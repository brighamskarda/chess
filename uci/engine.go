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

import "io"

const CommandBufferSize int = 128

// ChessEngine is the interface your chess engine should implement. Once you have a chess engine that implements this interface you can pass it to [NewUciEngine]. [UciEngineManager] will do all work of communicating with a client through the Universal Chess Interface (UCI).
type ChessEngine interface {
	start()
	stop()
}

// UciEngineManager wraps a [ChessEngine] and handles all the communication with the client programming running the chess engine. Any command the client wishes to send to the engine will automatically be parsed and the proper functions will be called on the ChessEngine. Should the ChessEngine have communication it wishes to send back to the client, that will also be handled automatically.
//
// With this type, making a UCI compatible chess engine should be trivial as your chess engine need only implement the [ChessEngine] interface and not worry about any of the communication.
//
// This struct can buffer up to [CommandBufferSize] commands from the client, at which point it can receive no more and the client program may have unexpected behavior due to blocking on output. This should be more then enough for most situations. If you hit this limit you should consider redesigning your engine to be more responsive.
//
// There is no buffer on commands being sent out to the client (except for possible OS level pipe buffers).
type UciEngineManager struct {
}

// Creates a new UciEngine that is ready to use. Call [UciEngineManager.Start] on it to start execution of the chess engine.
func NewUciEngine(stdin io.Reader, stdout io.Writer, engine ChessEngine) *UciEngineManager {
	return nil
}

// Start the chess engine. From this point on it will be reading from stdin and writing to stdout. This function will not return until the chess engine is told to quit. In many cases this is likely the last thing you need to call in your main function.
func (ucie *UciEngineManager) Start() {

}
