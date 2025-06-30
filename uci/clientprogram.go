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

//go:build !windows && !unix

package uci

import (
	"fmt"
	"io"
	"os/exec"
)

// clientProgram should be a uci compatible chess engine that is already running. [Client] will send and receive commands from it via Read and Write.
// When finished with a clientProgram be sure to call Wait() to free any resources.
type clientProgram struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	cmd    *exec.Cmd
}

// newClientProgram starts program with the specified settings. program should be a path to a uci compatible chess engine. Stdout will be directed to towards Read(), and Stdin will receive from Write(). Stderr will be directed towards ReadErr(). If the program path is invalid, or unable to successfully run then an error is returned.
//
//   - On windows the program will be started in a job object. This reduces the likelihood of orphaning child processes when calling Kill()
//   - Likewise on unix-like operating systems (linux, apple, etc.) the program is started in a process group to help prevent orphaned children.
//   - On other operating systems Kill() just ends the parent process.
func newClientProgram(program string, settings ClientSettings) (*clientProgram, error) {
	cmd := exec.Command(program, settings.Args...)
	cmd.Env = settings.Env
	cmd.Dir = settings.WorkDir
	cp := clientProgram{cmd: cmd}
	var err error
	cp.stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not start new uci engine: %w", err)
	}
	cp.stdin, err = cmd.StdinPipe()
	if err != nil {
		cp.stdout.Close()
		return nil, fmt.Errorf("could not start new uci engine: %w", err)
	}
	cp.stderr, err = cmd.StderrPipe()
	if err != nil {
		cp.stdout.Close()
		cp.stderr.Close()
		return nil, fmt.Errorf("could not start new uci engine: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		cp.stdin.Close()
		cp.stdout.Close()
		cp.stderr.Close()
		return nil, fmt.Errorf("could not start new uci engine: %w", err)
	}

	return &cp, nil
}

// Terminate asks the program to gracefully exit. Returns an error if the request was not sent successfully.
// Call [clientProgram.Wait] after this function to clean up resources.
//
// On windows this is implemented by sending a CTRL_BREAK_EVENT to the process. If you are calling this function from a GUI in windows it may be necessary to attach your GUI to a [console].
//
// On unix-like operating systems SIGTERM is sent to the process group.
//
// On other systems this function does nothing.
//
// [console]: https://learn.microsoft.com/en-us/windows/console/attachconsole
func (cp *clientProgram) Terminate() error {
	return nil
}

// Kill immediately stops the program. Returns an error if the request was not sent successfully
// Call wait after this function to clean up resources.
func (cp *clientProgram) Kill() error {
	return cp.cmd.Process.Kill()
}

// Wait waits for the program to finish and cleans up associated resources.
// It may return an error if the program did not exit successfully (like returning exit code 1), or there were io errors.
// Wait should only be called once. Ensure the io.WriteCloser is closed, and both readers are flushed to prevent blocking.
func (cp *clientProgram) Wait() error {
	return cp.cmd.Wait()
}

// Read reads from the program's stdout.
func (cp *clientProgram) Read(p []byte) (int, error) {
	return cp.stdout.Read(p)
}

// Write writes to the program's stdin.
func (cp *clientProgram) Write(p []byte) (int, error) {
	return cp.stdin.Write(p)
}

// ReadErr reads from the program's stderr.
func (cp *clientProgram) ReadErr(p []byte) (int, error) {
	return cp.stderr.Read(p)
}
