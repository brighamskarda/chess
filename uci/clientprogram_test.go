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

package uci

import (
	"strings"
	"testing"
	"time"
)

func TestNewUciEngine_WriteRead(t *testing.T) {
	engine, err := newClientProgram(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Fatalf("could not start program: %v", err)
	}
	defer engine.Wait()
	defer engine.Kill()

	testString := "Hello World\n"
	_, err = engine.Write([]byte(testString))
	if err != nil {
		t.Fatalf("could not write to program: %v", err)
	}

	buf := make([]byte, 20)
	i, err := engine.Read(buf)
	if err != nil {
		t.Fatalf("could not read from program: %v", err)
	}

	if string(buf[:i]) != testString {
		t.Fatalf("did not get correct response: expected %q, got %q", testString, buf[:i])
	}
}

func TestNewUciEngine_WriteStderr(t *testing.T) {
	stderr := strings.Builder{}
	engine, err := newClientProgram(dummyBinaryPath, ClientSettings{Stderr: &stderr})
	if err != nil {
		t.Fatalf("could not start program: %v", err)
	}
	defer engine.Wait()
	defer engine.Kill()

	testString := "Hello World\n"
	_, err = engine.Write([]byte(testString))
	if err != nil {
		t.Fatalf("could not write to program: %v", err)
	}

	time.Sleep(250 * time.Millisecond)

	if result := stderr.String(); result != testString {
		t.Fatalf("did not get correct response: expected %q, got %q", testString, result)
	}
}

func TestNewUciEngine_Terminate(t *testing.T) {
	engine, err := newClientProgram(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Fatalf("could not start program: %v", err)
	}

	if err := engine.Terminate(); err != nil {
		t.Errorf("terminate failed: %v", err)
	}

	engine.Wait()

	if engine.cmd.ProcessState == nil {
		t.Error("process did not terminate")
	}
}

func TestNewUciEngine_Kill(t *testing.T) {
	engine, err := newClientProgram(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Fatalf("could not start program: %v", err)
	}

	if err := engine.Kill(); err != nil {
		t.Errorf("kill failed: %v", err)
	}

	engine.Wait()
	if engine.cmd.ProcessState == nil {
		t.Error("process did not die")
	}
}

func TestNewUciEngine_WaitBlocks(t *testing.T) {
	engine, err := newClientProgram(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Fatalf("could not start program: %v", err)
	}
	defer engine.Kill()

	done := make(chan struct{})

	go func() {
		engine.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Wait() returned before Kill() was called")
	case <-time.After(100 * time.Millisecond):
	}

	if err := engine.Kill(); err != nil {
		t.Fatalf("failed to kill engine: %v", err)
	}

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait() did not return after Kill()")
	}
}
