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

// ChessEngine is an interface that defines common functions a chess engine should have.
//
// By following this interface,
// users can plug their chess engine into a [UciEngineBroker]
// to automatically gain all the benefits of having a UCI compliant chess engine.
type ChessEngine interface {
	// Initialize will be the first function called on the chess engine.
	//
	// Initialize will be called exactly once and should be used to initialize any internal state the engine relies on.
	//
	// A function is provided via which the engine can send info commands to the UCI chess client.
	// The function should be stored and used for the duration of the program.
	// The [InfoCmd] documentation should give a good idea of what kind of info can be sent.
	// Sending commands with only [InfoCmd.StringMsg] can be very useful for debugging.
	//
	// Keep in mind that the function is not buffered,
	// so sending info commands during a move search can slow it down significantly.
	// Implementing a buffered channel to asynchronously send info commands is a good idea.
	//
	// Keep in mind that if Initialize takes too long, the GUI may kill the engine.
	Initialize(func(*InfoCmd))

	// CopyProtection provides the engine an opportunity to do copyright checks.
	//
	// CopyProtection will be called after [ChessEngine.Initialize] and
	// should return true if copy protection checks succeeded,
	// or false if the engine could not perform its copy protection checks.
	CopyProtection() bool

	// Register provides and opportunity for the engine to register itself.
	//
	// If the engine is unable to register itself,
	// the user may expect reduced functionality.
	//
	// This function will by called after [ChessEngine.Initialize] and
	// should return true if it can successfully register.
	// The first time this function is called
	// the RegisterCmd will be nil,
	// and the engine may try to register without it.
	// If this function returns false (registration failed),
	// future calls to it will have a RegisterCmd included.
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
	// These options will be send to the GUI so the user may modify them.
	// There are no required options in the UCI specification.
	Options() []OptionCmd

	// SetDebug will receive a true when the client requests debug mode.
	//
	// This function can be called asynchronously at any time.
	// The engine should default to normal operations (debug = false).
	// When debug mode is on, the engine should send out additional infos to aid development.
	SetDebug(bool)

	// SetOption sets engine parameters.
	//
	// SetOption will not be called while the engine is searching.
	// If the engine does not support the given option then it can just ignore it.
	//
	// A type switch on OptionCmd with the following types will likely be necessary to implement this function.
	//   - [SetCheckOptionCmd]
	//   - [SetSpinOptionCmd]
	//   - [SetStringOptionCmd] (which double as setting a combo option)
	//   - [SetButtonOptionCmd]
	SetOption(SetOptionCmd)

	// Quit can be called asynchronously any time after Initialize.
	//
	// This is the broker's way of asking the engine to nicely stop its operations.
	// Failure to do so promptly may result in the engine being forcibly stopped.
	// Quit should not return until all cleanup is complete.
	//
	// The following are some (but not all) actions that should be taken to ensure a smooth shutdown.
	//	* Stop ongoing searches
	// 	* Close Open files
	// 	* Release remote software licenses
	//
	// Quit will only be called once.
	Quit()
}
