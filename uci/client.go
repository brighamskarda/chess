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
	"io"
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
	// Wait should only be called once. Ensure the io.WriteCloser is closed, and both readers are flushed to prevent blocking.
	Wait() error
	// Read reads from the program's stdout.
	Read(p []byte) (int, error)
	// Write writes to the program's stdin.
	Write(p []byte) (int, error)
	// ReadErr reads from the program's stderr.
	ReadErr(p []byte) (int, error)
}

// Client is the side of UCI that handles game state and sends commands to the [Engine]. Use this if you are developing a chess program that interacts with engines.
type Client struct {
	clientProgram clientProgram
	logger        *concurrentWriter
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

// NewClient takes in the path to UCI compliant chess engine and returns a [Client] that will allow you to interface with it. If it could not start the program, an error is returned.
func NewClient(program string, settings ClientSettings) (*Client, error) {
	c := &Client{}

	// c.setUpLogger(settings.Logger)

	// Setup cancel context
	// c.ctx, c.cancel = context.WithCancel(context.Background())

	// // Setup Engine
	// c.engine = exec.CommandContext(c.ctx, program, settings.Args...)

	// // Setup stdin and stdout
	// var err error
	// c.engineWriter, err = c.engine.StdinPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create new client: %w", err)
	// }
	// c.engineReader, err = c.engine.StdoutPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create new client: %w", err)
	// }

	// // Set env and working dir
	// c.engine.Env = settings.Env
	// c.engine.Dir = settings.WorkDir

	// // Start the engine
	// err = c.engine.Start()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create new client: %w", err)
	// }

	return c, nil
}

func (c *Client) setUpLogger(w io.Writer) {
	if w == nil {
		w = io.Discard
	}
	c.logger = &concurrentWriter{w: w}
}

// Quit sends the command to the engine to shutdown as soon as possible. After this is called, c should no longer be used and all resources for it will be freed. timeout is the amount of time the process should wait before forcibly shutting down the engine. If timeout is set to 0 then waiting for the engine to close gracefully may never end. An error is returned if the program was not exited gracefully. This error may be ignored if you don't care about the exit status of the engine.
func (c *Client) Quit(timeout time.Duration) error {
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

	// done := make(chan error)
	// go func() {
	// 	c.engineWriter.Write([]byte("quit\n"))
	// 	c.engineWriter.Close()
	// 	done <- c.engine.Wait()
	// }()

	// select {
	// case err := <-done:
	// 	return err
	// case <-ctx.Done():
	// 	c.cancel()
	// }

	return nil
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
