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
	"bytes"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNewGameDate(t *testing.T) {
	g := NewGame()
	date := time.Now()
	dateString := fmt.Sprintf("%d.%02d.%02d", date.Year(), date.Month(), date.Day())
	if g.Date != dateString {
		t.Errorf("incorrect result: expected %s, got %s", dateString, g.Date)
	}
}

func TestNewGameFromFEN(t *testing.T) {
	fen := "rnb1kbnr/pppp1ppp/8/3qp3/8/2P3B1/PP1PPPPP/RNBQK1NR b KQkq - 0 1"
	g, err := NewGameFromFEN(fen)
	if err != nil {
		t.Errorf("expected err to be nil")
	}
	if s, _ := g.Position().MarshalText(); string(s) != fen {
		t.Errorf("position does not match: expected %s, got %s", fen, s)
	}
	if g.OtherTags["SetUp"] != "1" {
		t.Errorf("SetUp tag not set to 1")
	}
	if g.OtherTags["FEN"] != fen {
		t.Errorf("FEN tag not set")
	}
}

func TestPgnMoveCopy(t *testing.T) {
	myPgnMove := PgnMove{
		Move:              Move{A1, B1, NoPieceType},
		NumericAnnotation: 255,
		PostCommentary:    []string{"my comment"},
		Variations: [][]PgnMove{
			[]PgnMove{
				PgnMove{
					Move:              Move{},
					NumericAnnotation: 0,
					PostCommentary:    []string{},
					Variations: [][]PgnMove{
						[]PgnMove{
							PgnMove{
								Move:              Move{},
								NumericAnnotation: 0,
								PostCommentary:    []string{},
								Variations:        [][]PgnMove{},
							},
						},
					},
				},
			},
			[]PgnMove{
				PgnMove{
					Move:              Move{},
					NumericAnnotation: 0,
					PostCommentary:    []string{},
					Variations:        [][]PgnMove{},
				},
			},
		},
	}

	myPgnMoveCopy := myPgnMove.Copy()
	if len(myPgnMove.Variations) != len(myPgnMoveCopy.Variations) {
		t.Errorf("copy failed, Variations lengths do not match")
	}
	if len(myPgnMove.Variations[0][0].Variations) != len(myPgnMoveCopy.Variations[0][0].Variations) {
		t.Errorf("deep copy failed, sub variation lengths do not match")
	}

	myPgnMoveCopy.Variations[0][0].NumericAnnotation = 199
	if myPgnMove.Variations[0][0].NumericAnnotation == myPgnMoveCopy.Variations[0][0].NumericAnnotation {
		t.Errorf("copy failed, changing value is seen in both copies")
	}
}

func TestIsCheckMate(t *testing.T) {
	testCases := []struct {
		fen      string
		expected bool
	}{
		{fen: DefaultFEN, expected: false},
		{fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", expected: false},
		{fen: "8/8/8/8/8/6k1/6q1/6K1 w - - 0 1", expected: true},
		{fen: "8/8/8/8/8/8/6k1/5K2 w - - 0 1", expected: false},
		{fen: "7r/2R5/p7/4Np2/3Pk3/4P2P/P3NP2/R3K3 w Q - 7 34", expected: false},
		{fen: "7r/2R5/p7/4Np2/3Pk3/2N1P2P/P4P2/R3K3 b Q - 8 34", expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.fen, func(t *testing.T) {
			game, _ := NewGameFromFEN(tc.fen)
			actual := game.IsCheckmate()
			if tc.expected != actual {
				t.Errorf("incorrect result for %s: expected %t, got %t", tc.fen, tc.expected, actual)
			}
		})
	}
}

func TestIsStaleMate(t *testing.T) {
	testCases := []struct {
		fen      string
		expected bool
	}{
		{fen: DefaultFEN, expected: false},
		{fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", expected: false},
		{fen: "5bnr/4p1pq/4Qpkr/7p/7P/4P3/PPPP1PP1/RNB1KBNR b KQ - 2 10", expected: true},
		{fen: "8/1k4p1/7p/5p1P/1q2p3/8/2K5/7r w - - 2 61", expected: true},
		{fen: "8/1k4p1/7p/5p1P/1q2p3/2K5/8/3r4 w - - 0 60", expected: false},
		{fen: "7K/2p2k2/1p4q1/1P6/7p/6bP/8/8 w - - 10 63", expected: true},
		{fen: "8/2p3rk/pb4pp/8/8/P3q2P/4K3/3B1R2 w - - 1 34", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.fen, func(t *testing.T) {
			game, _ := NewGameFromFEN(tc.fen)
			actual := game.IsStalemate()
			if tc.expected != actual {
				t.Errorf("incorrect result for %s: expected %t, got %t", tc.fen, tc.expected, actual)
			}
		})
	}
}

func TestGameMove(t *testing.T) {
	g := NewGame()
	g.Result = WhiteWins
	err := g.Move(Move{E2, E4, NoPieceType})
	if err != nil {
		t.Errorf("got error")
	}
	if s, _ := g.Position().MarshalText(); string(s) != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("incorrect position, got %s", string(s))
	}
	if g.Result != NoResult {
		t.Errorf("incorrect result, got %v", g.Result)
	}
}

func TestGameMoveCheckmate(t *testing.T) {
	g, _ := NewGameFromFEN("8/2p3rk/pb4pp/8/8/P6P/4K3/2qB1R2 b - - 0 33")
	err := g.Move(Move{C1, E3, NoPieceType})
	if err != nil {
		t.Errorf("got error")
	}
	if g.Result != BlackWins {
		t.Errorf("incorrect result, got %v", g.Result)
	}
}

func TestGameMoveUCI(t *testing.T) {
	g := NewGame()
	g.Result = WhiteWins
	err := g.MoveUCI("E2E4")
	if err != nil {
		t.Errorf("got error")
	}
	if s, _ := g.Position().MarshalText(); string(s) != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("incorrect position, got %s", string(s))
	}
	if g.Result != NoResult {
		t.Errorf("incorrect result, got %v", g.Result)
	}
}

func TestAnnotateMove(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	if g.AnnotateMove(0, 3) != nil {
		t.Error("annotate move gave an error")
	}
	if g.AnnotateMove(1, 2) != nil {
		t.Error("annotate move gave an error")
	}
	moveHistory := g.MoveHistory()
	if moveHistory[0].NumericAnnotation != 3 {
		t.Errorf("for move 0 expected 3, got %d", moveHistory[0].NumericAnnotation)
	}
	if moveHistory[1].NumericAnnotation != 2 {
		t.Errorf("for move 1 expected 23, got %d", moveHistory[1].NumericAnnotation)
	}
}

func TestAnnotateMoveError(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	if g.AnnotateMove(-1, 3) == nil {
		t.Error("annotate move gave no error")
	}
	if g.AnnotateMove(2, 2) == nil {
		t.Error("annotate move gave no error")
	}
	moveHistory := g.MoveHistory()
	if moveHistory[0].NumericAnnotation != 0 {
		t.Errorf("for move 0 expected 0, got %d", moveHistory[0].NumericAnnotation)
	}
	if moveHistory[1].NumericAnnotation != 0 {
		t.Errorf("for move 1 expected 0, got %d", moveHistory[1].NumericAnnotation)
	}
}

func TestMoveSAN(t *testing.T) {
	g := NewGame()
	if g.MoveSAN("e4") != nil {
		t.Fatal()
	}
	if g.MoveSAN("Nc6") != nil {
		t.Fatal()
	}

	if s, _ := g.Position().MarshalText(); string(s) != "r1bqkbnr/pppppppp/2n5/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 1 2" {
		t.Errorf("MoveSAN did not work correctly, ending position was %q", string(s))
	}
}

func TestPositionPly(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	ply := 0
	expected := DefaultFEN
	actual, _ := g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 1
	expected = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	actual, _ = g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 2
	expected = "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"
	actual, _ = g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}
}

func TestPositionPly_AltStart(t *testing.T) {
	g, err := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN1 w Qkq - 0 1")
	if err != nil {
		t.Fatal()
	}
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	ply := 0
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN1 w Qkq - 0 1"
	actual, _ := g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 1
	expected = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBN1 b Qkq e3 0 1"
	actual, _ = g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 2
	expected = "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBN1 w Qkq d6 0 2"
	actual, _ = g.PositionPly(ply).MarshalText()
	if expected != string(actual) {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}
}

func TestPositionPlyError(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	if g.PositionPly(-1) != nil {
		t.Errorf("did not get nil")
	}
	if g.PositionPly(3) != nil {
		t.Errorf("did not get nil")
	}
}

func TestCommentAfterMove(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentAfterMove(0, "comment 1")
	g.CommentAfterMove(1, "comment 2")
	moveHistory := g.MoveHistory()
	if moveHistory[0].PostCommentary[0] != "comment 1" {
		t.Errorf("for move 0 got %q", moveHistory[0].PostCommentary)
	}
	if moveHistory[1].PostCommentary[0] != "comment 2" {
		t.Errorf("for move 1 got %q", moveHistory[1].PostCommentary)
	}
}

func TestCommentAfterMoveErr(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	if g.CommentAfterMove(0, "comment 1") != nil {
		t.Errorf("got error when none expected")
	}
	if g.CommentAfterMove(-1, "fff") == nil {
		t.Errorf("did not get error")
	}
	if g.CommentAfterMove(2, "fff") == nil {
		t.Errorf("did not get error")
	}
}

func TestCommentBeforeMove(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentBeforeMove(0, "comment 1")
	g.CommentBeforeMove(1, "comment 2")
	moveHistory := g.MoveHistory()
	if moveHistory[0].PreCommentary[0] != "comment 1" {
		t.Errorf("for move 0 got %q", moveHistory[0].PreCommentary)
	}
	if moveHistory[1].PreCommentary[0] != "comment 2" {
		t.Errorf("for move 1 got %q", moveHistory[1].PreCommentary)
	}
}

func TestCommentBeforeMoveErr(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	if g.CommentBeforeMove(0, "comment 1") != nil {
		t.Errorf("got error when none expected")
	}
	if g.CommentBeforeMove(-1, "fff") == nil {
		t.Errorf("did not get error")
	}
	if g.CommentBeforeMove(2, "fff") == nil {
		t.Errorf("did not get error")
	}
}

func TestDeleteCommentBefore(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentBeforeMove(0, "comment 1")
	g.CommentBeforeMove(1, "comment 2")
	g.CommentBeforeMove(0, "comment 3")
	g.CommentBeforeMove(1, "comment 4")
	moveHistory := g.MoveHistory()
	if len(moveHistory[0].PreCommentary) != 2 || len(moveHistory[1].PreCommentary) != 2 {
		t.Errorf("comment before not working")
	}
	g.DeleteCommentBefore(0, 1)
	moveHistory = g.MoveHistory()
	if len(moveHistory[0].PreCommentary) != 1 {
		t.Errorf("len of moveHistory[0].PreCommentary incorrect, got %d", len(moveHistory[0].PreCommentary))
	}
	if moveHistory[0].PreCommentary[0] != "comment 1" {
		t.Errorf("deleted wrong comment from move 0")
	}

	g.DeleteCommentBefore(1, 0)
	moveHistory = g.MoveHistory()
	if len(moveHistory[1].PreCommentary) != 1 {
		t.Errorf("len of moveHistory[1].PreCommentary incorrect, got %d", len(moveHistory[1].PreCommentary))
	}
	if moveHistory[1].PreCommentary[0] != "comment 4" {
		t.Errorf("deleted wrong comment from move 1")
	}
}

func TestDeleteCommentBeforeError(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentBeforeMove(0, "comment 1")
	g.CommentBeforeMove(1, "comment 2")

	if g.DeleteCommentBefore(0, 0) != nil {
		t.Error("got error when none expected")
	}
	if g.DeleteCommentBefore(-1, 0) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentBefore(0, 0) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentBefore(1, 1) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentBefore(1, -1) == nil {
		t.Error("did not get error")
	}
}

func TestDeleteCommentAfter(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentAfterMove(0, "comment 1")
	g.CommentAfterMove(1, "comment 2")
	g.CommentAfterMove(0, "comment 3")
	g.CommentAfterMove(1, "comment 4")
	moveHistory := g.MoveHistory()
	if len(moveHistory[0].PostCommentary) != 2 || len(moveHistory[1].PostCommentary) != 2 {
		t.Errorf("comment before not working")
	}
	g.DeleteCommentAfter(0, 1)
	moveHistory = g.MoveHistory()
	if len(moveHistory[0].PostCommentary) != 1 {
		t.Errorf("len of moveHistory[0].PreCommentary incorrect, got %d", len(moveHistory[0].PostCommentary))
	}
	if moveHistory[0].PostCommentary[0] != "comment 1" {
		t.Errorf("deleted wrong comment from move 0")
	}

	g.DeleteCommentAfter(1, 0)
	moveHistory = g.MoveHistory()
	if len(moveHistory[1].PostCommentary) != 1 {
		t.Errorf("len of moveHistory[1].PreCommentary incorrect, got %d", len(moveHistory[1].PostCommentary))
	}
	if moveHistory[1].PostCommentary[0] != "comment 4" {
		t.Errorf("deleted wrong comment from move 1")
	}
}

func TestDeleteCommentAfterError(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fatal()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fatal()
	}

	g.CommentAfterMove(0, "comment 1")
	g.CommentAfterMove(1, "comment 2")

	if g.DeleteCommentAfter(0, 0) != nil {
		t.Error("got error when none expected")
	}
	if g.DeleteCommentAfter(-1, 0) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentAfter(0, 0) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentAfter(1, 1) == nil {
		t.Error("did not get error")
	}
	if g.DeleteCommentAfter(1, -1) == nil {
		t.Error("did not get error")
	}
}

func makeGameWithVariation() *Game {
	g := NewGame()
	g.Move(Move{E2, E4, NoPieceType})
	g.Move(Move{D7, D5, NoPieceType})

	g.MakeVariation(1, []PgnMove{{
		Move:              Move{B8, C6, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}})
	return g
}

func TestMakeVariation(t *testing.T) {
	g := NewGame()
	g.Move(Move{E2, E4, NoPieceType})
	g.Move(Move{D7, D5, NoPieceType})

	err := g.MakeVariation(1, []PgnMove{{
		Move:              Move{B8, C6, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}})

	if err != nil {
		t.Error("did not expect error")
	}
	moveHistory := g.MoveHistory()
	if len(moveHistory[1].Variations) != 1 {
		t.Errorf("variation not added correctly")
	}
	if len(moveHistory[1].Variations[0]) != 2 {
		t.Errorf("variation not added correctly")
	}
}

func TestMakeVariationError(t *testing.T) {
	g := NewGame()
	g.Move(Move{E2, E4, NoPieceType})
	g.Move(Move{D7, D5, NoPieceType})

	variationToAdd := []PgnMove{{
		Move:              Move{B8, C6, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}}

	err := g.MakeVariation(0, variationToAdd)
	if err == nil {
		t.Error("did not get error")
	}

	err = g.MakeVariation(2, variationToAdd)
	if err == nil {
		t.Error("did not get error")
	}

	err = g.MakeVariation(-1, variationToAdd)
	if err == nil {
		t.Error("did not get error")
	}

	moveHistory := g.MoveHistory()
	for i, m := range moveHistory {
		if len(m.Variations) != 0 {
			t.Errorf("move %d has a more than 0 variations", i)
		}
	}
}

func TestDeleteVariation(t *testing.T) {
	g := makeGameWithVariation()
	if g.DeleteVariation(1, 0) != nil {
		t.Error("DeleteVariation produced an error.")
	}
	moveHistory := g.MoveHistory()
	if len(moveHistory[1].Variations) != 0 {
		t.Errorf("variation not deleted correctly")
	}
}

func TestDeleteVariationError(t *testing.T) {
	g := makeGameWithVariation()
	if g.DeleteVariation(-1, 0) == nil {
		t.Error("DeleteVariation produced no error.")
	}
	if g.DeleteVariation(0, 0) == nil {
		t.Error("DeleteVariation produced no error.")
	}
	if g.DeleteVariation(2, 0) == nil {
		t.Error("DeleteVariation produced no error.")
	}
	if g.DeleteVariation(1, -1) == nil {
		t.Error("DeleteVariation produced no error.")
	}
	if g.DeleteVariation(1, 1) == nil {
		t.Error("DeleteVariation produced no error.")
	}
	moveHistory := g.MoveHistory()
	if len(moveHistory[1].Variations) != 1 {
		t.Errorf("variation incorrectly deleted")
	}
}

func TestGetVariation(t *testing.T) {
	g := makeGameWithVariation()

	g.MakeVariation(1, []PgnMove{{
		Move:              Move{G8, F6, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}})

	g.MakeVariation(0, []PgnMove{{
		Move:              Move{H2, H4, NoPieceType},
		NumericAnnotation: 0,
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	}})

	newG, err := g.GetVariation(1, 0)
	if err != nil {
		t.Errorf("returned error when none expected")
	}

	expectedPosition := "r1bqkbnr/pppppppp/2n5/8/3PP3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 2"
	actualPosition, _ := newG.Position().MarshalText()
	if expectedPosition != string(actualPosition) {
		t.Errorf("variation does not match expected position: expected %q, got %q", expectedPosition, actualPosition)
	}

	moveHistory := newG.MoveHistory()
	if len(moveHistory[1].Variations) != 2 {
		t.Errorf("alt variations not preserved")
	}
	if len(moveHistory[1].Variations[0]) != 1 || moveHistory[1].Variations[0][0].Move != (Move{D7, D5, NoPieceType}) {
		t.Errorf("original move not preserved")
	}
	if len(moveHistory[1].Variations[1]) != 2 {
		t.Errorf("alt variation not preserved")
	}
	if len(moveHistory[0].Variations) != 1 {
		t.Errorf("move 0 variation not preserved")
	}
}

func TestGetVariationError(t *testing.T) {
	g := makeGameWithVariation()

	_, err := g.GetVariation(0, 0)
	if err == nil {
		t.Error("expected error for 0, 0")
	}

	_, err = g.GetVariation(-1, 0)
	if err == nil {
		t.Error("expected error for -1, 0")
	}

	_, err = g.GetVariation(2, 0)
	if err == nil {
		t.Error("expected error for 2, 0")
	}

	_, err = g.GetVariation(1, -1)
	if err == nil {
		t.Error("expected error for 1, -1")
	}

	_, err = g.GetVariation(1, 1)
	if err == nil {
		t.Error("expected error for 1, 1")
	}
}

func TestGameMarshalText(t *testing.T) {
	g := NewGame()
	moves := []string{
		"d4",
		"e6",
		"e4",
		"d5",
		"exd5",
		"exd5",
		"Nf3",
		"Nf6",
		"Ne5",
		"Qe7",
		"f4",
		"Bg4",
		"Be2",
		"Bd7",
		"g4",
		"Ne4",
		"c4",
		"Qh4+",
		"Kf1",
		"Qf2#",
	}

	for _, m := range moves {
		g.MoveSAN(m)
	}

	g.CommentBeforeMove(0, "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20")
	g.AnnotateMove(4, 2)
	g.AnnotateMove(5, 10)
	g.CommentAfterMove(19, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				PostCommentary:    []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.Date = "2025.04.09"
	g.OtherTags["WhiteElo"] = "1090"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "0-1"]
[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20} 1. d4 e6
2. e4 (2. f4 g5 {Another variation here} {another comment here} (2... h5!) 3.
h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7. Be2 Bd7
8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.}
0-1`
	actual, _ := g.MarshalText()
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameMarshalText_AltStart(t *testing.T) {
	g, _ := NewGameFromFEN("r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16")
	moves := []string{
		"Qd3",
		"Qg4+",
		"Kf7",
	}

	for _, m := range moves {
		if g.MoveSAN(m) != nil {
			t.Fatal()
		}
	}
	g.Date = "2025.04.09"
	g.OtherTags["WhiteElo"] = "1090"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "*"]
[FEN "r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16"]
[SetUp "1"]
[WhiteElo "1090"]

16... Qd3 17. Qg4+ Kf7 *`
	actual, _ := g.MarshalText()
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameMarshalText_NumericAnnotationGlyphs(t *testing.T) {
	g := NewGame()
	moves := []string{
		"d4",
		"e6",
		"e4",
		"d5",
		"exd5",
		"exd5",
		"Nf3",
	}

	for _, m := range moves {
		g.MoveSAN(m)
	}

	g.AnnotateMove(1, 1)
	g.AnnotateMove(2, 2)
	g.AnnotateMove(3, 3)
	g.AnnotateMove(4, 4)
	g.AnnotateMove(5, 5)
	g.AnnotateMove(6, 6)

	g.Date = "2025.04.09"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "*"]

1. d4 e6! 2. e4? d5!! 3. exd5?? exd5!? 4. Nf3?! *`
	actual, _ := g.MarshalText()
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameMarshalText2(t *testing.T) {
	g := NewGame()
	moves := []string{
		"d4",
		"e6",
		"e4",
		"d5",
		"exd5",
		"exd5",
		"Nf3",
	}

	for _, m := range moves {
		g.MoveSAN(m)
	}

	g.AnnotateMove(1, 1)
	g.AnnotateMove(2, 2)
	g.AnnotateMove(3, 3)
	g.AnnotateMove(4, 4)
	g.AnnotateMove(5, 5)
	g.AnnotateMove(6, 6)

	g.Date = "2025.04.09"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "*"]

1. d4 e6! 2. e4? d5!! 3. exd5?? exd5!? 4. Nf3?! *`
	actual, err := g.MarshalText()
	if err != nil {
		t.Errorf("incorrect result: expected err to be nil.")
	}
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameMarshalTextReduced(t *testing.T) {
	g := NewGame()
	moves := []string{
		"d4",
		"e6",
		"e4",
		"d5",
		"exd5",
		"exd5",
		"Nf3",
		"Nf6",
		"Ne5",
		"Qe7",
		"f4",
		"Bg4",
		"Be2",
		"Bd7",
		"g4",
		"Ne4",
		"c4",
		"Qh4+",
		"Kf1",
		"Qf2#",
	}

	for _, m := range moves {
		g.MoveSAN(m)
	}

	g.CommentBeforeMove(0, "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20")
	g.AnnotateMove(4, 2)
	g.AnnotateMove(5, 10)
	g.CommentAfterMove(19, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				PostCommentary:    []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.Date = "2025.04.09"
	g.OtherTags["WhiteElo"] = "1090"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "0-1"]

1. d4 e6 2. e4 d5 3. exd5 exd5 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7. Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2# 0-1`
	actual, _ := g.MarshalTextReduced()
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameMarshalTextReduced_AltStart(t *testing.T) {
	g, _ := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/1NBQKBNR b Kkq - 0 1")
	moves := []string{
		"e6",
		"e4",
		"d5",
		"exd5",
		"exd5",
		"Nf3",
		"Nf6",
		"Ne5",
		"Qe7",
		"f4",
		"Bg4",
		"Be2",
		"Bd7",
		"g4",
		"Ne4",
		"c4",
		"Qh4+",
		"Kf1",
		"Qf2#",
	}

	for _, m := range moves {
		g.MoveSAN(m)
	}

	g.CommentBeforeMove(0, "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20")
	g.AnnotateMove(3, 2)
	g.AnnotateMove(4, 10)
	g.CommentAfterMove(18, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(1, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				PostCommentary:    []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(1, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			PostCommentary:    []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.Date = "2025.04.09"
	g.OtherTags["WhiteElo"] = "1090"

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "0-1"]
[FEN "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/1NBQKBNR b Kkq - 0 1"]
[SetUp "1"]

1... e6 2. e4 d5 3. exd5 exd5 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7. Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2# 0-1`
	actual, _ := g.MarshalTextReduced()
	if string(actual) != expected {
		t.Errorf(`incorrect result: expected 
"""
%s
"""

got 
"""
%s
"""`, expected, actual)
	}
}

func TestGameUnmarshal(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "idc"]
[Site "ur mom's house"]
[Date "2025.04.09"]
[Round "2"]
[White "phil"]
[Black "donna"]
[Result "0-1"]
[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}
1. d4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7.
Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.}
0-1`))

	if err != nil {
		t.Fatalf("err != nil: %s", err)
	}

	if g.Event != "idc" {
		t.Errorf("event incorrect")
	}

	if g.Site != "ur mom's house" {
		t.Errorf("site incorrect")
	}

	if g.Date != "2025.04.09" {
		t.Errorf("date incorrect")
	}

	if g.Round != "2" {
		t.Errorf("round incorrect")
	}

	if g.White != "phil" {
		t.Errorf("white player incorrect")
	}

	if g.Black != "donna" {
		t.Errorf("black player incorrect")
	}

	if g.Result != BlackWins {
		t.Errorf("result incorrect")
	}

	if g.OtherTags["WhiteElo"] != "1090" {
		t.Errorf("whiteElo incorrect")
	}

	moveHis := g.MoveHistory()
	if moveHis[0].PreCommentary[0] != "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20" {
		t.Errorf("commentary incorrect")
	}
	if len(moveHis) != 20 {
		t.Errorf("moveHistory incorrect length")
	}

	if len(moveHis[0].Variations) != 0 {
		t.Errorf("moveHistory has variation where none are present")
	}

	if len(moveHis[2].Variations) != 2 {
		t.Errorf("ply 2 does not have 2 variations")
	}

	if len(moveHis[2].Variations[0][1].PostCommentary) != 2 {
		t.Errorf("ply 2 missing commentary")
	}

	if len(moveHis[2].Variations[0]) != 3 {
		t.Errorf("variation length incorrect")
	}

	if moveHis[4].NumericAnnotation != 2 {
		t.Errorf("ply 4 missing numeric annotation")
	}

	if moveHis[5].NumericAnnotation != 10 {
		t.Errorf("ply 4 missing numeric annotation")
	}
}

func TestGameUnmarshalComments(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "idc"]
[Site "ur mom's house"]
[Date "2025.04.09"]
[Round "2"]
[White "phil"]
[Black "donna"]
[Result "0-1"]
%[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}
1. d4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4; this is a comment 
7. Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need ;this comment to 
be even 
longer than before, preferably longer than 80% characters for some testing.}
0-1`))

	if err != nil {
		t.Fatalf("err != nil: %s", err)
	}
	moveHis := g.MoveHistory()
	if moveHis[11].PostCommentary[0] != "this is a comment" {
		t.Errorf("semicolon move comment not parsed")
	}
	if moveHis[19].PostCommentary[0] != `Black wins by checkmate. Now I need ;this comment to 
be even 
longer than before, preferably longer than 80% characters for some testing.` {
		t.Errorf("multiline commentary not parsed")
	}
	if _, ok := g.OtherTags["WhiteElo"]; ok != false {
		t.Errorf("dev escape not working")
	}
}

func TestGameUnmarshalPreComments(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "idc"]
[Site "ur mom's house"]
[Date "2025.04.09"]
[Round "2"]
[White "phil"]
[Black "donna"]
[Result "0-1"]
%[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}{precomment 2}
1. d4 e6 2. e4 {This is a post comment} ({ This is a pre comment}{This is another pre comment }2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4; this is a comment 
7. Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need ;this comment to 
be even 
longer than before, preferably longer than 80% characters for some testing.}
0-1`))

	if err != nil {
		t.Fatalf("err != nil: %s", err)
	}
	moveHis := g.MoveHistory()
	if len(moveHis[0].PreCommentary) != 2 {
		t.Errorf("start of game pre commentary not working")
	}
	if moveHis[0].PreCommentary[0] != "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20" {
		t.Errorf("start of game pre commentary comment 1 incorrect")
	}
	if moveHis[0].PreCommentary[1] != "precomment 2" {
		t.Errorf("start of game pre commentary comment 2 incorrect")
	}
	if len(moveHis[2].PreCommentary) != 0 {
		t.Errorf("found pre commentary on move 2, should not exist")
	}
	if moveHis[2].PostCommentary[0] != "This is a post comment" {
		t.Errorf("post commentary on move 2 incorrect, got %s", moveHis[2].PostCommentary[0])
	}
	if len(moveHis[2].Variations[0][0].PreCommentary) != 2 {
		t.Errorf("pre commentary on first variation incorrect length")
	}
	variation := moveHis[2].Variations[0]
	if variation[0].PreCommentary[0] != "This is a pre comment" {
		t.Errorf("variation pre comment 1 incorrect, got %s", variation[0].PreCommentary[0])
	}
	if variation[0].PreCommentary[1] != "This is another pre comment" {
		t.Errorf("variation pre comment 1 incorrect, got %s", variation[0].PreCommentary[1])
	}
}

func TestGameUnmarshalAltStart(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "*"]
[FEN "r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16"]
[SetUp "1"]
[WhiteElo "1090"]

16... Qd3 17. Qg4+ Kf7 *`))

	if err != nil {
		t.Errorf("err != nil")
	}

	if g.OtherTags["FEN"] != "r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16" {
		t.Errorf("fen not set")
	}

	if g.OtherTags["SetUp"] != "1" {
		t.Errorf("setup != 1")
	}

	if s, _ := g.PositionPly(0).MarshalText(); string(s) != "r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16" {
		t.Errorf("alt start position not set")
	}

	if s, _ := g.PositionPly(1).MarshalText(); string(s) != "r6r/ppp3pp/2n1Nnk1/4p3/2Q5/B2q4/P4PPP/RN3RK1 w - - 1 17" {
		t.Errorf("alt ply 1 incorrect")
	}
}

func TestGameUnmarshalFailure(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "idc"]
[Site "ur mom's house"]
[Date "2025.04.09"]
[Round "2"]
[White "phil"]
[Black "donna"]
[Result "0-1"]
[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}
1. dfs4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4; this is a comment 
7. Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need ;this comment to 
be even longer than before, preferably longer than 80% characters for some testing.}
0-1`))

	if err == nil {
		t.Errorf("err == nil")
	}
	if !compareGames(g, NewGame()) {
		t.Errorf("game was modified on failure")
	}
}

// compareGames returns true if they are the same.
func compareGames(g1 *Game, g2 *Game) bool {
	return g1.Event == g2.Event &&
		g1.Site == g2.Site &&
		g1.Date == g2.Date &&
		g1.Round == g2.Round &&
		g1.White == g2.White &&
		g1.Black == g2.Black &&
		g1.Result == g2.Result &&
		maps.Equal(g1.OtherTags, g2.OtherTags) &&
		compareMoveHistories(g1.moveHistory, g2.moveHistory)
}

func compareMoveHistories(mh1 []PgnMove, mh2 []PgnMove) bool {
	if len(mh1) != len(mh2) {
		return false
	}
	for i := range mh1 {
		if !slices.Equal(mh1[i].PostCommentary, mh2[i].PostCommentary) ||
			!slices.Equal(mh1[i].PreCommentary, mh2[i].PreCommentary) ||
			mh1[i].NumericAnnotation != mh2[i].NumericAnnotation ||
			mh1[i].Move != mh2[i].Move ||
			len(mh1[i].Variations) != len(mh2[i].Variations) {
			return false
		}
		for j := range mh1[i].Variations {
			if !compareMoveHistories(mh1[i].Variations[j], mh2[i].Variations[j]) {
				return false
			}
		}
	}
	return true
}

func TestGameUnmarshal_randomBraces(t *testing.T) {
	g := NewGame()
	g.UnmarshalText([]byte("[] \"]\"]\n\nRandom game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}\n1. d4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...\nh5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7.\nBe2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#\n{"))
	// Just make sure it doesn't panic
}

func TestBlackStart_altPos(t *testing.T) {
	g := &Game{}
	err := g.UnmarshalText([]byte(`[Event "Hamguy123's Study: Chapter 1"]
[Result "*"]
[Variant "From Position"]
[FEN "rnbqk2r/pppp1ppp/4pn2/8/1bPP4/2N5/PPQ1PPPP/R1B1KBNR b KQkq - 0 1"]
[ECO "?"]
[Opening "?"]
[StudyName "Hamguy123's Study"]
[ChapterName "Chapter 1"]
[SetUp "1"]
[UTCDate "2025.05.22"]
[UTCTime "01:06:42"]
[Annotator "https://lichess.org/@/Hamguy123"]
[ChapterURL "https://lichess.org/study/lfAbjJ1u/fSzsIBlu"]

1... Nd5 2. e3 Qf6 3. Bd2 O-O 4. O-O-O Nc6 5. Be2 b6 6. Nh3 Rb8 7. Rhf1 *
`))
	if err != nil {
		t.Fatal()
	}
}

func TestParsePgn(t *testing.T) {
	f, err := os.Open("./testdata/SaintLouis2023.pgn")
	if err != nil {
		t.Errorf("issue reading test file at \"./testdata/SaintLouis2023.pgn\"")
	}
	defer f.Close()
	games, err := ParsePGN(f)

	if err != nil {
		t.Fatalf("got error parsing pgn file: %v", err)
	}

	if len(games) != 91 {
		t.Fatalf("did not parse all 91 games.")
	}

	if len(games[0].MoveHistory()) != 89 {
		t.Errorf("first game move length not 89, got %d", len(games[0].MoveHistory()))
	}

	if games[0].Result != WhiteWins {
		t.Errorf("first game result incorrect, got %v", games[0].Result)
	}

	if games[0].Black != "Sevian,Samuel" {
		t.Errorf("first game black incorrect, got %s", games[0].Black)
	}

	if len(games[90].MoveHistory()) != 137 {
		t.Errorf("first game move length not 137, got %d", len(games[90].MoveHistory()))
	}

	if games[90].Result != Draw {
		t.Errorf("first game result incorrect, got %v", games[90].Result)
	}

	if games[90].Black != "Lee,Alice" {
		t.Errorf("first game black incorrect, got %s", games[90].Black)
	}
}

func TestParsePgn_BadGame(t *testing.T) {
	f, err := os.Open("./testdata/SaintLouis2023.pgn")
	if err != nil {
		t.Errorf("issue reading test file at \"./testdata/SaintLouis2023.pgn\"")
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("issue reading test file at \"./testdata/SaintLouis2023.pgn\"")
	}
	b[726] = 'k'
	br := bytes.NewReader(b)
	games, err := ParsePGN(br)

	if err == nil {
		t.Fatalf("did not get error parsing pgn file")
	}

	if len(games) != 90 {
		t.Fatalf("did not parse 90 games, got %d", len(games))
	}

	if len(games[0].MoveHistory()) != 143 {
		t.Errorf("first game move length not 143, got %d", len(games[0].MoveHistory()))
	}
}

func TestThreeFoldDraw_Basic(t *testing.T) {
	pgn := []byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]
[Link "https://www.chess.com/analysis?tab=analysis"]

1. e4 e5 2. Ke2 Ke7 3. Ke1 Ke8 4. Ke2 Ke7 5. Ke1 Ke8 6. Ke2 Ke7 *`)

	g := Game{}
	err := g.UnmarshalText(pgn)
	if err != nil {
		t.Error("issue unmarshaling game")
	}
	if !g.CanClaimDrawThreeFold() {
		t.Error("could not claim three fold draw")
	}
}

func TestThreeFoldDraw_Long(t *testing.T) {
	pgn := []byte(`[Event "rated blitz game"]
[Site "https://lichess.org/uI3XYntv"]
[Date "2025.06.10"]
[White "pawnbroker222"]
[Black "Saint-Denis"]
[Result "*"]
[GameId "uI3XYntv"]
[UTCDate "2025.06.10"]
[UTCTime "17:10:39"]
[WhiteElo "1730"]
[BlackElo "1782"]
[WhiteRatingDiff "+1"]
[BlackRatingDiff "-1"]
[Variant "Standard"]
[TimeControl "180+2"]
[ECO "A11"]
[Opening "English Opening: Caro-Kann Defensive System"]
[Termination "Normal"]

1. c4 c6 2. Nc3 d5 3. cxd5 cxd5 4. d4 Nf6 5. g3 Nc6 6. Bg2 Bg4 7. Nf3 a6 8. O-O e6 9. e3 Bd6 10. a3 Rc8 11. Bd2 O-O 12. Na4 Ne4 13. Nc5 Nxd2 14. Qxd2 Bxc5 15. dxc5 Re8 16. b4 e5 17. h3 Bxf3 18. Bxf3 e4 19. Bg2 Ne5 20. Rfd1 Nd3 21. Bf1 Qd7 22. Kg2 Re6 23. Bxd3 exd3 24. Qxd3 Rh6 25. Rh1 Qc6 26. Qd4 Re8 27. Rad1 Rh5 28. Qf4 h6 29. Rd2 d4+ 30. Kh2 dxe3 31. fxe3 Qe6 32. g4 Rg5 33. Re1 f5 34. gxf5 Rxf5 35. Qd6 Qe4 36. Qd4 Qe5+ 37. Qxe5 Rfxe5 38. Rde2 g5 39. Kg3 Kf7 40. e4 Kf6 41. Kg4 Kg6 42. Re3 h5+ 43. Kg3 R8e7 44. Kf3 Rf7+ 45. Kg3 Rf4 46. R1e2 h4+ 47. Kg2 g4 48. hxg4 Rxg4+ 49. Kh2 Kg5 50. Kh3 Rf4 51. Rg2+ Kh5 52. Rg8 Rfxe4 53. Rh8+ Kg5 54. Rg8+ Kf6 55. Rxe4 Rxe4 56. Rf8+ Ke7 57. Rf3 Kd7 58. Rb3 Kc6 59. Rd3 Kb5 60. Rd7 Kc6 61. Rd3 Rc4 62. Kg2 Kc7 63. Kf2 b6 64. cxb6+ Kxb6 65. Ke2 Kb5 66. Kd2 Ka4 67. Ke2 a5 68. bxa5 Kxa5 69. Kf2 Ka4 70. Kg2 Re4 71. Kh3 Kb5 72. Rf3 Ra4 73. Rc3 Kb6 74. Kh2 Kb5 75. Kh3 Kb6 76. Rd3 Kb5 77. Re3 Kc6 78. Rf3 Kd6 79. Rf6+ Ke5 80. Rf3 Ke6 81. Re3+ Kf5 82. Rf3+ Kg5 83. Rb3 Kh5 84. Rb4 Rxa3+ 85. Kh2 h3 86. Rc4 Kg5 87. Kh1 Rf3 88. Kh2 Re3 89. Ra4 Kh5 90. Rb4 Rf3 91. Ra4 Kg5 92. Rb4 Kf5 93. Rc4 Ke5 94. Rb4 Kf5 95. Rc4 Kg5 96. Rd4 Kh5 97. Re4 Kg5 98. Rd4 Kh5 99. Rc4 Kg5 100. Rb4 Kh5 101. Rc4 Kg5 *`)

	g := Game{}
	err := g.UnmarshalText(pgn)
	if err != nil {
		t.Error("issue unmarshaling game")
	}
	if !g.CanClaimDrawThreeFold() {
		t.Error("could not claim three fold draw")
	}
}

func TestThreeFoldDraw_EnPassantDifference(t *testing.T) {
	pgn := []byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

1. e4 e6 2. Ke2 Ke7 3. Ke1 Ke8 4. e5 Ke7 5. Ke2 f5 6. Ke1 Ke8 7. Ke2 Ke7 8. Ke1
Ke8 9. Ke2 Ke7 *`)

	g := Game{}
	err := g.UnmarshalText(pgn)
	if err != nil {
		t.Error("issue unmarshaling game")
	}
	if g.CanClaimDrawThreeFold() {
		t.Error("could claim three fold draw")
	}

	if g.Move(Move{E2, E1, NoPieceType}) != nil {
		t.Error("issue performing move")
	}
	if !g.CanClaimDrawThreeFold() {
		t.Error("could not claim three fold draw")
	}
}

func TestThreeFoldDraw_CastleDifference(t *testing.T) {
	pgn := []byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

1. e4 d5 2. Nf3 dxe4 3. Ng5 Nf6 4. Bc4 Nc6 5. Bxf7+ Kd7 6. Be6+ Ke8 7. Bf7+ Kd7
8. Be6+ Ke8 9. Bf7+ *`)

	g := Game{}
	err := g.UnmarshalText(pgn)
	if err != nil {
		t.Error("issue unmarshaling game")
	}
	if g.CanClaimDrawThreeFold() {
		t.Error("could claim three fold draw")
	}

	if g.Move(Move{E8, D7, NoPieceType}) != nil {
		t.Error("issue performing move")
	}
	if !g.CanClaimDrawThreeFold() {
		t.Error("could not claim three fold draw")
	}
}

func TestThreeFoldRepetition_NoLegalEnPassant(t *testing.T) {
	pgn := []byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

1. e4 e5 2. Nf3 Nf6 3. Ng1 Ng8 4.Nf3 Nf6 5. Ng1 Ng8 *`)
	game := &Game{}
	game.UnmarshalText(pgn)

	if game.CanClaimDrawThreeFold() == false {
		t.Error("thought three fold draw was legal when it was not")
	}

	pgn = []byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

1. e4 c5 2. e5 Qb6 3. Nf3 Qe6 4.Ng1 d5 5.Nf3 Nc6 6.Ng1 Nb8 7.Nf3 Nc6 8.Ng1 Nb8 *`)
	game = &Game{}
	game.UnmarshalText(pgn)

	if game.CanClaimDrawThreeFold() == false {
		t.Error("thought three fold draw was legal when it was not")
	}
}

// TestRealPGNs ignores files that are in subdirectories. Useful for ignoring large files.
func TestRealPGNs(t *testing.T) {
	testdir := "./testdata/extra_pgns"
	files, err := os.ReadDir(testdir)
	if err != nil {
		t.Fatalf("could not read test directory ./testdata/extra_pgns: %v", err)
	}

	for _, fileEntry := range files {
		if fileEntry.IsDir() || fileEntry.Name() == ".gitkeep" {
			continue
		}

		t.Logf("parsing file %s", fileEntry.Name())
		file, err := os.Open(testdir + "/" + fileEntry.Name())
		if err != nil {
			t.Errorf("could not read file %s", fileEntry.Name())
		}
		defer file.Close()
		pgn, err := io.ReadAll(file)
		if err != nil {
			t.Errorf("issues reading file %s", fileEntry.Name())
		}

		pgnReader := bytes.NewReader(pgn)
		games, errs := ParsePGN(pgnReader)
		if errs != nil {
			t.Errorf("error parsing file %s: %s", fileEntry.Name(), errs.Error())
		}

		_, err = createPgn(games)
		if err != nil {
			t.Errorf("error regenerating pgn for file %s: %s", fileEntry.Name(), err)
		}
	}
}

func createPgn(games []*Game) (string, error) {
	sb := strings.Builder{}
	for i, g := range games {
		pgn, err := g.MarshalText()
		if err != nil {
			return sb.String(), fmt.Errorf("error generating pgn for game %d: %w", i, err)
		}
		sb.Write(pgn)
		sb.WriteRune('\n')
	}
	return sb.String(), nil
}

func TestNewGameFromFEN_missingKing(t *testing.T) {
	_, err := NewGameFromFEN("rnbq1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 0 1")
	if err == nil {
		t.Error("did not get error")
	}
}

func TestNewGameFromFEN_setResult(t *testing.T) {
	g, err := NewGameFromFEN("rnbqkbnr/pppppQpp/8/8/2B5/8/PPPPPPPP/RNB1K1NR b KQkq - 0 1")
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	if g.Result != WhiteWins {
		t.Error("result not set")
	}
}

func TestGameUnmarshal_NoMoves(t *testing.T) {
	g := NewGame()
	err := g.UnmarshalText([]byte(`[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2000.01.01"]
[Round "1"]
[White "Gopher 1"]
[Black "Gopher 2"]
[Result "0-1"]

0-1`))
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}

	if g.Result != BlackWins || g.White != "Gopher 1" {
		t.Error("game not modified")
	}
}

func TestGameMarshal_NoMoves(t *testing.T) {
	g := NewGame()
	g.Date = "2025.04.09"
	g.Result = Draw
	actual, err := g.MarshalText()
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}

	expected := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "1/2-1/2"]

1/2-1/2`
	if string(actual) != expected {
		t.Errorf("output did not match expected. Expected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

func TestParsePgn_NotAPgn(t *testing.T) {
	pgn := []byte(`This text represents a file that is not a pgn`)
	_, err := ParsePGN(bytes.NewReader(pgn))
	if err == nil {
		t.Error("did not get error for reader that does not represent a pgn")
	}
}

func TestParsePgn_PartialPgn(t *testing.T) {
	pgn := []byte(`[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2025.04.09"]
[Round "1"]
[White "?"]
[Black "?"]
[Result "1/2-1/2"]`)
	_, err := ParsePGN(bytes.NewReader(pgn))
	if err == nil {
		t.Error("did not get error for partial pgn")
	}
}

func BenchmarkParsePgn(b *testing.B) {
	file, err := os.Open("./testdata/SaintLouis2023.pgn")
	if err != nil {
		b.Fatalf("issue reading test file at \"./testdata/SaintLouis2023.pgn\"")
	}
	defer file.Close()

	pgn, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("issue reading test file at \"./testdata/SaintLouis2023.pgn\"")
	}

	r := bytes.NewReader(pgn)
	for b.Loop() {
		ParsePGN(r)
		r.Reset(pgn)
	}
}

func BenchmarkCanClaimThreeFold(b *testing.B) {
	pgn := []byte(`[Event "rated blitz game"]
[Site "https://lichess.org/Cb0mTytV"]
[Date "2025.06.21"]
[White "llimllib"]
[Black "MATSER-3000"]
[Result "1/2-1/2"]
[GameId "Cb0mTytV"]
[UTCDate "2025.06.21"]
[UTCTime "01:51:13"]
[WhiteElo "1717"]
[BlackElo "1708"]
[WhiteRatingDiff "+0"]
[BlackRatingDiff "+0"]
[Variant "Standard"]
[TimeControl "300+0"]
[ECO "C66"]
[Opening "Ruy Lopez: Berlin Defense, Improved Steinitz Defense"]
[Termination "Time forfeit"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 Nf6 4. O-O d6 5. c3 Be7 6. Re1 O-O 7. d3 h6 8. Nbd2 a6 9. Ba4 Be6 10. d4 exd4 11. cxd4 Bg4 12. h3 Bh5 13. Nf1 b5 14. Bb3 Na5 15. Bc2 c5 16. Ng3 Bg6 17. e5 dxe5 18. dxe5 Bxc2 19. Qxc2 Nd5 20. Nf5 Qc7 21. Nxe7+ Nxe7 22. Be3 c4 23. Rad1 Nac6 24. Bc5 Ng6 25. Bxf8 Rxf8 26. e6 Nb4 27. Qe4 Ne7 28. exf7+ Rxf7 29. a3 Nbc6 30. Ne5 Rf6 31. Nxc6 Nxc6 32. Qe8+ Kh7 33. Rd7 Qf4 34. Re2 Qc1+ 35. Kh2 Qf4+ 36. g3 Qf3 37. Qe4+ Qxe4 38. Rxe4 Rxf2+ 39. Kg1 Rxb2 40. Rd6 Na5 41. Rxa6 Nb3 42. Re7 Rb1+ 43. Kf2 Nd4 44. Raa7 Nf5 45. Rf7 Nd6 46. Rxg7+ Kh8 47. Rh7+ Kg8 48. Rxh6 Ne4+ 49. Kf3 Ng5+ 50. Kg4 Nf7 51. Rc6 Ne5+ 52. Kf4 Nxc6 53. Ra8+ Kf7 54. Rc8 Rf1+ 55. Kg4 Ne5+ 56. Kh4 Ng6+ 57. Kh5 Rf5+ 58. Kg4 Kf6 59. h4 Ne5+ 60. Kh3 Rf3 61. Rb8 Rxa3 62. Rxb5 c3 63. Rc5 Ra2 64. h5 c2 65. h6 Kg6 66. Kh4 Ra4+ 67. g4 Nf3+ 68. Kg3 Nd4 69. g5 Ra3+ 70. Kf4 Ne2+ 71. Kg4 c1=Q 72. Rxc1 Nxc1 73. Kf4 Ra5 74. Kg4 Rxg5+ 75. Kf4 Kxh6 76. Ke4 Kg6 77. Ke3 Rf5 78. Ke4 Kg5 79. Kd4 Rf4+ 80. Ke3 1/2-1/2`)
	game := &Game{}
	if game.UnmarshalText(pgn) != nil {
		b.Fatalf("pgn not parsed")
	}
	for b.Loop() {
		game.CanClaimDrawThreeFold()
	}
}

// This fuzz test is quite slow.
func FuzzParsePgn(f *testing.F) {
	file, err := os.Open("./testdata/SaintLouis2023-first3.pgn")
	if err != nil {
		f.Fatalf("issue reading test file at \"./testdata/SaintLouis2023-first3.pgn\"")
	}
	defer file.Close()

	pgn, err := io.ReadAll(file)
	if err != nil {
		f.Fatalf("issue reading test file at \"./testdata/SaintLouis2023-first3.pgn\"")
	}

	f.Add(pgn)
	f.Fuzz(func(t *testing.T, pgn []byte) {
		pgnReader := bytes.NewReader(pgn)
		ParsePGN(pgnReader)
		// Just make sure it doesn't panic.
	})
}
func FuzzGameUnmarshal(f *testing.F) {
	inputs := []string{"[]\n\n0",
		"[a]\n\n0",
		"[ \"]\n\n0",
		"%0\n\n(", "%0\n\n{}",
		"[]\n\n}{ s y . w I r",
		"%\n\nAA1x0 0-1"}

	for _, s := range inputs {
		f.Add([]byte(s))
	}

	f.Add([]byte(`[Event "idc"]
[Site "ur mom's house"]
[Date "2025.04.09"]
[Round "2"]
[White "phil"]
[Black "donna"]
[Result "0-1"]
[WhiteElo "1090"]

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}
1. d4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7.
Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.}
0-1`))
	f.Fuzz(func(t *testing.T, pgn []byte) {
		// Just make sure it doesn't panic
		g := NewGame()
		g.UnmarshalText(pgn)
	})
}

func FuzzGameUnmarshal_altStart(f *testing.F) {
	f.Add([]byte(`[Event "Hamguy123's Study: Chapter 1"]
[Result "*"]
[Variant "From Position"]
[FEN "rnbqk2r/pppp1ppp/4pn2/8/1bPP4/2N5/PPQ1PPPP/R1B1KBNR b KQkq - 0 1"]
[ECO "?"]
[Opening "?"]
[StudyName "Hamguy123's Study"]
[ChapterName "Chapter 1"]
[SetUp "1"]
[UTCDate "2025.05.22"]
[UTCTime "01:06:42"]
[Annotator "https://lichess.org/@/Hamguy123"]
[ChapterURL "https://lichess.org/study/lfAbjJ1u/fSzsIBlu"]

1... Nd5 2. e3 Qf6 3. Bd2 O-O 4. O-O-O Nc6 5. Be2 b6 6. Nh3 Rb8 7. Rhf1 *
`))
	f.Fuzz(func(t *testing.T, pgn []byte) {
		// Just make sure it doesn't panic
		g := NewGame()
		g.UnmarshalText(pgn)
	})
}
