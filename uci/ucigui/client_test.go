// Copyright (C) 2025 Brigham Skarda

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

package ucigui

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewClient_ErrorOnInvalidBinary(t *testing.T) {
	_, err := NewClient("./dkfdks.exe", ClientSettings{})
	if err == nil {
		t.Error("did not get error on invalid binary")
	}
}

func TestNewClient_NoErrorOnValidBinary(t *testing.T) {
	_, err := NewClient(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Errorf("%v", err)
	}
}

type clientProgramMock struct {
	stdinReader  *io.PipeReader
	stdinWriter  *io.PipeWriter
	stdoutReader *io.PipeReader
	stdoutWriter *io.PipeWriter
	stderrReader *io.PipeReader
	stderrWriter *io.PipeWriter
}

func (cp *clientProgramMock) Terminate() error {
	return nil
}
func (cp *clientProgramMock) Kill() error {
	cp.stdinReader.Close()
	cp.stdinWriter.Close()
	cp.stdoutReader.Close()
	cp.stdoutWriter.Close()
	cp.stderrReader.Close()
	cp.stderrWriter.Close()
	return nil
}
func (cp *clientProgramMock) Wait() error {
	return nil
}
func (cp *clientProgramMock) Write(p []byte) (int, error) {
	return cp.stdinWriter.Write(p)
}
func (cp *clientProgramMock) Read(p []byte) (int, error) {
	return cp.stdoutReader.Read(p)
}
func (cp *clientProgramMock) ReadErr(p []byte) (int, error) {
	return cp.stderrReader.Read(p)
}
func (cp *clientProgramMock) CloseStdin() error {
	return cp.stdinWriter.Close()
}

func newDummyClientProgram() *clientProgramMock {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()
	return &clientProgramMock{
		stdinReader:  stdinReader,
		stdinWriter:  stdinWriter,
		stdoutReader: stdoutReader,
		stdoutWriter: stdoutWriter,
		stderrReader: stderrReader,
		stderrWriter: stderrWriter,
	}
}

func TestClient_StdoutToLogger(t *testing.T) {
	dummyClient := newDummyClientProgram()
	defer dummyClient.Kill()
	testLogger := strings.Builder{}
	_, err := newClientFromClientProgram(dummyClient, ClientSettings{Logger: &testLogger})
	if err != nil {
		t.Fatalf("couldn't make client: %v", err)
	}

	dummyClient.stdoutWriter.Write([]byte("line1\n"))
	dummyClient.stdoutWriter.Write([]byte("line2"))
	dummyClient.stdoutWriter.Write([]byte("line3\n"))

	time.Sleep(10 * time.Millisecond)

	expected := "<<< line1\n<<< line2line3\n"
	loggerOutput := testLogger.String()
	if loggerOutput != expected {
		t.Errorf("logger output does not match: expected %v, got %v", expected, loggerOutput)
	}
}

func TestClient_StderrToLogger(t *testing.T) {
	dummyClient := newDummyClientProgram()
	defer dummyClient.Kill()
	testLogger := strings.Builder{}
	_, err := newClientFromClientProgram(dummyClient, ClientSettings{Logger: &testLogger})
	if err != nil {
		t.Fatalf("couldn't make client: %v", err)
	}

	dummyClient.stderrWriter.Write([]byte("line1\n"))
	dummyClient.stderrWriter.Write([]byte("line2"))
	dummyClient.stderrWriter.Write([]byte("line3\n"))

	time.Sleep(10 * time.Millisecond)

	expected := "!<! line1\n!<! line2line3\n"
	loggerOutput := testLogger.String()
	if loggerOutput != expected {
		t.Errorf("logger output does not match: expected %v, got %v", expected, loggerOutput)
	}
}

func TestClient_QuitOnRealProgram(t *testing.T) {
	cp, err := newClientProgram(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	c, err := newClientFromClientProgram(cp, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	c.Quit(time.Second, time.Second)
}

type clientProgramMock_DelayedWait struct {
	TerminateCalled bool
	KillCalled      bool
	TimeToDelay     time.Duration
}

func (cp *clientProgramMock_DelayedWait) Terminate() error {
	cp.TerminateCalled = true
	return nil
}
func (cp *clientProgramMock_DelayedWait) Kill() error {
	cp.KillCalled = true
	return nil
}
func (cp *clientProgramMock_DelayedWait) Wait() error {
	time.Sleep(cp.TimeToDelay)
	return nil
}
func (cp *clientProgramMock_DelayedWait) Write(p []byte) (int, error) {
	return 0, nil
}
func (cp *clientProgramMock_DelayedWait) Read(p []byte) (int, error) {
	return 0, nil
}
func (cp *clientProgramMock_DelayedWait) ReadErr(p []byte) (int, error) {
	return 0, nil
}
func (cp *clientProgramMock_DelayedWait) CloseStdin() error {
	return nil
}

func TestClient_QuitProcess(t *testing.T) {
	cp := &clientProgramMock_DelayedWait{TimeToDelay: time.Second}
	c, err := newClientFromClientProgram(cp, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = c.Quit(100*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !cp.TerminateCalled {
		t.Error("Terminate() not called")
	}

	if !cp.KillCalled {
		t.Error("Kill() not called")
	}
}

func TestClient_QuitSendsQuit(t *testing.T) {
	cp := newDummyClientProgram()
	c, err := newClientFromClientProgram(cp, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	buf := make([]byte, 20)
	var nRead int
	go func() {
		nRead, err = cp.stdinReader.Read(buf)
	}()

	c.Quit(100*time.Millisecond, 100*time.Millisecond)

	if err != nil {
		t.Errorf("%v", err)
	}

	expected := "quit\n"
	got := string(buf[:nRead])
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestClient_QuitLogs(t *testing.T) {
	cp := newDummyClientProgram()
	testLogger := strings.Builder{}
	c, err := newClientFromClientProgram(cp, ClientSettings{Logger: &testLogger})
	if err != nil {
		t.Fatalf("%v", err)
	}

	c.Quit(time.Millisecond, time.Millisecond)

	expected := ">>> quit\n"
	got := testLogger.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
