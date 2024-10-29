package chess

import (
	"errors"
	"strings"
	"unicode"
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

// String returns the uppercase ascii letter representation of the piece. NO-PIECE-TYPE, or INVALID PIECE TYPE otherwise.
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

func isValidPieceType(pt PieceType) bool {
	return pt <= 6
}

func parsePieceType(r rune) (PieceType, error) {
	r = unicode.ToLower(r)
	switch r {
	case 'p':
		return Pawn, nil
	case 'r':
		return Rook, nil
	case 'n':
		return Knight, nil
	case 'b':
		return Bishop, nil
	case 'q':
		return Queen, nil
	case 'k':
		return King, nil
	default:
		return NoPieceType, errors.New("can't parse piece type")
	}
}

type Piece struct {
	Color Color
	Type  PieceType
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

// String gives a one letter ascii character representing the piece. " " for no piece. INVALID PIECE otherwise.
func (p Piece) String() string {
	if !isValidPiece(p) {
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

func isValidPiece(p Piece) bool {
	if !isValidPieceType(p.Type) || !isValidColor(p.Color) {
		return false
	}
	if p.Type == NoPieceType || p.Color == NoColor {
		if p.Type != NoPieceType || p.Color != NoColor {
			return false
		}
	}
	return true
}

// ParsePiece attempts to parse a piece from a given rune. Currently only supports ascii characters (no piece symbols). Uppercase is white, lowercase is black.
func ParsePiece(r rune) (Piece, error) {
	pieceType, err := parsePieceType(r)
	if err != nil {
		return NoPiece, errors.New("can't parse piece type")
	}
	var color Color
	if unicode.IsUpper(r) {
		color = White
	} else {
		color = Black
	}
	return Piece{color, pieceType}, nil
}
