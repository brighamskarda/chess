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
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"
)

// ClientSettings determines how [NewClient] should execute its command.
type ClientSettings struct {
	// Args is the arguments that should be passed to the engine. May be nil.
	Args []string
	// Env is the environment variables that the engine should run with. If nil, it will run with environment of the parent process. If empty it will run with an empty environment. Entries should be formatted as "VAR_NAME=VALUE".
	Env []string
	// WorkDir is the working directory to run the engine from. If empty it will run in the working directory of the parent process.
	WorkDir string
	// Logger is where all communication between the engine and the client will take place.
	//
	// `>>>` Three arrows right indicates a line that was sent to stdin for the engine.
	//
	// `<<<` Three arrows left indicates a line that was received from the engine's stdout.
	//
	// `!<!` This pattern indicates a line that was received from the engine's stderr.
	Logger io.Writer
}

// clientProgram should be a uci compatible chess engine that is already running.
// When finished with a clientProgram be sure to call Wait() to free any resources.
type clientProgram interface {
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
	Terminate() error
	// Kill immediately stops the program. Returns an error if the request was not sent successfully
	// Call wait after this function to clean up resources.
	Kill() error
	// Wait waits for the program to finish and cleans up associated resources.
	// It may return an error if the program did not exit successfully (like returning exit code 1), or there were io errors.
	// Wait should only be called once. Ensure that stdin is closed by calling CloseStdin(), and both readers are flushed to prevent blocking.
	Wait() error
	// Read reads from the program's stdout.
	Read(p []byte) (int, error)
	// ReadErr reads from the program's stderr.
	ReadErr(p []byte) (int, error)
	// Write writes to the program's stdin.
	Write(p []byte) (int, error)
	// CloseStdin closes the underlying pipe to stdin. You cannot use Write() after this is called.
	CloseStdin() error
}

type concurrentWriter struct {
	m sync.Mutex
	w io.Writer
}

func (cw *concurrentWriter) Write(p []byte) (int, error) {
	cw.m.Lock()
	defer cw.m.Unlock()
	return cw.w.Write(p)
}

// Client is the side of UCI that handles game state and sends commands to the [Engine]. Use this if you are developing a chess program that interacts with engines.
type Client struct {
	clientProgram clientProgram
	logger        *concurrentWriter
	infoBuf       *concurrentCircBuf[*Info]
}

// NewClient takes in the path to UCI compliant chess engine and returns a [Client] that will allow you to interface with it. If it could not start the program, an error is returned.
func NewClient(program string, settings ClientSettings) (*Client, error) {
	cp, err := newClientProgram(program, settings)
	if err != nil {
		return nil, fmt.Errorf("could not make new client: %w", err)
	}

	return newClientFromClientProgram(cp, settings)
}

func newClientFromClientProgram(cp clientProgram, settings ClientSettings) (*Client, error) {
	c := &Client{
		clientProgram: cp,
		infoBuf:       newCircBuf[*Info](128),
	}

	c.setUpLogger(settings.Logger)
	c.startReadLoop()

	return c, nil
}

func (c *Client) setUpLogger(w io.Writer) {
	if w == nil {
		w = io.Discard
	}
	c.logger = &concurrentWriter{w: w}
}

func (c *Client) startReadLoop() {
	go c.stdoutReadLoop()
	go stderrReadLoop(c.logger, readErrWrapper{cp: c.clientProgram})
}

type readErrWrapper struct{ cp clientProgram }

func (rew readErrWrapper) Read(p []byte) (int, error) {
	return rew.cp.ReadErr(p)
}

func stderrReadLoop(cw *concurrentWriter, r io.Reader) {
	scnr := bufio.NewScanner(r)
	prefix := []byte("!<! ")
	originalPrefixLen := len(prefix)
	for scnr.Scan() {
		prefix = append(prefix, scnr.Bytes()...)
		prefix = append(prefix, '\n')
		_, err := cw.Write(prefix)
		if err != nil {
			break
		}
		prefix = prefix[:originalPrefixLen]
	}
	for scnr.Scan() {
	}
}

func (c *Client) stdoutReadLoop() {
	scnr := bufio.NewScanner(c.clientProgram)
	prefix := []byte("<<< ")
	originalPrefixLen := len(prefix)
	for scnr.Scan() {
		line := scnr.Bytes()

		// send to logger
		prefix = append(prefix, line...)
		prefix = append(prefix, '\n')
		c.logger.Write(prefix)
		prefix = prefix[:originalPrefixLen]

		// parse command for rest of Client
		c.handleCommand(line)
	}
}

func (c *Client) handleCommand(line []byte) {
	command := parseCommand(line)
	switch command.commandType() {
	case info:
		c.infoBuf.Push(command.(*Info))
	}
}

func parseCommand(line []byte) command {
	commandType := findCommandType(line)
	switch commandType {
	case unknown:
		return basicCommand{cmdType: unknown, msg: string(line)}
	case info:
		return parseInfoCommand(line)
	}
	panic(fmt.Sprintf("could not parse command, unexpected command type %d", commandType))
}

func findCommandType(line []byte) commandType {
	scnr := bufio.NewScanner(bytes.NewBuffer(slices.Clone(line)))
	scnr.Split(bufio.ScanWords)
	for scnr.Scan() {
		token := strings.ToLower(scnr.Text())
		switch token {
		case "id":
			return id
		case "uciok":
			return uciok
		case "readyok":
			return readyok
		case "bestmove":
			return bestmove
		case "copyprotection":
			return copyprotection
		case "registration":
			return registration
		case "info":
			return info
		case "option":
			return option
		}
	}
	return unknown
}

// send is the proper way to send messages to the engine as it will also send the messages to the client's logger.
func (c *Client) send(p []byte) error {
	prefix := []byte(">>> ")
	prefix = append(prefix, p...)
	c.logger.Write(prefix)
	_, err := c.clientProgram.Write(p)
	return err
}

// ReadInfo returns an [Info] received from the engine. This is implemented as a circular buffer with size 128, meaning that only the last 128 infos are stored. If the engine sends more, the oldest info is deleted to make room. It is a good idea to setup a continuous read loop if you don't want to miss any infos.
//
// ReadInfo blocks if no infos are available. The returned info is safe to modify.
func (c *Client) ReadInfo() *Info {
	return c.infoBuf.Next()
}

// Quit sends the quit command to the engine to shutdown as soon as possible.  timeout1 is the amount of time the process should wait before resorting to other measures. Once the timeout1 is reached a request will be sent to the engine to gracefully shutdown (SIGTERM on unix). timeout2 will then begin, and once it is reached the engine will be shutdown forcibly.
//
// This library makes use of process groups (unix) and job objects (windows) to help ensure that the engine don't leave behind stray processes on forcible exits.
//
// Any errors encountered during the exit process will be reported, though outside of debugging they can likely be ignored as this function will only return once the engine is dead.
//
// After this is called, c should no longer be used and all resources for it will be freed.
func (c *Client) Quit(timeout1 time.Duration, timeout2 time.Duration) error {
	var errs error

	ctx, cancel := context.WithTimeout(context.Background(), timeout1)
	defer cancel()

	done := make(chan error)
	go func() {
		err := c.send([]byte("quit\n"))
		time.Sleep(timeout1 / 5) // gives the engine time to read quit before the pipe is closed.
		err = errors.Join(err, c.clientProgram.CloseStdin())
		done <- errors.Join(err, c.clientProgram.Wait())
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		ctx2, cancel2 := context.WithTimeout(context.Background(), timeout2)
		defer cancel2()
		errs = c.clientProgram.Terminate()

		select {
		case err := <-done:
			return errors.Join(errs, err)
		case <-ctx2.Done():
			errs = errors.Join(errs, c.clientProgram.Kill())
		}
	}

	errs = errors.Join(errs, <-done)
	return errs
}

// TODO make sure that the constant read loop flushed stdout when it gets the cancel context.

// Uci should be the first function you always call. It tells the program on the other end of the writer to enter uci mode.
//
// tell engine to use the uci (universal chess interface),
// this will be sent once as a first command after program boot
// to tell the engine to switch to uci mode.
// After receiving the uci command the engine must identify itself with the "id" command
// and send the "option" commands to tell the GUI which engine settings the engine supports if any.
// After that the engine should send "uciok" to acknowledge the uci mode.
// If no uciok is sent within a certain time period, the engine task will be killed by the GUI.
// func (c *Client) Uci(timeout) error {
// 	n, err := c.engineWriter.Write([]byte("uci\n"))
// 	if err != nil {

// 	}
// }
