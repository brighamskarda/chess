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
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// inputCommandBufferSize is the number of parsed commands that should be buffered while waiting for another command to execute. Similar to [inputLineBufferSize], having a buffer here and increase throughput.
const inputCommandBufferSize = 16

// outputCommandBufferSize is the number of commands to buffer for the output stream. Having a buffer can increase command throughput.
const outputCommandBufferSize = 16

// errorReportingLocation is where users should report internal library errors that are printed out.
const errorReportingLocation = "https://github.com/brighamskarda/chess/issues"

// UciEngineBroker will automatically handle all UCI communication for a [ChessEngine]. Having a broker significantly simplifies the development of a chess engine as the engine can be developed without worrying about the complexities of the Universal Chess Interface (UCI).
//
// ChessEngines should do all of their logging through info strings, but this broker
type UciEngineBroker struct {
	// Engine is where the actual logic of the chess engine is contained. When the broker receives commands from the UCI client, it will translate those commands into the appropriate calls to the engine.
	Engine ChessEngine
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

	inputCommands := make(chan clientToEngineCmd, inputCommandBufferSize)
	outputCommands := make(chan engineToClientCmd, outputCommandBufferSize)
	broker.outputCommands = outputCommands
	broker.ctx, broker.ctxCancel = context.WithCancel(context.Background())

	go broker.signalListener()
	go broker.commandInputLoop(inputCommands)
	go broker.commandOutputLoop(outputCommands)
	broker.executeCommands(inputCommands)
}

// printError wraps writes to Error in a mutex lock in case a non-concurrent writer is provided.
func (broker *UciEngineBroker) printError(err string) {
	broker.errorLocker.Lock()
	fmt.Fprintln(broker.Error, "UciEngineBroker error:", err)
	broker.errorLocker.Unlock()
}

// signalListener ensures that the uci engine broker context is cancelled when a sigterm or sigint is received. Should work on windows and linux.
func (broker *UciEngineBroker) signalListener() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	broker.ctxCancel()
}

// readLines reads lines from the brokers input and calls broker.ctxCancel() if there is an error reading. It is common practice for UCI chess engine to shutdown once stdin has been closed.
func (broker *UciEngineBroker) readLines(line chan<- []byte) {
	defer close(line)
	bufReader := bufio.NewReader(broker.Input)

	l, err := bufReader.ReadBytes('\n')
	for ; err == nil; l, err = bufReader.ReadBytes('\n') {
		select {
		case <-broker.ctx.Done():
			return
		case line <- l:
		}
	}

	broker.ctxCancel()
}

// commandInputLoop reads from lines and parses them as [clientToEngineCmd]s. These commands are then sent to the provided output channel. Onces all lines are read and the channel is closed, the commands channel is closed.
//
// As per the UCI specification unknown commands are simply ignored and new commands will continue to be parsed.
func (broker *UciEngineBroker) commandInputLoop(commands chan<- clientToEngineCmd) {
	defer close(commands)
	defer broker.Input.Close()

	line := make(chan []byte)
	go broker.readLines(line)

Loop:
	for {
		select {
		case <-broker.ctx.Done():
			break Loop
		case value, ok := <-line:
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

// commandOutputLoop marshals the engineToClientCommands and outputs them to the writer until the channel is closed. If there is an error when writing to the output the broker context is cancelled since this is a pretty good sign that something has gone wrong and the engine should shut down.
func (broker *UciEngineBroker) commandOutputLoop(commands <-chan engineToClientCmd) {
	defer cleanOutChannel(commands)
	defer broker.ctxCancel()
	defer broker.Output.Close()

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

				_, err = broker.Output.Write(text)
				if err != nil {
					broker.printError(fmt.Sprintf("error encountered while trying to write to output, closing output writer and shutting down: %v", err))
					break Loop
				}
			} else {
				break Loop
			}
		}
	}
}

// cleanOutChannel ensures that a channel is read from until it is closed to prevent deadlocks
func cleanOutChannel[T any](ch <-chan T) {
	for {
		_, ok := <-ch
		if !ok {
			break
		}
	}
}

// executeCommands is the core of the UciEngineBroker. It takes all the commands that are being parsed and translates them into function calls to the engine.
func (broker *UciEngineBroker) executeCommands(inputCommands <-chan clientToEngineCmd) {
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

// doCommand calls different command handlers based on the underlying type of the cmd.
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
