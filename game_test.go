package chess

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIsValidResult(t *testing.T) {
	if(isValidResult(5)) {
		t.Error("incorrect result: input 5: expected false, got true")
	}
	if(!isValidResult(Draw)) {
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
	if game.result != NoResult {
		t.Errorf("result incorrect: expected 0, got %v", game.result)
	}
}