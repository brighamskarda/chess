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

package chess

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Result represents the result of a chess [Game].
type Result uint8

const (
	NoResult Result = iota
	WhiteWins
	BlackWins
	Draw
)

// MarshalText returns "1-0" for [WhiteWins], "0-1" for [BlackWins], "1/2-1/2" for [Draw], and * otherwise. err is always nil.
func (r Result) MarshalText() (text []byte, err error) {
	switch r {
	case WhiteWins:
		return []byte("1-0"), nil
	case BlackWins:
		return []byte("0-1"), nil
	case Draw:
		return []byte("1/2-1/2"), nil
	default:
		return []byte("*"), nil
	}
}

// UnmarshalText sets [WhiteWins] for "1-0", [BlackWins] for "0-1", [Draw] for "1/2-1/2", and [NoResult] for "*". Error is returned if text is not one of these four values.
func (r *Result) UnmarshalText(text []byte) error {
	switch string(text) {
	case "1-0":
		*r = WhiteWins
		return nil
	case "0-1":
		*r = BlackWins
		return nil
	case "1/2-1/2":
		*r = Draw
		return nil
	case "*":
		*r = NoResult
		return nil
	default:
		return fmt.Errorf("could not parse result %q", text)
	}
}

// PgnMove is an expanded move struct used in [Game]. It provides fields for Numeric Annotation Glyphs, commentary and Recursive Annotation Variation (RAV - move variations).
type PgnMove struct {
	Move              Move
	NumericAnnotation uint8
	Commentary        []string
	// Variations supports multiple variations. Hence it is a 2d slice. The first move in a variation should replace the current move.
	Variations [][]PgnMove
}

// Copy provides a copy of the PgnMove. This is a deep copy so the variations are separate.
func (m PgnMove) Copy() PgnMove {
	newPgnMove := PgnMove{
		Move:              m.Move,
		NumericAnnotation: m.NumericAnnotation,
		Commentary:        slices.Clone(m.Commentary),
		Variations:        make([][]PgnMove, 0, len(m.Variations)),
	}
	for _, variation := range m.Variations {
		newVariation := make([]PgnMove, 0, len(variation))
		for _, move := range variation {
			newVariation = append(newVariation, move.Copy())
		}
		newPgnMove.Variations = append(newPgnMove.Variations, newVariation)
	}
	return newPgnMove
}

// Game represents all parts of the PGN game notation standard found here. https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm.
//
// It ensures that only legal moves are performed and keeps track of move history. It also provides utilities for determining if a draw can be claimed.
//
// This library does not support chess 960, largely because special castling rights are not implemented. Otherwise starting games from arbitrary positions is supported.
//
// The field descriptions given are pulled from the PGN specification.
type Game struct {
	pos         *Position
	moveHistory []PgnMove
	// moves are the current legal moves, nil if they haven't been generated. If it is an empty slice then there are no legal moves.
	moves []Move

	// Event should be reasonably descriptive. Abbreviations are to be avoided unless absolutely necessary. A consistent event naming should be used to help facilitate database scanning. If the name of the event is unknown, a single question mark should appear as the tag value.
	Event string
	// Site should include city and region names along with a standard name for the country. The use of the IOC (International Olympic Committee) three letter names is suggested for those countries where such codes are available. If the site of the event is unknown, a single question mark should appear as the tag value. A comma may be used to separate a city from a region. No comma is needed to separate a city or region from the IOC country code. A later section of this document gives a list of three letter nation codes along with a few additions for "locations" not covered by the IOC.
	Site string
	// Date gives the starting date for the game. (Note: this is not necessarily the same as the starting date for the event.) The date is given with respect to the local time of the site given in the Event tag. The Date tag value field always uses a standard ten character format: "YYYY.MM.DD". The first four characters are digits that give the year, the next character is a period, the next two characters are digits that give the month, the next character is a period, and the final two characters are digits that give the day of the month. If the any of the digit fields are not known, then question marks are used in place of the digits.
	Date string
	// Round gives the playing round for the game. In a match competition, this value is the number of the game played. If the use of a round number is inappropriate, then the field should be a single hyphen character. If the round is unknown, a single question mark should appear as the tag value.
	//
	// Some organizers employ unusual round designations and have multipart playing rounds and sometimes even have conditional rounds. In these cases, a multipart round identifier can be made from a sequence of integer round numbers separated by periods. The leftmost integer represents the most significant round and succeeding integers represent round numbers in descending hierarchical order.
	Round string
	// 	White is the name of the player or players of the white pieces. The names are given as they would appear in a telephone directory. The family or last name appears first. If a first name or first initial is available, it is separated from the family name by a comma and a space. Finally, one or more middle initials may appear. (Wherever a comma appears, the very next character should be a space. Wherever an initial appears, the very next character should be a period.) If the name is unknown, a single question mark should appear as the tag value.
	//
	// The intent is to allow meaningful ASCII sorting of the tag value that is independent of regional name formation customs. If more than one person is playing the white pieces, the names are listed in alphabetical order and are separated by the colon character between adjacent entries. A player who is also a computer program should have appropriate version information listed after the name of the program.
	//
	// The format used in the FIDE Rating Lists is appropriate for use for player name tags.
	White string
	// Black is the name of the player or players of the black pieces. The names are given here as they are for the White tag value.
	Black string
	// Result field is the result of the game. It is always exactly the same as the game termination marker that concludes the associated movetext. It is always one of four possible values: "1-0" (White wins), "0-1" (Black wins), "1/2-1/2" (drawn game), and "*" (game still in progress, game abandoned, or result otherwise unknown). Note that the digit zero is used in both of the first two cases; not the letter "O".
	//
	// Note that this field is not just a string, a type was provided to make it slightly easier to use for computers.
	Result Result
	// OtherTags is intended for custom PGN game tags. Some examples are provided here: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c9
	OtherTags map[string]string

	// Commentary represents a game level comment. It will appear before the first move.
	Commentary string
}

// NewGame returns a fresh game of chess with the starting position initialized. Tags are set as follows:
//
// * Event - ?
//
// * Site - https://github.com/brighamskarda/chess
//
// * Date - <CurrentDate>
//
// * Round - 1
//
// * White - ?
//
// * Black - ?
//
// * Result - NoResult
func NewGame() *Game {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	date := time.Now()
	return &Game{
		pos:         pos,
		moveHistory: []PgnMove{},
		moves:       nil,
		Event:       "?",
		Site:        "https://github.com/brighamskarda/chess",
		Date:        date.Format("2006.01.02"),
		Round:       "1",
		White:       "?",
		Black:       "?",
		Result:      NoResult,
		OtherTags:   map[string]string{},
		Commentary:  "",
	}
}

// NewGameFromFEN starts a game specified from the provided fen string. Returns an error if [Position.UnmarshalText] could not parse the FEN. Values are set the same as [NewGame]. The SetUp and FEN tags are also filled into OtherTags.
func NewGameFromFEN(fen string) (*Game, error) {
	pos := &Position{}
	err := pos.UnmarshalText([]byte(fen))
	if err != nil {
		return nil, fmt.Errorf("could not make game: %w", err)
	}
	date := time.Now()
	return &Game{
		pos:         pos,
		moveHistory: []PgnMove{},
		moves:       nil,
		Event:       "?",
		Site:        "https://github.com/brighamskarda/chess",
		Date:        date.Format("2006.01.02"),
		Round:       "1",
		White:       "?",
		Black:       "?",
		Result:      NoResult,
		OtherTags: map[string]string{
			"SetUp": "1",
			"FEN":   fen,
		},
		Commentary: "",
	}, nil
}

// UnmarshalText is capable of unmarshaling a single game in pgn format. See also [ParsePGN]
func (g *Game) UnmarshalText(text []byte) error {
	// Be sure to not read lines beginning with %. These are comments.
	// Semicolons are commentary and go to the end of the line.
	return nil
}

// ParsePGN reads to the end of the provided reader and provides a list of the games parsed from the PGN.
func ParsePGN(pgn io.Reader) ([]*Game, error) {
	// Be sure to not read lines beginning with %. These are comments.
	// Semicolons are commentary and go to the end of the line.

	// Implement some fuzz testing with this to make sure is always returns error and never panics.
	return nil, nil
}

// Copy returns a copy of the game.
func (g *Game) Copy() *Game {
	return &Game{
		pos:         g.pos.Copy(),
		moveHistory: g.MoveHistory(),
		moves:       slices.Clone(g.moves),
		Event:       g.Event,
		Site:        g.Site,
		Date:        g.Date,
		Round:       g.Round,
		White:       g.White,
		Black:       g.Black,
		Result:      g.Result,
		OtherTags:   maps.Clone(g.OtherTags),
		Commentary:  g.Commentary,
	}
}

// LegalMoves returns a copy the current legal moves, cached for performance.
func (g *Game) LegalMoves() []Move {
	return slices.Clone(g.legalMoves())
}

// legalMoves returns the struct with the current legal moves. It is not a copy, don't modify it. Cached for performance.
func (g *Game) legalMoves() []Move {
	if g.moves == nil {
		g.moves = LegalMoves(g.pos)
	}
	return g.moves
}

// IsCheckmate returns true if the side to move is in check and there are no legal moves.
func (g *Game) IsCheckmate() bool {
	if len(g.legalMoves()) == 0 && g.pos.IsCheck() {
		return true
	}
	return false
}

// IsStalemate returns true if the side to move is not in check and has no legal moves.
func (g *Game) IsStalemate() bool {
	if len(g.legalMoves()) == 0 && !g.pos.IsCheck() {
		return true
	}
	return false
}

// CanClaimDraw returns true if the side to move can claim a draw either due to the 50 move rule, or three fold repetition.
func (g *Game) CanClaimDraw() bool {
	return true
}

// CanClaimDrawThreeFold returns true if the side to move can claim a draw due to three fold repetition.
func (g *Game) CanClaimDrawThreeFold() bool {
	return true
}

// Move performs the given move m only if it is legal. Otherwise an error is produced.
//
// If Result is set then it will be set to [NoResult]. If the game ends (in stalemate or checkmate) then Result will also be set appropriately.
func (g *Game) Move(m Move) error {
	if !slices.Contains(g.legalMoves(), m) {
		return errors.New("illegal move")
	}
	g.pos.Move(m)
	g.moves = nil
	g.moveHistory = append(g.moveHistory, PgnMove{
		Move:              m,
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	})
	g.setResult()
	return nil
}

// setResult sets the game result, defaults to NoResult.
func (g *Game) setResult() {
	if g.IsStalemate() {
		g.Result = Draw
	} else if g.IsCheckmate() {
		if g.pos.SideToMove == White {
			g.Result = BlackWins
		} else if g.pos.SideToMove == Black {
			g.Result = WhiteWins
		}
	} else {
		g.Result = NoResult
	}
}

// MoveUCI parses and performs a UCI chess move. https://www.wbec-ridderkerk.nl/html/UCIProtocol.html
//
// Errors are returned if m could not be parsed or the move was illegal.
func (g *Game) MoveUCI(m string) error {
	move, err := ParseUCIMove(m)
	if err != nil {
		return err
	}
	return g.Move(move)
}

// MoveSAN parses and performs a SAN (Standard Algebraic Notation) chess move. https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.2.3
//
// Errors are returned if m could not be parsed or the move was illegal.
func (g *Game) MoveSAN(m string) error {
	move, err := ParseSANMove(m, g.pos)
	if err != nil {
		return err
	}
	return g.Move(move)
}

// Position returns a copy of the current position.
func (g *Game) Position() *Position {
	return g.pos.Copy()
}

// PositionPly returns a copy of the position at a certain ply (half move). 0 returns the initial game position.
//
// If a negative number is provided, or ply goes beyond the number of moves played nil is returned.
func (g *Game) PositionPly(ply int) *Position {
	pos := &Position{}
	if g.OtherTags["SetUp"] == "1" {

		err := pos.UnmarshalText([]byte(g.OtherTags["FEN"]))
		if err != nil {
			panic("game somehow got invalid FEN starting position, can get position at ply")
		}
	} else {
		pos.UnmarshalText([]byte(DefaultFEN))
	}
	for _, m := range g.moveHistory[:ply] {
		pos.Move(m.Move)
	}
	return pos
}

// MoveHistory returns a copy of all the moves played this game with their annotations, commentary and variations. Will not return nil.
func (g *Game) MoveHistory() []PgnMove {
	moveHistoryCopy := make([]PgnMove, 0, len(g.moveHistory))
	for _, move := range g.moveHistory {
		moveHistoryCopy = append(moveHistoryCopy, move.Copy())
	}
	return moveHistoryCopy
}

// AnnotateMove applies a numeric annotation glyph (NAG) to the specified move number. NAG's can be found here: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c10
//
// plyNum starts at 0 for the first move. Any previous nag is overwritten.
func (g *Game) AnnotateMove(plyNum int, nag uint8) {
	g.moveHistory[plyNum].NumericAnnotation = nag
}

// CommentMove appends a comment to the specified move.
//
// plyNum starts at 0 for the first move.
func (g *Game) CommentMove(plyNum int, comment string) {
	g.moveHistory[plyNum].Commentary = append(g.moveHistory[plyNum].Commentary, comment)
}

// DeleteComment deletes a comment from the specified move.
//
// plyNum and commentNum start at 0 for the first move.
func (g *Game) DeleteComment(plyNum int, commentNum int) {
	g.moveHistory[0].Commentary = slices.Delete(g.moveHistory[0].Commentary, commentNum, commentNum+1)
}

// MakeVariation adds a set of variation moves to the specified move. The variation should begin with a move that replaces the current move. Variation moves must be legal.
//
// plyNum starts at 0 for the first move. moves is not copied, it is simply added to the variations list for the appropriate pgn move.
func (g *Game) MakeVariation(plyNum int, moves []PgnMove) {
	g.moveHistory[plyNum].Variations = append(g.moveHistory[plyNum].Variations, moves)
}

// DeleteVariation deletes the variation at the plyNum. Give the index you want to delete (starting at 0).
func (g *Game) DeleteVariation(plyNum int, variationNum int) {
	g.moveHistory[plyNum].Variations = slices.Delete(g.moveHistory[plyNum].Variations, variationNum, variationNum+1)
}

// GetVariation returns a new game where the specified variation is followed. All other specified variations are preserved, and the main line is also preserved as a variation.
func (g *Game) GetVariation(plyNum int, variationNum int) *Game {
	newMoveHistory := g.MoveHistory()[0:plyNum]
	// Make the specified variation the main move history
	for _, m := range g.moveHistory[plyNum].Variations[variationNum] {
		newMoveHistory = append(newMoveHistory, m.Copy())
	}

	// Make the main move history a variation
	variation := []PgnMove{}
	for _, m := range g.moveHistory[plyNum:] {
		variation = append(variation, m.Copy())
	}
	newMoveHistory[plyNum].Variations = append(newMoveHistory[plyNum].Variations, variation)

	// Add other move variations
	for varNum, variation := range g.moveHistory[plyNum].Variations {
		if varNum == variationNum {
			continue
		}
		variationCopy := []PgnMove{}
		for _, m := range variation {
			variationCopy = append(variationCopy, m.Copy())
		}
		newMoveHistory[plyNum].Variations = append(newMoveHistory[plyNum].Variations, variationCopy)
	}

	newGame := g.Copy()
	newGame.pos = newGame.PositionPly(0)
	newGame.moveHistory = []PgnMove{}
	newGame.moves = nil
	for _, m := range newMoveHistory {
		newGame.Move(m.Move)
	}

	newGame.moveHistory = newMoveHistory

	return newGame
}

// MarshalText implements [encoding.TextMarshaler]. It provides the game as a valid PGN that can be written to a file. Multiple PGNs can be written to the same file. Just be sure to separate them with a new line.
//
// The seven tag roster will appear in order, then all other tags will appear in alphabetical order for consistency. err is always nil.
func (g *Game) MarshalText() (text []byte, err error) {
	return []byte(g.String()), nil
}

// String provides the same functionality as [Game.MarshalText].
func (g *Game) String() string {
	lines := make([]string, 0, 10)
	g.addTags(&lines)
	lines = append(lines, "")
	if g.Commentary != "" {
		lines = append(lines, fmt.Sprintf("{%s}", g.Commentary))
	}
	g.addMoveText(&lines)
	pgn := strings.Builder{}
	for i, l := range lines {
		pgn.WriteString(l)
		if i != len(lines)-1 {
			pgn.WriteString("\n")
		}
	}
	return pgn.String()
}

func (g *Game) addTags(lines *[]string) {
	*lines = append(*lines, fmt.Sprintf("[Event %q]", g.Event))
	*lines = append(*lines, fmt.Sprintf("[Site %q]", g.Site))
	*lines = append(*lines, fmt.Sprintf("[Date %q]", g.Date))
	*lines = append(*lines, fmt.Sprintf("[Round %q]", g.Round))
	*lines = append(*lines, fmt.Sprintf("[White %q]", g.White))
	*lines = append(*lines, fmt.Sprintf("[Black %q]", g.Black))
	rstStr, err := g.Result.MarshalText()
	if err != nil {
		rstStr = []byte("*")
	}
	*lines = append(*lines, fmt.Sprintf("[Result %q]", rstStr))
	g.addOtherTags(lines)
}

func (g *Game) addOtherTags(lines *[]string) {
	keys := make([]string, 0, len(g.OtherTags))
	for k := range g.OtherTags {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		*lines = append(*lines, fmt.Sprintf("[%s %q]", k, g.OtherTags[k]))
	}
}

func (g *Game) addMoveText(lines *[]string) {
	currPos := g.PositionPly(0)
	currentLine := strings.Builder{}
	currentLine.Grow(80)

	includeBlackMoveNum := currPos.SideToMove == Black
	for _, m := range g.moveHistory {
		if currPos.SideToMove == White {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "."
			appendToPgnLine(moveNum, &currentLine, lines)
		}
		if includeBlackMoveNum {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "..."
			appendToPgnLine(moveNum, &currentLine, lines)
		}

		sanMove := " " + m.Move.StringSAN(currPos)
		nag := nagString(m.NumericAnnotation)
		if m.NumericAnnotation <= 6 {
			sanMove += nag
			appendToPgnLine(sanMove, &currentLine, lines)
		} else {
			appendToPgnLine(sanMove, &currentLine, lines)
			appendToPgnLine(" "+nag, &currentLine, lines)
		}

		for _, comment := range m.Commentary {
			cmtStr := fmt.Sprintf(" {%s}", comment)
			appendToPgnLine(cmtStr, &currentLine, lines)
		}

		for _, variation := range m.Variations {
			appendVariation(currPos.Copy(), variation, &currentLine, lines)
		}

		currPos.Move(m.Move)
		if (len(m.Commentary) > 0 || len(m.Variations) > 0) && currPos.SideToMove == Black {
			includeBlackMoveNum = true
		} else {
			includeBlackMoveNum = false
		}
	}
	result, _ := g.Result.MarshalText()
	appendToPgnLine(" "+string(result), &currentLine, lines)
	*lines = append(*lines, currentLine.String())
}

// appendToPgnLine appends string s to currentLine. If currentLine would be longer than 80, then it is appended to lines and currentLine is reset to the value of s.
func appendToPgnLine(s string, currentLine *strings.Builder, lines *[]string) {
	if len(s)+currentLine.Len() > 80 {
		*lines = append(*lines, currentLine.String())
		currentLine.Reset()
	}
	if currentLine.Len() == 0 {
		s = strings.TrimSpace(s)
	}
	currentLine.WriteString(s)
}

func nagString(nag uint8) string {
	switch nag {
	case 0:
		return ""
	case 1:
		return "!"
	case 2:
		return "?"
	case 3:
		return "!!"
	case 4:
		return "??"
	case 5:
		return "!?"
	case 6:
		return "?!"
	default:
		return "$" + strconv.FormatUint(uint64(nag), 10)
	}
}

func appendVariation(currPos *Position, moves []PgnMove, currentLine *strings.Builder, lines *[]string) {
	includeBlackMoveNum := currPos.SideToMove == Black
	for i, m := range moves {
		if currPos.SideToMove == White {
			moveNum := strconv.FormatUint(uint64(currPos.FullMove), 10) + "."
			if i == 0 {
				moveNum = "(" + moveNum
			}
			moveNum = " " + moveNum
			appendToPgnLine(moveNum, currentLine, lines)
		}
		if includeBlackMoveNum {
			moveNum := strconv.FormatUint(uint64(currPos.FullMove), 10) + "..."
			if i == 0 {
				moveNum = "(" + moveNum
			}
			moveNum = " " + moveNum
			appendToPgnLine(moveNum, currentLine, lines)
		}

		sanMove := " " + m.Move.StringSAN(currPos)
		nag := nagString(m.NumericAnnotation)
		if m.NumericAnnotation <= 6 {
			sanMove += nag
			if i == len(moves)-1 && len(m.Commentary) == 0 && len(m.Variations) == 0 {
				sanMove += ")"
			}
			appendToPgnLine(sanMove, currentLine, lines)
		} else {
			appendToPgnLine(sanMove, currentLine, lines)
			if i == len(moves)-1 && len(m.Commentary) == 0 && len(m.Variations) == 0 {
				nag += ")"
			}
			appendToPgnLine(" "+nag, currentLine, lines)
		}

		for j, comment := range m.Commentary {
			cmtStr := fmt.Sprintf(" {%s}", comment)
			if i == len(moves)-1 && j == len(m.Commentary)-1 && len(m.Variations) == 0 {
				cmtStr += ")"
			}
			appendToPgnLine(cmtStr, currentLine, lines)
		}

		for j, variation := range m.Variations {
			appendVariation(currPos.Copy(), variation, currentLine, lines)
			if i == len(moves)-1 && j == len(m.Variations)-1 {
				appendToPgnLine(")", currentLine, lines)
			}
		}

		currPos.Move(m.Move)
		if (m.NumericAnnotation > 6 || len(m.Commentary) > 0 || len(m.Variations) > 0) && currPos.SideToMove == Black {
			includeBlackMoveNum = true
		} else {
			includeBlackMoveNum = false
		}
	}
}

// ReducedString provides the game as a valid PGN following these rules: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c3.2.4
func (g *Game) ReducedString() string {
	return ""
}
