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
	"io"
	"log/slog"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/brighamskarda/chess/v2"
)

func makeOsPipe(t *testing.T) (*os.File, *os.File) {
	reader, writer, err := os.Pipe()
	t.Cleanup(func() {
		reader.Close()
		writer.Close()
	})
	if err != nil {
		t.Fatalf("failed to make pipe: %v", err)
	}
	return reader, writer
}

func makeUciEngineBroker(reader io.ReadCloser, writer io.WriteCloser) *UciEngineBroker {
	return &UciEngineBroker{
		Input:  reader,
		Output: writer,
		Log:    slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})),
		Engine: &mockEngine{},
	}
}

// DiscardWriteCloser wraps io.Writer to add a Close() method
type DiscardWriteCloser struct {
	io.Writer
}

// Close implements io.WriteCloser by doing nothing
func (d DiscardWriteCloser) Close() error {
	return nil
}

// TestEngineShutsDownWhenStdinIsClose ensures that best practice is followed
// by making sure the engine shuts down when it detects stdin has been closed.
func TestEngineShutsDownWhenStdinIsClosed(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	_, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	startReturned := make(chan struct{})
	go func() {
		broker.Start(t.Context())
		startReturned <- struct{}{}
	}()

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}

	time.Sleep(time.Second)

	err = stdinW.Close()
	if err != nil {
		t.Errorf("problem closing stdin: %v", err)
	}

	// Finish this test by ensuring startReturned receives a value within a second of the stdinW being closed.
	select {
	case <-startReturned:
		// Success: The broker.Start() returned as expected after stdin was closed.
		if broker.Engine.(*mockEngine).quit != 1 && broker.Engine.(*mockEngine).initialize == 1 {
			t.Errorf("broker did not call Quit on the engine exactly 1 time, was called %v times", broker.Engine.(*mockEngine).quit)
		}
	case <-time.After(time.Second):
		// Failure: The broker did not shut down within the 1-second timeout.
		t.Errorf("engine did not shut down within 1 second of stdin being closed")
	}
}

func TestEngineShutsDownWhenOutputIsClosed(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	startReturned := make(chan struct{})
	go func() {
		broker.Start(t.Context())
		startReturned <- struct{}{}
	}()

	// Close the reader end (simulating the client/GUI closing the pipe)
	stdoutR.Close()

	// Send a command that forces the engine to write an output
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("failed to write to stdout: %v", err)
	}

	// The commandOutputLoop should encounter a write error and call ctxCancel()
	select {
	case <-startReturned:
		// Success: broker.Start() returned because the broken output pipe triggered shutdown
		if broker.Engine.(*mockEngine).quit != 1 && broker.Engine.(*mockEngine).initialize == 1 {
			t.Errorf("broker did not call Quit on the engine exactly 1 time")
		}
	case <-time.After(time.Second * 2):
		t.Errorf("broker did not shut down after the output pipe was closed")
	}
}

func TestEngineShutsDownWhenContextCancelled(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	_, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	ctx, cancel := context.WithCancel(t.Context())

	startReturned := make(chan struct{})
	go func() {
		broker.Start(ctx)
		startReturned <- struct{}{}
	}()

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}

	cancel()

	// Finish this test by ensuring startReturned receives a value within a second of the context being cancelled.
	select {
	case <-startReturned:
		// Success: The broker.Start() returned as expected after stdin was closed.
		if broker.Engine.(*mockEngine).quit != 1 && broker.Engine.(*mockEngine).initialize == 1 {
			t.Errorf("broker did not call Quit on the engine exactly 1 time, was called %v times", broker.Engine.(*mockEngine).quit)
		}
	case <-time.After(time.Second):
		// Failure: The broker did not shut down within the 1-second timeout.
		t.Errorf("engine did not shut down within 1 second of context cancellation")
	}
}

func TestEngineShutsDownWhenQuit(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	startReturned := make(chan struct{})
	go func() {
		broker.Start(t.Context())
		startReturned <- struct{}{}
	}()

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}

	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "uciok\n" {
			break
		}
	}

	_, err = stdinW.WriteString("quit\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}

	// Finish this test by ensuring startReturned receives a value within a second of the context being cancelled.
	select {
	case <-startReturned:
		// Success: The broker.Start() returned as expected after stdin was closed.
		if broker.Engine.(*mockEngine).quit != 1 && broker.Engine.(*mockEngine).initialize == 1 {
			t.Errorf("broker did not call Quit on the engine exactly 1 time, was called %v times", broker.Engine.(*mockEngine).quit)
		}
	case <-time.After(time.Second):
		// Failure: The broker did not shut down within the 1-second timeout.
		t.Errorf("engine did not shut down within 1 second of quit being called")
	}
}

func TestNoOtherCommandsParsedBeforeInitialize(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())

	_, err := stdinW.WriteString("debug on\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("quit\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("uci\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}

	output := bufio.NewReader(stdoutR)
	testOutput(output, "id name mockEngine v0.1\n", t)

	if broker.Engine.(*mockEngine).debug != 0 {
		t.Errorf("debug was called before initialize")
	}
	if broker.Engine.(*mockEngine).quit != 0 {
		t.Errorf("quit was called before initialize")
	}
}

// startNewUciBroker starts a new uci engine broker and returns the stdin and stdout pipes that the client would see.
func startNewUciBroker(t *testing.T) (stdinW *os.File, stdoutR *os.File) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)
	go broker.Start(t.Context())
	return
}

// testOutput reads a line from the reader and compares it to an expected output.
func testOutput(reader *bufio.Reader, expected string, t *testing.T) {
	got, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("got err %q, expected %q", err, expected)
		return
	}

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
		return
	}

	t.Logf("Received: %q", expected)
}

// TestEngineInitialization tests that the broker outputs id, author, options, and uciok after receiving a uci command.
func TestEngineInitialization(t *testing.T) {
	stdinW, stdoutR := startNewUciBroker(t)

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	output := bufio.NewReader(stdoutR)

	testOutput(output, "id name mockEngine v0.1\n", t)
	testOutput(output, "id author Brigham Skarda\n", t)
	testOutput(output, "option name checkOpt type check default true\n", t)
	testOutput(output, "option name spinOpt type spin default 3 min 1 max 5\n", t)
	testOutput(output, "option name comboOpt type combo default one var one var two var three\n", t)
	testOutput(output, "option name stringOpt type string default sss\n", t)
	testOutput(output, "option name buttonOpt type button\n", t)
	testOutput(output, "uciok\n", t)
}

func TestEngineQuitsOnInterrupt(t *testing.T) {
	t.Log("Unfortunately there is no great way to test signals on windows, this is a test that needs to be done by hand.")
}

func TestEngineDebugMode(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("debug on\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "readyok\n" {
			break
		}
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).debug
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Debug to be called %v times, but was called %v times", expectedVal, actualVal)
	}
	expectedBool := true
	actualBool := broker.Engine.(*mockEngine).debugState
	if expectedBool != actualBool {
		t.Errorf("expected Engine.DebugState to be %v , but was %v", expectedBool, actualBool)
	}

	_, err = stdinW.WriteString("debug off\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	stdoutR.Read(make([]byte, 256))

	expectedVal = 2
	actualVal = broker.Engine.(*mockEngine).debug
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Debug to be called %v times, but was called %v times", expectedVal, actualVal)
	}
	expectedBool = false
	actualBool = broker.Engine.(*mockEngine).debugState
	if expectedBool != actualBool {
		t.Errorf("expected Engine.DebugState to be %v , but was %v", expectedBool, actualBool)
	}
}

func TestIsReady(t *testing.T) {
	stdinW, stdoutR := startNewUciBroker(t)

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for registration ok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "registration ok\n" {
			break
		}
	}

	testOutput(out, "readyok\n", t)
}

func TestSetOption(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("setoption name Nullmove value true\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for readyok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "readyok\n" {
			break
		}
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).setOption
	if expectedVal != actualVal {
		t.Errorf("expected Engine.SetOption to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

func TestCopyProtection(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for uciok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "uciok\n" {
			break
		}
	}

	testOutput(out, "copyprotection checking\n", t)
	testOutput(out, "copyprotection ok\n", t)

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).copyProtection
	if expectedVal != actualVal {
		t.Errorf("expected Engine.CopyProtection to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

type mockEngineBadCopyProtect struct {
	mockEngine
}

func (engine *mockEngineBadCopyProtect) CopyProtection() bool {
	engine.copyProtection++
	return false
}

func TestCopyProtectionBad(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := &UciEngineBroker{
		Input:  stdinR,
		Output: stdoutW,
		Log:    slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})),
		Engine: &mockEngineBadCopyProtect{},
	}

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for uciok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "uciok\n" {
			break
		}
	}

	testOutput(out, "copyprotection checking\n", t)
	testOutput(out, "copyprotection error\n", t)

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngineBadCopyProtect).copyProtection
	if expectedVal != actualVal {
		t.Errorf("expected Engine.CopyProtection to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

func TestRegistration(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for copyprotection ok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "copyprotection ok\n" {
			break
		}
	}

	testOutput(out, "registration checking\n", t)
	testOutput(out, "registration ok\n", t)

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).register
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Register to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

type mockEngineBadRegister struct {
	mockEngine
}

func (engine *mockEngineBadRegister) Register(cmd *RegisterCmd) bool {
	engine.register++
	if cmd == nil {
		return false
	} else {
		return true
	}
}

func TestRegistrationBad(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := &UciEngineBroker{
		Input:  stdinR,
		Output: stdoutW,
		Log:    slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})),
		Engine: &mockEngineBadRegister{},
	}

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	// wait for copyprotection ok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "copyprotection ok\n" {
			break
		}
	}

	testOutput(out, "registration checking\n", t)
	testOutput(out, "registration error\n", t)

	_, err = stdinW.WriteString("register name bs code 123\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	testOutput(out, "registration checking\n", t)
	testOutput(out, "registration ok\n", t)

	expectedVal := 2
	actualVal := broker.Engine.(*mockEngineBadRegister).register
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Register to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

func TestSetStartPosition(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("position startpos moves e2e4 d7d5\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "readyok\n" {
			break
		}
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).setPosition
	if expectedVal != actualVal {
		t.Errorf("expected Engine.SetPosition to be called %v times, but was called %v times", expectedVal, actualVal)
	}

	expectedStr := chess.DefaultFEN
	actualStr, _ := broker.Engine.(*mockEngine).position.MarshalText()
	if expectedStr != string(actualStr) {
		t.Errorf("expected position to be %q, but was %q", expectedStr, actualStr)
	}

	expectedMoves := []chess.Move{{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType},
		{FromSquare: chess.D7, ToSquare: chess.D5, Promotion: chess.NoPieceType}}
	actualMoves := broker.Engine.(*mockEngine).moveHistory
	if !slices.Equal(expectedMoves, actualMoves) {
		t.Errorf("expected move history to be %v, but was %v", expectedMoves, actualMoves)
	}
}

func TestSetPosition(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	testPos := "1r6/5pp1/R1R4p/1r1pP3/2pkQPP1/7P/1P6/2K5 w - - 0 41"

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("position fen " + testPos + " moves e2e4 d7d5\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "readyok\n" {
			break
		}
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).setPosition
	if expectedVal != actualVal {
		t.Errorf("expected Engine.SetPosition to be called %v times, but was called %v times", expectedVal, actualVal)
	}

	expectedStr := testPos
	actualStr, _ := broker.Engine.(*mockEngine).position.MarshalText()
	if expectedStr != string(actualStr) {
		t.Errorf("expected position to be %q, but was %q", expectedStr, actualStr)
	}

	expectedMoves := []chess.Move{{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType},
		{FromSquare: chess.D7, ToSquare: chess.D5, Promotion: chess.NoPieceType}}
	actualMoves := broker.Engine.(*mockEngine).moveHistory
	if !slices.Equal(expectedMoves, actualMoves) {
		t.Errorf("expected move history to be %v, but was %v", expectedMoves, actualMoves)
	}
}

func TestUciNewGame(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("ucinewgame\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	out := bufio.NewReader(stdoutR)
	for {
		text, _ := out.ReadString('\n')
		if text == "readyok\n" {
			break
		}
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).newGame
	if expectedVal != actualVal {
		t.Errorf("expected Engine.NewGame to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}

func TestEvaluate(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("go mate 3\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	readyChan := make(chan struct{})
	out := bufio.NewReader(stdoutR)

	// wait for readyok to indicate the commands have been processed
	go func() {
		for {
			text, _ := out.ReadString('\n')
			if text == "readyok\n" {
				break
			}
		}
		readyChan <- struct{}{}
	}()

	select {
	case <-readyChan:
		break
	case <-time.After(mockEngineEvaluateTime / 2):
		t.Errorf("broker is blocking on go command, it shouldn't")
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).evaluate
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Evaluate to be called %v times, but was called %v times", expectedVal, actualVal)
	}

	testOutput(out, "bestmove e2e4 ponder d7d5\n", t)
}

func TestEvaluateThenStop(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start(t.Context())
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("go mate 3\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("stop\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	readyChan := make(chan struct{})
	out := bufio.NewReader(stdoutR)

	// wait for bestmove
	go func() {
		for {
			text, _ := out.ReadString('\n')
			if text == "bestmove e2e4 ponder d7d5\n" {
				break
			}
		}
		readyChan <- struct{}{}
	}()

	select {
	case <-readyChan:
		break
	case <-time.After(mockEngineEvaluateTime / 2):
		t.Errorf("stop did not stop the engine, it shouldn't")
	}

	expectedVal := 1
	actualVal := broker.Engine.(*mockEngine).stop
	if expectedVal != actualVal {
		t.Errorf("expected Engine.Stop to be called %v times, but was called %v times", expectedVal, actualVal)
	}
}
