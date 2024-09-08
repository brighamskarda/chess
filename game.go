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

func parseResult(s string) Result {
	switch s {
	case "1-0":
		return WhiteWins
	case "0-1":
		return BlackWins
	case "1/2-1/2":
		return Draw
	default:
		return NoResult
	}
}

func (r Result) String() string {
	switch r {
	case WhiteWins:
		return "1-0"
	case BlackWins:
		return "0-1"
	case Draw:
		return "1/2-1/2"
	default:
		return "*"
	}
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
	}
}

// Move performs the given move. If move m is not legal g remains unchanged and an error is returned.
// If the move is legal the result tag is set to * (NoResult). If the position ends in checkmate
// or stalemate the result tag is updated accordingly.
func (g *Game) Move(m Move) error {
	legalMoves := GenerateLegalMoves(g.position)
	if !slices.Contains(legalMoves, m) {
		return errors.New("m is not a legal move")
	}
	g.position.Move(m)
	g.moveHistory = append(g.moveHistory, m)
	if IsCheckMate(g.position) {
		if g.position.Turn == Black {
			g.SetResult(WhiteWins)
		}
		if g.position.Turn == White {
			g.SetResult(BlackWins)
		}
	} else if IsStaleMate(g.position) {
		g.SetResult(Draw)
	} else {
		g.SetResult(NoResult)
	}
	return nil
}

// Returns a copy of current game.
func (g *Game) Copy() *Game {
	positionCopy := *g.position
	gameCopy := &Game{
		position:    &positionCopy,
		moveHistory: slices.Clone(g.moveHistory),
		tags:        maps.Clone(g.tags),
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

func (g *Game) IsCheckMate() bool {
	return IsCheckMate(g.position)
}

func (g *Game) IsStaleMate() bool {
	return IsStaleMate(g.position)
}

func (g *Game) GetResult() Result {
	return parseResult(g.tags["Result"])
}

func (g *Game) SetResult(r Result) {
	g.tags["Result"] = r.String()
}

func (g *Game) ValidMoves() []Move {
	return GenerateLegalMoves(g.position)
}

func (g *Game) GetTag(t string) (string, error) {
	s, ok := g.tags[t]
	if !ok {
		return s, fmt.Errorf("game does not contain tag \"%s\"", t)
	}
	return s, nil
}

// SetTag sets any tag for the game so that it will show up in the pgn file. The Result tag cannot be
// set with this function, please use [Game.SetResult] instead.
func (g *Game) SetTag(tag string, value string) {
	if tag == "Result" {
		return
	}
	g.tags[tag] = value
}

// Remove tag will remove any pgn tag except the 7 required tags specified [here].
//
// [here]: http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.1.1
func (g *Game) RemoveTag(tag string) {
	requiredTags := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}
	if slices.Contains(requiredTags, tag) {
		return
	}
	delete(g.tags, tag)
}

func (g *Game) GetAllTags() map[string]string {
	return maps.Clone(g.tags)
}
