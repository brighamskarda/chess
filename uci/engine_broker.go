// Copyright (C) 2026 Brigham Skarda
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
)

// errorReportingLocation is where users should report internal library errors that are printed out.
const errorReportingLocation = "https://github.com/brighamskarda/chess/issues"

// UciEngineBroker will automatically handle all UCI communication for a [ChessEngine].
//
// Having a broker significantly simplifies the development of a chess engine
// as the engine can be developed without worrying about the complexities of the Universal Chess Interface (UCI).
type UciEngineBroker struct {
	// Engine is where the actual logic of the chess engine is contained.
	//
	// When the broker receives commands from the UCI client,
	// it will translate those commands into the appropriate calls to the engine.
	Engine ChessEngine

	// Input is the source from which the UciEngineBroker will read commands from the UCI client.
	//
	// In most cases this should be [os.Stdin]
	Input io.Reader

	// Output is the destination to which the engine commands will be sent to the UCI client.
	//
	// In most cases this should be [os.Stdout].
	Output io.Writer

	// outputLocker ensures only one go routine writes to Error at a time.
	outputLocker sync.Mutex

	// Log is where the broker can output information outside of normal engine communication.
	//
	// Log is an optional field.
	// But Logs provide information on how the broker is running,
	// and report errors that are encountered.
	// It is recommended that a logger with at least a level of [slog.Error] is provided.
	// Here is a simple logger you can provided that logs errors to [os.Stderr].
	//
	//		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	//
	// Log should not write to the same location as [UciEngineBroker.Output].
	// [os.Stderr] and external files are great places to log to.
	Log *slog.Logger

	// ctx indicates if the engine should keep running, or if it should shutdown.
	ctx context.Context

	// ctxCancel should be called when the program needs to shutdown.
	//
	// It will close ctx resulting in all parts of the uci broker and engine to shutdown.
	ctxCancel context.CancelCauseFunc

	// quitWg makes sure the engine has shutdown before the broker stops.
	quitWg sync.WaitGroup

	// DisableSignalHandling removes the default signaling functionality.
	//
	// By default the broker will automatically shutdown
	// when it receives the appropriate signal from the OS.
	// This functionality always desirable so this flag is provided.
	DisableSignalHandling bool

	// isInitialized indicates if Initialize() has been called on the engine yet.
	isInitialized bool
}

// Start the UciEngineBroker.
//
// This function will not return until
// the UCI client requests the engine to shutdown,
// the context is cancelled,
// or there is an error.
// Until then, it will read stdin for commands from the UCI client,
// and it will send commands from the engine back to the UCI client via stdout.
//
// The provided context will also be passed into the to [UciEngineBroker.Log] whenever it is called.
//
// Start returns an error if the broker is stopped for any reason besides
// the context being cancelled,
// or the quit command being received from the engine.
//
// A UciEngineBroker should only be started once.
// To restart a chess engine make a new UciEngineBroker
func (broker *UciEngineBroker) Start(ctx context.Context) error {
	// Make sure the error logger isn't nil.
	if broker.Log == nil {
		broker.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	broker.Log.InfoContext(ctx, "starting UCI engine broker")

	// setup cancellation context.
	broker.ctx, broker.ctxCancel = context.WithCancelCause(ctx)
	// setup signal handling
	if !broker.DisableSignalHandling {
		broker.Log.DebugContext(broker.ctx, "setting up os.Signal handling")
		go broker.terminationListener()
	}

	// Start executing commands received from the client. Runs a loop until cancellation.
	broker.executeCommands()
	broker.quitWg.Wait()

	// See if we stopped executing commands due to an error.
	err := context.Cause(broker.ctx)
	if err != context.Canceled {
		return fmt.Errorf("UCI engine broker stopped with error: %w", err)
	}
	return nil
}

// sendCommand sends an engine command to the UCI client.
func (broker *UciEngineBroker) sendCommand(cmd engineToClientCmd) {
	text, err := cmd.marshalText()
	if err != nil {
		broker.Log.ErrorContext(broker.ctx, "failed to send command", slog.Any("command", cmd), slog.Any("error", err))
		return
	}

	broker.outputLocker.Lock()
	_, err = broker.Output.Write(text)
	broker.outputLocker.Unlock()
	if err != nil {
		broker.Log.ErrorContext(broker.ctx, "failed to send command", slog.Any("error", err))
		broker.Log.ErrorContext(broker.ctx, "shutting down UCI engine broker, got an error when trying to output commands")
		broker.ctxCancel(fmt.Errorf("output writer error: %w", err))
		return
	}

	broker.Log.DebugContext(broker.ctx, "sent command to client", slog.Any("text", text))
}

// terminationListener calls the cancel function when it receives an os request to shutdown.
//
// Specifically, it listens for os.Interrupt or syscall.SIGTERM.
func (broker *UciEngineBroker) terminationListener() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-ch:
		broker.ctxCancel(fmt.Errorf("os.Signal %v received", s))
	case <-broker.ctx.Done():
		// cleanup go routine.
	}
}

// readLines reads lines from the brokers input and calls broker.ctxCancel() if there is an error reading.
//
// It is common practice for UCI chess engine to shutdown once stdin has been closed.
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

	broker.ctxCancel(fmt.Errorf("input reader error: %w", err))
}

// executeCommands takes commands from the Input and translates them to functions calls to the engine.
func (broker *UciEngineBroker) executeCommands() {
	inputLines := make(chan []byte)
	go broker.readLines(inputLines)
Loop:
	for {
		if broker.ctx.Err() != nil {
			break Loop
		}

		select {
		case <-broker.ctx.Done():
			break Loop
		case line, ok := <-inputLines:
			if !ok {
				break Loop
			}
			cmd, err := parseClientToEngineCmd(line)
			if err != nil {
				broker.Log.ErrorContext(broker.ctx, "skipping client command", slog.Any("error", err))
			}
			broker.doCommand(cmd)
		}
	}
}

// doCommand calls different command handlers based on the underlying type of the cmd.
func (broker *UciEngineBroker) doCommand(cmd clientToEngineCmd) {
	if !broker.isInitialized && reflect.TypeOf(cmd) != reflect.TypeFor[*uciCmd]() {
		broker.Log.WarnContext(broker.ctx, "skipping invalid first command, expected uciCmd", slog.Any("got", reflect.TypeOf(cmd)))
		return
	}

	switch c := cmd.(type) {
	case *uciCmd:
		broker.handleUciCommand()
	case *debugCmd:
		broker.handleDebugCommand(c.on)
	case *isReadyCmd:
		broker.handleIsReadyCommand()
	case SetOptionCmd:
		broker.Engine.SetOption(c)
	case *quitCmd:
		broker.ctxCancel(nil)
	default:
		broker.Log.ErrorContext(broker.ctx, fmt.Sprintf("command with unknown type received, "+
			"this indicates an internal library error, please report such errors to %v", errorReportingLocation),
			slog.Any("unknownCommand", reflect.TypeOf(cmd)))

	}
}

func (broker *UciEngineBroker) handleUciCommand() {
	init := sync.OnceFunc(func() {
		broker.Engine.Initialize(func(infoCmd *InfoCmd) {
			broker.sendCommand(infoCmd)
		})
	},
	)

	// Increment the wait group so that the program doesn't exit until Quit has finished.
	broker.quitWg.Add(1)
	context.AfterFunc(broker.ctx, func() {
		defer broker.quitWg.Done()
		init() // make sure initialization is finished before calling quit.
		broker.Engine.Quit()
	})

	init()
	broker.isInitialized = true

	// send out the engine name
	broker.sendCommand(&idCmd{
		isAuthor: false,
		id:       broker.Engine.Name(),
	})

	// send out the engine author
	broker.sendCommand(&idCmd{
		isAuthor: true,
		id:       broker.Engine.Author(),
	})

	// send out the engine options
	for _, opt := range broker.Engine.Options() {
		broker.sendCommand(opt)
	}

	// send uciok
	broker.sendCommand(&uciokCmd{})
}

func (broker *UciEngineBroker) handleDebugCommand(debug bool) {
	broker.Engine.SetDebug(debug)
}

func (broker *UciEngineBroker) handleIsReadyCommand() {
	broker.sendCommand(&readyOkCmd{})
}
