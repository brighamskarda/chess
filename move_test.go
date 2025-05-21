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
	"testing"
)

func TestMoveString(t *testing.T) {
	expected := "a1b2"
	actual := Move{A1, B2, NoPieceType}.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
	expected = "h2c1q"
	actual = Move{H2, C1, Queen}.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_basicPawnMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	m := Move{E2, E4, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "e4"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{A7, A6, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "a6"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_pawnCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/ppp1p1pp/8/3p4/4PpP1/8/PPPP1P1P/RNBQKBNR b KQkq g3 0 1"))
	m := Move{F4, G3, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "fxg3"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{E4, D5, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "exd5"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_pawnPromotion(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pPpppp1p/8/8/8/8/PPPPP1pP/RNBQKB1R w KQkq - 0 1"))
	m := Move{B7, A8, Knight}
	actual := m.StringSAN(pos)
	expected := "bxa8=N"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{G2, G1, Queen}
	actual = m.StringSAN(pos)
	expected = "g1=Q"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_BasicMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/1Pp1pp1p/8/8/8/4P3/PPPP2pP/RNBQKB1R w KQkq - 0 1"))
	m := Move{F1, D3, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Bd3"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{D8, D4, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "Qd4"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_BasicCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/1Pp1pp1p/8/8/8/4P3/PPPP2pP/RNBQKB1R w KQkq - 0 1"))
	m := Move{F1, G2, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Bxg2"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{A8, A2, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "Rxa2"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_FileDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/1Pp1pp1p/Pr6/3N4/8/4P3/1PPP2pP/RNBQKB1R w KQkq - 0 1"))
	m := Move{D5, C3, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Ndc3"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{B6, A6, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "Rbxa6"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_RankDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/1Pp1pp1p/P7/8/r2N4/4P3/1PPN2pP/R1BQKB1R w KQkq - 0 1"))
	m := Move{D4, B3, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "N4b3"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{A4, A6, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "R4xa6"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_SquareDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("qnb1kbnr/1Pp1pp1p/P1q5/8/qBqB4/4P3/1BPB2pP/R1BQKB1R w KQk - 0 1"))
	m := Move{D2, C3, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Bd2c3"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(m)
	m = Move{A4, A6, NoPieceType}
	actual = m.StringSAN(pos)
	expected = "Qa4xa6"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

// No disambiguation needed due to check.
func TestMoveStringSAN_NoDisambiguation(t *testing.T) {
	// This test is a great example of the stupidity of SAN
	pos := &Position{}
	pos.UnmarshalText([]byte("3k4/3n1n2/8/8/8/8/3R4/3K4 b - - 0 1"))
	m := Move{F7, F5, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Nf5"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_CheckSymbol(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("3k4/8/8/8/8/3b4/8/3K4 b - - 0 1"))
	m := Move{D3, E2, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Be2+"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestMoveStringSAN_CheckmateSymbol(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("8/8/8/3p4/8/1K2PN2/p3Q3/7k w - - 0 74"))
	m := Move{E2, H2, NoPieceType}
	actual := m.StringSAN(pos)
	expected := "Qh2#"
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestParseUCIMove(t *testing.T) {
	expected := Move{A1, A2, NoPieceType}
	actual, err := ParseUCIMove("a1a2")
	if expected != actual {
		t.Errorf("incorrect result: expected %v, got %v", expected, actual)
	}
	if err != nil {
		t.Errorf("incorrect result for \"a1a2\": expected err to be nil")
	}

	expected = Move{H2, C1, Queen}
	actual, err = ParseUCIMove("h2c1q")
	if expected != actual {
		t.Errorf("incorrect result: expected %v, got %v", expected, actual)
	}
	if err != nil {
		t.Errorf("incorrect result for \"h2c1q\": expected err to be nil")
	}
}

func TestParseUCIMoveErr(t *testing.T) {
	_, err := ParseUCIMove("a1c")
	if err == nil {
		t.Error("Expected err to be nil.")
	}
}

func TestParseSANMove_PawnMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	move, err := ParseSANMove("e4", pos)
	if err != nil {
		t.Errorf("got error on e4")
	}
	expected := Move{E2, E4, NoPieceType}
	if move != expected {
		t.Errorf("got error on e4: expected %v, got %v", expected, move)
	}

	pos.Move(move)
	move, err = ParseSANMove("A6", pos)
	if err != nil {
		t.Errorf("got error on A6")
	}
	expected = Move{A7, A6, NoPieceType}
	if move != expected {
		t.Errorf("got error on A6: expected %v, got %v", expected, move)
	}
}

func TestParseSANMove_PawnPromotion(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("8/1P2k3/8/8/8/8/4K1p1/8 w - - 0 1"))
	move, err := ParseSANMove("b8=q", pos)
	if err != nil {
		t.Errorf("got error on b8=q")
	}
	expected := Move{B7, B8, Queen}
	if move != expected {
		t.Errorf("got error on b8=q: expected %v, got %v", expected, move)
	}

	pos.Move(move)
	move, err = ParseSANMove("G1=N", pos)
	if err != nil {
		t.Errorf("got error on G1=N")
	}
	expected = Move{G2, G1, Knight}
	if move != expected {
		t.Errorf("got error on G1=N: expected %v, got %v", expected, move)
	}
}

func TestParseSANMove_PawnCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("8/4k3/2r5/1P6/5Pp1/8/3K4/8 b - f3 0 1"))
	s := "gxf3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{G4, F3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "bxc6"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{B5, C6, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_Castle(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k3/8/8/8/8/8/8/4K2R w Kq - 0 1"))
	s := "O-O"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{E1, G1, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "O-O-O"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{E8, C8, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_Basic(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/8/8/8/8/8/8/2N1K2R w Kq - 0 1"))
	s := "nD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "Bh7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, H7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_BasicCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/5R2/8/8/8/3r4/8/2N1K2R w Kq - 0 1"))
	s := "nxD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BXf7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_FileDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/8/4b3/8/8/8/8/2N1NK1R w Kq - 0 1"))
	s := "ncD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BGf7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_FileDisambiguationCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/5R2/4b3/8/8/3r4/8/2N1NK1R w Kq - 0 1"))
	s := "ncxD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BGxf7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_RankDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/8/6b1/2N5/8/8/8/2N2K1R w Kq - 0 1"))
	s := "n1D3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "B8f7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_RankDisambiguationCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/5R2/6b1/2N5/8/3p4/8/2N2K1R w Kq - 0 1"))
	s := "n1xD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "B8Xf7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_SquareDisambiguation(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/8/4b1b1/2N5/8/8/8/2N1NK1R w Kq - 0 1"))
	s := "nC1D3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BG8f7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_SquareDisambiguationCapture(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3k1b1/5R2/4b1b1/2N5/8/3r4/8/2N1NK1R w Kq - 0 1"))
	s := "nC1xD3"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BG8xf7"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_SquareDisambiguationCaptureCheck(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("4K1b1/5R2/4b1b1/2N1k3/8/3r4/8/2N1N3 w - - 0 1"))
	s := "nC1xD3+"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BG8xf7+"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_SquareDisambiguationCaptureCheckmate(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("4K1b1/5R2/4b1b1/2N1k3/8/3r4/8/2N1N3 w - - 0 1"))
	s := "nC1xD3#"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C1, D3, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	pos.Move(move)
	s = "BG8xf7#"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{G8, F7, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_DisambiguateLetterB(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("7k/1p1b4/2N5/8/8/8/8/7K b - - 0 1"))
	s := "bxc6"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{B7, C6, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}

	s = "Bxc6"
	move, err = ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected = Move{D7, C6, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_DisambiguateMoveThatResultsInCheck(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("k3r3/8/8/8/8/8/2N1N3/4K3 w - - 0 1"))
	s := "Nd4"
	move, err := ParseSANMove(s, pos)
	if err != nil {
		t.Errorf("got error on %s", s)
	}
	expected := Move{C2, D4, NoPieceType}
	if move != expected {
		t.Errorf("got error on %s: expected %v, got %v", s, expected, move)
	}
}

func TestParseSANMove_miscInput(t *testing.T) {

	p := Position{}
	p.UnmarshalText([]byte("rnbqkbnr/pppp1pp1/8/8/8/8/1PPP1PPP/RNBQKBNR w KQkq - 0 1"))

	inputs := []string{
		"RAX0", "Bꀀ0", "\U0006bac10", "P2XA1",
	}

	for _, san := range inputs {
		// Just make sure it doesn't panic
		ParseSANMove(san, p.Copy())
	}
}

func FuzzParseSANMove(f *testing.F) {
	legalMoves := []string{
		"Ra3", "Ra4", "b3", "b4", "c3", "c4", "d3", "d4",
		"Ke2", "Qf3", "Bc4", "f4", "g3", "g4", "h3", "h4",
		"Na3", "Nc3", "Nf3", "Nh3",
	}

	for _, m := range legalMoves {
		f.Add(m)
	}

	p := Position{}
	p.UnmarshalText([]byte("rnbqkbnr/pppp1pp1/8/8/8/8/1PPP1PPP/RNBQKBNR w KQkq - 0 1"))
	f.Fuzz(func(t *testing.T, san string) {
		// Just make sure it doesn't panic
		ParseSANMove(san, p.Copy())
	})
}
