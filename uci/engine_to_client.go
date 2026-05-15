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

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/brighamskarda/chess/v2"
)

// engineToClientCmd is an interface under which all uci commands from the engine to the client will be contained.
type engineToClientCmd interface {
	marshalText() ([]byte, error)
}

// idCmd
//   - name <x>
//     this must be sent after receiving the "uci" command to identify the engine,
//     e.g. "id name Shredder X.Y\n"
//   - author <x>
//     this must be sent after receiving the "uci" command to identify the engine,
//     e.g. "id author Stefan MK\n"
type idCmd struct {
	// isAuthor is true if this is an author id command.
	// If it is false then it is assumed to be the name id command.
	isAuthor bool
	// id is either the name of the engine, or the author of the engine.
	id string
}

func (cmd *idCmd) marshalText() ([]byte, error) {
	text := bytes.Buffer{}
	text.WriteString("id ")
	if cmd.isAuthor {
		text.WriteString("author ")
	} else {
		text.WriteString("name ")
	}
	text.WriteString(cmd.id)
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// uciokCmd
// Must be sent after the id and optional options to tell the GUI that the engine
// has sent all infos and is ready in uci mode.
type uciokCmd struct{}

func (cmd *uciokCmd) marshalText() ([]byte, error) {
	return []byte("uciok\n"), nil
}

// readyOkCmd
// This must be sent when the engine has received an "isready" command and has
// processed all input and is ready to accept new commands now.
// It is usually sent after a command that can take some time to be able to wait for the engine,
// but it can be used anytime, even when the engine is searching,
// and must always be answered with "isready".
type readyOkCmd struct{}

func (cmd *readyOkCmd) marshalText() ([]byte, error) {
	return []byte("readyok\n"), nil
}

// BestMove gives results of an engine's evaluation.
type BestMove struct {
	// Move is the best move according to the engine's evaluation.
	Move chess.Move
	// PonderMove is an optional field that
	// indicates the move the engine would like to ponder on
	// (if given the opportunity).
	//
	// PonderMove should be a move the engine thinks the opponent will play
	// in response to Move.
	PonderMove Optional[chess.Move]
}

func (cmd *BestMove) marshalText() ([]byte, error) {
	text := bytes.Buffer{}
	text.WriteString("bestmove ")

	move, err := cmd.Move.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("could not marshal bestMoveCmd: %w", err)
	}
	text.Write(move)

	if cmd.PonderMove.HasValue() {
		text.WriteString(" ponder ")

		ponderMove, err := cmd.PonderMove.Value().MarshalText()
		if err != nil {
			return nil, fmt.Errorf("could not marshal bestMoveCmd: %w", err)
		}
		text.Write(ponderMove)
	}

	text.WriteByte('\n')
	return text.Bytes(), nil
}

// copyProtectionCmd - this is needed for copyprotected engines.
//
// After the uciok command the engine can tell the GUI,
// that it will check the copy protection now. This is done by "copyprotection checking".
// If the check is ok the engine should send "copyprotection ok", otherwise "copyprotection error".
// If there is an error the engine should not function properly but should not quit alone.
// If the engine reports "copyprotection error" the GUI should not use this engine
// and display an error message instead!
// The code in the engine can look like this
//
//	TellGUI("copyprotection checking\n");
//	 // ... check the copy protection here ...
//	 if(ok)
//	    TellGUI("copyprotection ok\n");
//	else
//	   TellGUI("copyprotection error\n");
type copyProtectionCmd uint8

const (
	copyProtectChecking copyProtectionCmd = iota
	copyProtectOk
	copyProtectError
)

func (cmd *copyProtectionCmd) marshalText() ([]byte, error) {
	switch *cmd {
	case copyProtectChecking:
		return []byte("copyprotection checking\n"), nil
	case copyProtectOk:
		return []byte("copyprotection ok\n"), nil
	case copyProtectError:
		return []byte("copyprotection error\n"), nil
	default:
		return nil, fmt.Errorf("could not marshal copyprotection command: invalid value %v", *cmd)
	}
}

// registrationCmd - this is needed for engines that need a username and/or a code to function with all features.
//
// Analog to the "copyprotection" command the engine can send "registration checking"
// after the uciok command followed by either "registration ok" or "registration error".
// Also after every attempt to register the engine it should answer with "registration checking"
// and then either "registration ok" or "registration error".
// In contrast to the "copyprotection" command, the GUI can use the engine after the engine has
// reported an error, but should inform the user that the engine is not properly registered
// and might not use all its features.
// In addition the GUI should offer to open a dialog to
// enable registration of the engine. To try to register an engine the GUI can send
// the "register" command.
// The GUI has to always answer with the "register" command	if the engine sends "registration error"
// at engine startup (this can also be done with "register later")
// and tell the user somehow that the engine is not registered.
// This way the engine knows that the GUI can deal with the registration procedure and the user
// will be informed that the engine is not properly registered.
type registrationCmd uint8

const (
	registerChecking registrationCmd = iota
	registerOk
	registerError
)

func (cmd *registrationCmd) marshalText() ([]byte, error) {
	switch *cmd {
	case registerChecking:
		return []byte("registration checking\n"), nil
	case registerOk:
		return []byte("registration ok\n"), nil
	case registerError:
		return []byte("registration error\n"), nil
	default:
		return nil, fmt.Errorf("could not marshal registration command: invalid value %v", *cmd)
	}
}

// InfoCmd can be used to send information to the UCI client.
//
// Ideally this should be done whenever one of the infos fields has changed.
// The engine can send any combination of information at once.
//
// I suggest to start sending "currmove", "currmovenumber", "currline" and "refutation"
// only after one second of evaluation
// to avoid the heavy traffic often associated
// with the start of an evaluation.
type InfoCmd struct {
	// Depth <x> - search Depth in plies
	Depth Optional[int]

	// SelDepth - selective search depth in plies.
	//
	// If the engine sends SelDepth there must also be a "depth" present in the same command.
	SelDepth Optional[int]

	// Time Time searched in milliseconds, this should be sent together with the pv.
	Time Optional[int]

	// Nodes Nodes searched, the engine should send this info regularly.
	Nodes Optional[int]

	// Pv <move1> ... <movei> - the best line found. The primary variation.
	Pv []chess.Move

	// MultiPv <num> - this for the multi pv mode.
	//
	// For the best move/pv add "MultiPv 1" in the string when you send the pv.
	// In k-best mode always send all k variants in k strings together.
	MultiPv Optional[int]

	// Score
	//
	// 	- cp <x>
	// 		the Score from the engine's point of view in centipawns.
	// 	- mate <y>
	// 		mate in y moves, not plies.
	// 		If the engine is getting mated use negative values for y.
	// 	- lowerbound
	//       the Score is just a lower bound.
	// 	- upperbound
	// 	   the Score is just an upper bound.
	Score Optional[InfoScore]

	// CurrMove - currently searching this move.
	CurrMove Optional[chess.Move]

	// CurrMoveNumber - currently searching move number x, for the first move x should be 1 not 0.
	CurrMoveNumber Optional[int]

	// HashFull - the hash is x permill full, the engine should send this info regularly.
	HashFull Optional[int]

	// Nps - x nodes per second searched, the engine should send this info regularly.
	Nps Optional[int]

	// TbHits - x positions where found in the endgame table bases.
	TbHits Optional[int]

	// SbHits - x positions where found in the shredder endgame databases.
	//
	// This isn't really used.
	SbHits Optional[int]

	// CpuLoad - the cpu usage of the engine is x permill.
	CpuLoad Optional[int]

	// StringMsg - a generic string message.
	//
	// This is very useful for debugging/development
	// as you can effectively display anything you want.
	StringMsg Optional[string]

	// Refutation <move1> <move2> ... <movei>
	//
	// Move <move1> is refuted by the line <move2> ... <movei>, i can be any number >= 1.
	// Example: after move d1h5 is searched, the engine can send "info Refutation d1h5 g6h5"
	// if g6h5 is the best answer after d1h5 or if g6h5 refutes the move d1h5.
	// if there is no Refutation for d1h5 found, the engine should just send "info Refutation d1h5"
	// The engine should only send this if the option "UCI_ShowRefutations" is set to true.
	Refutation []chess.Move

	// CurrLine <cpunr> <move1> ... <movei>
	//
	// This is the current line the engine is calculating.
	// If the engine is just using one cpu, <cpunr> can be omitted.
	// If <cpunr> is greater than 1, always send all lines together.
	// The engine should only send this if the option "UCI_ShowCurrLine" is set to true.
	CurrLine Optional[CurrentLine]
}

func (cmd *InfoCmd) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	text.WriteString("info")

	cmd.marshalDepth(text)
	cmd.marshalSelDepth(text)
	cmd.marshalTime(text)
	cmd.marshalNodes(text)
	cmd.marshalPv(text)
	cmd.marshalMultiPv(text)
	cmd.marshalScore(text)
	cmd.marshalCurrMove(text)
	cmd.marshalCurrMoveNumber(text)
	cmd.marshalHashFull(text)
	cmd.marshalNps(text)
	cmd.marshalTbHits(text)
	cmd.marshalSbHits(text)
	cmd.marshalCpuLoad(text)
	cmd.marshalRefutation(text)
	cmd.marshalCurrLine(text)
	// it is important that the string message is marshaled last as everything after it is considered part of the string to the client.
	cmd.marshalString(text)

	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *InfoCmd) marshalDepth(text *bytes.Buffer) {
	if cmd.Depth.HasValue() {
		text.WriteString(" depth ")
		text.WriteString(strconv.Itoa(cmd.Depth.Value()))
	}
}

func (cmd *InfoCmd) marshalSelDepth(text *bytes.Buffer) {
	if cmd.SelDepth.HasValue() {
		text.WriteString(" seldepth ")
		text.WriteString(strconv.Itoa(cmd.SelDepth.Value()))
	}
}

func (cmd *InfoCmd) marshalTime(text *bytes.Buffer) {
	if cmd.Time.HasValue() {
		text.WriteString(" time ")
		text.WriteString(strconv.Itoa(cmd.Time.Value()))
	}
}

func (cmd *InfoCmd) marshalNodes(text *bytes.Buffer) {
	if cmd.Nodes.HasValue() {
		text.WriteString(" nodes ")
		text.WriteString(strconv.Itoa(cmd.Nodes.Value()))
	}
}

func marshalMoveList(moves []chess.Move, text *bytes.Buffer) {
	for _, m := range moves {
		moveText, _ := m.MarshalText()
		text.WriteByte(' ')
		text.Write(moveText)
	}
}

func (cmd *InfoCmd) marshalPv(text *bytes.Buffer) {
	if cmd.Pv != nil {
		text.WriteString(" pv")
		marshalMoveList(cmd.Pv, text)
	}
}

func (cmd *InfoCmd) marshalMultiPv(text *bytes.Buffer) {
	if cmd.MultiPv.HasValue() {
		text.WriteString(" multipv ")
		text.WriteString(strconv.Itoa(cmd.MultiPv.Value()))
	}
}

func (cmd *InfoCmd) marshalScore(text *bytes.Buffer) {
	if cmd.Score.HasValue() {
		text.WriteString(" score ")

		value := cmd.Score.Value()

		if value.IsMate {
			text.WriteString("mate ")
		} else {
			text.WriteString("cp ")
		}

		text.WriteString(strconv.Itoa(value.Score))

		if value.IsLowerBound {
			text.WriteString(" lowerbound")
		}
		if value.IsUpperBound {
			text.WriteString(" upperbound")
		}
	}
}

func (cmd *InfoCmd) marshalCurrMove(text *bytes.Buffer) {
	if cmd.CurrMove.HasValue() {
		text.WriteString(" currmove ")
		moveText, _ := cmd.CurrMove.Value().MarshalText()
		text.Write(moveText)
	}
}

func (cmd *InfoCmd) marshalCurrMoveNumber(text *bytes.Buffer) {
	if cmd.CurrMoveNumber.HasValue() {
		text.WriteString(" currmovenumber ")
		text.WriteString(strconv.Itoa(cmd.CurrMoveNumber.Value()))
	}
}

func (cmd *InfoCmd) marshalHashFull(text *bytes.Buffer) {
	if cmd.HashFull.HasValue() {
		text.WriteString(" hashfull ")
		text.WriteString(strconv.Itoa(cmd.HashFull.Value()))
	}
}

func (cmd *InfoCmd) marshalNps(text *bytes.Buffer) {
	if cmd.Nps.HasValue() {
		text.WriteString(" nps ")
		text.WriteString(strconv.Itoa(cmd.Nps.Value()))
	}
}

func (cmd *InfoCmd) marshalTbHits(text *bytes.Buffer) {
	if cmd.TbHits.HasValue() {
		text.WriteString(" tbhits ")
		text.WriteString(strconv.Itoa(cmd.TbHits.Value()))
	}
}

func (cmd *InfoCmd) marshalSbHits(text *bytes.Buffer) {
	if cmd.SbHits.HasValue() {
		text.WriteString(" sbhits ")
		text.WriteString(strconv.Itoa(cmd.SbHits.Value()))
	}
}

func (cmd *InfoCmd) marshalCpuLoad(text *bytes.Buffer) {
	if cmd.CpuLoad.HasValue() {
		text.WriteString(" cpuload ")
		text.WriteString(strconv.Itoa(cmd.CpuLoad.Value()))
	}
}

func (cmd *InfoCmd) marshalRefutation(text *bytes.Buffer) {
	if cmd.Refutation != nil {
		text.WriteString(" refutation")
		marshalMoveList(cmd.Refutation, text)
	}
}

func (cmd *InfoCmd) marshalCurrLine(text *bytes.Buffer) {
	if cmd.CurrLine.HasValue() {
		currLine := cmd.CurrLine.Value()

		text.WriteString(" currline")
		if currLine.CpuNr.HasValue() {
			text.WriteByte(' ')
			text.WriteString(strconv.Itoa(currLine.CpuNr.Value()))
		}

		marshalMoveList(currLine.Moves, text)
	}
}

func (cmd *InfoCmd) marshalString(text *bytes.Buffer) {
	if cmd.StringMsg.HasValue() {
		text.WriteString(" string ")
		text.WriteString(cmd.StringMsg.Value())
	}
}

// CurrentLine is used in [InfoCmd] to show the current line the engine is calculating.
type CurrentLine struct {
	// CpuNr is the CPU number that is evaluating this line.
	CpuNr Optional[int]
	// Moves is the line of moves being evaluated.
	Moves []chess.Move
}

// InfoScore is used in [InfoCmd] to show the engine's current evaluation of the position.
type InfoScore struct {
	// Score is either the number of plies until mate,
	// or the engine's evaluation of the position in centipawns.
	Score int
	// IsMate is true if the score represents how many plies until mate.
	// Otherwise score is assumed to be the engines evaluation in centipawns.
	IsMate bool
	// IsLowerBound indicates that this score is a lower bound.
	// Should be false if IsUpperBound is set.
	IsLowerBound bool
	// IsUpperBound indicates that this score is an upper bound.
	// Should be false if IsLowerBound is set.
	IsUpperBound bool
}

// Option represents an option that the chess engine supports.
//
// Options are sent to the UCI client so it knows what it can modify in the engine.
// There is no need to make any types implementing this interface
// as all valid options can be represented using the structs provided in this module.
// See [CheckOption], [SpinOption], [ComboOption], [StringOption], and [ButtonOption]
//
// When creating an option it is highly recommended to avoid the words
// "value" and "name" in its parameters as these can confuse UCI parsers.
//
// The following options are predefined in the UCI chess standard.
// Options with these names should not be used for anything else.
// Furthermore any options not listed below that are still prepended with "UCI_" will likely be ignored by the GUI.
//   - <name> = Hash, type is spin
//     The value in MB for memory for hash tables can be changed,
//     this should be answered with the first "setoptions" command at program boot
//     if the engine has sent the appropriate "option name Hash" command,
//     which should be supported by all engines!
//     So the engine should use a very small hash first as default.
//   - <name> = NalimovPath, type string
//     This is the path on the hard disk to the Nalimov compressed format.
//     Multiple directories can be concatenated with ";"
//   - <name> = NalimovCache, type spin
//     This is the size in MB for the cache for the nalimov table bases
//     These last two options should also be present in the initial options exchange dialog
//     when the engine is booted if the engine supports it
//   - <name> = Ponder, type check
//     This means that the engine is able to ponder.
//     The GUI will send this whenever pondering is possible or not.
//     Note: The engine should not start pondering on its own if this is enabled, this option is only
//     needed because the engine might change its time management algorithm when pondering is allowed.
//   - <name> = OwnBook, type check
//     This means that the engine has its own book which is accessed by the engine itself.
//     if this is set, the engine takes care of the opening book and the GUI will never
//     execute a move out of its book for the engine. If this is set to false by the GUI,
//     the engine should not access its own book.
//   - <name> = MultiPV, type spin
//     The engine supports multi best line or k-best mode. the default value is 1
//   - <name> = UCI_ShowCurrLine, type check, should be false by default,
//     The engine can show the current line it is calculating. see "info currline" above.
//   - <name> = UCI_ShowRefutations, type check, should be false by default,
//     The engine can show a move and its refutation in a line. see "info refutations" above.
//   - <name> = UCI_LimitStrength, type check, should be false by default,
//     The engine is able to limit its strength to a specific Elo number,
//     This should always be implemented together with "UCI_Elo".
//   - <name> = UCI_Elo, type spin
//     The engine can limit its strength in Elo within this interval.
//     If UCI_LimitStrength is set to false, this value should be ignored.
//     If UCI_LimitStrength is set to true, the engine should play with this specific strength.
//     This should always be implemented together with "UCI_LimitStrength".
//   - <name> = UCI_AnalyseMode, type check
//     The engine wants to behave differently when analysing or playing a game.
//     For example when playing it can use some kind of learning.
//     This is set to false if the engine is playing a game, otherwise it is true.
//   - <name> = UCI_Opponent, type string
//     With this command the GUI can send the name, title, elo and if the engine is playing a human
//     or computer to the engine.
//     The format of the string has to be [GM|IM|FM|WGM|WIM|none] [<elo>|none] [computer|human] <name>
//     Examples:
//     "setoption name UCI_Opponent value GM 2800 human Gary Kasparov"
//     "setoption name UCI_Opponent value none none computer Shredder"
//   - <name> = UCI_EngineAbout, type string
//     With this command, the engine tells the GUI information about itself, for example a license text,
//     usually it doesn't make sense that the GUI changes this text with the setoption command.
//     Example:
//     "option name UCI_EngineAbout type string default Shredder by Stefan Meyer-Kahlen, see www.shredderchess.com"
//   - <name> = UCI_ShredderbasesPath, type string
//     this is either the path to the folder on the hard disk containing the Shredder endgame databases or
//     the path and filename of one Shredder endgame datbase.
//   - <name> = UCI_SetPositionValue, type string
//     the GUI can send this to the engine to tell the engine to use a certain value in centipawns from white's
//     point of view if evaluating this specifix position.
//     The string can have the formats:
//     <value> + <fen> | clear + <fen> | clearall
type Option interface {
	engineToClientCmd
	// optionName should return the name/id of the option.
	optionName() string
}

// CheckOption represents an option of type Check that the engine supports.
//
// Check options are simple true or false check boxes.
//
// See [Option] for more info.
type CheckOption struct {
	// Name is the name of the option.
	//
	// Certain options are expected to be of a specific type.
	// See [Option] for more info.
	Name string
	// DefaultValue is the default setting the engine uses for this option.
	DefaultValue bool
}

func (cmd *CheckOption) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "check")
	text.WriteString(" default ")
	text.WriteString(strconv.FormatBool(cmd.DefaultValue))
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *CheckOption) optionName() string {
	return cmd.Name
}

// SpinOption represents an option of type Spin that the engine supports.
//
// Spin options are like a spin wheel that can represent numbers
// between a defined minimum and maximum value.
//
// See [Option] for more info.
type SpinOption struct {
	// Name is the name of the option.
	//
	// Certain options are expected to be of a specific type.
	// See [Option] for more info.
	Name string
	// DefaultValue is the default setting the engine uses for this option.
	//
	// It should be between the minimum and maximum values.
	DefaultValue int
	// Min is the minimum value (inclusive) that this option can be set to.
	Min int
	// Max is the maximum value (inclusive) that this option can be set to.
	Max int
}

func (cmd *SpinOption) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "spin")
	text.WriteString(" default ")
	text.WriteString(strconv.Itoa(cmd.DefaultValue))
	text.WriteString(" min ")
	text.WriteString(strconv.Itoa(cmd.Min))
	text.WriteString(" max ")
	text.WriteString(strconv.Itoa(cmd.Max))
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *SpinOption) optionName() string {
	return cmd.Name
}

// ComboOption represents an option of type Combo that the engine supports.
//
// Combo options can represent a predefined subset of strings.
//
// See [Option] for more info.
type ComboOption struct {
	// Name is the name of the option.
	//
	// Certain options are expected to be of a specific type.
	// See [Option] for more info.
	Name string
	// DefaultValue is the default setting the engine uses for this option.
	//
	// It should be one of the Variants.
	DefaultValue string
	// Variants are the predefined possible values for the option.
	//
	// All Variants should be unique, and an empty string is a valid variant.
	Variants []string
}

func (cmd *ComboOption) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "combo")
	text.WriteString(" default ")
	text.WriteString(cmd.DefaultValue)
	for _, v := range cmd.Variants {
		text.WriteString(" var ")
		text.WriteString(v)
	}
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *ComboOption) optionName() string {
	return cmd.Name
}

// ButtonOption represents an option of type Button that the engine supports.
//
// Button options have no parameters.
// Like a button, they are simply pressed.
//
// See [Option] for more info.
type ButtonOption struct {
	// Name is the name of the option.
	//
	// Certain options are expected to be of a specific type.
	// See [Option] for more info.
	Name string
}

func (cmd *ButtonOption) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "button")
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *ButtonOption) optionName() string {
	return cmd.Name
}

// StringOption represents an option of type String that the engine supports.
//
// String options can represent any string.
//
// See [Option] for more info.
type StringOption struct {
	// Name is the name of the option.
	//
	// Certain options are expected to be of a specific type.
	// See [Option] for more info.
	Name string
	// DefaultValue is the default setting the engine uses for this option.
	//
	// It may be an empty string.
	DefaultValue string
}

func (cmd *StringOption) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "string")
	text.WriteString(" default ")
	if len(cmd.DefaultValue) > 0 {
		text.WriteString(cmd.DefaultValue)
	} else {
		text.WriteString("<empty>")
	}
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *StringOption) optionName() string {
	return cmd.Name
}

// startOptionCmd starts marshalling an option command.
//
// It is the same for every type. "option name <id> type <t>"
func startOptionCmd(text *bytes.Buffer, name string, optionType string) {
	text.WriteString("option name ")
	text.WriteString(name)
	text.WriteString(" type ")
	text.WriteString(optionType)
}
