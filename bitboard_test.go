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

func TestRookAttacks_emptyBoard(t *testing.T) {
	var bb Bitboard
	var expected Bitboard
	var actual Bitboard

	bb = 0x1
	expected = 0x1010101010101fe
	actual = bb.rookAttacks(0x1)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x8000000000000000
	expected = 0x7f80808080808080
	actual = bb.rookAttacks(0x8000000000000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x200004000000
	expected = 0x2424df24fb242424
	actual = bb.rookAttacks(0x200004000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_rightBlocked(t *testing.T) {
	var bb Bitboard = 0x4000000

	var expected Bitboard = 0x40404041b040404
	actual := bb.rookAttacks(0x50000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_leftBlocked(t *testing.T) {
	var bb Bitboard = 0x4000000

	var expected Bitboard = 0x4040404fa040404
	actual := bb.rookAttacks(0x2000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_upBlocked(t *testing.T) {
	var bb Bitboard = 0x4000000

	var expected Bitboard = 0x404fb040404
	actual := bb.rookAttacks(0x40000000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_downBlocked(t *testing.T) {
	var bb Bitboard = 0x4000000

	var expected Bitboard = 0x4040404fb040400
	actual := bb.rookAttacks(0x400)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_blockedAllWays(t *testing.T) {
	var bb Bitboard = 0x100000000000

	var expected Bitboard = 0x10681010000000
	actual := bb.rookAttacks(0x104a0010001000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestRookAttacks_onSameRankFile(t *testing.T) {
	var bb Bitboard = 0x100014000000

	var expected Bitboard = 0x1414ff14ff141414
	actual := bb.rookAttacks(0x100014000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestKnightAttacksCorners(t *testing.T) {
	var bb Bitboard = 0x100000000000000
	var expected Bitboard = 0x4020000000000

	actual := bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x8000000000000000
	expected = 0x20400000000000
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x80
	expected = 0x402000
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x1
	expected = 0x20400
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestKnightAttacksOneFromCorners(t *testing.T) {
	var bb Bitboard = 0x2000000000000
	var expected Bitboard = 0x800080500000000

	actual := bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x40000000000000
	expected = 0x100010a000000000
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x4000
	expected = 0xa0100010
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x200
	expected = 0x5080008
	actual = bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestKnightAttacksCenter(t *testing.T) {
	var bb Bitboard = 0x200000000000
	var expected Bitboard = 0x5088008850000000

	actual := bb.knightAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_emptyBoard(t *testing.T) {
	var bb Bitboard
	var expected Bitboard

	bb = 0x81
	expected = 0x8142241818244200
	actual := bb.bishopAttacks(0x81)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x8100000000000000
	expected = 0x42241818244281
	actual = bb.bishopAttacks(0x8100000000000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x24000000000000
	expected = 0x5a005a9924428100
	actual = bb.bishopAttacks(0x24000000000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_NEBlocked(t *testing.T) {
	var bb Bitboard = 0x200
	var expected Bitboard = 0x8050005

	actual := bb.bishopAttacks(0x200008000200)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_SEBlocked(t *testing.T) {
	var bb Bitboard = 0x2000000000000
	var expected Bitboard = 0x500050800000000

	actual := bb.bishopAttacks(0x800200000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_SWBlocked(t *testing.T) {
	var bb Bitboard = 0x40000000000000
	var expected Bitboard = 0xa000a01000000000

	actual := bb.bishopAttacks(0x1000040000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_NWBlocked(t *testing.T) {
	var bb Bitboard = 0x4000
	var expected Bitboard = 0x10a000a0

	actual := bb.bishopAttacks(0x40010000000)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_blockedAllWays(t *testing.T) {
	var bb Bitboard = 0x10000000
	var expected Bitboard = 0x402800284400

	actual := bb.bishopAttacks(0x400810004400)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestBishopAttacks_onSameRankFile(t *testing.T) {
	var bb Bitboard = 0x440010004400
	var expected Bitboard = 0x11aa44aa11aa44aa

	actual := bb.bishopAttacks(0x440010004400)
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}

func TestKingAttacks(t *testing.T) {
	var bb Bitboard
	var expected Bitboard
	var actual Bitboard

	bb = 0x8100000000000081
	expected = 0x42c300000000c342
	actual = bb.kingAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}

	bb = 0x800000000
	expected = 0x1c141c000000
	actual = bb.kingAttacks()
	if expected != actual {
		t.Errorf("incorrect result: expected\n%s\n\ngot\n%s", expected.String(), actual.String())
	}
}
