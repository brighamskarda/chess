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

import "strings"

// PieceType represents the type of a piece like a rook or a queen. See also [Piece]
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

// String returns a single lowercase letter representation of the piece type if valid, else "-".
func (pt PieceType) String() string {
	switch pt {
	case NoPieceType:
		return "-"
	case Pawn:
		return "p"
	case Rook:
		return "r"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Queen:
		return "q"
	case King:
		return "k"
	default:
		return "-"
	}
}

func parsePieceType(s string) PieceType {
	switch s {
	case "p":
		return Pawn
	case "n":
		return Knight
	case "b":
		return Bishop
	case "r":
		return Rook
	case "q":
		return Queen
	case "k":
		return King
	default:
		return NoPieceType
	}
}

// Piece represents a chess piece with type and color. The zero value is no piece.
type Piece struct {
	Type  PieceType
	Color Color
}

var (
	NoPiece = Piece{Type: NoPieceType, Color: NoColor}

	WhitePawn   = Piece{Type: Pawn, Color: White}
	WhiteRook   = Piece{Type: Rook, Color: White}
	WhiteKnight = Piece{Type: Knight, Color: White}
	WhiteBishop = Piece{Type: Bishop, Color: White}
	WhiteQueen  = Piece{Type: Queen, Color: White}
	WhiteKing   = Piece{Type: King, Color: White}

	BlackPawn   = Piece{Type: Pawn, Color: Black}
	BlackRook   = Piece{Type: Rook, Color: Black}
	BlackKnight = Piece{Type: Knight, Color: Black}
	BlackBishop = Piece{Type: Bishop, Color: Black}
	BlackQueen  = Piece{Type: Queen, Color: Black}
	BlackKing   = Piece{Type: King, Color: Black}
)

// String returns a single letter representation of the piece type if valid, else "-".
//
// White pieces are uppercase and black pieces are lowercase.
func (p Piece) String() string {
	if p.Color == White {
		return strings.ToUpper(p.Type.String())
	} else if p.Color == Black {
		return p.Type.String()
	} else {
		return "-"
	}
}
