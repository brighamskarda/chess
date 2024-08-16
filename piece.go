package chess

import (
	"strings"
)

type PieceType uint8

const (
	NoPieceType PieceType = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

func (pt PieceType) String() string {
	switch pt {
	case NoPieceType:
		return "NO-PIECE-TYPE"
	case Pawn:
		return "P"
	case Rook:
		return "R"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "INVALID PIECE TYPE"
	}
}

func pieceTypeIsValid(pt PieceType) bool {
	return pt <= 6
}

type Piece struct {
	Color Color
	Type  PieceType
}

func (p Piece) String() string {
	if !pieceIsValid(p) {
		return "INVALID PIECE"
	}
	if p.Color == NoColor || p.Type == NoPieceType {
		return " "
	}
	pieceStr := p.Type.String()
	if p.Color == Black {
		return strings.ToLower(pieceStr)
	}
	return pieceStr
}

func pieceIsValid(p Piece) bool {
	if !pieceTypeIsValid(p.Type) || !colorIsValid(p.Color) {
		return false
	}
	return true
}

var (
	NoPiece Piece = Piece{NoColor, NoPieceType}

	WhitePawn   Piece = Piece{White, Pawn}
	WhiteRook   Piece = Piece{White, Rook}
	WhiteKnight Piece = Piece{White, Knight}
	WhiteBishop Piece = Piece{White, Bishop}
	WhiteQueen  Piece = Piece{White, Queen}
	WhiteKing   Piece = Piece{White, King}

	BlackPawn   Piece = Piece{Black, Pawn}
	BlackRook   Piece = Piece{Black, Rook}
	BlackKnight Piece = Piece{Black, Knight}
	BlackBishop Piece = Piece{Black, Bishop}
	BlackQueen  Piece = Piece{Black, Queen}
	BlackKing   Piece = Piece{Black, King}
)
