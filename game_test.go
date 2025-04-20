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
		Commentary:        "my comment",
		Variations: [][]PgnMove{
			[]PgnMove{
				PgnMove{
					Move:              Move{},
					NumericAnnotation: 0,
					Commentary:        "",
					Variations: [][]PgnMove{
						[]PgnMove{
							PgnMove{
								Move:              Move{},
								NumericAnnotation: 0,
								Commentary:        "",
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
					Commentary:        "",
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
