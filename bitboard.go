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

import "strconv"

// Bitboard is a 64 bit representation of a chess board. Each bit corresponds to a square with the the least significant bit (rightmost bit if using bit shifts) is A1, then B1, all the way up to H8.
//
// Most commonly there is a bitboard for each type of piece on the board with positive bits indicating squares occupied by that piece type. Bitboard are also used to represent all occupied squares, and squares that are attacked by certain pieces.
type Bitboard uint64

// Bit returns 1 if the bit at the specified index is set, otherwise 0.
func (bb Bitboard) Bit(index uint8) uint8 {
	if bb&(1<<index) > 0 {
		return 1
	} else {
		return 0
	}
}

// SetBit returns a copy of the bitboard with the specified bit set.
func (bb Bitboard) SetBit(index uint8) Bitboard {
	return bb | 1<<index
}

// ClearBit returns a copy of the bitboard with the specified bit cleared.
func (bb Bitboard) ClearBit(index uint8) Bitboard {
	return bb &^ (1 << index)
}

// Square returns 1 if the bit at the specified square is set, otherwise 0.
func (bb Bitboard) Square(s Square) uint8 {
	if s.File > FileH || s.Rank > Rank8 || s.File == NoFile || s.Rank == NoRank {
		return 0
	}
	return bb.Bit(squareToIndex(s))
}

// SetSquare returns a copy of the bitboard with the specified square set.
func (bb Bitboard) SetSquare(s Square) Bitboard {
	return bb.SetBit(squareToIndex(s))
}

// ClearSquare returns a copy of the bitboard with the specified square set.
func (bb Bitboard) ClearSquare(s Square) Bitboard {
	return bb.ClearBit(squareToIndex(s))
}

func squareToIndex(s Square) uint8 {
	return (uint8(s.File) - 1) + (uint8(s.Rank)-1)*8
}

// String gives a string representing the bitboard as if you were looking at a chess board from white's perspective.
func (bb Bitboard) String() string {
	s := ""
	for r := Rank8; r != NoRank; r -= 1 {
		for f := FileA; f <= FileH; f += 1 {
			s += strconv.Itoa(int(bb.Square(Square{f, r})))
		}
		if r != Rank1 {
			s += "\n"
		}
	}
	return s
}
