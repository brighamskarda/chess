package chess

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"
)

// Game is guaranteed to always represent a valid game of chess. Invalid positions are not allowed,
// but the move history may not represent an entire game (undoing all moves may not lead to the starting
// chess position). Game can be used to parse and generate PGNs.
//
// The zero value for Game is not valid and should not be used.
type Game struct {
	position    *Position
	moveHistory []Move
	tags        map[string]string
	result      Result
}

type Result byte

const (
	NoResult Result = iota
	WhiteWins
	BlackWins
	Draw
)

func isValidResult(r Result) bool {
	return r <= Draw
}

// NewGame returns a [*Game] representing the starting position for a game of chess.
func NewGame() *Game {
	position, _ := ParseFen(DefaultFen)
	currentDate := time.Now()
	tags := map[string]string{
		"Event":  "Golang chess match",
		"Site":   "github.com/brighamskarda/chess",
		"Date":   fmt.Sprintf("%4d.%2d.%2d", currentDate.Year(), currentDate.Month(), currentDate.Day()),
		"Round":  "1",
		"White":  "White Player",
		"Black":  "Black Player",
		"Result": "*",
	}
	return &Game{
		position:    position,
		moveHistory: []Move{},
		tags:        tags,
		result:      NoResult,
	}
}

// Move performs the given move. If move m is not legal g remains unchanged and an error is returned.
func (g *Game) Move(m Move) error {
	legalMoves := GenerateLegalMoves(g.position)
	if !slices.Contains(legalMoves, m) {
		return errors.New("m is not a legal move")
	}
	g.position.Move(m)
	g.moveHistory = append(g.moveHistory, m)
	return nil
}

// Returns a copy of current game.
func (g *Game) Copy() *Game {
	positionCopy := *g.position
	gameCopy := &Game{
		position:    &positionCopy,
		moveHistory: slices.Clone(g.moveHistory),
		tags:        maps.Clone(g.tags),
		result:      g.result,
	}
	return gameCopy
}

func (g *Game) String() string {
	return g.position.String()
}

// Position returns a copy of the game's position
func (g *Game) Position() *Position {
	var pos Position = *g.position
	return &pos
}

func (g *Game) Turn() Color {
	return g.position.Turn
}

func (g *Game) WhiteKingSideCastle() bool {
	return g.position.WhiteKingSideCastle
}

func (g *Game) WhiteQueenSideCastle() bool {
	return g.position.WhiteQueenSideCastle
}

func (g *Game) BlackKingSideCastle() bool {
	return g.position.BlackKingSideCastle
}

func (g *Game) BlackQueenSideCastle() bool {
	return g.position.BlackQueenSideCastle
}

func (g *Game) HalfMove() uint16 {
	return g.position.HalfMove
}

func (g *Game) FullMove() uint16 {
	return g.position.FullMove
}

func (g *Game) EnPassant() Square {
	return g.position.EnPassant
}
