package chess

import (
	"testing"
)

func TestNoPieceToString(t *testing.T) {
	pieceToTest := Piece{
		Color: NoColor,
		Type:  NoPieceType,
	}
	if pieceToTest.String() != " " {
		t.Error("No Piece does not equal \" \"")
	}

	pieceToTest.Color = White
	if pieceToTest.String() != "INVALID PIECE" {
		t.Error("No Piece does not equal \"INVALID PIECE\" when piece.Color == White")
	}

	pieceToTest.Color = NoColor
	pieceToTest.Type = Pawn
	if pieceToTest.String() != "INVALID PIECE" {
		t.Error("No Piece does not equal \"INVALID PIECE\" when piece.Type == Pawn")
	}
}

func TestPieceToString(t *testing.T) {
	pieceToTest := Piece{
		Color: White,
		Type:  Pawn,
	}
	if pieceToTest.String() != "P" {
		t.Error("White pawn does not equal \"P\"")
	}

	pieceToTest.Color = Black
	if pieceToTest.String() != "p" {
		t.Error("Black pawn does not equal \"p\"")
	}

	pieceToTest.Type = Bishop
	if pieceToTest.String() != "b" {
		t.Error("Black bishop does not equal \"b\"")
	}
}
