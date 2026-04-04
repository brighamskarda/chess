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
	"encoding"
	"fmt"
	"strconv"

	"github.com/brighamskarda/chess/v2"
)

// clientToEngineCmd is an interface under which all uci commands from the engine to the client will be contained.
type engineToClientCmd interface {
	encoding.TextMarshaler
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

func (cmd *idCmd) MarshalText() ([]byte, error) {
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

func (cmd *uciokCmd) MarshalText() ([]byte, error) {
	return []byte("uciok\n"), nil
}

// readyokCmd
// This must be sent when the engine has received an "isready" command and has
// processed all input and is ready to accept new commands now.
// It is usually sent after a command that can take some time to be able to wait for the engine,
// but it can be used anytime, even when the engine is searching,
// and must always be answered with "isready".
type readyokCmd struct{}

func (cmd *readyokCmd) MarshalText() ([]byte, error) {
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

func (cmd *bestMoveCmd) MarshalText() ([]byte, error) {
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

func (cmd *copyprotectionCmd) MarshalText() ([]byte, error) {
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

func (cmd *registrationCmd) MarshalText() ([]byte, error) {
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

func (cmd *infoCmd) MarshalText() ([]byte, error) {
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

// checkOptionCmd - a checkbox that can either be true or false
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
type checkOptionCmd struct {
	// name <id>
	// The option has the name id.
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// Usually those options should not be displayed in the normal engine options window of the GUI but
	// get a special treatment. "Pondering" for example should be set automatically when pondering is
	// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// automatically by the GUI.
	//
	// There are global variables defined for these predefined options.
	name string
	// defaultValue - the default value of this parameter is x
	defaultValue bool
}

func (cmd *checkOptionCmd) MarshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.name, "check")
	text.WriteString(" default ")
	text.WriteString(strconv.FormatBool(cmd.defaultValue))
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// spinOptionCmd - a spin wheel that can be an integer in a certain range
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
type spinOptionCmd struct {
	// name <id>
	// The option has the name id.
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// Usually those options should not be displayed in the normal engine options window of the GUI but
	// get a special treatment. "Pondering" for example should be set automatically when pondering is
	// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// automatically by the GUI.
	//
	// There are global variables defined for these predefined options.
	name string
	// defaultValue - the default value of this parameter is x
	defaultValue int
	// min - the minimum value of this parameter is x
	min int
	// max - the maximum value of this parameter is x
	max int
}

func (cmd *spinOptionCmd) MarshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.name, "spin")
	text.WriteString(" default ")
	text.WriteString(strconv.Itoa(cmd.defaultValue))
	text.WriteString(" min ")
	text.WriteString(strconv.Itoa(cmd.min))
	text.WriteString(" max ")
	text.WriteString(strconv.Itoa(cmd.max))
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// comboOptionCmd - a combo box that can have different predefined strings as a value
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
type comboOptionCmd struct {
	// name <id>
	// The option has the name id.
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// Usually those options should not be displayed in the normal engine options window of the GUI but
	// get a special treatment. "Pondering" for example should be set automatically when pondering is
	// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// automatically by the GUI.
	//
	// There are global variables defined for these predefined options.
	name string
	// defaultValue - the default value of this parameter is x
	defaultValue string
	// variants - the predefined possible values for the parameter
	variants []string
}

func (cmd *comboOptionCmd) MarshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.name, "combo")
	text.WriteString(" default ")
	text.WriteString(cmd.defaultValue)
	for _, v := range cmd.variants {
		text.WriteString(" var ")
		text.WriteString(v)
	}
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// buttonOptionCmd - a button that can be pressed to send a command to the engine
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
type buttonOptionCmd struct {
	// name <id>
	// The option has the name id.
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// Usually those options should not be displayed in the normal engine options window of the GUI but
	// get a special treatment. "Pondering" for example should be set automatically when pondering is
	// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// automatically by the GUI.
	//
	// There are global variables defined for these predefined options.
	name string
}

func (cmd *buttonOptionCmd) MarshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.name, "button")
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// stringOptionCmd -a text field that has a string as a value, an empty string has the value "<empty>"
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
type stringOptionCmd struct {
	// name <id>
	// The option has the name id.
	// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
	// Usually those options should not be displayed in the normal engine options window of the GUI but
	// get a special treatment. "Pondering" for example should be set automatically when pondering is
	// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// automatically by the GUI.
	//
	// There are global variables defined for these predefined options.
	name string
	// defaultValue - the default value of this parameter is x. An empty string will be encoded to <empty>.
	defaultValue string
}

func (cmd *stringOptionCmd) MarshalText() ([]byte, error) {
	text := &bytes.Buffer{}
	startOptionCmd(text, cmd.name, "string")
	text.WriteString(" default ")
	if len(cmd.defaultValue) > 0 {
		text.WriteString(cmd.defaultValue)
	} else {
		text.WriteString("<empty>")
	}
	text.WriteByte('\n')
	return text.Bytes(), nil
}

// startOptionCmd starts marshalling an option command. It is the same for every type. "option name <id> type <t>"
func startOptionCmd(text *bytes.Buffer, name string, optionType string) {
	text.WriteString("option name ")
	text.WriteString(name)
	text.WriteString(" type ")
	text.WriteString(optionType)
}
