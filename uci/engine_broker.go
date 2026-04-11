// Copyright (C) 2026 Brigham Skarda

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
)

// inputLineBufferSize is the number of lines that should be buffered from the [UciEngineBroker.Input]. Having a buffer can increase throughput.
const inputLineBufferSize = 16

// inputCommandBufferSize is the number of parsed commands that should be buffered while waiting for another command to execute. Similar to [inputLineBufferSize], having a buffer here and increase throughput.
const inputCommandBufferSize = 4

// outputCommandBufferSize is the number of commands to buffer for the output stream. Having a buffer can increase command throughput.
const outputCommandBufferSize = 16

// errorReportingLocation is where users should report internal library errors that are printed out.
const errorReportingLocation = "https://github.com/brighamskarda/chess/issues"

// UciEngineBroker will automatically handle all UCI communication for a [ChessEngine]. Having a broker significantly simplifies the development of a chess engine as the engine can be developed without worrying about the complexities of the Universal Chess Interface (UCI).
//
// ChessEngines should do all of their logging through info strings, but this broker
type UciEngineBroker struct {
	// Input is the source from which the UciEngineBroker will read commands from the UCI client.
	//
	// In most cases this should be [os.Stdin]
	Input io.ReadCloser
	// Output is the destination to which the engine commands will be sent to the UCI client.
	//
	// In most cases this should be [os.Stdout]
	Output io.WriteCloser
	// Error is where the UciEngineBroker will log errors that shouldn't be sent to the client. In most cases it should be pretty empty as errors will only occur if the client is sending invalid or malformed UCI commands.
	//
	// In most cases this should be [os.Stderr]
	Error io.Writer
	// errorLocker ensures only one go routine writes to Error at a time.
	errorLocker sync.Mutex
	// Engine is where the actual logic of the chess engine is contained. When the broker receives commands from the UCI client, it will translate those commands into the appropriate calls to the engine.
	Engine ChessEngine
	// ctx indicates if the engine should keep running, or if it should shutdown.
	ctx context.Context
	// ctxCancel should be called when the program needs to shutdown. It will close ctx resulting in all parts of the uci broker and engine to shutdown.
	ctxCancel context.CancelFunc
	// outputCommands is the queue of commands to be sent back to the client.
	outputCommands chan<- engineToClientCmd
}

// Starts the UciEngineBroker. This function will not return until the UCI client requests the engine to shutdown. Until then it will read stdin for commands from the UCI client, and it will send command from the engine back the the UCI client via stdout.
func (broker *UciEngineBroker) Start() {
	// TODO set up listeners for shutdown events like SIGTERM

	inputLines := make(chan []byte, inputLineBufferSize)
	inputCommands := make(chan clientToEngineCmd, inputCommandBufferSize)
	outputCommands := make(chan engineToClientCmd, outputCommandBufferSize)
	broker.outputCommands = outputCommands
	broker.ctx, broker.ctxCancel = context.WithCancel(context.Background())

	go broker.inputParserLoop(broker.Input, inputLines)
	go broker.clientCommandParserLoop(inputLines, inputCommands)
	go broker.commandOutputLoop(broker.Output, outputCommands)
	broker.executeCommands(inputCommands)
}

// printError wraps writes to Error in a mutex lock in case a non-concurrent writer is provided.
func (broker *UciEngineBroker) printError(err string) {
	broker.errorLocker.Lock()
	fmt.Fprintln(broker.Error, "UciEngineBroker error:", err)
	broker.errorLocker.Unlock()
}

// inputParserLoop reads lines of text from the input, and puts them into lines. The function closes the channel and reader, then exits when input returns an error or when broker.ctx is cancelled. The rest of the data will be sent on an error, even if it doesn't end in new line.
func (broker *UciEngineBroker) inputParserLoop(input io.ReadCloser, lines chan<- []byte) {
	defer input.Close()
	defer close(lines)

	bufReader := bufio.NewReader(input)

	line := make(chan []byte)

	go func() {
		defer close(line)

		l, err := bufReader.ReadBytes('\n')
		for ; err == nil; l, err = bufReader.ReadBytes('\n') {
			select {
			case <-broker.ctx.Done():
				return
			case line <- l:
			}
		}

		if len(l) > 0 {
			select {
			case <-broker.ctx.Done():
				return
			case line <- l:
			}
		}
	}()

Loop:
	for {
		select {
		case <-broker.ctx.Done():
			break Loop
		case value, ok := <-line:
			if ok {
				lines <- value
			} else {
				break Loop
			}
		}
	}
}

// clientCommandParserLoop reads from lines and parses them as [clientToEngineCmd]s. These commands are then sent to the provided output channel. Onces all lines are read and the channel is closed, the commands channel is closed.
//
// As per the UCI specification unknown commands are simply ignored and new commands will continue to be parsed.
func (broker *UciEngineBroker) clientCommandParserLoop(lines <-chan []byte, commands chan<- clientToEngineCmd) {
	defer close(commands)
Loop:
	for {
		select {
		case <-broker.ctx.Done():
			break Loop
		case value, ok := <-lines:
			if ok {
				cmd, err := parseClientToEngineCmd(value)
				if err != nil {
					broker.errorLocker.Lock()
					fmt.Fprintln(broker.Error, err)
					broker.errorLocker.Unlock()
					continue Loop
				}

				select {
				case <-broker.ctx.Done():
					break Loop
				case commands <- cmd:
				}
			} else {
				break Loop
			}
		}
	}
}

// commandOutputLoop marshals the engineToClientCommands and outputs them to the writer until the channel is closed.
func (broker *UciEngineBroker) commandOutputLoop(output io.WriteCloser, commands <-chan engineToClientCmd) {
	defer output.Close()
	defer broker.ctxCancel()

Loop:
	for {
		select {
		case <-broker.ctx.Done():
			break Loop
		case cmd, ok := <-commands:
			if ok {
				text, err := cmd.marshalText()
				if err != nil {
					broker.printError(fmt.Sprintf("engine to client command encountered an error while marshaling: %q This indicates an internal library error. Please report such errors to %v", err, errorReportingLocation))
					continue Loop
				}

				_, err = output.Write(text)
				if err != nil {
					broker.printError(fmt.Sprintf("error encountered while trying to write to output, closing output writer and shutting down: %v", err))
					broker.ctxCancel()
				}
			} else {
				break Loop
			}
		}
	}
}

// executeCommands is the core of the UciEngineBroker. It takes all the commands that are being parsed and translates them into function calls to the engine.
func (broker *UciEngineBroker) executeCommands(inputCommands <-chan clientToEngineCmd) {
	defer broker.ctxCancel()
	defer close(broker.outputCommands)

	for {
		select {
		case <-broker.ctx.Done():
			broker.Engine.Quit()
			return
		case cmd, ok := <-inputCommands:
			if ok {
				broker.doCommand(cmd)
			} else {
				broker.ctxCancel()
			}
		}
	}
}

func (broker *UciEngineBroker) doCommand(cmd clientToEngineCmd) {
	switch cmd.(type) {
	case *uciCmd:
		broker.handleUciCommand()
	default:
		broker.printError(fmt.Sprintf("command with unknown type %T received in UciEngineBroker. This indicates an internal library error. Please report such errors to %v", cmd, errorReportingLocation))
	}
}

func (broker *UciEngineBroker) handleUciCommand() {
	broker.Engine.Initialize()

	// send out the engine name
	broker.outputCommands <- &idCmd{
		isAuthor: false,
		id:       broker.Engine.Name(),
	}

	// send out the engine author
	broker.outputCommands <- &idCmd{
		isAuthor: true,
		id:       broker.Engine.Author(),
	}

	// send out the engine options
	for _, opt := range broker.Engine.Options() {
		broker.outputCommands <- opt
	}

	// send uciok
	broker.outputCommands <- &uciokCmd{}

}
