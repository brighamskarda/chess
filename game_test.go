package chess

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIsValidResult(t *testing.T) {
	if isValidResult(5) {
		t.Error("incorrect result: input 5: expected false, got true")
	}
	if !isValidResult(Draw) {
		t.Error("incorrect result: input Draw: expected true, got false")
	}
}

func TestNewGame(t *testing.T) {
	game := NewGame()
	if *game.position != *getDefaultPosition() {
		t.Error("default position incorrect: ", cmp.Diff(getDefaultPosition(), game.position))
	}
	if game.tags["Event"] != "Golang chess match" {
		t.Errorf(`default game tag "Event" incorrect: expected Golang chess match, got %s`, game.tags["Event"])
	}
	if game.tags["Site"] != "github.com/brighamskarda/chess" {
		t.Errorf(`default game tag "Site" incorrect: expected github.com/brighamskarda/chess, got %s`, game.tags["Event"])
	}
	currentDate := time.Now()
	currentDateString := fmt.Sprintf("%4d.%2d.%2d", currentDate.Year(), currentDate.Month(), currentDate.Day())
	if game.tags["Date"] != currentDateString {
		t.Errorf(`default game tag "Date" incorrect: expected %s, got %s`, currentDateString, game.tags["Event"])
	}
	if game.tags["Round"] != "1" {
		t.Errorf(`default game tag "Round" incorrect: expected 1, got %s`, game.tags["Round"])
	}
	if game.tags["White"] != "White Player" {
		t.Errorf(`default game tag "White" incorrect: expected White Player, got %s`, game.tags["White"])
	}
	if game.tags["Black"] != "Black Player" {
		t.Errorf(`default game tag "Black" incorrect: expected Black Player, got %s`, game.tags["Black"])
	}
	if game.tags["Result"] != "*" {
		t.Errorf(`default game tag "Result" incorrect: expected *, got %s`, game.tags["Result"])
	}
	if len(game.moveHistory) != 0 {
		t.Errorf("move history not empty: expected [], got %v", game.moveHistory)
	}
}

func TestGameMove(t *testing.T) {
	game := NewGame()
	err := game.Move(Move{E2, E4, NoPieceType})
	if err != nil {
		t.Errorf("incorrect result: input E2E4: expected nil, got %v", err)
	}
	tempGame := game.Copy()
	err = game.Move(Move{D2, D4, NoPieceType})
	if err == nil {
		t.Errorf("incorrect result: input D2D4: expected error, got nil")
	}
	if !reflect.DeepEqual(game, tempGame) {
		t.Error("incorrect result: after invalid move D2D4 game state changed: ", cmp.Diff(game, tempGame))
	}
	err = game.Move(Move{E7, E5, NoPieceType})
	if err != nil {
		t.Errorf("incorrect result: input E7E5: expected nil, got %v", err)
	}
	expectedMoveHistory := []Move{
		{E2, E4, NoPieceType},
		{E7, E5, NoPieceType},
	}
	if !reflect.DeepEqual(expectedMoveHistory, game.moveHistory) {
		t.Errorf("incorrect result: moveHistory incorrect: expected %v, got %v", expectedMoveHistory, game.moveHistory)
	}
}

func TestGameMoveUpdatesResultToNoResult(t *testing.T) {
	game := NewGame()
	game.SetResult(WhiteWins)
	if game.GetResult() != WhiteWins {
		t.Errorf("game.SetResult did not work")
	}
	game.Move(Move{E2, E4, NoPieceType})
	if game.GetResult() != NoResult {
		t.Errorf("game.Move did not updateResult")
	}
	game.SetResult(WhiteWins)
	game.Move(Move{D2, D4, NoPieceType})
	if game.GetResult() == NoResult {
		t.Errorf("game.Move updates result when it shouldn't")
	}
}

func TestGameMoveUpdatesResultForCheckmate(t *testing.T) {
	game := NewGame()
	game.position, _ = ParseFen("r4r1k/pp3Q2/6pB/4P3/6P1/8/P5PP/R4R1K w - - 3 22")
	game.Move(Move{F7, G7, NoPieceType})
	if game.GetResult() != WhiteWins {
		t.Errorf("game.Move did not update result for white winning")
	}

	game = NewGame()
	game.position, _ = ParseFen("K7/4kpp1/5b1p/8/8/6P1/1q3P1P/2q5 b - - 9 44")
	game.Move(Move{C1, A1, NoPieceType})
	if game.GetResult() != BlackWins {
		t.Errorf("game.Move did not update result for black winning")
	}
}

func TestGameMoveUpdatesResultForStalemate(t *testing.T) {
	game := NewGame()
	game.position, _ = ParseFen("8/8/8/8/4K2p/7k/6p1/6B1 w - - 0 50")
	game.Move(Move{E4, F3, NoPieceType})
	if game.GetResult() != Draw {
		t.Errorf("game.Move did not update result for black to move on stalemate")
	}

	game = NewGame()
	game.position, _ = ParseFen("8/7p/3b2p1/4kp2/4p3/2pn3r/8/3K4 b - - 3 50")
	game.Move(Move{H3, H2, NoPieceType})
	if game.GetResult() != Draw {
		t.Errorf("game.Move did not update result for white to move on stalemate")
	}
}

func TestGetandSetTags(t *testing.T) {
	game := NewGame()
	tag, err := game.GetTag("Round")
	if err != nil {
		t.Errorf("game.GetTag(\"Round\") returned err: %v", err)
	}
	if tag != "1" {
		t.Errorf(`game.GetTag("Round") incorrect result: expected 1, got %s`, tag)
	}

	_, err = game.GetTag("fjdkslfjd")
	if err == nil {
		t.Errorf(`game.GetTag("fjdkslfjd") did not return error`)
	}

	game.SetTag("Round", "2")
	tag, err = game.GetTag("Round")
	if err != nil {
		t.Errorf("game.GetTag(\"Round\") returned err: %v", err)
	}
	if tag != "2" {
		t.Errorf(`game.GetTag("Round") incorrect result: expected 2, got %s`, tag)
	}

	game.SetTag("Result", "1/2-1/2")
	tag, _ = game.GetTag("Result")
	if tag != "*" {
		t.Errorf(`game.SetTag("Result", "1/2-1/2") set result when it shouldn't have`)
	}
}

func TestRemoveTag(t *testing.T) {
	game := NewGame()
	game.RemoveTag("Site")
	tag, err := game.GetTag("Site")
	if tag != "github.com/brighamskarda/chess" || err != nil {
		t.Errorf("game.RemoveTag removed a tag it shouldn't have")
	}

	game.SetTag("hi", "lol")
	game.RemoveTag("hi")
	_, err = game.GetTag("hi")
	if err == nil {
		t.Errorf("game.RemoveTag did not work")
	}
}

func TestGetAllTags(t *testing.T) {
	game := NewGame()
	game.SetTag("hi", "lol")
	allTags := game.GetAllTags()
	if len(allTags) != 8 {
		t.Errorf("did not get all tags")
	}
}

func TestSetPosition(t *testing.T) {
	game := NewGame()
	game.Move(Move{E2, E4, NoPieceType})
	game.SetResult(WhiteWins)
	position, _ := ParseFen("8/7p/3b2p1/2p2p2/3kpn2/7r/8/5K2 b - - 1 46")
	err := game.SetPosition(position)
	if err != nil {
		t.Error("SetPosition game an error")
	}
	if *game.position != *position {
		t.Error("SetPosition did not set the position")
	}
	if tag, _ := game.GetTag("SetUp"); tag != "1" {
		t.Error("SetPosition did not set tag 'SetUp'")
	}
	if tag, _ := game.GetTag("FEN"); tag != "8/7p/3b2p1/2p2p2/3kpn2/7r/8/5K2 b - - 1 46" {
		t.Error("SetPosition did not set tag 'FEN'")
	}
	if len(game.moveHistory) != 0 {
		t.Error("SetPosition did not clear move history")
	}
	if game.GetResult() != NoResult {
		t.Error("SetPosition did not reset result")
	}

	position, _ = ParseFen("8/7p/3b2p1/4kp2/4p3/2pn4/7r/3K4 w - - 4 51")
	game.SetPosition(position)
	if game.GetResult() != Draw {
		t.Error("SetPosition did not set result")
	}
}

func TestSetInvalidPosition(t *testing.T) {
	game := NewGame()
	position, _ := ParseFen("8/7p/3b2p1/4kp2/4p3/2pn4/7r/3KK3 w - - 4 51")
	err := game.SetPosition(position)
	if err == nil {
		t.Errorf("game accepted an invalid position")
	}
}
