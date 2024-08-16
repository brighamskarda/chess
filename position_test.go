package chess

import (
	"testing"
)

const emptyBoardStr string = `8        
7        
6        
5        
4        
3        
2        
1        
 ABCDEFGH`

const defaultBoardStr string = `8rnbqkbnr
7pppppppp
6        
5        
4        
3        
2PPPPPPPP
1RNBQKBNR
 ABCDEFGH`

func TestBoardString(t *testing.T) {
	pos := Position{}
	if pos.String() != emptyBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.String(), emptyBoardStr)
	}
	pos = getDefaultPosition()
	if pos.String() != defaultBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.String(), defaultBoardStr)
	}
}

func getDefaultPosition() Position {
	pos := Position{}
	pos.Board[0] = BlackRook
	pos.Board[1] = BlackKnight
	pos.Board[2] = BlackBishop
	pos.Board[3] = BlackQueen
	pos.Board[4] = BlackKing
	pos.Board[5] = BlackBishop
	pos.Board[6] = BlackKnight
	pos.Board[7] = BlackRook
	for i := 8; i < 16; i++ {
		pos.Board[i] = BlackPawn
	}
	for i := 48; i < 56; i++ {
		pos.Board[i] = WhitePawn
	}
	pos.Board[56] = WhiteRook
	pos.Board[57] = WhiteKnight
	pos.Board[58] = WhiteBishop
	pos.Board[59] = WhiteQueen
	pos.Board[60] = WhiteKing
	pos.Board[61] = WhiteBishop
	pos.Board[62] = WhiteKnight
	pos.Board[63] = WhiteRook

	pos.Turn = White

	pos.WhiteKingSideCastle = true
	pos.WhiteQueenSideCastle = true
	pos.BlackKingSideCastle = true
	pos.BlackQueenSideCastle = true

	pos.EnPassant = nil

	pos.HalfMove = 0
	pos.FullMove = 1

	return pos
}

func TestPositionPieceAt(t *testing.T) {
	pos := getDefaultPosition()
	if pos.PieceAt(C1) != WhiteBishop {
		t.Errorf("pos.At(C1) != WhiteBishop. Actual %s", pos.PieceAt(C1).String())
	}
}

func TestPositionSetPieceAt(t *testing.T) {
	pos := Position{}
	pos.SetPieceAt(D4, BlackPawn)
	if pos.Board[35] != BlackPawn {
		t.Fail()
	}
}

func TestParseFen(t *testing.T) {
	pos, err := ParseFen(DefaultFen)
	if err != nil {
		t.Error("ParseFen set error nil when fen is valid")
	}
	if pos != getDefaultPosition() {
		t.Errorf("ParseFen incorrect output. Actual:%+v\nExpected%+v", pos, getDefaultPosition())
	}

	pos, err = ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16")
	if err != nil {
		t.Error("ParseFen set error nil when fen is valid")
	}
	updatedPos := getDefaultPosition()
	updatedPos.Board[0] = WhitePawn
	updatedPos.Turn = Black
	updatedPos.BlackKingSideCastle = false
	*updatedPos.EnPassant = E6
	updatedPos.HalfMove = 32
	updatedPos.FullMove = 16
	if pos != updatedPos {
		t.Errorf("ParseFen incorrect output. Actual:%+v\nExpected%+v", pos, updatedPos)
	}
}

func TestParseFenInvalid(t *testing.T) {
	_, err := ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3216")
	if err == nil {
		t.Error("ParseFen failed to set error for Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3216")
	}

	_, err = ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3a2 16")
	if err == nil {
		t.Error("ParseFen failed to set error for Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3a2 16")
	}
}

func BenchmarkParseFen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseFen(DefaultFen)
	}
}
