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

import "testing"

func TestParseFEN_PositionString(t *testing.T) {
	pos := &Position{}
	err := pos.UnmarshalText([]byte(DefaultFEN))
	fen := pos.String()
	if fen != DefaultFEN {
		t.Errorf("incorrect result: expected %q, got %q", DefaultFEN, fen)
	}
	if err != nil {
		t.Errorf("incorrect result for default fen: expected error to be nil.")
	}

	input := "1k1r3r/ppq2Rp1/2p1p1p1/4N1b1/Q2P4/2P4P/PP6/R3K3 b Q - 0 23"
	pos = &Position{}
	err = pos.UnmarshalText([]byte(input))
	fen = pos.String()
	if fen != input {
		t.Errorf("incorrect result: expected %q, got %q", input, fen)
	}
	if err != nil {
		t.Errorf("incorrect result for random fen: expected error to be nil.")
	}
}

func TestParseFENError(t *testing.T) {
	input := "1k1r3r/pq2Rp1/2p1p1p1/4N1b1/Q2P4/2P4P/PP6/R3K3 b Q - 0 23"
	pos := &Position{}
	err := pos.UnmarshalText([]byte(input))
	if err == nil {
		t.Errorf("incorrect result for missing piece: expected error")
	}

	input = "1k1r3r/ppq2Rp1/2p1p1p1/3N1b1/Q2P4/2P4P/PP6/R3K3 b Q - 0 23"
	pos = &Position{}
	err = pos.UnmarshalText([]byte(input))
	if err == nil {
		t.Errorf("incorrect result for wrong number in FEN: expected error")
	}

	input = "1k1r3r/ppq2Rp1/2p1p1p1/4N1b1/Q2P4/2P4P/PP6/R3K3 b Q - 0 "
	pos = &Position{}
	err = pos.UnmarshalText([]byte(input))
	if err == nil {
		t.Errorf("incorrect result: expected error due to missing field")
	}
}

func TestParseFEN_EmptyFields(t *testing.T) {
	input := "rn1qk2r/pbppppbp/1p3np1/8/4P3/3P1NP1/PPP2PBP/RNBQ1RK1 b - - 24 6"
	pos := &Position{}
	err := pos.UnmarshalText([]byte(input))
	fen := pos.String()
	if err != nil {
		t.Errorf("incorrect result: expected error to be nil")
	}
	if fen != input {
		t.Errorf("incorrect result: expected %q, got %q", input, fen)
	}
}

func TestPositionPrettyString_Default(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	expected := `8rnbqkbnr
7pppppppp
6--------
5--------
4--------
3--------
2PPPPPPPP
1RNBQKBNR
 ABCDEFGH

Side To Move: White
Castle Rights: KQkq
En Passant Square: -
Half Move: 0
Full Move: 1`
	actual := pos.PrettyString(true, true)

	if actual != expected {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected, actual)
	}

	actual = pos.PrettyString(false, true)
	expected = `1RNBKQBNR
2PPPPPPPP
3--------
4--------
5--------
6--------
7pppppppp
8rnbkqbnr
 HGFEDCBA

Side To Move: White
Castle Rights: KQkq
En Passant Square: -
Half Move: 0
Full Move: 1`

	if actual != expected {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected, actual)
	}

	actual = pos.PrettyString(false, false)
	expected = `1RNBKQBNR
2PPPPPPPP
3--------
4--------
5--------
6--------
7pppppppp
8rnbkqbnr
 HGFEDCBA`

	if actual != expected {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected, actual)
	}
}

func TestPositionPrettyString_NotDefault(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rn1qk2r/pbppppbp/1p3np1/8/4P3/3P1NP1/PPP2PBP/RNBQ1RK1 b kq e3 0 6"))
	expected := `8rn-qk--r
7pbppppbp
6-p---np-
5--------
4----P---
3---P-NP-
2PPP--PBP
1RNBQ-RK-
 ABCDEFGH

Side To Move: Black
Castle Rights: kq
En Passant Square: E3
Half Move: 0
Full Move: 6`
	actual := pos.PrettyString(true, true)

	if actual != expected {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected, actual)
	}
}

func TestPositionPrettyString_NotDefault_EmptyFields(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rn1qk2r/pbppppbp/1p3np1/8/4P3/3P1NP1/PPP2PBP/RNBQ1RK1 b - - 24 6"))
	expected := `1-KR-QBNR
2PBP--PPP
3-PN-P---
4---P----
5--------
6-pn---p-
7pbppppbp
8r--kq-nr
 HGFEDCBA

Side To Move: Black
Castle Rights: -
En Passant Square: -
Half Move: 24
Full Move: 6`
	actual := pos.PrettyString(false, true)

	if actual != expected {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected, actual)
	}
}

func TestIsCheck(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	if pos.IsCheck() {
		t.Error("incorrect result: default position: expected false, got true")
	}

	fen := "r2q1n1k/5Qb1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R w - - 0 24"
	pos = &Position{}
	pos.UnmarshalText([]byte(fen))
	if pos.IsCheck() {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "r2q1nQk/6b1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R b - - 1 24"
	pos = &Position{}
	pos.UnmarshalText([]byte(fen))
	if !pos.IsCheck() {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "rnbq2nr/ppp1b1kN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R4RK1 w - - 6 26"
	pos = &Position{}
	pos.UnmarshalText([]byte(fen))
	if pos.IsCheck() {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "rnbq2nr/ppp1bRkN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R5K1 b - - 7 26"

	pos.UnmarshalText([]byte(fen))
	if !pos.IsCheck() {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}
}
func TestIsCheckPawn(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.SetPiece(BlackPawn, F2)
	if !pos.IsCheck() {
		t.Error("incorrect result for black pawn on f2: expected true, got false")
	}

	pos.SetPiece(NoPiece, F2)
	pos.SetPiece(BlackPawn, E2)
	if pos.IsCheck() {
		t.Error("incorrect result for black pawn on e2: expected false, got true")
	}

	pos.SetPiece(NoPiece, E2)
	pos.SetPiece(BlackPawn, D2)
	if !pos.IsCheck() {
		t.Error("incorrect result for black pawn on d2: expected true, got false")
	}

	pos.SetPiece(NoPiece, D2)
	pos.SetPiece(WhitePawn, F7)
	pos.SideToMove = Black
	if !pos.IsCheck() {
		t.Error("incorrect result for white pawn on f7: expected true, got false")
	}
}

func TestIsCheckRook(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.SetPiece(BlackRook, E8)
	pos.SetPiece(NoPiece, E2)
	if pos.IsCheck() {
		t.Error("incorrect result for black rock on e8: blocked by own piece: expected false, got true")
	}

	pos.SetPiece(WhitePawn, E2)
	pos.SetPiece(NoPiece, E7)
	if pos.IsCheck() {
		t.Error("incorrect result for black rock on e8: blocked by opponent piece: expected false, got true")
	}

	pos.SetPiece(NoPiece, E2)
	if !pos.IsCheck() {
		t.Error("incorrect result for black rook on e8: expected true, got false")
	}

	pos.SetPiece(BlackKing, E8)
	pos.SetPiece(BlackRook, D1)
	if !pos.IsCheck() {
		t.Error("incorrect result for black rook on d1: expected true, got false")
	}

	pos.SetPiece(WhiteRook, E1)
	pos.SideToMove = Black
	if !pos.IsCheck() {
		t.Error("incorrect result for white rook on e1: expected true, got false")
	}
}

func TestIsCheckKnight(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.SetPiece(BlackKnight, D3)
	if !pos.IsCheck() {
		t.Error("incorrect result for black knight: expected true, got false")
	}

	pos.SetPiece(NoPiece, D3)
	pos.SetPiece(WhiteKnight, D6)
	if pos.IsCheck() {
		t.Error("incorrect result for white knight: was still whites turn: expected false, got true")
	}

	pos.SideToMove = Black
	if !pos.IsCheck() {
		t.Error("incorrect result for white knight: expected true, got false")
	}
}

func TestIsCheckBishop(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black
	pos.SetPiece(BlackKing, A1)
	pos.SetPiece(WhiteBishop, H8)
	if !pos.IsCheck() {
		t.Error("incorrect result for black bishop on H8: expected true, got false")
	}

	pos.SetPiece(BlackQueen, E5)
	if pos.IsCheck() {
		t.Error("incorrect result for black bishop on H8: blocked by own queen: expected false, got true")
	}

	pos = &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKing, A8)
	pos.SetPiece(BlackBishop, H1)
	if !pos.IsCheck() {
		t.Error("incorrect result for white bishop on H1: expected true, got false")
	}
}

func TestIsCheckQueen(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKing, D4)
	pos.SetPiece(BlackQueen, B6)
	if !pos.IsCheck() {
		t.Error("incorrect result for black queen on diagonal: expected true, got false")
	}

	pos.SideToMove = Black
	pos.SetPiece(BlackKing, D4)
	pos.SetPiece(WhiteQueen, B6)
	if !pos.IsCheck() {
		t.Error("incorrect result for white queen on diagonal: expected true, got false")
	}

	pos = &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKing, D4)
	pos.SetPiece(BlackQueen, D2)
	if !pos.IsCheck() {
		t.Error("incorrect result for black queen on vertical: expected true, got false")
	}

	pos.SideToMove = Black
	pos.SetPiece(BlackKing, D4)
	pos.SetPiece(WhiteQueen, D2)
	if !pos.IsCheck() {
		t.Error("incorrect result for white queen on vertical: expected true, got false")
	}

	pos.SetPiece(NoPiece, D2)
	pos.SetPiece(WhiteQueen, C2)
	if pos.IsCheck() {
		t.Error("incorrect result for white queen on horse jump: expected false, got true")
	}
}

func TestIsCheckKing(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKing, E4)
	pos.SetPiece(BlackKing, D4)
	if !pos.IsCheck() {
		t.Error("incorrect result for black king on d4 horizontal: expected true, got false")
	}

	pos.SideToMove = Black
	if !pos.IsCheck() {
		t.Error("incorrect result for white king on e4 horizontal: expected true, got false")
	}

	pos.SetPiece(NoPiece, E4)
	pos.SetPiece(WhiteKing, F4)
	if pos.IsCheck() {
		t.Error("incorrect result for white king two spaces away: expected false, got true")
	}
}

func TestRandomMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.Move(Move{B1, D5, NoPieceType})
	expected := "rnbqkbnr/pppppppp/8/3N4/8/8/PPPPPPPP/R1BQKBNR b KQkq - 1 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(Move{E8, G8, NoPieceType})
	expected = "rnbq1bkr/pppppppp/8/3N4/8/8/PPPPPPPP/R1BQKBNR w KQ - 0 2"
	actual = pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestWhiteEnPassantMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/ppppp1pp/8/8/5p2/8/PPPPPPPP/RNBQKBNR w KQkq - 53 1"))
	pos.Move(Move{E2, E4, NoPieceType})
	expected := "rnbqkbnr/ppppp1pp/8/8/4Pp2/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(Move{F4, E3, NoPieceType})
	expected = "rnbqkbnr/ppppp1pp/8/8/8/4p3/PPPP1PPP/RNBQKBNR w KQkq - 0 2"
	actual = pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestBlackEnPassantMove(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 9 1"))
	pos.Move(Move{F7, F5, NoPieceType})
	expected := "rnbqkbnr/ppppp1pp/8/4Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 2"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	pos.Move(Move{E5, F6, NoPieceType})
	expected = "rnbqkbnr/ppppp1pp/5P2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2"
	actual = pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestWhiteKSCastle(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1"))
	pos.Move(Move{E1, G1, NoPieceType})
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 1 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestWhiteQSCastle(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3KBNR w KQkq - 0 1"))
	pos.Move(Move{E1, C1, NoPieceType})
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/2KR1BNR b kq - 1 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestBlackKSCastle(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqk2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"))
	pos.Move(Move{E8, G8, NoPieceType})
	expected := "rnbq1rk1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 1 2"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestBlackQSCastle(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"))
	pos.Move(Move{E8, C8, NoPieceType})
	expected := "2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 1 2"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestPromotion(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.Move(Move{E1, F3, Knight})
	expected := "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQ1BNR b kq - 1 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestCastleRightsRemoved(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	pos.Move(Move{H1, H4, NoPieceType})
	expected := "rnbqkbnr/pppppppp/8/8/7R/8/PPPPPPPP/RNBQKBN1 b Qkq - 1 1"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestEnPassantRemoved(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"))
	pos.Move(Move{F8, B3, NoPieceType})
	expected := "rnbqk1nr/pppppppp/8/8/4P3/1b6/PPPP1PPP/RNBQKBNR w KQkq - 1 2"
	actual := pos.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestPositionMarshalers(t *testing.T) {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	fen, _ := pos.MarshalText()
	if string(fen) != DefaultFEN {
		t.Errorf("position marshalers not working")
	}
}

func FuzzPositionUnmarshal(f *testing.F) {
	f.Add(DefaultFEN)

	f.Fuzz(func(t *testing.T, fen string) {
		pos := &Position{}
		pos.UnmarshalText([]byte(fen))
		// Make sure it doesn't panic
	})
}
