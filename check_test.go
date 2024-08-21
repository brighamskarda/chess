package chess

import (
	"testing"
)

func TestIsCheck(t *testing.T) {
	pos := getDefaultPosition()
	if IsCheck(&pos) {
		t.Error("incorrect result: default position: expected false, got true")
	}

	fen := "r2q1n1k/5Qb1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R w - - 0 24"
	pos, _ = ParseFen(fen)
	if IsCheck(&pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "r2q1nQk/6b1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R b - - 1 24"
	pos, _ = ParseFen(fen)
	if !IsCheck(&pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "rnbq2nr/ppp1b1kN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R4RK1 w - - 6 26"
	pos, _ = ParseFen(fen)
	if IsCheck(&pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "rnbq2nr/ppp1bRkN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R5K1 b - - 7 26"
	pos, _ = ParseFen(fen)
	if !IsCheck(&pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}
}

func TestIsCheckPawn(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(F2, BlackPawn)
	if !IsCheck(&pos) {
		t.Error("incorrect result for black pawn on f2: expected true, got false")
	}

	pos.SetPieceAt(F2, NoPiece)
	pos.SetPieceAt(E2, BlackPawn)
	if IsCheck(&pos) {
		t.Error("incorrect result for black pawn on e2: expected false, got true")
	}

	pos.SetPieceAt(E2, NoPiece)
	pos.SetPieceAt(D2, BlackPawn)
	if !IsCheck(&pos) {
		t.Error("incorrect result for black pawn on d2: expected true, got false")
	}

	pos.SetPieceAt(D2, NoPiece)
	pos.SetPieceAt(F7, WhitePawn)
	pos.Turn = Black
	if !IsCheck(&pos) {
		t.Error("incorrect result for white pawn on f7: expected true, got false")
	}
}

func TestIsCheckRook(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(E8, BlackRook)
	pos.SetPieceAt(E2, NoPiece)
	if IsCheck(&pos) {
		t.Error("incorrect result for black rock on e8: blocked by own piece: expected false, got true")
	}

	pos.SetPieceAt(E2, WhitePawn)
	pos.SetPieceAt(E7, NoPiece)
	if IsCheck(&pos) {
		t.Error("incorrect result for black rock on e8: blocked by opponent piece: expected false, got true")
	}

	pos.SetPieceAt(E2, NoPiece)
	if !IsCheck(&pos) {
		t.Error("incorrect result for black rook on e8: expected true, got false")
	}

	pos.SetPieceAt(E8, BlackKing)
	pos.SetPieceAt(D1, BlackRook)
	if !IsCheck(&pos) {
		t.Error("incorrect result for black rook on d1: expected true, got false")
	}

	pos.SetPieceAt(E1, WhiteRook)
	pos.Turn = Black
	if !IsCheck(&pos) {
		t.Error("incorrect result for white rook on e1: expected true, got false")
	}
}

func TestIsCheckKnight(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(D3, BlackKnight)
	if !IsCheck(&pos) {
		t.Error("incorrect result for black knight: expected true, got false")
	}

	pos.SetPieceAt(D3, NoPiece)
	pos.SetPieceAt(D6, WhiteKnight)
	if IsCheck(&pos) {
		t.Error("incorrect result for white knight: was still whites turn: expected false, got true")
	}

	pos.Turn = Black
	if !IsCheck(&pos) {
		t.Error("incorrect result for white knight: expected true, got false")
	}
}

func TestIsCheckBishop(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(A1, BlackKing)
	pos.SetPieceAt(H8, WhiteBishop)
	if !IsCheck(pos) {
		t.Error("incorrect result for black bishop on H8: expected true, got false")
	}

	pos.SetPieceAt(E6, BlackQueen)
	if IsCheck(pos) {
		t.Error("incorrect result for black bishop on H8: blocked by own queen: expected false, got true")
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(A8, WhiteKing)
	pos.SetPieceAt(H1, BlackBishop)
	if !IsCheck(pos) {
		t.Error("incorrect result for white bishop on H1: expected true, got false")
	}
}

func TestIsCheckQueen(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(D4, WhiteKing)
	pos.SetPieceAt(B6, BlackQueen)
	if !IsCheck(pos) {
		t.Error("incorrect result for black queen on diagonal: expected true, got false")
	}

	pos.Turn = Black
	pos.SetPieceAt(D4, BlackKing)
	pos.SetPieceAt(B6, WhiteQueen)
	if !IsCheck(pos) {
		t.Error("incorrect result for white queen on diagonal: expected true, got false")
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(D4, WhiteKing)
	pos.SetPieceAt(D2, BlackQueen)
	if !IsCheck(pos) {
		t.Error("incorrect result for black queen on vertical: expected true, got false")
	}

	pos.Turn = Black
	pos.SetPieceAt(D4, BlackKing)
	pos.SetPieceAt(D2, WhiteQueen)
	if !IsCheck(pos) {
		t.Error("incorrect result for white queen on vertical: expected true, got false")
	}

	pos.SetPieceAt(D2, NoPiece)
	pos.SetPieceAt(C2, WhiteQueen)
	if IsCheck(pos) {
		t.Error("incorrect result for white queen on horse jump: expected false, got true")
	}
}

func TestIsCheckKing(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(E4, WhiteKing)
	pos.SetPieceAt(D4, BlackKing)
	if !IsCheck(pos) {
		t.Error("incorrect result for black king on d4 horizontal: expected true, got false")
	}

	pos.Turn = Black
	if !IsCheck(pos) {
		t.Error("incorrect result for white king on e4 horizontal: expected true, got false")
	}

	pos.SetPieceAt(E4, NoPiece)
	pos.SetPieceAt(F4, WhiteKing)
	if IsCheck(pos) {
		t.Error("incorrect result for white king two spaces away: expected false, got true")
	}
}
