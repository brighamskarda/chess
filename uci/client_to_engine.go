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
	"encoding"
	"fmt"
	"strconv"
	"unicode"

	"github.com/brighamskarda/chess/v2"
)

// clientToEngineCmd is an interface under which all uci commands from the client to the engine will be contained.
type clientToEngineCmd interface {
	encoding.TextUnmarshaler
	// getCmdText returns the full message received to construct the command.
	//
	// The message should be copied so the user can modify it.
	getCmdText() []byte
}

// baseClientCommand provides fields that should be in all uci commands.
type baseClientCommand struct {
	// message represents the unmodified message that was received or sent.
	message []byte
}

// getCmdText is a base implementation of the function for the [clientToEngineCmd] interface.
func (cmd *baseClientCommand) getCmdText() []byte {
	return bytes.Clone(cmd.message)
}

// initBaseEngineCommand initializes the base engine command in the same way for any type that contains it.
//
// It ensures the text is copied so it can't be modified by outside sources.
func (cmd *baseClientCommand) initBaseEngineCommand(text []byte) {
	cmd.message = bytes.Clone(text)
}

// uciCmd tells engine to enter uci mode uci.
//
// This will be sent once as a first command after program boot to tell the engine to switch to uciCmd mode.
// After receiving the uciCmd command the engine must identify itself with the "id" command and send the "option" commands
// to tell the GUI which engine settings the engine supports if any.
// After that the engine should send "uciok" to acknowledge the uciCmd mode.
// If no uciok is sent within a certain time period, the engine task will be killed by the GUI.
type uciCmd struct {
	baseClientCommand
}

func (cmd *uciCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "uci" {
		return fmt.Errorf("could not unmarshal uci command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// debugCmd [ on | off ]
//
// Switch the debugCmd mode of the engine on and off.
// In debugCmd mode the engine should send additional infos to the GUI, e.g. with the "info string" command,
// to help debugging, e.g. the commands that the engine has received etc.
// This mode should be switched off by default and this command can be sent any time, also when the engine is thinking.
type debugCmd struct {
	baseClientCommand
	// on is true if the engine should be set to debug mode.
	on bool
}

func (cmd *debugCmd) UnmarshalText(text []byte) error {
	words := bytes.Fields(text)
	if len(words) != 2 {
		return fmt.Errorf("could not unmarshal debug command %q: expected only two fields", text)
	}
	if string(words[0]) != "debug" {
		return fmt.Errorf("could not unmarshal debug command %q: expected first field to be \"debug\"", text)
	}

	switch string(words[1]) {
	case "on":
		cmd.on = true
	case "off":
		cmd.on = false
	default:
		return fmt.Errorf("could not unmarshal debug command %q: should be set to \"on\" or \"off\".", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// isReadyCmd is used to synchronize the engine with the GUI.
//
// When the GUI has sent a command or multiple commands that can take some time to complete,
// this command can be used to wait for the engine to be ready again or to ping the engine to find out if it is still alive.
// E.g. this should be sent after setting the path to the tablebases as this can take some time.
// This command is also required once before the engine is asked to do any search to wait for the engine to finish initializing.
// This command must always be answered with "readyok"
// and can be sent also when the engine is calculating in which case the engine should also immediately answer with "readyok" without stopping the search.
type isReadyCmd struct {
	baseClientCommand
}

func (cmd *isReadyCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "isready" {
		return fmt.Errorf("could not unmarshal isready command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// validateSetOption ensure the first two words are setoption and name.
//
// It returns an error if they are not.
func validateSetOption(text []byte) error {
	words := bytes.Fields(text)
	if len(words) < 3 {
		return fmt.Errorf("could not unmarshal setoption command %q: set option commands should be at least three words", text)
	}
	if string(words[0]) != "setoption" || string(words[1]) != "name" {
		return fmt.Errorf("could not unmarshal setoption command %q: does not begin with \"setoption name\"", text)
	}
	if string(words[2]) == "value" {
		return fmt.Errorf("could not unmarshal setoption command %q: 'value' is not a valid option name", text)

	}

	return nil
}

// getOptionName parses the option id from the text.
func getOptionName(text []byte) string {
	startIndex := bytes.Index(text, []byte("name")) + 5
	endIndex := bytes.Index(text, []byte("value"))
	if endIndex == -1 {
		endIndex = len(text)
	}
	return string(bytes.TrimSpace(text[startIndex:endIndex]))
}

// getOptionValue parses the option value from the text.
func getOptionValue(text []byte) string {
	startIndex := bytes.Index(text, []byte("value"))
	if startIndex == -1 {
		return ""
	}
	return string(bytes.TrimSpace(text[startIndex+6:]))
}

// SetOptionCmd is sent to change one of an engine's configurable options.
//
// When passed to the engine this will be one of four types.
//   - [SetCheckOptionCmd]
//   - [SetSpinOptionCmd]
//   - [SetStringOptionCmd] (which double as setting a combo option)
//   - [SetButtonOptionCmd]
type SetOptionCmd interface {
	clientToEngineCmd
	// optionName returns the name/id of the option.
	OptionName() string
}

// SetCheckOptionCmd is an option with the check type.
//
//   - check - a checkbox that can either be true or false
type SetCheckOptionCmd struct {
	baseClientCommand
	// name is the name/id of the option being set.
	name string
	// Is true if the "checkbox" for this command is set.
	Checkbox bool
}

func (cmd *SetCheckOptionCmd) UnmarshalText(text []byte) error {
	err := validateSetOption(text)
	if err != nil {
		return err
	}

	cmd.name = getOptionName(text)
	if cmd.name == "" {
		return fmt.Errorf("could not unmarshal setOption command %q: empty option name", text)
	}

	value := getOptionValue(text)

	switch value {
	case "true":
		cmd.Checkbox = true
	case "false":
		cmd.Checkbox = false
	default:
		return fmt.Errorf("could not unmarshal setOption command %q: check type options should only have values of true or false", text)
	}

	cmd.initBaseEngineCommand(text)
	return nil
}

func (cmd *SetCheckOptionCmd) OptionName() string {
	return cmd.name
}

// SetSpinOptionCmd is an option with the spin type.
//
//   - spin - a spin wheel that can be an integer in a certain range
type SetSpinOptionCmd struct {
	baseClientCommand
	// name is the name/id of the option being set.
	name string
	// value represents the numeric integer value for the spin option.
	Value int
}

func (cmd *SetSpinOptionCmd) UnmarshalText(text []byte) error {
	err := validateSetOption(text)
	if err != nil {
		return err
	}

	cmd.name = getOptionName(text)
	if cmd.name == "" {
		return fmt.Errorf("could not unmarshal setOption command %q: empty option name", text)
	}

	value := getOptionValue(text)

	i, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return fmt.Errorf("could not unmarshal setOption command %q: spin type options should only represent integer values", text)
	}

	cmd.Value = int(i)
	cmd.initBaseEngineCommand(text)
	return nil
}

func (cmd *SetSpinOptionCmd) OptionName() string {
	return cmd.name
}

// SetStringOptionCmd is an option with the string **or combo** type.
//
//   - string - a text field that has a string as a value, an empty string has the value "<empty>"
type SetStringOptionCmd struct {
	baseClientCommand
	// name is the name/id of the option being set.
	name string
	// Value is the string content assigned to this option. If <empty> was sent then this will be an empty string.
	Value string
}

func (cmd *SetStringOptionCmd) UnmarshalText(text []byte) error {
	err := validateSetOption(text)
	if err != nil {
		return err
	}

	cmd.name = getOptionName(text)
	if cmd.name == "" {
		return fmt.Errorf("could not unmarshal setOption command %q: empty option name", text)
	}

	value := getOptionValue(text)

	switch value {
	case "":
		return fmt.Errorf("could not unmarshal setStringOption command %q: empty value", text)
	case "<empty>":
		cmd.Value = ""
	default:
		cmd.Value = value
	}

	cmd.initBaseEngineCommand(text)
	return nil
}

func (cmd *SetStringOptionCmd) OptionName() string {
	return cmd.name
}

// SetButtonOptionCmd is an option with the button type.
//
//   - button - a button that can be pressed to send a command to the engine
type SetButtonOptionCmd struct {
	baseClientCommand
	// name is the name/id of the button to be "pressed".
	name string
}

func (cmd *SetButtonOptionCmd) UnmarshalText(text []byte) error {
	err := validateSetOption(text)
	if err != nil {
		return err
	}

	cmd.name = getOptionName(text)
	if cmd.name == "" {
		return fmt.Errorf("could not unmarshal setOption command %q: empty option name", text)
	}

	cmd.initBaseEngineCommand(text)
	return nil
}

func (cmd *SetButtonOptionCmd) OptionName() string {
	return cmd.name
}

// RegisterCmd is the command to register an engine or to tell it that registration
// will be done later.
//
// The following tokens are allowed:
//   - later
//     the user doesn't want to register the engine now.
//   - name <x>
//     the engine should be registered with the name <x>
//   - code <y>
//     the engine should be registered with the code <y>
//
// Example:
//
//	"register later"
//	"register name Stefan MK code 4359874324"
type RegisterCmd struct {
	baseClientCommand
	// Later indicates the user doesn't want to register the engine now.
	Later bool
	// Name indicates the name the engine should be registered with.
	Name Optional[string]
	// Code indicates the registration code.
	Code Optional[string]
}

func (cmd *RegisterCmd) UnmarshalText(text []byte) error {
	wordsIterator := bytes.FieldsSeq(text)
	firstTwoWords := make([]string, 0, 2)
	wordsIterator(func(val []byte) bool {
		firstTwoWords = append(firstTwoWords, string(val))
		return len(firstTwoWords) < 2
	})

	if len(firstTwoWords) != 2 || firstTwoWords[0] != "register" {
		return fmt.Errorf("could not unmarshal register command %q: should start with \"register <later|name|code>\"", text)
	}

	switch string(firstTwoWords[1]) {
	case "later":
		cmd.Later = true
		cmd.Name = OptionalEmpty[string]()
		cmd.Code = OptionalEmpty[string]()
	case "name", "code":
		cmd.Later = false
		nameIndex := bytes.Index(text, []byte("name"))
		codeIndex := bytes.Index(text, []byte("code"))
		if codeIndex == -1 {
			cmd.Name = OptionalOf(string(bytes.TrimSpace(text[nameIndex+5:])))
			cmd.Code = OptionalEmpty[string]()
		} else if nameIndex == -1 {
			cmd.Name = OptionalEmpty[string]()
			cmd.Code = OptionalOf(string(bytes.TrimSpace(text[codeIndex+5:])))
		} else if nameIndex < codeIndex {
			cmd.Name = OptionalOf(string(bytes.TrimSpace(text[nameIndex+5 : codeIndex])))
			cmd.Code = OptionalOf(string(bytes.TrimSpace(text[codeIndex+5:])))
		} else if nameIndex > codeIndex {
			cmd.Name = OptionalOf(string(bytes.TrimSpace(text[nameIndex+5:])))
			cmd.Code = OptionalOf(string(bytes.TrimSpace(text[codeIndex+5 : nameIndex])))
		}
	default:
		return fmt.Errorf("could not unmarshal register command %q: should start with \"register <later|name|code>\"", text)
	}

	if (cmd.Name.HasValue() && len(cmd.Name.Value()) == 0) || (cmd.Code.HasValue() && len(cmd.Code.Value()) == 0) {
		return fmt.Errorf("could not unmarshal register command %q: name and code register commands should have a value", text)

	}

	cmd.initBaseEngineCommand(text)
	return nil
}

// uciNewGameCmd is sent to the engine when the next search (started with "position" and "go") will be from
// a different game.
//
// This can be a new game the engine should play or a new game it should analyse but
// also the next position from a testsuite with positions only.
// If the GUI hasn't sent a "ucinewgame" before the first "position" command, the engine shouldn't
// expect any further ucinewgame commands as the GUI is probably not supporting the ucinewgame command.
// So the engine should not rely on this command even though all new GUIs should support it.
// As the engine's reaction to "ucinewgame" can take some time the GUI should always send "isready"
// after "ucinewgame" to wait for the engine to finish its operation.
type uciNewGameCmd struct {
	baseClientCommand
}

func (cmd *uciNewGameCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "ucinewgame" {
		return fmt.Errorf("could not unmarshal ucinewgame command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// positionCmd [fen <fenstring> | startpos ]  moves <move1> .... <movei>
//
// Set up the position described in fenstring on the internal board and
// play the moves on the internal chess board.
// if the game was played  from the start position the string "startpos" will be sent
// Note: no "new" command is needed. However, if this position is from a different game than
// the last position sent to the engine, the GUI should have sent a "ucinewgame" inbetween.
type positionCmd struct {
	baseClientCommand
	// position is the position that should be set by the chess engine.
	position *chess.Position
	// moves are the moves that should be applied to the position.
	// If "moves" was not sent, then this is nil. If moves was sent, but was not followed by an moves then this is an empty slice.
	moves []chess.Move
}

func (cmd *positionCmd) UnmarshalText(text []byte) error {
	words := bytes.Fields(text)
	if len(words) < 2 ||
		string(words[0]) != "position" ||
		(string(words[1]) != "fen" && string(words[1]) != "startpos") {
		return fmt.Errorf("could not unmarshal position command %q: must start with \"position [fen | startpos]\"", text)
	}

	if string(words[1]) == "startpos" {
		pos := chess.Position{}
		pos.UnmarshalText([]byte(chess.DefaultFEN))
		cmd.position = &pos
	} else {
		if err := cmd.parseFen(text); err != nil {
			return fmt.Errorf("could not unmarshal position command %q: %w", text, err)
		}
	}

	err := cmd.parseMoves(text)
	if err != nil {
		return fmt.Errorf("could not unmarshal position command %q: %w", text, err)
	}

	cmd.initBaseEngineCommand(text)
	return nil
}

// parseFen finds the starting and ending of the fen string in the position command text and parses it.
//
// Does not support startpos. Returns an error if parsing was not successful.
func (cmd *positionCmd) parseFen(text []byte) error {
	startIndex := bytes.Index(text, []byte("fen")) + 4
	endIndex := bytes.Index(text, []byte("moves"))
	if endIndex == -1 {
		endIndex = len(text)
	}

	pos := chess.Position{}
	if err := pos.UnmarshalText(text[startIndex:endIndex]); err != nil {
		return err
	}

	cmd.position = &pos
	return nil
}

// parseMoves parses the moves from the position command.
//
// If moves does not exist, sets the moves to nil.
// If there was a problem unmarshaling any of the moves an error will be returned.
func (cmd *positionCmd) parseMoves(text []byte) error {
	startIndex := bytes.Index(text, []byte("moves"))
	if startIndex == -1 {
		cmd.moves = nil
		return nil
	}
	startIndex += 6

	fields := bytes.Fields(text[startIndex:])
	cmd.moves = make([]chess.Move, 0, len(fields))
	for _, f := range fields {
		move := chess.Move{}
		if err := move.UnmarshalText(f); err != nil {
			return err
		}
		cmd.moves = append(cmd.moves, move)
	}
	return nil
}

// goCmd start calculating on the current position set up with the "position" command.
//
// There are a number of commands that can follow this command, all will be sent in the same string.
// If one command is not sent its value should be interpreted as it would not influence the search.
//
// In this struct nil means that the command was not present.
type goCmd struct {
	baseClientCommand
	// searchmoves <move1> .... <movei>
	//
	// restrict search to this moves only
	// Example: After "position startpos" and "go infinite searchmoves e2e4 d2d4"
	// the engine should only search the two moves e2e4 and d2d4 in the initial position.
	searchMoves []chess.Move
	// ponder start searching in pondering mode.
	//
	// Do not exit the search in ponder mode, even if it's mate!
	// This means that the last move sent in in the position string is the ponder move.
	// The engine can do what it wants to do, but after a "ponderhit" command
	// it should execute the suggested move to ponder on. This means that the ponder move sent by
	// the GUI can be interpreted as a recommendation about which move to ponder. However, if the
	// engine decides to ponder on a different move, it should not display any mainlines as they are
	// likely to be misinterpreted by the GUI because the GUI expects the engine to ponder
	// on the suggested move.
	ponder bool
	// wtime - white has x msec left on the clock
	wtime Optional[int]
	// btime - black has x msec left on the clock
	btime Optional[int]
	// winc - white increment per move in mseconds if x > 0
	winc Optional[int]
	// binc - black increment per move in mseconds if x > 0
	binc Optional[int]
	// movestogo - there are x moves to the next time control, this will only be sent if x > 0,
	// if you don't get this and get the wtime and btime it's sudden death
	movestogo Optional[int]
	// depth - search x plies only
	depth Optional[int]
	// nodes - search x nodes only
	nodes Optional[int]
	// mate - search for a mate in x moves
	mate Optional[int]
	// movetime search exactly x mseconds
	movetime Optional[int]
	// infinite search until the "stop" command.
	//
	// Do not exit the search without being told so in this mode!
	infinite bool
}

func (cmd *goCmd) UnmarshalText(text []byte) error {
	*cmd = goCmd{}

	fields := bytes.Fields(text)
	if string(fields[0]) != "go" {
		return fmt.Errorf("could not unmarshal go command %q: expected \"go\" to be the first field", text)
	}

	if err := cmd.parseFields(fields); err != nil {
		return fmt.Errorf("could not unmarshal go command %q: %w", text, err)
	}

	cmd.initBaseEngineCommand(text)
	return nil
}

func (cmd *goCmd) parseFields(fields [][]byte) error {
	for i := 1; i < len(fields); i++ {
		f := string(fields[i])
		switch f {
		case "searchmoves":
			cmd.searchMoves = parseSearchMoves(fields[i+1:])
			i += len(cmd.searchMoves)
		case "ponder":
			cmd.ponder = true
		case "wtime", "btime", "winc", "binc", "movestogo", "depth", "nodes", "mate", "movetime":
			if err := cmd.parseFieldWithValue(fields[i:]); err != nil {
				return err
			}
			i++
		case "infinite":
			cmd.infinite = true
		default:
			return fmt.Errorf("could not unmarshal unknown command %q", f)
		}
	}
	return nil
}

// parseFieldWithValue parses goCmd fields that have a number after them.
//
// Expects the first field to be the name, and the second field to be the value.
// Returns an error if there is no second field, or there was a problem parsing it.
// A first field is expected though.
func (cmd *goCmd) parseFieldWithValue(fields [][]byte) error {
	fieldName := string(fields[0])
	if len(fields) < 2 {
		return fmt.Errorf("could not unmarshal %v: no value present", fieldName)
	}
	valueString := string(fields[1])
	v, err := strconv.Atoi(valueString)
	if err != nil {
		return fmt.Errorf("could not unmarshal %v %q: %w", fieldName, valueString, err)
	}

	switch fieldName {
	case "wtime":
		cmd.wtime = OptionalOf(v)
	case "btime":
		cmd.btime = OptionalOf(v)
	case "winc":
		cmd.winc = OptionalOf(v)
	case "binc":
		cmd.binc = OptionalOf(v)
	case "movestogo":
		cmd.movestogo = OptionalOf(v)
	case "depth":
		cmd.depth = OptionalOf(v)
	case "nodes":
		cmd.nodes = OptionalOf(v)
	case "mate":
		cmd.mate = OptionalOf(v)
	case "movetime":
		cmd.movetime = OptionalOf(v)
	default:
		return fmt.Errorf("could not unmarshal %v %q: unknown command, this indicates an error uci library, not user code.", fieldName, valueString)
	}
	return nil
}

// parseSearchMoves parses the fields into chess moves until a field fails to be parsed.
//
// The return value will never be null, though it may be empty.
// the number of fields parsed can be extrapolated from the length of the return value.
func parseSearchMoves(fields [][]byte) []chess.Move {
	moves := []chess.Move{}
	for _, f := range fields {
		var m chess.Move
		if m.UnmarshalText(f) != nil {
			// found an invalid move, this means it is probably a new command
			return moves
		}

		moves = append(moves, m)
	}
	return moves
}

// stopCmd - stop calculating as soon as possible
//
// don't forget the "bestmove" and possibly the "ponder" token when finishing the search
type stopCmd struct {
	baseClientCommand
}

func (cmd *stopCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "stop" {
		return fmt.Errorf("could not unmarshal stop command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// ponderhitCmd - the user has played the expected move.
//
// This will be sent if the engine was told to ponder on the same move the user has played.
// The engine should continue searching but switch from pondering to normal search.
type ponderhitCmd struct {
	baseClientCommand
}

func (cmd *ponderhitCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "ponderhit" {
		return fmt.Errorf("could not unmarshal ponderhit command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// quitCmd the program as soon as possible
type quitCmd struct {
	baseClientCommand
}

func (cmd *quitCmd) UnmarshalText(text []byte) error {
	if string(bytes.TrimSpace(text)) != "quit" {
		return fmt.Errorf("could not unmarshal quit command %q", text)
	}
	cmd.initBaseEngineCommand(text)
	return nil
}

// unmarshalOptionCommand first tries to unmarshal the check option, then the spin option, then the string option, then the button option.
func unmarshalOptionCommand(text []byte) (clientToEngineCmd, error) {
	// check option
	cmd, err1 := unmarshalClientToEngineCmd[SetCheckOptionCmd](text)
	if err1 == nil {
		return cmd, nil
	}

	// spin option
	cmd, err2 := unmarshalClientToEngineCmd[SetSpinOptionCmd](text)
	if err2 == nil {
		return cmd, nil
	}

	// string/combo option
	cmd, err3 := unmarshalClientToEngineCmd[SetStringOptionCmd](text)
	if err3 == nil {
		return cmd, nil
	}

	// button option
	cmd, err4 := unmarshalClientToEngineCmd[SetButtonOptionCmd](text)
	if err4 == nil {
		return cmd, nil
	}

	return nil, fmt.Errorf("could not unmarshal setoption command: "+
		"4 separate parses where attempted and all failed, each of there errors are provided:"+
		"\n\t%w\n\t%w\n\t%w\n\t%w", err1, err2, err3, err4)
}

// commandSet is used to lookup the specific parses for different commands based on the first word in the command.
var commandSet = map[string]func([]byte) (clientToEngineCmd, error){
	"uci":        unmarshalClientToEngineCmd[uciCmd],
	"debug":      unmarshalClientToEngineCmd[debugCmd],
	"isready":    unmarshalClientToEngineCmd[isReadyCmd],
	"setoption":  unmarshalOptionCommand,
	"register":   unmarshalClientToEngineCmd[RegisterCmd],
	"ucinewgame": unmarshalClientToEngineCmd[uciNewGameCmd],
	"position":   unmarshalClientToEngineCmd[positionCmd],
	"go":         unmarshalClientToEngineCmd[goCmd],
	"stop":       unmarshalClientToEngineCmd[stopCmd],
	"ponderhit":  unmarshalClientToEngineCmd[ponderhitCmd],
	"quit":       unmarshalClientToEngineCmd[quitCmd],
}

// parseClientToEngineCmd attempts to parse any client command and return the corresponding clientToEngineCmd.
//
// Parsing is based on the first word in the command.
// If the first word not a valid command, the it will continue on to the second word, and so on.
// If no valid command is found, or parsing of the specified command failed, an error is returned.
func parseClientToEngineCmd(text []byte) (clientToEngineCmd, error) {
	cursor := 0
	for cursor < len(text) {
		// 1. Find the start of the next field (first non-space character)
		start := bytes.IndexFunc(text[cursor:], func(r rune) bool {
			return !unicode.IsSpace(r)
		})

		// If no more non-space characters are found, we're done
		if start == -1 {
			break
		}
		start += cursor // Adjust 'start' to be relative to the beginning of 'text'

		// 2. Find the end of this field (the next space character)
		end := bytes.IndexFunc(text[start:], unicode.IsSpace)

		if end == -1 {
			// No more spaces, so the field goes to the end of the slice
			end = len(text)
		} else {
			end += start
		}

		// 3. Extract the field and check the map
		field := string(text[start:end])
		if unmarshalFunc, ok := commandSet[field]; ok {
			// Lazy Match: Return as soon as we hit a valid command
			return unmarshalFunc(text[start:])
		}

		// 4. Move the cursor forward to search the next field
		cursor = end
	}

	return nil, fmt.Errorf("could not parse client to engine command %q: no valid command keyword was found", text)
}

// A generic function that will parse an engine command an return the value or an error, all at once.
//
// (no need for a separate declaration and call to UnmarshalText)
func unmarshalClientToEngineCmd[T any, PT interface {
	*T
	clientToEngineCmd
}](text []byte) (clientToEngineCmd, error) {
	var result PT = new(T)
	err := result.UnmarshalText(text)
	if err != nil {
		return nil, err
	}
	return result, nil
}
