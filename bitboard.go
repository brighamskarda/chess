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

import (
	"math/bits"
	"strconv"
)

// Bitboard is a 64 bit representation of a chess board. Each bit corresponds to a square with the the least significant bit (rightmost bit if using bit shifts) is A1, then B1, all the way up to H8.
//
// Most commonly there is a bitboard for each type of piece on the board with positive bits indicating squares occupied by that piece type. Bitboards are also used to represent all occupied squares, and squares that are attacked by certain pieces.
type Bitboard uint64

// Bit returns 1 if the bit at the specified index is set, otherwise 0. If index is too high 0 is always returned.
func (bb Bitboard) Bit(index uint8) uint8 {
	if bb&(1<<index) > 0 {
		return 1
	} else {
		return 0
	}
}

// SetBit returns a copy of the bitboard with the specified bit set. Nothing happens if index is too high.
func (bb Bitboard) SetBit(index uint8) Bitboard {
	return bb | 1<<index
}

// ClearBit returns a copy of the bitboard with the specified bit cleared. Nothing happens if index is too high.
func (bb Bitboard) ClearBit(index uint8) Bitboard {
	return bb & ^(1 << index)
}

// Square returns 1 if the bit at the specified square is set, otherwise 0.
func (bb Bitboard) Square(s Square) uint8 {
	if !squareOnBoard(s) {
		return 0
	}
	return bb.Bit(squareToIndex(s))
}

// SetSquare returns a copy of the bitboard with the specified square set. Nothing is different if the square is invalid.
func (bb Bitboard) SetSquare(s Square) Bitboard {
	if !squareOnBoard(s) {
		return bb
	}
	return bb.SetBit(squareToIndex(s))
}

// ClearSquare returns a copy of the bitboard with the specified square cleared. Nothing is different if the square is invalid.
func (bb Bitboard) ClearSquare(s Square) Bitboard {
	if !squareOnBoard(s) {
		return bb
	}
	return bb.ClearBit(squareToIndex(s))
}

func squareToIndex(s Square) uint8 {
	return (uint8(s.File) - 1) + (uint8(s.Rank)-1)*8
}

func indexToSquare(index int) Square {
	if index >= 64 || index < 0 {
		return NoSquare
	}
	file := File((index % 8) + 1)
	rank := Rank((index / 8) + 1)
	return Square{file, rank}
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

// WhitePawnAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of white pawns.
func (bb Bitboard) WhitePawnAttacks() Bitboard {
	return bb.pawnAttacksNE() | bb.pawnAttacksNW()
}

func (bb Bitboard) pawnAttacksNE() Bitboard {
	return (bb << 9) & 0xFEFEFEFEFEFEFEFE
}

func (bb Bitboard) pawnAttacksNW() Bitboard {
	return (bb << 7) & 0x7F7F7F7F7F7F7F7F
}

// BlackPawnAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of black pawns.
func (bb Bitboard) BlackPawnAttacks() Bitboard {
	return bb.pawnAttacksSE() | bb.pawnAttacksSW()
}

func (bb Bitboard) pawnAttacksSE() Bitboard {
	return (bb >> 7) & 0xFEFEFEFEFEFEFEFE
}

func (bb Bitboard) pawnAttacksSW() Bitboard {
	return (bb >> 9) & 0x7F7F7F7F7F7F7F7F
}

// RookAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of rooks. Occupied should indicate all squares on the board occupied by either color, including the rooks that are moving.
func (bb Bitboard) RookAttacks(occupied Bitboard) Bitboard {
	var movesRight Bitboard = (occupied ^
		(occupied | 0x101010101010101 - 2*bb)) & ^Bitboard(0x101010101010101)

	occupied_reverse := bits.Reverse64(uint64(occupied))
	var movesLeft Bitboard = Bitboard(bits.Reverse64((occupied_reverse ^ (occupied_reverse | 0x101010101010101 - 2*bits.Reverse64(uint64(bb)))) & ^uint64(0x101010101010101)))

	ccRook := bb.rotate90CC()
	ccOccupied := occupied.rotate90CC()

	var movesDown Bitboard = ((ccOccupied ^ ((ccOccupied | 0x101010101010101) - 2*ccRook)) & ^Bitboard(0x101010101010101)).rotate90C()

	cRook := bb.rotate90C()
	cOccupied := occupied.rotate90C()

	var movesUp Bitboard = ((cOccupied ^ ((cOccupied | 0x101010101010101) - 2*cRook)) & ^Bitboard(0x101010101010101)).rotate90CC()
	return movesLeft | movesRight | movesDown | movesUp
}

// Knight Attacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of knights.
func (bb Bitboard) knightAttacks() Bitboard {
	return ((bb << 17) & 0xfefefefefefefefe) |
		((bb << 10) & 0xfcfcfcfcfcfcfcfc) |
		((bb >> 6) & 0xfcfcfcfcfcfcfcfc) |
		((bb >> 15) & 0xfefefefefefefefe) |
		((bb >> 17) & 0x7f7f7f7f7f7f7f7f) |
		((bb >> 10) & 0x3f3f3f3f3f3f3f3f) |
		((bb << 6) & 0x3f3f3f3f3f3f3f3f) |
		((bb << 15) & 0x7f7f7f7f7f7f7f7f)
}

// BishopAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of bishops. Occupied should indicate all squares on the board occupied by either color, including the bishops that are moving.
func (bb Bitboard) BishopAttacks(occupied Bitboard) Bitboard {
	var attacks Bitboard = 0
	for bb != 0 {
		singleBishop := bb & -bb
		bb ^= singleBishop
		attacks |= singleBishop.diagonalAttacks(occupied) | singleBishop.antiDiagonalAttacks(occupied)
	}
	return attacks
}

func (bb Bitboard) diagonalAttacks(occupied Bitboard) Bitboard {
	diagonalMask := getDiagonalMask(bb)
	forward := occupied & diagonalMask
	reverse := Bitboard(bits.ReverseBytes64(uint64(forward)))

	forward -= 2 * bb
	reverse -= 2 * Bitboard(bits.ReverseBytes64(uint64(bb)))

	forward ^= Bitboard(bits.ReverseBytes64(uint64(reverse)))
	forward &= diagonalMask

	return forward
}

func (bb Bitboard) antiDiagonalAttacks(occupied Bitboard) Bitboard {
	diagonalMask := getAntiDiagonalMask(bb)
	forward := occupied & diagonalMask
	reverse := Bitboard(bits.ReverseBytes64(uint64(forward)))

	forward -= 2 * bb
	reverse -= 2 * Bitboard(bits.ReverseBytes64(uint64(bb)))

	forward ^= Bitboard(bits.ReverseBytes64(uint64(reverse)))
	forward &= diagonalMask

	return forward
}

var diagonalMasks []Bitboard

func getDiagonalMask(bb Bitboard) Bitboard {
	if diagonalMasks == nil {
		initializeDiagonalMasks()
	}
	return diagonalMasks[bits.TrailingZeros64(uint64(bb))]
}

func initializeDiagonalMasks() {
	diagonalMasks = make([]Bitboard, 0, 64)
	for i := range 64 {
		rank := Rank(i / 8)
		file := File(i % 8)
		diff := int(file) - int(rank)

		var mask Bitboard = 0

		for r, f := 0, diff; r < 8; r, f = r+1, f+1 {
			if f >= 0 && f < 8 {
				mask = mask.SetSquare(Square{File: File(f + 1), Rank: Rank(r + 1)})
			}
		}
		diagonalMasks = append(diagonalMasks, mask)
	}
}

var antiDiagonalMasks []Bitboard

func getAntiDiagonalMask(bb Bitboard) Bitboard {
	if antiDiagonalMasks == nil {
		initializeAntiDiagonalMasks()
	}
	return antiDiagonalMasks[bits.TrailingZeros64(uint64(bb))]
}

func initializeAntiDiagonalMasks() {
	antiDiagonalMasks = make([]Bitboard, 0, 64)
	for i := range 64 {
		rank := Rank(i / 8)
		file := File(i % 8)
		sum := int(file) + int(rank)

		var mask Bitboard = 0

		for r, f := 0, sum; r < 8 && f >= 0; r, f = r+1, f-1 {
			if f < 8 {
				mask = mask.SetSquare(Square{File: File(f + 1), Rank: Rank(r + 1)})
			}
		}
		antiDiagonalMasks = append(antiDiagonalMasks, mask)
	}
}

// QueenAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of queens. Occupied should indicate all squares on the board occupied by either color, including the queens that are moving.
func (bb Bitboard) QueenAttacks(occupied Bitboard) Bitboard {
	var attacks Bitboard = 0
	attacks |= bb.RookAttacks(occupied)
	attacks |= bb.BishopAttacks(occupied)
	return attacks
}

// KingAttacks returns a bitboard indicating all the squares attacked by this bitboard assuming its a bitboard of kings.
func (bb Bitboard) KingAttacks() Bitboard {
	return ((bb >> 9) & 0x7f7f7f7f7f7f7f7f) |
		(bb >> 8) |
		((bb >> 7) & 0xfefefefefefefefe) |
		((bb >> 1) & 0x7f7f7f7f7f7f7f7f) |
		((bb << 1) & 0xfefefefefefefefe) |
		((bb << 7) & 0x7f7f7f7f7f7f7f7f) |
		(bb << 8) |
		((bb << 9) & 0xfefefefefefefefe)
}

// flipVert flips a bitboard vertically.
func (bb Bitboard) flipVert() Bitboard {
	return Bitboard(bits.ReverseBytes64(uint64(bb)))
}

// flipDiagA1H8 flips a board diagonally over the A1H8 line.
func (bb Bitboard) flipDiagA1H8() Bitboard {
	t := 0x0f0f0f0f00000000 & (bb ^ (bb << 28))
	bb ^= t ^ (t >> 28)
	t = 0x3333000033330000 & (bb ^ (bb << 14))
	bb ^= t ^ (t >> 14)
	t = 0x5500550055005500 & (bb ^ (bb << 7))
	bb ^= t ^ (t >> 7)
	return bb
}

// rotate90CC rotates a bitboard 90 degrees counter clockwise.
func (bb Bitboard) rotate90CC() Bitboard {
	return bb.flipVert().flipDiagA1H8()
}

// rotate90C rotates a bitboard 90 degrees clockwise.
func (bb Bitboard) rotate90C() Bitboard {
	return bb.flipDiagA1H8().flipVert()
}
