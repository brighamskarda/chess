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

import "math/bits"

// PseudoLegalMoves are moves that are legal except they leave one's king in check. Returns nil if moves could not be generated (for example if pos.SideToMove was not set). Returns an empty slice if move generation was successful, but no moves were found.
func PseudoLegalMoves(pos *Position) []Move {
	if pos.SideToMove != White && pos.SideToMove != Black {
		return nil
	}
	moves := make([]Move, 0)
	moves = append(moves, pawnMoves(pos)...)
	moves = append(moves, rookMoves(pos)...)
	moves = append(moves, knightMoves(pos)...)
	moves = append(moves, bishopMoves(pos)...)
	moves = append(moves, queenMoves(pos)...)
	moves = append(moves, kingMoves(pos)...)
	return moves
}

func pawnMoves(pos *Position) []Move {
	if pos.SideToMove == White {
		return whitePawnMoves(pos)
	}

	if pos.SideToMove == Black {
		return blackPawnMoves(pos)
	}

	return nil
}

func whitePawnMoves(pos *Position) []Move {
	moves := make([]Move, 0)

	whitePawns := pos.Bitboard(WhitePawn)
	occupied := pos.OccupiedBitboard()
	enemies := pos.ColorBitboard(Black) | (1 << squareToIndex(pos.EnPassant))

	moves = append(moves, whitePawnsMoveForward(whitePawns, occupied)...)
	moves = append(moves, whitePawnsMoveForward2(whitePawns, occupied)...)
	moves = append(moves, whitePawnsTakeNE(whitePawns, enemies)...)
	moves = append(moves, whitePawnsTakeNW(whitePawns, enemies)...)
	includeWhitePawnPromotions(&moves)
	return moves
}

func whitePawnsMoveForward(whitePawns Bitboard, occupied Bitboard) []Move {
	moves := make([]Move, 0)
	moveForward := (whitePawns << 8) &^ occupied

	for moveForward != 0 {
		squareIndex := bits.TrailingZeros64(uint64(moveForward))
		moveForward ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File, square.Rank - 1}, square, NoPieceType})
	}
	return moves
}

func whitePawnsMoveForward2(whitePawns Bitboard, occupied Bitboard) []Move {
	moves := make([]Move, 0)
	moveForward2 := (whitePawns << 16) &^ occupied

	for moveForward2 != 0 {
		squareIndex := bits.TrailingZeros64(uint64(moveForward2))
		moveForward2 ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		if square.Rank == Rank4 && occupied.Square(Square{square.File, Rank3}) == 0 {
			moves = append(moves, Move{Square{square.File, square.Rank - 2}, square, NoPieceType})
		}
	}
	return moves
}

func whitePawnsTakeNE(whitePawns Bitboard, enemies Bitboard) []Move {
	moves := make([]Move, 0)
	neAttacks := whitePawns.pawnAttacksNE() & enemies

	for neAttacks != 0 {
		squareIndex := bits.TrailingZeros64(uint64(neAttacks))
		neAttacks ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File - 1, square.Rank - 1}, square, NoPieceType})
	}
	return moves
}

func whitePawnsTakeNW(whitePawns Bitboard, enemies Bitboard) []Move {
	moves := make([]Move, 0)
	neAttacks := whitePawns.pawnAttacksNW() & enemies

	for neAttacks != 0 {
		squareIndex := bits.TrailingZeros64(uint64(neAttacks))
		neAttacks ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File + 1, square.Rank - 1}, square, NoPieceType})
	}
	return moves
}

func includeWhitePawnPromotions(moves *[]Move) {
	numMoves := len(*moves)
	for i := range numMoves {
		if (*moves)[i].ToSquare.Rank == Rank8 {
			moveCopy := (*moves)[i]
			(*moves)[i].Promotion = Queen
			moveCopy.Promotion = Rook
			*moves = append(*moves, moveCopy)
			moveCopy.Promotion = Knight
			*moves = append(*moves, moveCopy)
			moveCopy.Promotion = Bishop
			*moves = append(*moves, moveCopy)
		}
	}
}

func blackPawnMoves(pos *Position) []Move {
	moves := make([]Move, 0)

	blackPawns := pos.Bitboard(BlackPawn)
	occupied := pos.OccupiedBitboard()
	enemies := pos.ColorBitboard(White) | (1 << squareToIndex(pos.EnPassant))

	moves = append(moves, blackPawnsMoveForward(blackPawns, occupied)...)
	moves = append(moves, blackPawnsMoveForward2(blackPawns, occupied)...)
	moves = append(moves, blackPawnsTakeSE(blackPawns, enemies)...)
	moves = append(moves, blackPawnsTakeSW(blackPawns, enemies)...)
	includeBlackPawnPromotions(&moves)
	return moves
}

func blackPawnsMoveForward(blackPawns Bitboard, occupied Bitboard) []Move {
	moves := make([]Move, 0)
	moveForward := (blackPawns >> 8) &^ occupied

	for moveForward != 0 {
		squareIndex := bits.TrailingZeros64(uint64(moveForward))
		moveForward ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File, square.Rank + 1}, square, NoPieceType})
	}
	return moves
}

func blackPawnsMoveForward2(blackPawns Bitboard, occupied Bitboard) []Move {
	moves := make([]Move, 0)
	moveForward2 := (blackPawns >> 16) &^ occupied

	for moveForward2 != 0 {
		squareIndex := bits.TrailingZeros64(uint64(moveForward2))
		moveForward2 ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		if square.Rank == Rank5 && occupied.Square(Square{square.File, Rank6}) == 0 {
			moves = append(moves, Move{Square{square.File, square.Rank + 2}, square, NoPieceType})
		}
	}
	return moves
}

func blackPawnsTakeSE(blackPawns Bitboard, enemies Bitboard) []Move {
	moves := make([]Move, 0)
	seAttacks := blackPawns.pawnAttacksSE() & enemies

	for seAttacks != 0 {
		squareIndex := bits.TrailingZeros64(uint64(seAttacks))
		seAttacks ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File - 1, square.Rank + 1}, square, NoPieceType})
	}
	return moves
}

func blackPawnsTakeSW(blackPawns Bitboard, enemies Bitboard) []Move {
	moves := make([]Move, 0)
	swAttacks := blackPawns.pawnAttacksSW() & enemies

	for swAttacks != 0 {
		squareIndex := bits.TrailingZeros64(uint64(swAttacks))
		swAttacks ^= 1 << squareIndex
		square := indexToSquare(squareIndex)
		moves = append(moves, Move{Square{square.File + 1, square.Rank + 1}, square, NoPieceType})
	}
	return moves
}

func includeBlackPawnPromotions(moves *[]Move) {
	numMoves := len(*moves)
	for i := 0; i < numMoves; i++ {
		if (*moves)[i].ToSquare.Rank == Rank1 {
			moveCopy := (*moves)[i]
			(*moves)[i].Promotion = Queen
			moveCopy.Promotion = Rook
			*moves = append(*moves, moveCopy)
			moveCopy.Promotion = Knight
			*moves = append(*moves, moveCopy)
			moveCopy.Promotion = Bishop
			*moves = append(*moves, moveCopy)
		}
	}
}

func rookMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	var rooks Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		rooks = pos.Bitboard(WhiteRook)
		allies = pos.ColorBitboard(White)
	case Black:
		rooks = pos.Bitboard(BlackRook)
		allies = pos.ColorBitboard(Black)
	default:
		return moves
	}

	for rooks != 0 {
		singleRookIndex := bits.TrailingZeros64(uint64(rooks))
		singleRookBitboard := Bitboard(1 << singleRookIndex)
		rooks ^= singleRookBitboard
		rookAttacks := singleRookBitboard.RookAttacks(occupied)
		rookAttacks &^= allies
		singleRookSquare := indexToSquare(singleRookIndex)
		for rookAttacks != 0 {
			attackIndex := bits.TrailingZeros64(uint64(rookAttacks))
			rookAttacks ^= 1 << attackIndex
			attackSquare := indexToSquare(attackIndex)
			moves = append(moves, Move{singleRookSquare, attackSquare, NoPieceType})
		}
	}

	return moves
}

func knightMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	var knights Bitboard
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		knights = pos.Bitboard(WhiteKnight)
		allies = pos.ColorBitboard(White)
	case Black:
		knights = pos.Bitboard(BlackKnight)
		allies = pos.ColorBitboard(Black)
	default:
		return moves
	}

	for knights != 0 {
		singleKnightIndex := bits.TrailingZeros64(uint64(knights))
		singleKnightBitboard := Bitboard(1 << singleKnightIndex)
		knights ^= singleKnightBitboard
		knightAttacks := singleKnightBitboard.KnightAttacks()
		knightAttacks &^= allies
		singleKnightSquare := indexToSquare(singleKnightIndex)
		for knightAttacks != 0 {
			attackIndex := bits.TrailingZeros64(uint64(knightAttacks))
			knightAttacks ^= 1 << attackIndex
			attackSquare := indexToSquare(attackIndex)
			moves = append(moves, Move{singleKnightSquare, attackSquare, NoPieceType})
		}
	}

	return moves
}

func bishopMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	var bishops Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		bishops = pos.Bitboard(WhiteBishop)
		allies = pos.ColorBitboard(White)
	case Black:
		bishops = pos.Bitboard(BlackBishop)
		allies = pos.ColorBitboard(Black)
	default:
		return moves
	}

	for bishops != 0 {
		singleBishopIndex := bits.TrailingZeros64(uint64(bishops))
		singleBishopBitboard := Bitboard(1 << singleBishopIndex)
		bishops ^= singleBishopBitboard
		bishopAttacks := singleBishopBitboard.BishopAttacks(occupied)
		bishopAttacks &^= allies
		singleBishopSquare := indexToSquare(singleBishopIndex)
		for bishopAttacks != 0 {
			attackIndex := bits.TrailingZeros64(uint64(bishopAttacks))
			bishopAttacks ^= 1 << attackIndex
			attackSquare := indexToSquare(attackIndex)
			moves = append(moves, Move{singleBishopSquare, attackSquare, NoPieceType})
		}
	}

	return moves
}

func queenMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	var queens Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		queens = pos.Bitboard(WhiteQueen)
		allies = pos.ColorBitboard(White)
	case Black:
		queens = pos.Bitboard(BlackQueen)
		allies = pos.ColorBitboard(Black)
	default:
		return moves
	}

	for queens != 0 {
		singleQueenIndex := bits.TrailingZeros64(uint64(queens))
		singleQueenBitboard := Bitboard(1 << singleQueenIndex)
		queens ^= singleQueenBitboard
		// Combine rook and bishop attacks for queen moves
		queenAttacks := singleQueenBitboard.RookAttacks(occupied) | singleQueenBitboard.BishopAttacks(occupied)
		queenAttacks &^= allies
		singleQueenSquare := indexToSquare(singleQueenIndex)
		for queenAttacks != 0 {
			attackIndex := bits.TrailingZeros64(uint64(queenAttacks))
			queenAttacks ^= 1 << attackIndex
			attackSquare := indexToSquare(attackIndex)
			moves = append(moves, Move{singleQueenSquare, attackSquare, NoPieceType})
		}
	}

	return moves
}

func kingMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	var kings Bitboard
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		kings = pos.Bitboard(WhiteKing)
		allies = pos.ColorBitboard(White)
	case Black:
		kings = pos.Bitboard(BlackKing)
		allies = pos.ColorBitboard(Black)
	default:
		return moves
	}

	for kings != 0 {
		singleKingIndex := bits.TrailingZeros64(uint64(kings))
		singleKingBitboard := Bitboard(1 << singleKingIndex)
		kings ^= singleKingBitboard
		kingAttacks := singleKingBitboard.KingAttacks()
		kingAttacks &^= allies
		singleKingSquare := indexToSquare(singleKingIndex)
		for kingAttacks != 0 {
			attackIndex := bits.TrailingZeros64(uint64(kingAttacks))
			kingAttacks ^= 1 << attackIndex
			attackSquare := indexToSquare(attackIndex)
			moves = append(moves, Move{singleKingSquare, attackSquare, NoPieceType})
		}
	}

	moves = append(moves, castleMoves(pos)...)

	return moves
}

func castleMoves(pos *Position) []Move {
	moves := make([]Move, 0)
	occupied := pos.OccupiedBitboard()

	if pos.SideToMove == White {
		attacked := pos.getAttackedSquares(Black)
		if pos.WhiteKsCastle &&
			pos.Piece(E1) == WhiteKing &&
			pos.Piece(H1) == WhiteRook &&
			occupied.Square(F1) == 0 &&
			occupied.Square(G1) == 0 &&
			attacked.Square(E1) == 0 &&
			attacked.Square(F1) == 0 &&
			attacked.Square(G1) == 0 {
			moves = append(moves, Move{E1, G1, NoPieceType})
		}
		if pos.WhiteQsCastle &&
			pos.Piece(E1) == WhiteKing &&
			pos.Piece(A1) == WhiteRook &&
			occupied.Square(B1) == 0 &&
			occupied.Square(C1) == 0 &&
			occupied.Square(D1) == 0 &&
			attacked.Square(E1) == 0 &&
			attacked.Square(D1) == 0 &&
			attacked.Square(C1) == 0 {
			moves = append(moves, Move{E1, C1, NoPieceType})
		}
	} else if pos.SideToMove == Black {
		attacked := pos.getAttackedSquares(White)
		if pos.BlackKsCastle &&
			pos.Piece(E8) == BlackKing &&
			pos.Piece(H8) == BlackRook &&
			occupied.Square(F8) == 0 &&
			occupied.Square(G8) == 0 &&
			attacked.Square(E8) == 0 &&
			attacked.Square(F8) == 0 &&
			attacked.Square(G8) == 0 {
			moves = append(moves, Move{E8, G8, NoPieceType})
		}
		if pos.BlackQsCastle &&
			pos.Piece(E8) == BlackKing &&
			pos.Piece(A8) == BlackRook &&
			occupied.Square(B8) == 0 &&
			occupied.Square(C8) == 0 &&
			occupied.Square(D8) == 0 &&
			attacked.Square(E8) == 0 &&
			attacked.Square(D8) == 0 &&
			attacked.Square(C8) == 0 {
			moves = append(moves, Move{E8, C8, NoPieceType})
		}
	}

	return moves
}

// LegalMoves returns all legal moves for pos. Returns nil if moves could not be generated (for example if pos.SideToMove was not set). Returns an empty slice if move generation was successful, but no moves were found.
func LegalMoves(pos *Position) []Move {
	pseudoLegalMoves := PseudoLegalMoves(pos)
	legalMoves := make([]Move, 0)
	for _, m := range pseudoLegalMoves {
		tempPos := pos.Copy()
		tempPos.Move(m)
		tempPos.SideToMove = pos.SideToMove
		if !tempPos.IsCheck() {
			legalMoves = append(legalMoves, m)
		}
	}
	return legalMoves
}
