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

// PgnMove is an expanded move struct used in [Game]. It provides fields for Numeric Annotation Glyphs, commentary and Recursive Annotation Variation (RAV - move variations).
type PgnMove struct {
	Move              Move
	NumericAnnotation uint8
	Commentary        string
	// Variations supports multiple variations. Hence it is a 2d slice.
	Variations [][]PgnMove
}

// Copy provides a copy of the PgnMove. This is a deep copy so the variations are separate.
func (m PgnMove) Copy() PgnMove {
	newPgnMove := PgnMove{
		Move:              m.Move,
		NumericAnnotation: m.NumericAnnotation,
		Commentary:        m.Commentary,
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
	pos, _ := ParseFEN(DefaultFEN)
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

// NewGameFromFEN starts a game specified from the provided fen string. Returns an error if [ParseFEN] could not parse the FEN. Values are set the same as [NewGame]. The SetUp and FEN tags are also filled into OtherTags.
func NewGameFromFEN(fen string) (*Game, error) {
	pos, err := ParseFEN(fen)
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

// ParsePGN reads to the end of the provided reader and provides a list of the games parsed from the PGN.
func ParsePGN(pgn io.Reader) ([]*Game, error) {
	// Be sure to not read lines beginning with %. These are comments.
	// Semicolons are commentary and go to the end of the line.
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
		Commentary:        "",
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
	var pos *Position
	if g.OtherTags["SetUp"] == "1" {
		var err error
		pos, err = ParseFEN(g.OtherTags["FEN"])
		if err != nil {
			panic("game somehow got invalid FEN starting position, can get position at ply")
		}
	} else {
		pos, _ = ParseFEN(DefaultFEN)
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

// CommentMove applies a comment to the specified move.
//
// plyNum starts at 0 for the first move. Any previous comment is overwritten.
func (g *Game) CommentMove(plyNum int, comment string) {
	g.moveHistory[plyNum].Commentary = comment
}

// MakeVariation adds a set of variation moves to the specified move. Variation moves must be legal.
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

// String provides the game as a valid PGN that can be written to a file. Multiple PGNs can be written to the same file. Just be sure to separate them with a new line.
func (g *Game) String() string {
	return ""
}

// ReducedString provides the game as a valid PGN following these rules: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c3.2.4
func (g *Game) ReducedString() string {
	return ""
}
