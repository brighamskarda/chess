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
