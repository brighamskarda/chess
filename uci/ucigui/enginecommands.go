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
	"bytes"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/brighamskarda/chess/v2"
)

// command represents a command received from the engine. It is necessary to type assert certain commands. But this should be avoided when possible.
type command interface {
	commandType() commandType
}

type commandType uint8

const (
	unknownCommandType commandType = iota
	id
	uciok
	readyok
	bestmove
	copyprotection
	registration
	info
	option
)

type basicCommand struct {
	cmdType commandType
	msg     string
}

func (bc basicCommand) commandType() commandType {
	return bc.cmdType
}

// Score is used in [Info] to represent the score reported by the engine.
type Score struct {
	// Cp - the score from the engine's point of view in centipawns.
	Cp *int
	// Mate - mate in y moves, not plies.
	// If the engine is getting mated use negative values for y.
	Mate *int
	// Lowerbound - the score is just a lower bound. True if present.
	Lowerbound bool
	// Upperbound - the score is just an upper bound. True if present.
	Upperbound bool
}

// Currline is used in [Info] to represent the line being evaluated by a certain cpu.
type Currline struct {
	// Cpunr - the number of the cpu.
	Cpunr *uint
	// Moves - the line being evaluated.
	Moves []chess.Move
}

// Info represents information sent from the engine to the client with the info command. Any fields not received are set to nil.
type Info struct {
	// Depth - search depth in plies.
	Depth *uint
	// Seldepth - selective search depth in plies,
	// if the engine sends seldepth there must also be a "depth" present in the same string.
	Seldepth *uint
	// Time - the time searched in ms, this should be sent together with the pv.
	Time *uint
	// Nodes - nodes searched, the engine should send this info regularly.
	Nodes *uint
	// Pv (Primary Variation) - the best line found.
	Pv []chess.Move
	// Multipv - this for the multi pv mode.
	// for the best move/pv add "multipv 1" in the string when you send the pv.
	// in k-best mode always send all k variants in k strings together.
	Multipv *uint
	// Score - see [Score].
	Score *Score
	// Currmove - currently searching this move.
	Currmove *chess.Move
	// Currmovenumber - currently searching move number x, for the first move x should be 1 not 0.
	Currmovenumber *uint
	// Hashfull - the hash is x permill (like percent, but in 1,000 parts) full, the engine should send this info regularly.
	Hashfull *uint
	// Nps - nodes per second searched, the engine should send this info regularly
	Nps *uint
	// Tbhits - x positions where found in the endgame table bases
	Tbhits *uint
	// Cpuload - the cpu usage of the engine is x permill (per 1,000 parts).
	CpuLoad *uint
	// String - any string str which will by displayed by the engine,
	// if there is a string command the rest of the line will be interpreted as a string.
	String *string
	// Refutation - move <move1> is refuted by the line <move2> ... <movei>, i can be any number >= 1.
	//    Example: after move d1h5 is searched, the engine can send
	//    "info refutation d1h5 g6h5"
	//    if g6h5 is the best answer after d1h5 or if g6h5 refutes the move d1h5.
	//    if there is no refutation for d1h5 found, the engine should just send
	//    "info refutation d1h5"
	// The engine should only send this if the option "UCI_ShowRefutations" is set to true.
	Refutation []chess.Move
	// Currline - this is the current line the engine is calculating. <cpunr> is the number of the cpu if
	//    the engine is running on more than one cpu. <cpunr> = 1,2,3....
	//    if the engine is just using one cpu, <cpunr> can be omitted.
	//    If <cpunr> is greater than 1, always send all k lines in k strings together.
	// The engine should only send this if the option "UCI_ShowCurrLine" is set to true.
	Currline *Currline
}

func (i *Info) commandType() commandType {
	return info
}

func parseInfoCommand(line []byte) *Info {
	info := &Info{}

	tokens := bytes.Fields(line)

	i := findTokenIndex(tokens, []byte("info")) + 1
	if i <= 0 {
		return info
	}

loop:
	for ; i < len(tokens)-1; i++ {
		switch strings.ToLower(string(tokens[i])) {
		case "depth":
			info.Depth = parseUintPointer(tokens[i+1])
			if info.Depth != nil {
				i++
			}
		case "seldepth":
			info.Seldepth = parseUintPointer(tokens[i+1])
			if info.Seldepth != nil {
				i++
			}
		case "time":
			info.Time = parseUintPointer(tokens[i+1])
			if info.Time != nil {
				i++
			}
		case "nodes":
			info.Nodes = parseUintPointer(tokens[i+1])
			if info.Nodes != nil {
				i++
			}
		case "pv":
			var numParsed int
			info.Pv, numParsed = parseMoveLine(tokens[i+1:])
			i += numParsed
		case "multipv":
			info.Multipv = parseUintPointer(tokens[i+1])
			if info.Nodes != nil {
				i++
			}
		case "score":
			var numParsed int
			info.Score, numParsed = parseScore(tokens[i+1:])
			i += numParsed
		case "currmove":
			i++
			info.Currmove = &chess.Move{}
			if info.Currmove.UnmarshalText(tokens[i]) != nil {
				i--
				info.Currmove = nil
			}
		case "currmovenumber":
			info.Currmovenumber = parseUintPointer(tokens[i+1])
			if info.Currmovenumber != nil {
				i++
			}
		case "hashfull":
			info.Hashfull = parseUintPointer(tokens[i+1])
			if info.Hashfull != nil {
				i++
			}
		case "nps":
			info.Nps = parseUintPointer(tokens[i+1])
			if info.Nps != nil {
				i++
			}
		case "tbhits":
			info.Tbhits = parseUintPointer(tokens[i+1])
			if info.Tbhits != nil {
				i++
			}
		case "cpuload":
			info.CpuLoad = parseUintPointer(tokens[i+1])
			if info.CpuLoad != nil {
				i++
			}
		case "refutation":
			var numParsed int
			info.Refutation, numParsed = parseMoveLine(tokens[i+1:])
			i += numParsed
		case "currline":
			var numParsed int
			info.Currline, numParsed = parseCurrLine(tokens[i+1:])
			i += numParsed
		case "string":
			s := string(bytes.TrimSpace(line[bytes.Index(bytes.ToLower(line), []byte("string"))+6:]))
			info.String = &s
			break loop
		}

	}

	return info
}

func findTokenIndex(tokens [][]byte, tokenToFind []byte) int {
	for i, token := range tokens {
		if bytes.EqualFold(token, tokenToFind) {
			return i
		}
	}
	return -1
}

func parseUintPointer(token []byte) *uint {
	i, err := strconv.ParseUint(string(token), 10, 0)
	if err != nil {
		return nil
	}
	val := uint(i)
	return &val
}

func parseIntPointer(token []byte) *int {
	i, err := strconv.ParseInt(string(token), 10, 0)
	if err != nil {
		return nil
	}
	val := int(i)
	return &val
}

func parseMoveLine(tokens [][]byte) (moves []chess.Move, numParsed int) {
	for _, t := range tokens {
		var newMove chess.Move
		if newMove.UnmarshalText(t) != nil {
			break
		}
		moves = append(moves, newMove)
		numParsed++
	}
	return
}

func parseScore(tokens [][]byte) (score *Score, numParsed int) {
	score = &Score{}

	if len(tokens) < 2 {
		return nil, numParsed
	}

	// Need to parse at least one of these to avoid returning nil
	switch strings.ToLower(string(tokens[0])) {
	case "cp":
		score.Cp = parseIntPointer(tokens[1])
		numParsed += 2
		if score.Cp == nil {
			return nil, numParsed - 1
		}
	case "mate":
		score.Mate = parseIntPointer(tokens[1])
		numParsed += 2
		if score.Mate == nil {
			return nil, numParsed - 1
		}
	default:
		return nil, numParsed
	}

	// Continue parsing other options
loop:
	for ; numParsed < len(tokens); numParsed++ {
		switch strings.ToLower(string(tokens[numParsed])) {
		case "cp":
			numParsed++
			if numParsed < len(tokens) {
				score.Cp = parseIntPointer(tokens[numParsed])
			}
			if score.Cp == nil {
				numParsed--
			}
		case "mate":
			numParsed++
			if numParsed < len(tokens) {
				score.Mate = parseIntPointer(tokens[numParsed])
			}
			if score.Mate == nil {
				numParsed--
			}
		case "lowerbound":
			score.Lowerbound = true
		case "upperbound":
			score.Upperbound = true
		default:
			break loop
		}
	}

	return
}

func parseCurrLine(tokens [][]byte) (currLine *Currline, numParsed int) {
	cpunr := parseUintPointer(tokens[0])
	if cpunr == nil {
		return
	}

	numParsed++

	line, n := parseMoveLine(tokens[1:])
	if n == 0 {
		return
	}
	numParsed += n

	currLine = &Currline{
		Cpunr: cpunr,
		Moves: line,
	}
	return
}

type OptionType uint8

const (
	_ OptionType = iota
	Check
	Spin
	Combo
	Button
	String
)

func parseOptionType(token []byte) OptionType {
	switch strings.ToLower(string(token)) {
	case "check":
		return Check
	case "spin":
		return Spin
	case "combo":
		return Combo
	case "button":
		return Button
	case "string":
		return String
	default:
		return 0
	}
}

// Option represent the configurable options sent by the chess engine. See the [UCI chess specification] for common options. As specified there, options with certain names will be verified to ensure they are of the right type. Options that break those conventions will not be parsed. Options without a type or name will also not be parsed.
//
// If default, min, max, or var are nil, it means they were not provided. As defined in the UCI specification, there are only certain combinations that make sense. Options that don't follow these rules may be parsed incorrectly. Every option that at least defines name and type will be parsed though, even if their attributes don't make sense or aren't parsed correctly.
//
//   - Check - May only have a default value of true or false.
//   - Spin - Can define min, max, and default. These should all be numbers.
//   - Combo - Must define at least one var. Can also have a default.
//   - String - Can only have a default.
//   - Button - Has no attributes.
//
// Options with the prefix "UCI_" that are not defined in the UCI specification are ignored, per the UCI specification.
//
// [UCI chess specification]: https://www.shredderchess.com/download/div/uci.zip
type Option struct {
	Name    string
	OType   OptionType
	Default *string
	Min     *int
	Max     *int
	Var     []string
}

func (o *Option) commandType() commandType {
	return option
}

func parseOptionCommand(line []byte) *Option {
	option := &Option{}

	tokens := bytes.Fields(line)

	i := findTokenIndex(tokens, []byte("option")) + 1
	if i <= 0 {
		return option
	}

	for ; i < len(tokens)-1; i++ {
		switch string(tokens[i]) {
		case "name":
			var numParsed int
			option.Name, numParsed = parseOptionName(line)
			i += numParsed
		case "type":
			option.OType = parseOptionType(tokens[i+1])
			if option.OType != 0 {
				i++
			}
		case "default":
			optionDefault, numParsed := parseOptionDefault(line, option.OType == String)
			option.Default = &optionDefault
			i += numParsed
		case "min":
			option.Min = parseIntPointer(tokens[i+1])
			if option.Min != nil {
				i++
			}
		case "max":
			option.Max = parseIntPointer(tokens[i+1])
			if option.Max != nil {
				i++
			}
		case "var":
			optionVar, numParsed := parseOptionVar(line, len(option.Var))
			option.Var = append(option.Var, optionVar)
			i += numParsed
		}
	}

	if !nameAndTypeExist(option) {
		return nil
	}

	if isPredefinedOption(option) {
		if !isLegalPredefOption(option) {
			return nil
		}
		return option
	}

	if strings.HasPrefix(option.Name, "UCI_") {
		return nil
	}

	return option
}

func nameAndTypeExist(o *Option) bool {
	return o.Name != "" && o.OType != 0
}

func parseOptionName(line []byte) (string, int) {
	startIndex := findTokenIndexWithWhiteSpace(line, "name")

	if startIndex < 0 {
		return "", 0
	}

	var optionName string
	nextOptionOffset := findNextOptionIndex(line[startIndex:])
	if nextOptionOffset == -1 {
		optionName = string(line[startIndex:])
	} else {
		optionName = string(line[startIndex : startIndex+nextOptionOffset])
	}
	optionName = strings.TrimSpace(optionName)
	return optionName, len(strings.Fields(optionName))
}

func parseOptionDefault(line []byte, isString bool) (string, int) {
	startIndex := findTokenIndexWithWhiteSpace(line, "default")

	if startIndex < 0 {
		return "", 0
	}

	if isString {
		defaultValue := strings.TrimSpace(string(line[startIndex:]))
		return defaultValue, len(strings.Fields(defaultValue))
	}

	var defaultValue string
	nextOptionOffset := findNextOptionIndex(line[startIndex:])
	if nextOptionOffset == -1 {
		defaultValue = string(line[startIndex:])
	} else {
		defaultValue = string(line[startIndex : startIndex+nextOptionOffset])
	}
	defaultValue = strings.TrimSpace(defaultValue)
	return defaultValue, len(strings.Fields(defaultValue))
}

func parseOptionVar(line []byte, varIndex int) (string, int) {

	startIndex := findTokenIndexWithWhiteSpace(line, "var")
	for range varIndex {
		startIndex += findTokenIndexWithWhiteSpace(line[startIndex:], "var")
	}

	if startIndex < 0 {
		return "", 0
	}

	var varValue string
	nextOptionOffset := findNextOptionIndex(line[startIndex:])
	if nextOptionOffset == -1 {
		varValue = string(line[startIndex:])
	} else {
		varValue = string(line[startIndex : startIndex+nextOptionOffset])
	}
	varValue = strings.TrimSpace(varValue)
	return varValue, len(strings.Fields(varValue))
}

func findNextOptionIndex(line []byte) int {
	nextOptionIndex := math.MaxInt

	options := [][]byte{[]byte(" type "), []byte(" default "), []byte(" min "), []byte(" max "), []byte(" var ")}

	for _, o := range options {
		possibleOption := bytes.Index(line, o)
		if possibleOption > 0 && possibleOption < nextOptionIndex {
			nextOptionIndex = possibleOption
		}
	}

	if nextOptionIndex == math.MaxInt {
		return -1
	}
	return nextOptionIndex
}

var predef map[string]struct{} = map[string]struct{}{
	"Hash":                  {},
	"NalimovPath":           {},
	"NalimovCache":          {},
	"Ponder":                {},
	"OwnBook":               {},
	"MultiPV":               {},
	"UCI_ShowCurrLine":      {},
	"UCI_ShowRefutations":   {},
	"UCI_LimitStrength":     {},
	"UCI_Elo":               {},
	"UCI_AnalyseMode":       {},
	"UCI_Opponent":          {},
	"UCI_EngineAbout":       {},
	"UCI_ShredderbasesPath": {},
	"UCI_SetPositionValue":  {},
}

func isPredefinedOption(o *Option) bool {
	_, ok := predef[o.Name]
	return ok
}

func isLegalPredefOption(o *Option) bool {
	switch o.Name {
	case "Hash", "NalimovCache", "MultiPV", "UCI_Elo":
		return o.OType == Spin
	case "NalimovPath", "UCI_Opponent", "UCI_EngineAbout", "UCI_ShredderbasesPath", "UCI_SetPositionValue":
		return o.OType == String
	case "Ponder", "OwnBook", "UCI_ShowCurrLine", "UCI_ShowRefutations", "UCI_LimitStrength", "UCI_AnalyseMode":
		return o.OType == Check
	default:
		return true
	}
}

func findTokenIndexWithWhiteSpace(line []byte, token string) int {
	re := regexp.MustCompile(`\s` + token + `\s`)
	matches := re.FindIndex(line)
	if matches != nil {
		return matches[1]
	}
	return -1
}

type idType uint8

const (
	_ idType = iota
	name
	author
)

type idCommand struct {
	idt   idType
	value string
}

func (i idCommand) commandType() commandType {
	return id
}

func parseIdCommand(line []byte) *idCommand {
	nameIndex := bytes.Index(line, []byte("name"))
	authorIndex := bytes.Index(line, []byte("author"))
	if nameIndex == -1 && authorIndex == -1 {
		return nil
	} else if nameIndex != -1 && authorIndex == -1 {
		return &idCommand{
			idt:   name,
			value: string(bytes.TrimSpace(line[nameIndex+4:])),
		}
	} else if nameIndex == -1 && authorIndex != -1 {
		return &idCommand{
			idt:   author,
			value: string(bytes.TrimSpace(line[authorIndex+6:])),
		}
	} else if nameIndex < authorIndex {
		return &idCommand{
			idt:   name,
			value: string(bytes.TrimSpace(line[nameIndex+4:])),
		}
	} else {
		return &idCommand{
			idt:   author,
			value: string(bytes.TrimSpace(line[authorIndex+6:])),
		}
	}
}

type bestMove struct {
	best   chess.Move
	ponder *chess.Move
}

func (bm bestMove) commandType() commandType {
	return bestmove
}

func parseBestMoveCommand(line []byte) *bestMove {
	best := bestMove{}

	tokens := bytes.Fields(line)

	bestMoveSet := false
	for i, t := range tokens {
		if bytes.EqualFold(t, []byte("ponder")) {
			return nil
		}
		if err := best.best.UnmarshalText(t); err == nil {
			if i < len(tokens)-1 {
				tokens = tokens[i+1:]
			} else {
				tokens = nil
			}
			bestMoveSet = true
			break
		}
	}

	if !bestMoveSet {
		return nil
	}

	for i, t := range tokens {
		if bytes.EqualFold(t, []byte("ponder")) && i < len(tokens)-1 {
			tokens = tokens[i+1:]
			break
		}
	}

	ponder := chess.Move{}
	for _, t := range tokens {
		if err := ponder.UnmarshalText(t); err == nil {
			best.ponder = &ponder
			break
		}
	}

	return &best
}

type copyProtection uint8

func (cp copyProtection) commandType() commandType {
	return copyprotection
}

const (
	_ copyProtection = iota
	checking
	ok
	cpError
)

func parseCopyProtection(line []byte) *copyProtection {
	tokens := bytes.Fields(line)

	var cp copyProtection
	for _, t := range tokens {
		if bytes.EqualFold(t, []byte("checking")) {
			cp = checking
			break
		}
		if bytes.EqualFold(t, []byte("ok")) {
			cp = ok
			break
		}
		if bytes.EqualFold(t, []byte("error")) {
			cp = cpError
			break
		}
	}

	if cp == 0 {
		return nil
	}
	return &cp
}
