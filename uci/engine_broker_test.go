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
	"io"
	"os"
	"testing"
	"time"
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
		Error:  os.Stderr,
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

// maliciousEngine stalls during a command execution
type maliciousEngine struct {
	mockEngine
}

// Imagine this is called by executeCommands when a "uci" command is received
func (e *maliciousEngine) Initialize() {
	// Simulated long-running task
	time.Sleep(time.Second * 5)
}

// TestEngineShutsDownWhenStdinIsClose ensures that best practice is followed by making sure the engine shuts down when it detects stdin has been closed.
func TestEngineShutsDownWhenStdinIsClosed(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	startReturned := make(chan struct{})
	go func() {
		broker.Start()
		startReturned <- struct{}{}
	}()

	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Errorf("problem writing to stdin: %v", err)
	}
	// Read something to ensure engine initializes.
	stdoutR.Read(make([]byte, 1))

	err = stdinW.Close()
	if err != nil {
		t.Errorf("problem closing stdin: %v", err)
	}

	// Finish this test by ensuring startReturned receives a value within a second of the stdinW being closed.
	select {
	case <-startReturned:
		// Success: The broker.Start() returned as expected after stdin was closed.
		if broker.Engine.(*mockEngine).quit != 1 {
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
		broker.Start()
		startReturned <- struct{}{}
	}()

	// Close the reader end (simulating the client/GUI closing the pipe)
	stdoutR.Close()

	// Send a command that forces the engine to write an output
	_, err := stdinW.WriteString("uci\n")
	if err != nil {
		t.Fatalf("failed to write to stdin: %v", err)
	}

	// The commandOutputLoop should encounter a write error and call ctxCancel()
	select {
	case <-startReturned:
		// Success: broker.Start() returned because the broken output pipe triggered shutdown
		if broker.Engine.(*mockEngine).quit != 1 {
			t.Errorf("broker did not call Quit on the engine exactly 1 time")
		}
	case <-time.After(time.Second * 2):
		t.Errorf("broker did not shut down after the output pipe was closed")
	}
}

// startNewUciBroker starts a new uci engine broker and returns the stdin and stdout pipes that the client would see.
func startNewUciBroker(t *testing.T) (stdinW *os.File, stdoutR *os.File) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)
	go broker.Start()
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

// TestEngineInitialization tests that the broker outputs id, author, options, and uciok after receiving a uci command.
func TestEngineDebugMode(t *testing.T) {
	stdinR, stdinW := makeOsPipe(t)
	stdoutR, stdoutW := makeOsPipe(t)
	broker := makeUciEngineBroker(stdinR, stdoutW)

	go broker.Start()

	_, err := stdinW.WriteString("debug on\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	_, err = stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}
	// wait for readyok to indicate the commands have been processed
	stdoutR.Read(make([]byte, 256))

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

// TestEngineInitialization tests that the broker outputs id, author, options, and uciok after receiving a uci command.
func TestIsReady(t *testing.T) {
	stdinW, stdoutR := startNewUciBroker(t)

	_, err := stdinW.WriteString("isready\n")
	if err != nil {
		t.Fatalf("error writing to stdin: %v", err)
	}

	output := bufio.NewReader(stdoutR)

	testOutput(output, "readyok\n", t)

}
