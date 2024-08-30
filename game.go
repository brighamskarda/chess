package chess

import (
	"fmt"
	"time"
)

// Game is guaranteed to always represent a valid game of chess. Invalid positions are not allowed,
// but the move history may not represent an entire game (undoing all moves may not lead to the starting
// chess position). Game can be used to parse and generate PGNs.
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

func NewGame() *Game {
	position, _ := ParseFen(DefaultFen)
	currentDate := time.Now()
	tags := map[string]string{
		"Event": "Golang chess match",
		"Site": "github.com/brighamskarda/chess",
		"Date": fmt.Sprintf("%4d.%2d.%2d", currentDate.Year(), currentDate.Month(), currentDate.Day()),
		"Round": "1",
		"White": "White Player",
		"Black": "Black Player",
		"Result": "*",
	}
	return &Game{
		position: position,
		moveHistory: []Move{},
		tags: tags,
		result: NoResult,
	}
}