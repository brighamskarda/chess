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
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ClientSettings determines how [NewClient] should execute its command.
type ClientSettings struct {
	// Args is the arguments that should be passed to the engine. May be nil.
	Args []string
	// Env is the environment variables that the engine should run with. If nil,
	// it will run with environment of the parent process. If empty it will run
	// with an empty environment. Entries should be formatted as "VAR_NAME=VALUE".
	Env []string
	// WorkDir is the working directory to run the engine from. If empty it will
	// run in the working directory of the parent process.
	WorkDir string
	// Logger is where all communication between the engine and the client will
	// be recorded. Logger is closely related to various IO routines, so if it
	// starts blocking for extended periods of time, timeouts may start to
	// occur. May be nil.
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

// Client is the GUI side of UCI that handles game state and sends commands to the engine. Use this if you are developing a chess program that interacts with engines.
//
// There are a mix of functions that can be called concurrently, and others that must be called sequentially. In general status related functions are safe to call concurrently (for ui updates), while functions that send commands to the engine should be called sequentially. [Client.Quit] is an exception to this rule. It can be called at anytime, and should only be called once. See individual method docs for more info.
type Client struct {
	clientProgram clientProgram
	logger        *concurrentWriter
	infoBuf       *concurrentCircBuf[*Info]
	commandBuf    *concurrentBuf[command]
	cpStatus      atomic.Uint32
	regStatus     atomic.Uint32
	engineName    atomic.Pointer[string]
	engineAuthor  atomic.Pointer[string]
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewClient takes in the path to a UCI compliant chess engine and returns a [Client] that will allow you to interface with it. If it could not start the program, an error is returned.
//
// To ensure no resources are leaked, be sure to always call [Client.Quit] on new clients once you are done with them. (calling `defer client.Quit()` right after NewClient can be a good idea.)
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
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.commandBuf = newConcBuf[command](c.ctx)

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
		cw.Write(prefix)
		prefix = prefix[:originalPrefixLen]
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
	if command == nil {
		return
	}
	switch command.commandType() {
	case unknownCommandType:
	case info:
		c.infoBuf.Push(command.(*Info))
	case copyprotection:
		c.cpStatus.Store(uint32(command.(CopyStatus)))
	case registration:
		c.regStatus.Store(uint32(command.(RegStatus)))
	default:
		c.commandBuf.Push(command)
	}
}

func parseCommand(line []byte) command {
	commandType := findCommandType(line)
	switch commandType {
	case unknownCommandType:
		return basicCommand{cmdType: unknownCommandType, msg: string(line)}
	case info:
		if parsedCommand := parseInfoCommand(line); parsedCommand != nil {
			return parsedCommand
		}
	case option:
		if parsedCommand := parseOptionCommand(line); parsedCommand != nil {
			return parsedCommand
		}
	case id:
		if parsedCommand := parseIdCommand(line); parsedCommand != nil {
			return *parsedCommand
		}
	case uciok:
		return basicCommand{
			cmdType: commandType,
			msg:     "",
		}
	case readyok:
		return basicCommand{
			cmdType: commandType,
			msg:     "",
		}
	case bestmove:
		if parsedCommand := parseBestMoveCommand(line); parsedCommand != nil {
			return *parsedCommand
		}
	case copyprotection:
		if parsedCommand := parseCopyProtection(line); parsedCommand != nil {
			return *parsedCommand
		}
	case registration:
		if parsedCommand := parseRegistration(line); parsedCommand != nil {
			return *parsedCommand
		}
	}
	return nil
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
	return unknownCommandType
}

// send is the proper way to send messages to the engine as it will also send the messages to the client's logger. Don't forget to add a new line.
func (c *Client) send(ctx context.Context, p []byte) error {
	done := make(chan error, 1)
	go func() {
		nWrit, err := c.clientProgram.Write(p)
		prefix := []byte(">>> ")
		prefix = append(prefix, p[:nWrit]...)
		c.logger.Write(prefix)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("problem sending message to engine: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("problem sending message to engine, context cancelled: %w", ctx.Err())
	}
}

// ReadInfo returns an [Info] received from the engine. This is implemented as a circular buffer with size 128, meaning that only the last 128 infos are stored. If the engine sends more, the oldest info is deleted to make room. It is a good idea to setup a continuous read loop if you don't want to miss any infos.
//
// ReadInfo blocks if no infos are available. The returned info is safe to modify.
//
// Safe to call concurrently.
func (c *Client) ReadInfo() *Info {
	return c.infoBuf.Next()
}

// Quit sends the quit command to the engine to shutdown as soon as possible. This function should always be called when a client is no longer being used to clean up resources.
//
// Once the timeout1 is reached a system-level request will be sent to the engine to gracefully shutdown (SIGTERM on unix). timeout2 will then begin, and once it is reached the engine will be shutdown forcibly.
//
// This library makes use of process groups (unix) and job objects (windows) to help ensure that the engine don't leave behind stray processes on forcible exits.
//
// Any errors encountered during the exit process will be reported, though outside of debugging they can likely be ignored as this function will only return once the engine is dead.
//
// After this is called, c should no longer be used and all resources for it will be freed.
//
// Safe to call concurrently. Should only be called once.
func (c *Client) Quit(timeout1 time.Duration, timeout2 time.Duration) error {
	var errs error
	defer c.cancel()

	timer1, cancel := context.WithTimeout(context.Background(), timeout1)
	defer cancel()

	done := make(chan error)
	go func() {
		err := c.send(timer1, []byte("quit\n"))
		time.Sleep(timeout1 / 5) // gives the engine time to read quit before the pipe is closed.
		err = errors.Join(err, c.clientProgram.CloseStdin())
		done <- errors.Join(err, c.clientProgram.Wait())
	}()

	select {
	case err := <-done:
		return err
	case <-timer1.Done():
		timer2, cancel2 := context.WithTimeout(context.Background(), timeout2)
		defer cancel2()
		errs = c.clientProgram.Terminate()

		select {
		case err := <-done:
			return errors.Join(errs, err)
		case <-timer2.Done():
			errs = errors.Join(errs, c.clientProgram.Kill())
		}
	}

	errs = errors.Join(errs, <-done)
	return errs
}

// Uci should be the first function you always call on a new client. It tells the program to enter uci mode. If successful [Client.Name] and [Client.Author] will be set, and the engine's options will be returned. Returns an error if the timeout is reached before receiving uciok. [Client.Quit] should be called after receiving an error.
//
// It is a good idea to call [Client.CopyrightStatus] and [Client.RegistrationStatus] a short time after this function to make sure the engine initialized correctly.
//
// Not safe for concurrent use.
func (c *Client) Uci(timeout time.Duration) ([]*Option, error) {
	timer, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	err := c.send(timer, []byte("uci\n"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize uci mode: %w", err)
	}

	nextCommand := make(chan command)
	go readCommandsUntilUciOk(timer, nextCommand, c.commandBuf)

	options := []*Option{}

	for {
		cmd := <-nextCommand
		if cmd == nil {
			return nil, errors.New("could not initialize uci mode, timeout reached before encountering uciok")
		}
		switch cmd.commandType() {
		case id:
			c.setId(cmd.(idCommand))
		case option:
			options = append(options, cmd.(*Option))
		case uciok:
			return options, nil
		}
	}
}

func readCommandsUntilUciOk(ctx context.Context, commands chan<- command, commandBuf *concurrentBuf[command]) {
	for {
		cmd, err := commandBuf.NextWithContext(ctx)
		if err != nil {
			commands <- nil
			break
		}
		commands <- cmd
		if cmd.commandType() == uciok {
			break
		}
	}
}

func (c *Client) setId(id idCommand) {
	switch id.idt {
	case author:
		c.engineAuthor.Store(&id.value)
	case name:
		c.engineName.Store(&id.value)
	}
}

// Name provides the id name sent by the engine after calling [Client.Uci]. Empty string if not set.
//
// Safe to call concurrently.
func (c *Client) Name() string {
	val := c.engineName.Load()
	if val == nil {
		return ""
	}
	return *val
}

// Author provides the author sent by the engine after calling [Client.Uci]. Empty string if not set.
//
// Safe to call concurrently.
func (c *Client) Author() string {
	val := c.engineAuthor.Load()
	if val == nil {
		return ""
	}
	return *val
}

// CopyrightStatus returns the current copyright status received from the engine.
//
//   - [CpUnknown] - No message has been received yet. The Engine may not have copy protection and is good to go.
//   - [CpChecking] - The engine is checking its copy protection. Check back in a few seconds.
//   - [CpOk] - The engine copy protection succeeded.
//   - [CpError] - The engine copy protection failed. You should call [Client.Quit] and make sure you configured the engine correctly.
//
// Safe to call concurrently.
func (c *Client) CopyrightStatus() CopyStatus {
	return CopyStatus(c.cpStatus.Load())
}

// RegistrationStatus returns the current registration status received from the engine.
//
//   - [RegUnknown] - No registration message has been received yet. The engine might not require registration and is ready.
//   - [RegChecking] - The engine is verifying its registration. Wait a few moments before proceeding.
//   - [RegOk] - The engine registration was successful.
//   - [RegError] - The engine registration failed. The engine may still work, but it is a good idea to call [Client.Register], even if it is just to send the later command.
//
// Safe to call concurrently.
func (c *Client) RegistrationStatus() RegStatus {
	return RegStatus(c.regStatus.Load())
}

// IsReady blocks until the engine responds with readyok. Returns false if timeout is reached. This is useful for synchronizing with the engine after calling functions like [Client.Uci], [Client.SetOption], and [Client.Register].
//
// Not safe for concurrent use.
func (c *Client) IsReady(timeout time.Duration) bool {
	timer, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	err := c.send(timer, []byte("isready\n"))
	if err != nil {
		return false
	}

	for {
		command, err := c.commandBuf.NextWithContext(timer)
		if err != nil {
			return false
		}
		if command.commandType() == readyok {
			return true
		}
	}
}

// Debug tells the engine to enter or exit debug mode. Returns an error if timeout expires before the command could be sent.
//
// Not safe for concurrent use.
func (c *Client) Debug(enabled bool, timeout time.Duration) error {
	timer, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	var message []byte
	if enabled {
		message = []byte("debug on\n")
	} else {
		message = []byte("debug off\n")
	}
	err := c.send(timer, message)
	if err != nil {
		return fmt.Errorf("could not send debug message: %w", err)
	}
	return nil
}

func (c *Client) Register() error {
	return nil
}
