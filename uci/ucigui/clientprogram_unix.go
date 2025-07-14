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

//go:build unix

package ucigui

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type unixClientProgram struct {
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
func newClientProgram(program string, settings ClientSettings) (clientProgram, error) {
	cmd := exec.Command(program, settings.Args...)
	cmd.Env = settings.Env
	cmd.Dir = settings.WorkDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0}
	cp := unixClientProgram{cmd: cmd}
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
		cp.stdin.Close()
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

func (cp *unixClientProgram) Terminate() error {
	return syscall.Kill(-cp.cmd.Process.Pid, syscall.SIGTERM)
}

func (cp *unixClientProgram) Kill() error {
	return syscall.Kill(-cp.cmd.Process.Pid, syscall.SIGKILL)
}

func (cp *unixClientProgram) Wait() error {
	return cp.cmd.Wait()
}

func (cp *unixClientProgram) Read(p []byte) (int, error) {
	return cp.stdout.Read(p)
}

func (cp *unixClientProgram) Write(p []byte) (int, error) {
	return cp.stdin.Write(p)
}

func (cp *unixClientProgram) ReadErr(p []byte) (int, error) {
	return cp.stderr.Read(p)
}

func (cp *unixClientProgram) CloseStdin() error {
	return cp.stdin.Close()
}
