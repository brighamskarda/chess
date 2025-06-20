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
	"sync"
)

// Huge thanks to Chess Programming and his example implementation of bitboards. It was extremely useful. https://www.youtube.com/watch?v=4ohJQ9pCkHI&t=1199s

var rookRelevantBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var bishopRelevantBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var rookMasks [64]uint64
var bishopMasks [64]uint64

var rookAttacks [64][4096]uint64
var bishopAttacks [64][512]uint64

var rookMagics = [64]uint64{
	0xa8002c000108020,
	0x6c00049b0002001,
	0x100200010090040,
	0x2480041000800801,
	0x280028004000800,
	0x900410008040022,
	0x280020001001080,
	0x2880002041000080,
	0xa000800080400034,
	0x4808020004000,
	0x2290802004801000,
	0x411000d00100020,
	0x402800800040080,
	0xb000401004208,
	0x2409000100040200,
	0x1002100004082,
	0x22878001e24000,
	0x1090810021004010,
	0x801030040200012,
	0x500808008001000,
	0xa08018014000880,
	0x8000808004000200,
	0x201008080010200,
	0x801020000441091,
	0x800080204005,
	0x1040200040100048,
	0x120200402082,
	0xd14880480100080,
	0x12040280080080,
	0x100040080020080,
	0x9020010080800200,
	0x813241200148449,
	0x491604001800080,
	0x100401000402001,
	0x4820010021001040,
	0x400402202000812,
	0x209009005000802,
	0x810800601800400,
	0x4301083214000150,
	0x204026458e001401,
	0x40204000808000,
	0x8001008040010020,
	0x8410820820420010,
	0x1003001000090020,
	0x804040008008080,
	0x12000810020004,
	0x1000100200040208,
	0x430000a044020001,
	0x280009023410300,
	0xe0100040002240,
	0x200100401700,
	0x2244100408008080,
	0x8000400801980,
	0x2000810040200,
	0x8010100228810400,
	0x2000009044210200,
	0x4080008040102101,
	0x40002080411d01,
	0x2005524060000901,
	0x502001008400422,
	0x489a000810200402,
	0x1004400080a13,
	0x4000011008020084,
	0x26002114058042,
}

var bishopMagics = [64]uint64{
	0x89a1121896040240,
	0x2004844802002010,
	0x2068080051921000,
	0x62880a0220200808,
	0x4042004000000,
	0x100822020200011,
	0xc00444222012000a,
	0x28808801216001,
	0x400492088408100,
	0x201c401040c0084,
	0x840800910a0010,
	0x82080240060,
	0x2000840504006000,
	0x30010c4108405004,
	0x1008005410080802,
	0x8144042209100900,
	0x208081020014400,
	0x4800201208ca00,
	0xf18140408012008,
	0x1004002802102001,
	0x841000820080811,
	0x40200200a42008,
	0x800054042000,
	0x88010400410c9000,
	0x520040470104290,
	0x1004040051500081,
	0x2002081833080021,
	0x400c00c010142,
	0x941408200c002000,
	0x658810000806011,
	0x188071040440a00,
	0x4800404002011c00,
	0x104442040404200,
	0x511080202091021,
	0x4022401120400,
	0x80c0040400080120,
	0x8040010040820802,
	0x480810700020090,
	0x102008e00040242,
	0x809005202050100,
	0x8002024220104080,
	0x431008804142000,
	0x19001802081400,
	0x200014208040080,
	0x3308082008200100,
	0x41010500040c020,
	0x4012020c04210308,
	0x208220a202004080,
	0x111040120082000,
	0x6803040141280a00,
	0x2101004202410000,
	0x8200000041108022,
	0x21082088000,
	0x2410204010040,
	0x40100400809000,
	0x822088220820214,
	0x40808090012004,
	0x910224040218c9,
	0x402814422015008,
	0x90014004842410,
	0x1000042304105,
	0x10008830412a00,
	0x2520081090008908,
	0x40102000a0a60140,
}

var sliderAttacksOnce sync.Once

func initSliderAttacks() {
	sliderAttacksOnce.Do(func() {
		for square := 0; square < 64; square++ {
			bishopMasks[square] = maskBishopAttacks(square)
			rookMasks[square] = maskRookAttacks(square)

			currentBishopMask := bishopMasks[square]
			currentRookMask := rookMasks[square]

			bishopBitCount := bits.OnesCount64(currentBishopMask)
			rookBitCount := bits.OnesCount64(currentRookMask)

			bishopOccupancyVariations := 1 << bishopBitCount
			rookOccupancyVariations := 1 << rookBitCount

			for count := 0; count < bishopOccupancyVariations; count++ {
				occupancy := setOccupancy(count, bishopBitCount, currentBishopMask)
				magic_index := (occupancy * bishopMagics[square]) >> (64 - bishopRelevantBits[square])
				bishopAttacks[square][magic_index] = bishopAttacksOnTheFly(square, occupancy)
			}

			for count := 0; count < rookOccupancyVariations; count++ {
				occupancy := setOccupancy(count, rookBitCount, currentRookMask)
				magic_index := (occupancy * rookMagics[square]) >> (64 - rookRelevantBits[square])
				rookAttacks[square][magic_index] = rookAttacksOnTheFly(square, occupancy)
			}
		}
	})
}

func maskBishopAttacks(square int) uint64 {
	var attacks uint64
	var file, rank int
	targetFile := square % 8
	targetRank := square / 8

	rank = targetRank + 1
	file = targetFile + 1
	for rank <= 6 && file <= 6 {
		attacks |= 1 << (rank*8 + file)
		rank++
		file++
	}

	rank = targetRank + 1
	file = targetFile - 1
	for rank <= 6 && file >= 1 {
		attacks |= 1 << (rank*8 + file)
		rank++
		file--
	}

	rank = targetRank - 1
	file = targetFile + 1
	for rank >= 1 && file <= 6 {
		attacks |= 1 << (rank*8 + file)
		rank--
		file++
	}

	rank = targetRank - 1
	file = targetFile - 1
	for rank >= 1 && file >= 1 {
		attacks |= 1 << (rank*8 + file)
		rank--
		file--
	}

	return attacks
}

func maskRookAttacks(square int) uint64 {
	var attacks uint64
	var file, rank int
	targetFile := square % 8
	targetRank := square / 8

	rank = targetRank + 1
	for rank <= 6 {
		attacks |= 1 << (rank*8 + targetFile)
		rank++
	}

	rank = targetRank - 1
	for rank >= 1 {
		attacks |= 1 << (rank*8 + targetFile)
		rank--
	}

	file = targetFile + 1
	for file <= 6 {
		attacks |= 1 << (targetRank*8 + file)
		file++
	}

	file = targetFile - 1
	for file >= 1 {
		attacks |= 1 << (targetRank*8 + file)
		file--
	}

	return attacks
}

func setOccupancy(index int, bitCount int, mask uint64) uint64 {
	var occupancy uint64

	for count := 0; count < bitCount; count++ {
		square := getLs1bIndex(mask)
		mask = popBit(mask, square)
		if index&(1<<count) != 0 {
			occupancy |= 1 << square
		}
	}

	return occupancy
}

func getLs1bIndex(bitboard uint64) int {
	if bitboard != 0 {
		return bits.OnesCount64((bitboard & -bitboard) - 1)
	} else {
		return -1
	}
}

func popBit(bitboard uint64, square int) uint64 {
	if getBit(bitboard, square) {
		bitboard ^= 1 << square
	}
	return bitboard
}

func getBit(bitboard uint64, square int) bool {
	return (bitboard & (1 << square)) != 0
}

func bishopAttacksOnTheFly(square int, block uint64) uint64 {
	var attacks uint64

	tr := square / 8
	tf := square % 8

	for r, f := tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for r, f := tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for r, f := tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for r, f := tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	return attacks
}

func rookAttacksOnTheFly(square int, block uint64) uint64 {
	var attacks uint64

	tr := square / 8
	tf := square % 8

	for r := tr + 1; r <= 7; r++ {
		attacks |= 1 << (r*8 + tf)
		if block&(1<<(r*8+tf)) != 0 {
			break
		}
	}

	for r := tr - 1; r >= 0; r-- {
		attacks |= 1 << (r*8 + tf)
		if block&(1<<(r*8+tf)) != 0 {
			break
		}
	}

	for f := tf + 1; f <= 7; f++ {
		attacks |= 1 << (tr*8 + f)
		if block&(1<<(tr*8+f)) != 0 {
			break
		}
	}

	for f := tf - 1; f >= 0; f-- {
		attacks |= 1 << (tr*8 + f)
		if block&(1<<(tr*8+f)) != 0 {
			break
		}
	}

	return attacks
}
