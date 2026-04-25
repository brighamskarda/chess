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

import "github.com/brighamskarda/chess/v2"

// ChessEngine is an interface that defines common functions a chess engine should have.
//
// By following this interface,
// users can plug their chess engine into a [UciEngineBroker]
// to automatically gain all the benefits of having a UCI compliant chess engine.
//
// All functions in this interface will by called synchronously with a few exceptions.
// 	- Quit can be called any time after Initialize.
// While the the engine is evaluating these functions can still be called.
//		- Stop
//		- PonderHit
//		- SetDebug
type ChessEngine interface {
	// Initialize will be the first function called on the chess engine.
	//
	// Initialize will be called exactly once and should be used to initialize any internal state the engine relies on.
	//
	// A function is provided via which the engine can send info commands to the UCI client.
	// The function should be stored and used for the duration of the program.
	// The [InfoCmd] documentation should give a good idea of what kind of info can be sent.
	// InfoCmd.StringMsg can be very useful for debugging.
	//
	// Keep in mind that the function is not buffered,
	// so sending info commands during a move search can slow it down significantly.
	// Implementing a buffered channel to asynchronously send info commands is a good idea.
	//
	// If Initialize takes too long, the UCI client may kill the engine.
	Initialize(func(*InfoCmd))

	// CopyProtection provides the engine an opportunity to do copyright checks.
	//
	// CopyProtection will be called after Initialize and
	// should return true if copy protection checks succeeded,
	// or false if the engine could not perform its copy protection checks.
	CopyProtection() bool

	// Register provides an opportunity for the engine to register itself.
	//
	// If the engine is unable to register itself,
	// the user may expect reduced functionality.
	//
	// This function will by called after [ChessEngine.Initialize] and
	// should return true if it can successfully register.
	// The first time this function is called
	// the RegisterCmd will be nil.
	// The engine may try to register without it.
	// If this function returns false (registration failed),
	// future calls to it will have a RegisterCmd included from the UCI client.
	Register(*RegisterCmd) bool

	// Name should return the name of the chess engine.
	//
	// It can contain spaces, but should not contain new lines.
	// A good name is something like "Super Powerful Engine v1.2"
	Name() string

	// Author should return the names of the chess engine developers.
	//
	// It can contain spaces, but should not contain new lines.
	// It may be appropriate to link to a separate authors file if there are a lot of authors.
	Author() string

	// Options should return the options supported by this engine.
	//
	// These options will be send to the UCI client so it knows what can be modified.
	// There are no required options in the UCI specification.
	Options() []Option

	// SetDebug will receive a true when the client requests debug mode.
	//
	// This function can be called asynchronously during evaluation.
	// The engine should default to normal operations (debug = false).
	// When debug mode is on, the engine should send out additional infos to aid development.
	SetDebug(bool)

	// SetOption sets engine parameters.
	//
	// If the engine does not support the given option then it can be ignored.
	//
	// A type switch on SetOption with the following types will likely be necessary to implement this function.
	//   - [SetCheckOption]
	//   - [SetSpinOption]
	//   - [SetStringOption] (which double as setting a combo option)
	//   - [SetButtonOption]
	SetOption(SetOption)

	// NewGame is called when the next position to evaluate is part of a different game.
	//
	// While it isn't strictly required to implement new game,
	// it is good practice to clear any cached values the engine was using.
	NewGame()

	// SetPosition tells the chess engine to setup a certain position for its next search.
	//
	// The position will always be provided,
	// but the list of moves may be empty or nil.
	// The moves should be applied to the position to get a final position.
	// By providing the move history
	// the engine can see if it is headed towards a draw by three-fold repetition.
	//
	// Both of the parameters may be consumed by this function.
	SetPosition(*chess.Position, []chess.Move)

	// Evaluate tells the engine to start evaluating on its current position.
	//
	// Various options are provided via the EvaluateCmd parameter
	// to modify how the engine behaves.
	//
	// **IF THE PONDER OPTION IS SET DO NOT RETURN UNTIL STOP OR PONDERHIT ARE CALLED.**
	//
	// This function should evaluate until
	// it finds what it thinks is the best move,
	// or Stop/Quit are called asynchronously,
	// at which point it should return its best move as soon as possible.
	//
	// Before returning, sending an InfoCmd with
	// the final stats of the evaluation is recommended.
	Evaluate(*EvaluateCmd) *BestMove

	// Stop asks the engine to stop its current move evaluation and return the best move.
	//
	// Stop may be called asynchronously any time the engine is evaluating.
	// If Stop is called and the engine is not evaluating, it should be ignored.
	Stop()

	// PonderHit tells the engine to switch from pondering to normal search.
	//
	// If the engine was told to ponder on a move
	// and the opponent plays that move
	// then PonderHit is called to indicate that the
	// engine was pondering on the correct move
	// and that it can now return its best move when its ready.
	//
	// If the opponent plays another move
	// then stop will be called instead and a new
	// position will be set before calling Evaluate again.
	//
	// If PonderHit is called and the engine is not pondering, it should be ignored.
	PonderHit()

	// Quit can be called asynchronously any time after Initialize.
	//
	// This is the broker's way of nicely asking the engine to stop its operations.
	// Failure to do so promptly may result in the engine being forcibly stopped.
	// Quit should not return until all cleanup is complete.
	//
	// The following are some (but not all) actions that should be taken to ensure a smooth shutdown.
	//	- Stop ongoing searches
	// 	- Close Open files
	// 	- Release remote software licenses
	//
	// Quit will only be called once.
	Quit()
}
