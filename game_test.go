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
	"fmt"
	"maps"
	"slices"
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
	if s := g.Position().String(); s != fen {
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
		Commentary:        []string{"my comment"},
		Variations: [][]PgnMove{
			[]PgnMove{
				PgnMove{
					Move:              Move{},
					NumericAnnotation: 0,
					Commentary:        []string{},
					Variations: [][]PgnMove{
						[]PgnMove{
							PgnMove{
								Move:              Move{},
								NumericAnnotation: 0,
								Commentary:        []string{},
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
					Commentary:        []string{},
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
	if g.Position().String() != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("incorrect position, got %s", g.Position().String())
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
	if g.Position().String() != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("incorrect position, got %s", g.Position().String())
	}
	if g.Result != NoResult {
		t.Errorf("incorrect result, got %v", g.Result)
	}
}

func TestAnnotateMove(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fail()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fail()
	}

	g.AnnotateMove(0, 3)
	g.AnnotateMove(1, 2)
	moveHistory := g.MoveHistory()
	if moveHistory[0].NumericAnnotation != 3 {
		t.Errorf("for move 0 expected 3, got %d", moveHistory[0].NumericAnnotation)
	}
	if moveHistory[1].NumericAnnotation != 2 {
		t.Errorf("for move 1 expected 23, got %d", moveHistory[1].NumericAnnotation)
	}
}

func TestMoveSAN(t *testing.T) {
	g := NewGame()
	if g.MoveSAN("e4") != nil {
		t.Fail()
	}
	if g.MoveSAN("Nc6") != nil {
		t.Fail()
	}

	if g.Position().String() != "r1bqkbnr/pppppppp/2n5/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 1 2" {
		t.Errorf("MoveSAN did not work correctly, ending position was %q", g.Position().String())
	}
}

func TestPositionPly(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fail()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fail()
	}

	ply := 0
	expected := DefaultFEN
	actual := g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 1
	expected = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	actual = g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 2
	expected = "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"
	actual = g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}
}

func TestPositionPly_AltStart(t *testing.T) {
	g, err := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN1 w Qkq - 0 1")
	if err != nil {
		t.Fail()
	}
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fail()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fail()
	}

	ply := 0
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN1 w Qkq - 0 1"
	actual := g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 1
	expected = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBN1 b Qkq e3 0 1"
	actual = g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}

	ply = 2
	expected = "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBN1 w Qkq d6 0 2"
	actual = g.PositionPly(ply).String()
	if expected != actual {
		t.Errorf("incorrect position for ply %d: expected %q, got %q", ply, expected, actual)
	}
}

func TestCommentMove(t *testing.T) {
	g := NewGame()
	if g.Move(Move{E2, E4, NoPieceType}) != nil {
		t.Fail()
	}
	if g.Move(Move{D7, D5, NoPieceType}) != nil {
		t.Fail()
	}

	g.CommentMove(0, "comment 1")
	g.CommentMove(1, "comment 2")
	moveHistory := g.MoveHistory()
	if moveHistory[0].Commentary[0] != "comment 1" {
		t.Errorf("for move 0 got %q", moveHistory[0].Commentary)
	}
	if moveHistory[1].Commentary[0] != "comment 2" {
		t.Errorf("for move 1 got %q", moveHistory[1].Commentary)
	}
}

func makeGameWithVariation() *Game {
	g := NewGame()
	g.Move(Move{E2, E4, NoPieceType})
	g.Move(Move{D7, D5, NoPieceType})

	g.MakeVariation(1, []PgnMove{{
		Move:              Move{B8, C6, NoPieceType},
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	}})
	return g
}
func TestMakeVariation(t *testing.T) {
	g := makeGameWithVariation()
	moveHistory := g.MoveHistory()
	if len(moveHistory[1].Variations) != 1 {
		t.Errorf("variation not added correctly")
	}
	if len(moveHistory[1].Variations[0]) != 2 {
		t.Errorf("variation not added correctly")
	}
}

func TestDeleteVariation(t *testing.T) {
	g := makeGameWithVariation()
	g.DeleteVariation(1, 0)
	moveHistory := g.MoveHistory()
	if len(moveHistory[1].Variations) != 0 {
		t.Errorf("variation not deleted correctly")
	}
}

func TestGetVariation(t *testing.T) {
	g := makeGameWithVariation()

	g.MakeVariation(1, []PgnMove{{
		Move:              Move{G8, F6, NoPieceType},
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	}, {
		Move:              Move{D2, D4, NoPieceType},
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	}})

	g.MakeVariation(0, []PgnMove{{
		Move:              Move{H2, H4, NoPieceType},
		NumericAnnotation: 0,
		Commentary:        []string{},
		Variations:        [][]PgnMove{},
	}})

	newG := g.GetVariation(1, 0)

	expectedPosition := "r1bqkbnr/pppppppp/2n5/8/3PP3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 2"
	actualPosition := newG.Position().String()
	if expectedPosition != actualPosition {
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

func TestGameString(t *testing.T) {
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

	g.Commentary = "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20"
	g.AnnotateMove(4, 2)
	g.AnnotateMove(5, 10)
	g.CommentMove(19, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				Commentary:        []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
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

{Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20}
1. d4 e6 2. e4 (2. f4 g5 {Another variation here} {another comment here} (2...
h5!) 3. h4) (2. a4) 2... d5 3. exd5? exd5 $10 4. Nf3 Nf6 5. Ne5 Qe7 6. f4 Bg4 7.
Be2 Bd7 8. g4 Ne4 9. c4 Qh4+ 10. Kf1 Qf2#
{Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.}
0-1`
	actual := g.String()
	if actual != expected {
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

func TestGameString_AltStart(t *testing.T) {
	g, _ := NewGameFromFEN("r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16")
	moves := []string{
		"Qd3",
		"Qg4+",
		"Kf7",
	}

	for _, m := range moves {
		if g.MoveSAN(m) != nil {
			t.Fail()
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
	actual := g.String()
	if actual != expected {
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

func TestGameString_NumericAnnotationGlyphs(t *testing.T) {
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
	actual := g.String()
	if actual != expected {
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

func TestGameReducedString(t *testing.T) {
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

	g.Commentary = "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20"
	g.AnnotateMove(4, 2)
	g.AnnotateMove(5, 10)
	g.CommentMove(19, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				Commentary:        []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(2, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
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
	actual := g.ReducedString()
	if actual != expected {
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

func TestGameReducedString_AltStart(t *testing.T) {
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

	g.Commentary = "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20"
	g.AnnotateMove(3, 2)
	g.AnnotateMove(4, 10)
	g.CommentMove(18, "Black wins by checkmate. Now I need this comment to be even longer than before, preferably longer than 80 characters for some testing.")
	g.MakeVariation(1, []PgnMove{
		PgnMove{
			Move:              Move{F2, F4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
		PgnMove{
			Move:              Move{G7, G5, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{"Another variation here", "another comment here"},
			Variations: [][]PgnMove{[]PgnMove{PgnMove{
				Move:              Move{H7, H5, NoPieceType},
				NumericAnnotation: 1,
				Commentary:        []string{},
				Variations:        [][]PgnMove{},
			}}},
		},
		PgnMove{
			Move:              Move{H2, H4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
			Variations:        [][]PgnMove{},
		},
	})
	g.MakeVariation(1, []PgnMove{
		PgnMove{
			Move:              Move{A2, A4, NoPieceType},
			NumericAnnotation: 0,
			Commentary:        []string{},
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
	actual := g.ReducedString()
	if actual != expected {
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

	if g.Commentary != "Random game I found on Lichess.com, https://lichess.org/YF5EBq7m#20" {
		t.Errorf("commentary incorrect")
	}

	moveHis := g.MoveHistory()
	if len(moveHis) != 20 {
		t.Errorf("moveHistory incorrect length")
	}

	if len(moveHis[0].Variations) != 0 {
		t.Errorf("moveHistory has variation where none are present")
	}

	if len(moveHis[2].Variations) != 2 {
		t.Errorf("ply 2 does not have 2 variations")
	}

	if len(moveHis[2].Variations[0][1].Commentary) != 2 {
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
	if moveHis[11].Commentary[0] != "this is a comment" {
		t.Errorf("semicolon move comment not parsed")
	}
	if moveHis[19].Commentary[0] != `Black wins by checkmate. Now I need ;this comment to 
be even 
longer than before, preferably longer than 80% characters for some testing.` {
		t.Errorf("multiline commentary not parsed")
	}
	if _, ok := g.OtherTags["WhiteElo"]; ok != false {
		t.Errorf("dev escape not working")
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

	if g.PositionPly(0).String() != "r2q3r/ppp3pp/2n1Nnk1/4p3/2Q5/B7/P4PPP/RN3RK1 b - - 0 16" {
		t.Errorf("alt start position not set")
	}

	if g.PositionPly(1).String() != "r6r/ppp3pp/2n1Nnk1/4p3/2Q5/B2q4/P4PPP/RN3RK1 w - - 1 17" {
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
		g1.Commentary == g2.Commentary &&
		compareMoveHistories(g1.moveHistory, g2.moveHistory)
}

func compareMoveHistories(mh1 []PgnMove, mh2 []PgnMove) bool {
	if len(mh1) != len(mh2) {
		return false
	}
	for i := range mh1 {
		if !slices.Equal(mh1[i].Commentary, mh2[i].Commentary) ||
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
		t.Fail()
	}
}

func TestGameUnmarshal_miscInputs(t *testing.T) {
	inputs := []string{"[]\n\n0",
		"[a]\n\n0",
		"[ \"]\n\n0",
		"%0\n\n(", "%0\n\n{}",
		"[]\n\n}{ s y . w I r",
		"%\n\nAA1x0 0-1"}
	for _, s := range inputs {
		g := NewGame()
		g.UnmarshalText([]byte(s))
		// Just make sure it doesn't panic
	}
}

func FuzzGameUnmarshal(f *testing.F) {
	f.Add(`[Event "idc"]
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
0-1`)
	f.Fuzz(func(t *testing.T, pgn string) {
		// Just make sure it doesn't panic
		g := NewGame()
		g.UnmarshalText([]byte(pgn))
	})
}

func FuzzGameUnmarshal_altStart(f *testing.F) {
	f.Add(`[Event "Hamguy123's Study: Chapter 1"]
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
`)
	f.Fuzz(func(t *testing.T, pgn string) {
		// Just make sure it doesn't panic
		g := NewGame()
		g.UnmarshalText([]byte(pgn))
	})
}
