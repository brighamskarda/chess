// Copyright (C) 2026 Brigham Skarda

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
	"bytes"
	"fmt"
	"strconv"

	"github.com/brighamskarda/chess/v2"
)

// clientToEngineCmd is an interface under which all uci commands from the engine to the client will be contained.
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
	// isAuthor is true if this is an author id command. If it is false then it is assumed to be the name id command.
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

// readyokCmd
// This must be sent when the engine has received an "isready" command and has
// processed all input and is ready to accept new commands now.
// It is usually sent after a command that can take some time to be able to wait for the engine,
// but it can be used anytime, even when the engine is searching,
// and must always be answered with "isready".
type readyokCmd struct{}

func (cmd *readyokCmd) marshalText() ([]byte, error) {
	return []byte("readyok\n"), nil
}

// bestMoveCmd - bestmove <move1> [ ponder <move2> ]
// the engine has stopped searching and found the move <move> best in this position.
// the engine can send the move it likes to ponder on. The engine must not start pondering automatically.
// this command must always be sent if the engine stops searching, also in pondering mode if there is a
// "stop" command, so for every "go" command a "bestmove" command is needed!
// Directly before that the engine should send a final "info" command with the final search information,
// the the GUI has the complete statistics about the last search.
type bestMoveCmd struct {
	// move is the best move the engine found at the end of its search
	move chess.Move
	// ponderMove is an optional move that can be sent to indicate what move the engine would like to ponder on.
	ponderMove *chess.Move
}

func (cmd *bestMoveCmd) marshalText() ([]byte, error) {
	text := bytes.Buffer{}
	text.WriteString("bestmove ")

	move, err := cmd.move.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("could not marshal bestMoveCmd: %w", err)
	}
	text.Write(move)

	if cmd.ponderMove != nil {
		text.WriteString(" ponder ")

		ponderMove, err := cmd.ponderMove.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("could not marshal bestMoveCmd: %w", err)
		}
		text.Write(ponderMove)
	}

	text.WriteByte('\n')
	return text.Bytes(), nil
}

// copyprotectionCmd - this is needed for copyprotected engines. After the uciok command the engine can tell the GUI,
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
type copyprotectionCmd uint8

const (
	copyprotectChecking copyprotectionCmd = iota
	copyprotectOk
	copyprotectError
)

func (cmd *copyprotectionCmd) marshalText() ([]byte, error) {
	switch *cmd {
	case copyprotectChecking:
		return []byte("copyprotection checking\n"), nil
	case copyprotectOk:
		return []byte("copyprotection ok\n"), nil
	case copyprotectError:
		return []byte("copyprotection error\n"), nil
	default:
		return nil, fmt.Errorf("could not marshal copyprotection command: invalid value %v", *cmd)
	}
}

// registrationCmd - this is needed for engines that need a username and/or a code to function with all features.
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

// infoCmd - the engine wants to send information to the GUI. This should be done whenever one of the info has changed.
// The engine can send only selected infos or multiple infos with one info command,
// e.g. "info currmove e2e4 currmovenumber 1" or
//
//	"info depth 12 nodes 123456 nps 100000".
//
// Also all infos belonging to the pv should be sent together
// e.g. "info depth 2 score cp 214 time 1242 nodes 2124 nps 34928 pv e2e4 e7e5 g1f3"
// I suggest to start sending "currmove", "currmovenumber", "currline" and "refutation" only after one second
// to avoid too much traffic.
type infoCmd struct {
	// depth <x> - search depth in plies
	depth Optional[int]

	// seldepth <x> - selective search depth in plies,
	// if the engine sends seldepth there must also be a "depth" present in the same string.
	seldepth Optional[int]

	// time <x> - the time searched in ms, this should be sent together with the pv.
	time Optional[int]

	// nodes <x> - x nodes searched, the engine should send this info regularly
	nodes Optional[int]

	// pv <move1> ... <movei> - the best line found
	pv []chess.Move

	// multipv <num> - this for the multi pv mode.
	// for the best move/pv add "multipv 1" in the string when you send the pv.
	// in k-best mode always send all k variants in k strings together.
	multipv Optional[int]

	// score
	// 	* cp <x>
	// 		the score from the engine's point of view in centipawns.
	// 	* mate <y>
	// 		mate in y moves, not plies.
	// 		If the engine is getting mated use negative values for y.
	// 	* lowerbound
	//       the score is just a lower bound.
	// 	* upperbound
	// 	   the score is just an upper bound.
	score Optional[infoScore]

	// currmove <move> - currently searching this move
	currmove Optional[chess.Move]

	// currmovenumber <x> - currently searching move number x,
	// for the first move x should be 1 not 0.
	currmovenumber Optional[int]

	// hashfull <x> - the hash is x permill full, the engine should send this info regularly
	hashfull Optional[int]

	// nps <x> - x nodes per second searched, the engine should send this info regularly
	nps Optional[int]

	// tbhits <x> - x positions where found in the endgame table bases
	tbhits Optional[int]

	// sbhits <x> - x positions where found in the shredder endgame databases
	sbhits Optional[int]

	// cpuload <x> - the cpu usage of the engine is x permill.
	cpuload Optional[int]

	// string <str> - any string str which will be displayed by the engine,
	// if there is a string command the rest of the line will be interpreted as <str>.
	stringMsg Optional[string]

	// refutation <move1> <move2> ... <movei>
	// move <move1> is refuted by the line <move2> ... <movei>, i can be any number >= 1.
	// Example: after move d1h5 is searched, the engine can send "info refutation d1h5 g6h5"
	// if g6h5 is the best answer after d1h5 or if g6h5 refutes the move d1h5.
	// if there is no refutation for d1h5 found, the engine should just send "info refutation d1h5"
	// The engine should only send this if the option "UCI_ShowRefutations" is set to true.
	refutation []chess.Move

	// currline <cpunr> <move1> ... <movei>
	// this is the current line the engine is calculating. <cpunr> is the number of the cpu
	// if the engine is running on more than one cpu. <cpunr> = 1,2,3....
	// if the engine is just using one cpu, <cpunr> can be omitted.
	// If <cpunr> is greater than 1, always send all k lines in k strings together.
	// The engine should only send this if the option "UCI_ShowCurrLine" is set to true.
	currline Optional[currentLine]
}

func (cmd *infoCmd) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	text.WriteString("info")

	cmd.marshalDepth(text)
	cmd.marshalSeldepth(text)
	cmd.marshalTime(text)
	cmd.marshalNodes(text)
	cmd.marshalPv(text)
	cmd.marshalMultipv(text)
	cmd.marshalScore(text)
	cmd.marshalCurrmove(text)
	cmd.marshalCurrmovenumber(text)
	cmd.marshalHashfull(text)
	cmd.marshalNps(text)
	cmd.marshalTbhits(text)
	cmd.marshalSbhits(text)
	cmd.marshalCpuload(text)
	cmd.marshalRefutation(text)
	cmd.marshalCurrline(text)
	// it is important that the string message is marshaled last as everything after it is considered part of the string to the client.
	cmd.marshalString(text)

	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *infoCmd) marshalDepth(text *bytes.Buffer) {
	if cmd.depth.HasValue() {
		text.WriteString(" depth ")
		text.WriteString(strconv.Itoa(cmd.depth.Value()))
	}
}

func (cmd *infoCmd) marshalSeldepth(text *bytes.Buffer) {
	if cmd.seldepth.HasValue() {
		text.WriteString(" seldepth ")
		text.WriteString(strconv.Itoa(cmd.seldepth.Value()))
	}
}

func (cmd *infoCmd) marshalTime(text *bytes.Buffer) {
	if cmd.time.HasValue() {
		text.WriteString(" time ")
		text.WriteString(strconv.Itoa(cmd.time.Value()))
	}
}

func (cmd *infoCmd) marshalNodes(text *bytes.Buffer) {
	if cmd.nodes.HasValue() {
		text.WriteString(" nodes ")
		text.WriteString(strconv.Itoa(cmd.nodes.Value()))
	}
}

func marshalMoveList(moves []chess.Move, text *bytes.Buffer) {
	for _, m := range moves {
		moveText, _ := m.MarshalText()
		text.WriteByte(' ')
		text.Write(moveText)
	}
}

func (cmd *infoCmd) marshalPv(text *bytes.Buffer) {
	if cmd.pv != nil {
		text.WriteString(" pv")
		marshalMoveList(cmd.pv, text)
	}
}

func (cmd *infoCmd) marshalMultipv(text *bytes.Buffer) {
	if cmd.multipv.HasValue() {
		text.WriteString(" multipv ")
		text.WriteString(strconv.Itoa(cmd.multipv.Value()))
	}
}

func (cmd *infoCmd) marshalScore(text *bytes.Buffer) {
	if cmd.score.HasValue() {
		text.WriteString(" score ")

		value := cmd.score.Value()

		if value.isMate {
			text.WriteString("mate ")
		} else {
			text.WriteString("cp ")
		}

		text.WriteString(strconv.Itoa(value.score))

		if value.isLowerbound {
			text.WriteString(" lowerbound")
		}
		if value.isUpperbound {
			text.WriteString(" upperbound")
		}
	}
}

func (cmd *infoCmd) marshalCurrmove(text *bytes.Buffer) {
	if cmd.currmove.HasValue() {
		text.WriteString(" currmove ")
		moveText, _ := cmd.currmove.Value().MarshalText()
		text.Write(moveText)
	}
}

func (cmd *infoCmd) marshalCurrmovenumber(text *bytes.Buffer) {
	if cmd.currmovenumber.HasValue() {
		text.WriteString(" currmovenumber ")
		text.WriteString(strconv.Itoa(cmd.currmovenumber.Value()))
	}
}

func (cmd *infoCmd) marshalHashfull(text *bytes.Buffer) {
	if cmd.hashfull.HasValue() {
		text.WriteString(" hashfull ")
		text.WriteString(strconv.Itoa(cmd.hashfull.Value()))
	}
}

func (cmd *infoCmd) marshalNps(text *bytes.Buffer) {
	if cmd.nps.HasValue() {
		text.WriteString(" nps ")
		text.WriteString(strconv.Itoa(cmd.nps.Value()))
	}
}

func (cmd *infoCmd) marshalTbhits(text *bytes.Buffer) {
	if cmd.tbhits.HasValue() {
		text.WriteString(" tbhits ")
		text.WriteString(strconv.Itoa(cmd.tbhits.Value()))
	}
}

func (cmd *infoCmd) marshalSbhits(text *bytes.Buffer) {
	if cmd.sbhits.HasValue() {
		text.WriteString(" sbhits ")
		text.WriteString(strconv.Itoa(cmd.sbhits.Value()))
	}
}

func (cmd *infoCmd) marshalCpuload(text *bytes.Buffer) {
	if cmd.cpuload.HasValue() {
		text.WriteString(" cpuload ")
		text.WriteString(strconv.Itoa(cmd.cpuload.Value()))
	}
}

func (cmd *infoCmd) marshalRefutation(text *bytes.Buffer) {
	if cmd.refutation != nil {
		text.WriteString(" refutation")
		marshalMoveList(cmd.refutation, text)
	}
}

func (cmd *infoCmd) marshalCurrline(text *bytes.Buffer) {
	if cmd.currline.HasValue() {
		currLine := cmd.currline.Value()

		text.WriteString(" currline")
		if currLine.cpunr.HasValue() {
			text.WriteByte(' ')
			text.WriteString(strconv.Itoa(currLine.cpunr.Value()))
		}

		marshalMoveList(currLine.moves, text)
	}
}

func (cmd *infoCmd) marshalString(text *bytes.Buffer) {
	if cmd.stringMsg.HasValue() {
		text.WriteString(" string ")
		text.WriteString(cmd.stringMsg.Value())
	}
}

// currentLine is used in [infoCmd] to show the current line the engine is calculating.
type currentLine struct {
	cpunr Optional[int]
	moves []chess.Move
}

// score is used in [infoCmd] to show the current line the engine is calculating.
type infoScore struct {
	// score is the score
	score int
	// isMate is true if the score represents how many plies until mate. Otherwise score is assumed to be the engines evaluation in centipawns.
	isMate bool
	// isLowerbound indicates that this score is a isLowerbound. Should be false if upperbound is set.
	isLowerbound bool
	// isUpperbound indicates that this score is an isUpperbound. Should be false if lowerbound is set.
	isUpperbound bool
}

// OptionCmd represents an option that the chess engine supports. This is sent to the client so it knows what options it can set in the engine.
//
// There is no need to make any types implementing this interface as all valid options can be represented using the structs provided in this module. See [CheckOptionCmd], [SpinOptionCmd], [ComboOptionCmd], [StringOptionCmd], and [ButtonOptionCmd]
//
// The following options are predefined in the UCI chess standard. Options with these names should not be used for anything else. Furthermore any options not listed below that are still prepended with "UCI_" will likely be ignored by the GUI.
//   - <id> = Hash, type is spin
//     the value in MB for memory for hash tables can be changed,
//     this should be answered with the first "setoptions" command at program boot
//     if the engine has sent the appropriate "option name Hash" command,
//     which should be supported by all engines!
//     So the engine should use a very small hash first as default.
//   - <id> = NalimovPath, type string
//     this is the path on the hard disk to the Nalimov compressed format.
//     Multiple directories can be concatenated with ";"
//   - <id> = NalimovCache, type spin
//     this is the size in MB for the cache for the nalimov table bases
//     These last two options should also be present in the initial options exchange dialog
//     when the engine is booted if the engine supports it
//   - <id> = Ponder, type check
//     this means that the engine is able to ponder.
//     The GUI will send this whenever pondering is possible or not.
//     Note: The engine should not start pondering on its own if this is enabled, this option is only
//     needed because the engine might change its time management algorithm when pondering is allowed.
//   - <id> = OwnBook, type check
//     this means that the engine has its own book which is accessed by the engine itself.
//     if this is set, the engine takes care of the opening book and the GUI will never
//     execute a move out of its book for the engine. If this is set to false by the GUI,
//     the engine should not access its own book.
//   - <id> = MultiPV, type spin
//     the engine supports multi best line or k-best mode. the default value is 1
//   - <id> = UCI_ShowCurrLine, type check, should be false by default,
//     the engine can show the current line it is calculating. see "info currline" above.
//   - <id> = UCI_ShowRefutations, type check, should be false by default,
//     the engine can show a move and its refutation in a line. see "info refutations" above.
//   - <id> = UCI_LimitStrength, type check, should be false by default,
//     The engine is able to limit its strength to a specific Elo number,
//     This should always be implemented together with "UCI_Elo".
//   - <id> = UCI_Elo, type spin
//     The engine can limit its strength in Elo within this interval.
//     If UCI_LimitStrength is set to false, this value should be ignored.
//     If UCI_LimitStrength is set to true, the engine should play with this specific strength.
//     This should always be implemented together with "UCI_LimitStrength".
//   - <id> = UCI_AnalyseMode, type check
//     The engine wants to behave differently when analysing or playing a game.
//     For example when playing it can use some kind of learning.
//     This is set to false if the engine is playing a game, otherwise it is true.
//   - <id> = UCI_Opponent, type string
//     With this command the GUI can send the name, title, elo and if the engine is playing a human
//     or computer to the engine.
//     The format of the string has to be [GM|IM|FM|WGM|WIM|none] [<elo>|none] [computer|human] <name>
//     Examples:
//     "setoption name UCI_Opponent value GM 2800 human Gary Kasparov"
//     "setoption name UCI_Opponent value none none computer Shredder"
//   - <id> = UCI_EngineAbout, type string
//     With this command, the engine tells the GUI information about itself, for example a license text,
//     usually it doesn't make sense that the GUI changes this text with the setoption command.
//     Example:
//     "option name UCI_EngineAbout type string default Shredder by Stefan Meyer-Kahlen, see www.shredderchess.com"
//   - <id> = UCI_ShredderbasesPath, type string
//     this is either the path to the folder on the hard disk containing the Shredder endgame databases or
//     the path and filename of one Shredder endgame datbase.
//   - <id> = UCI_SetPositionValue, type string
//     the GUI can send this to the engine to tell the engine to use a certain value in centipawns from white's
//     point of view if evaluating this specifix position.
//     The string can have the formats:
//     <value> + <fen> | clear + <fen> | clearall
type OptionCmd interface {
	engineToClientCmd
	// optionName should return the name/id of the option.
	optionName() string
}

// CheckOptionCmd - a checkbox that can either be true or false
//
// see [OptionCmd]
type CheckOptionCmd struct {
	// Name <id> - the option has the Name id.
	//
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// See [OptionCmd] for more info
	Name string
	// DefaultValue - the default value of this parameter is x
	DefaultValue bool
}

func (cmd *CheckOptionCmd) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "check")
	text.WriteString(" default ")
	text.WriteString(strconv.FormatBool(cmd.DefaultValue))
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *CheckOptionCmd) optionName() string {
	return cmd.Name
}

// SpinOptionCmd - a spin wheel that can be an integer in a certain range
//
// see [OptionCmd]
type SpinOptionCmd struct {
	// Name <id> - the option has the Name id.
	//
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// See [OptionCmd] for more info
	Name string
	// DefaultValue - the default value of this parameter is x
	DefaultValue int
	// Min - the minimum value of this parameter is x
	Min int
	// Max - the maximum value of this parameter is x
	Max int
}

func (cmd *SpinOptionCmd) marshalText() ([]byte, error) {
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

func (cmd *SpinOptionCmd) optionName() string {
	return cmd.Name
}

// ComboOptionCmd - a combo box that can have different predefined strings as a value
//
// see [OptionCmd]
type ComboOptionCmd struct {
	// Name <id> - the option has the Name id.
	//
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// See [OptionCmd] for more info
	Name string
	// DefaultValue - the default value of this parameter is x
	DefaultValue string
	// Variants - the predefined possible values for the parameter
	Variants []string
}

func (cmd *ComboOptionCmd) marshalText() ([]byte, error) {
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

func (cmd *ComboOptionCmd) optionName() string {
	return cmd.Name
}

// ButtonOptionCmd - a button that can be pressed to send a command to the engine
//
// see [OptionCmd]
type ButtonOptionCmd struct {
	// Name <id> - the option has the Name id.
	//
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// See [OptionCmd] for more info
	Name string
}

func (cmd *ButtonOptionCmd) marshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.Name, "button")
	text.WriteByte('\n')
	return text.Bytes(), nil
}

func (cmd *ButtonOptionCmd) optionName() string {
	return cmd.Name
}

// StringOptionCmd -a text field that has a string as a value, an empty string has the value "<empty>"
//
// see [OptionCmd]
type StringOptionCmd struct {
	// Name <id> - the option has the Name id.
	//
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// See [OptionCmd] for more info
	Name string
	// DefaultValue - the default value of this parameter is x. An empty string will be encoded to <empty>.
	DefaultValue string
}

func (cmd *StringOptionCmd) marshalText() ([]byte, error) {
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

func (cmd *StringOptionCmd) optionName() string {
	return cmd.Name
}

// startOptionCmd starts marshalling an option command. It is the same for every type. "option name <id> type <t>"
func startOptionCmd(text *bytes.Buffer, name string, optionType string) {
	text.WriteString("option name ")
	text.WriteString(name)
	text.WriteString(" type ")
	text.WriteString(optionType)
}
