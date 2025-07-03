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
	"bytes"
	"strconv"
	"strings"

	"github.com/brighamskarda/chess/v2"
)

// command represents a command received from the engine. It is necessary to type assert certain commands. But this should be avoided when possible.
type command interface {
	commandType() commandType
	message() string
}

type commandType uint8

const (
	unknown commandType = iota
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

func (bc basicCommand) message() string {
	return bc.msg
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

func (i *Info) message() string {
	return ""
}

func parseInfoCommand(line []byte) *Info {
	info := &Info{}

	tokens := bytes.Fields(line)

	i := findInfoIndex(tokens) + 1
	if i < 0 {
		return info
	}

loop:
	for ; i < len(tokens)-1; i++ {
		switch strings.ToLower(string(tokens[i])) {
		case "depth":
			i++
			info.Depth = parseUintPointer(tokens[i])
			if info.Depth == nil {
				i--
			}
		case "seldepth":
			i++
			info.Seldepth = parseUintPointer(tokens[i])
			if info.Seldepth == nil {
				i--
			}
		case "time":
			i++
			info.Time = parseUintPointer(tokens[i])
			if info.Time == nil {
				i--
			}
		case "nodes":
			i++
			info.Nodes = parseUintPointer(tokens[i])
			if info.Nodes == nil {
				i--
			}
		case "pv":
			var numParsed int
			info.Pv, numParsed = parseMoveLine(tokens[i+1:])
			i += numParsed
		case "multipv":
			i++
			info.Multipv = parseUintPointer(tokens[i])
			if info.Nodes == nil {
				i--
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
			i++
			info.Currmovenumber = parseUintPointer(tokens[i])
			if info.Currmovenumber == nil {
				i--
			}
		case "hashfull":
			i++
			info.Hashfull = parseUintPointer(tokens[i])
			if info.Hashfull == nil {
				i--
			}
		case "nps":
			i++
			info.Nps = parseUintPointer(tokens[i])
			if info.Nps == nil {
				i--
			}
		case "tbhits":
			i++
			info.Tbhits = parseUintPointer(tokens[i])
			if info.Tbhits == nil {
				i--
			}
		case "cpuload":
			i++
			info.CpuLoad = parseUintPointer(tokens[i])
			if info.CpuLoad == nil {
				i--
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

func findInfoIndex(tokens [][]byte) int {
	for i, token := range tokens {
		if bytes.EqualFold(token, []byte("info")) {
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
