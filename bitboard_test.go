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

import "testing"

func TestBitboardBit(t *testing.T) {
	var bb Bitboard = 0x8
	for index := range 66 {
		if bb.Bit(uint8(index)) != 0 && index != 3 {
			t.Errorf("incorrect result for index %d: expected 0", index)
		}
		if index == 3 && bb.Bit(uint8(index)) != 1 {
			t.Errorf("incorrect result for index 3: expected 1")
		}
	}
}

func TestBitboardSetBit(t *testing.T) {
	var bb Bitboard = 0
	bb = bb.SetBit(9)
	if bb != 0x200 {
		t.Errorf("incorrect result: expected 0x200, got 0x%X", uint64(bb))
	}
}

func TestBitboardClearBit(t *testing.T) {
	var bb Bitboard = 0xFFFFFFFFFFFFFFFF
	bb = bb.ClearBit(9)
	if bb != 0xFFFFFFFFFFFFFDFF {
		t.Errorf("incorrect result: expected 0xFFFFFFFFFFFFFFFF, got 0x%X", uint64(bb))
	}
}

func TestBitboardSquareFullBb(t *testing.T) {
	var bb Bitboard = 0xFFFFFFFFFFFFFFFF
	for _, square := range AllSquares {
		if bb.Square(square) != 1 {
			t.Errorf("incorrect result for full bitboard: expected 1")
		}
	}
	if bb.Square(NoSquare) != 0 {
		t.Errorf("incorrect result for full bitboard: expected NoSquare to be 0")
	}
	if bb.Square(Square{FileH, 9}) != 0 {
		t.Errorf("incorrect result for full bitboard: expected square h- to be 0")
	}
	if bb.Square(Square{0xFF, Rank8}) != 0 {
		t.Errorf("incorrect result for full bitboard: expected square -8 to be 0")
	}
}

func TestBitboardSquareEmptyBb(t *testing.T) {
	var bb Bitboard = 0
	for _, square := range AllSquares {
		if bb.Square(square) != 0 {
			t.Errorf("incorrect result for empty bitboard: expected 0")
		}
	}
}

func TestBitboardString(t *testing.T) {
	var bb Bitboard = 0xFFFFFFFFFFFFFFFF
	expected := `11111111
11111111
11111111
11111111
11111111
11111111
11111111
11111111`
	actual := bb.String()
	if actual != expected {
		t.Errorf("incorrect result for full bitboard: got \n%s", actual)
	}

	bb = 0
	expected = `00000000
00000000
00000000
00000000
00000000
00000000
00000000
00000000`
	actual = bb.String()
	if actual != expected {
		t.Errorf("incorrect result for empty bitboard: got \n%s", actual)
	}

	bb = 0x200
	expected = `00000000
00000000
00000000
00000000
00000000
00000000
01000000
00000000`
	actual = bb.String()
	if actual != expected {
		t.Errorf("incorrect result for square B2: got \n%s", actual)
	}
}
