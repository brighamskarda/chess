package chess

import (
	"testing"
)

func TestIsCheck(t *testing.T) {
	pos := getDefaultPosition()
	if IsCheck(pos) {
		t.Error("incorrect result: default&Position: expected false, got true")
	}

	fen := "r2q1n1k/5Qb1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R w - - 0 24"
	pos, _ = ParseFen(fen)
	if IsCheck(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "r2q1nQk/6b1/p2pB2p/2pPp1pP/2pr2N1/5PP1/PP6/2KR3R b - - 1 24"
	pos, _ = ParseFen(fen)
	if !IsCheck(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "rnbq2nr/ppp1b1kN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R4RK1 w - - 6 26"
	pos, _ = ParseFen(fen)
	if IsCheck(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "rnbq2nr/ppp1bRkN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R5K1 b - - 7 26"
	pos, _ = ParseFen(fen)
	if !IsCheck(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}
}

func TestIsCheckPawn(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(F2, BlackPawn)
	if !IsCheck(pos) {
		t.Error("incorrect result for black pawn on f2: expected true, got false")
	}

	pos.SetPieceAt(F2, NoPiece)
	pos.SetPieceAt(E2, BlackPawn)
	if IsCheck(pos) {
		t.Error("incorrect result for black pawn on e2: expected false, got true")
	}

	pos.SetPieceAt(E2, NoPiece)
	pos.SetPieceAt(D2, BlackPawn)
	if !IsCheck(pos) {
		t.Error("incorrect result for black pawn on d2: expected true, got false")
	}

	pos.SetPieceAt(D2, NoPiece)
	pos.SetPieceAt(F7, WhitePawn)
	pos.Turn = Black
	if !IsCheck(pos) {
		t.Error("incorrect result for white pawn on f7: expected true, got false")
	}
}

func TestIsCheckRook(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(E8, BlackRook)
	pos.SetPieceAt(E2, NoPiece)
	if IsCheck(pos) {
		t.Error("incorrect result for black rock on e8: blocked by own piece: expected false, got true")
	}

	pos.SetPieceAt(E2, WhitePawn)
	pos.SetPieceAt(E7, NoPiece)
	if IsCheck(pos) {
		t.Error("incorrect result for black rock on e8: blocked by opponent piece: expected false, got true")
	}

	pos.SetPieceAt(E2, NoPiece)
	if !IsCheck(pos) {
		t.Error("incorrect result for black rook on e8: expected true, got false")
	}

	pos.SetPieceAt(E8, BlackKing)
	pos.SetPieceAt(D1, BlackRook)
	if !IsCheck(pos) {
		t.Error("incorrect result for black rook on d1: expected true, got false")
	}

	pos.SetPieceAt(E1, WhiteRook)
	pos.Turn = Black
	if !IsCheck(pos) {
		t.Error("incorrect result for white rook on e1: expected true, got false")
	}
}

func TestIsCheckKnight(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(D3, BlackKnight)
	if !IsCheck(pos) {
		t.Error("incorrect result for black knight: expected true, got false")
	}

	pos.SetPieceAt(D3, NoPiece)
	pos.SetPieceAt(D6, WhiteKnight)
	if IsCheck(pos) {
		t.Error("incorrect result for white knight: was still whites turn: expected false, got true")
	}

	pos.Turn = Black
	if !IsCheck(pos) {
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

	pos.SetPieceAt(E5, BlackQueen)
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

func BenchmarkIsCheck(b *testing.B) {
	pos := getDefaultPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsCheck(pos)
	}
}

func TestIsCheckMate(t *testing.T) {
	fen := "rnbq2nr/ppp1bRkN/4p1B1/3PP1Qp/2P5/6P1/PP4PP/R5K1 b - - 7 26"
	pos, _ := ParseFen(fen)
	if !IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "2B5/1p4k1/2p2pp1/2n1p1Kp/2P1P2N/4R1P1/p2r2P1/7r w - - 0 46"
	pos, _ = ParseFen(fen)
	if !IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "k3b3/8/7p/5Pp1/7K/r7/8/6r1 w - g6 0 1"
	pos, _ = ParseFen(fen)
	if IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "k3b3/8/7p/5Pp1/8/r7/8/6r1 w - g6 0 1"
	pos, _ = ParseFen(fen)
	if IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "k7/8/1R5p/R4Pp1/8/8/8/6r1 b - - 0 1"
	pos, _ = ParseFen(fen)
	if !IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "k1K5/ppp5/1bP5/8/8/8/8/3r4 w - - 0 1"
	pos, _ = ParseFen(fen)
	if IsCheckMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}
}

func BenchmarkIsCheckMate(b *testing.B) {
	pos, _ := ParseFen("3rkbnr/1p1bp3/1q1p3p/p5pQ/3n4/PPR5/5PPP/6K1 b - - 2 2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsCheckMate(pos)
	}
}

func TestIsStaleMate(t *testing.T) {
	fen := "k1K5/ppp5/1bP5/8/8/8/8/3r4 w - - 0 1"
	pos, _ := ParseFen(fen)
	if IsStaleMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected false, got true", fen)
	}

	fen = "5bnr/4p1pq/4Qpkr/7p/7P/4P3/PPPP1PP1/RNB1KBNR b KQ - 2 10"
	pos, _ = ParseFen(fen)
	if !IsStaleMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}

	fen = "2b5/pp3kp1/3p3p/3P4/5b2/6R1/6PK/r7 w - - 0 32"
	pos, _ = ParseFen(fen)
	if !IsStaleMate(pos) {
		t.Errorf("incorrect result: fen = %s: expected true, got false", fen)
	}
}

func BenchmarkIsStaleMate(b *testing.B) {
	pos := getDefaultPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsStaleMate(pos)
	}
}
