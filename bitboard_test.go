package chess

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetBitBoard(t *testing.T) {
	pos, _ := ParseFen(DefaultFen)
	var expected BitBoard = 0b00001000_000000000_00000000_0000000_00000000_00000000_00000000_00000000
	if bb := pos.GetBitBoard(WhiteKing); bb != expected {
		t.Errorf("WhiteKing bitboard incorrect: expected %b, got %b", expected, bb)
	}

	expected = 0b00000000_000000000_00000000_0000000_00000000_00000000_00000000_00001000
	if bb := pos.GetBitBoard(BlackKing); bb != expected {
		t.Errorf("BlackKing bitboard incorrect: expected %b, got %b", expected, bb)
	}

	expected = 0b01000010_000000000_00000000_0000000_00000000_00000000_00000000_00000000
	if bb := pos.GetBitBoard(WhiteKnight); bb != expected {
		t.Errorf("WhiteKnight bitboard incorrect: expected %b, got %b", expected, bb)
	}

	expected = 0b00000000_000000000_00000000_0000000_00000000_00000000_11111111_00000000
	if bb := pos.GetBitBoard(BlackPawn); bb != expected {
		t.Errorf("BlackPawn bitboard incorrect: expected %b, got %b", expected, bb)
	}

	pos = &Position{}
	expected = 0b00000000_000000000_00000000_0000000_00000000_00000000_00000000_00000000
	if bb := pos.GetBitBoard(WhitePawn); bb != expected {
		t.Errorf("WhitePawn bitboard incorrect (empty board): expected %b, got %b", expected, bb)
	}
}

func TestGetBoard(t *testing.T) {
	pos, _ := ParseFen(DefaultFen)
	expected := Board{
		WhiteKings:   0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackKings:   0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
		WhiteQueens:  0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackQueens:  0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
		WhiteBishops: 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackBishops: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
		WhiteKnights: 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackKnights: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
		WhiteRooks:   0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackRooks:   0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
		WhitePawns:   0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
		BlackPawns:   0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
	}
	if board := pos.GetBoard(); board != expected {
		t.Errorf("Board incorrect: %v", cmp.Diff(expected, board))
	}

	pos = &Position{}
	expected = Board{}
	if board := pos.GetBoard(); board != expected {
		t.Errorf("Empty board incorrect: %v", cmp.Diff(expected, board))
	}
}
