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
	"slices"
)

// PseudoLegalMoves are moves that are legal except they leave one's king in check. Returns nil if moves could not be generated (for example if pos.SideToMove was not set). Returns an empty slice if move generation was successful, but no moves were found.
func PseudoLegalMoves(pos *Position) []Move {
	if pos.SideToMove != White && pos.SideToMove != Black {
		return nil
	}
	pawnMoves := pawnMoves(pos)
	rookMoves := rookMoves(pos)
	knightMoves := knightMoves(pos)
	bishopMoves := bishopMoves(pos)
	queenMoves := queenMoves(pos)
	kingMoves := kingMoves(pos)
	moves := make([]Move, 0, len(pawnMoves)+len(rookMoves)+len(knightMoves)+len(bishopMoves)+len(queenMoves)+len(kingMoves))
	moves = append(moves, pawnMoves...)
	moves = append(moves, rookMoves...)
	moves = append(moves, knightMoves...)
	moves = append(moves, bishopMoves...)
	moves = append(moves, queenMoves...)
	moves = append(moves, kingMoves...)
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
	occupied := pos.OccupiedBitboard()
	enemies := pos.ColorBitboard(Black) | (1 << squareToIndex(pos.EnPassant))

	moves := make([]Move, 0, 32)

	// Forward 1 square
	moveForward := (pos.whitePawns << 8) &^ occupied
	for mf := moveForward; mf != 0; {
		squareIndex := bits.TrailingZeros64(uint64(mf))
		mf ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File, to.Rank - 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Forward 2 squares
	moveForward2 := (pos.whitePawns << 16) &^ occupied
	for mf2 := moveForward2; mf2 != 0; {
		squareIndex := bits.TrailingZeros64(uint64(mf2))
		mf2 ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		mid := Square{to.File, to.Rank - 1}
		from := Square{to.File, to.Rank - 2}
		if occupied.Square(mid) == 0 && from.Rank == Rank2 {
			moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
		}
	}

	// Take NE
	neAttacks := pos.whitePawns.pawnAttacksNE() & enemies
	for ne := neAttacks; ne != 0; {
		squareIndex := bits.TrailingZeros64(uint64(ne))
		ne ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File - 1, to.Rank - 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Take NW
	nwAttacks := pos.whitePawns.pawnAttacksNW() & enemies
	for nw := nwAttacks; nw != 0; {
		squareIndex := bits.TrailingZeros64(uint64(nw))
		nw ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File + 1, to.Rank - 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Handle promotions
	numMoves := len(moves)
	for i := 0; i < numMoves; i++ {
		if moves[i].ToSquare.Rank == Rank8 {
			base := moves[i]
			moves[i].Promotion = Queen

			for _, promo := range []PieceType{Rook, Knight, Bishop} {
				copy := base
				copy.Promotion = promo
				moves = append(moves, copy)
			}
		}
	}

	return moves
}

func blackPawnMoves(pos *Position) []Move {
	occupied := pos.OccupiedBitboard()
	enemies := pos.ColorBitboard(White) | (1 << squareToIndex(pos.EnPassant))

	moves := make([]Move, 0, 32)

	// Forward 1 square
	moveForward := (pos.blackPawns >> 8) &^ occupied
	for mf := moveForward; mf != 0; {
		squareIndex := bits.TrailingZeros64(uint64(mf))
		mf ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File, to.Rank + 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Forward 2 squares
	moveForward2 := (pos.blackPawns >> 16) &^ occupied
	for mf2 := moveForward2; mf2 != 0; {
		squareIndex := bits.TrailingZeros64(uint64(mf2))
		mf2 ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		mid := Square{to.File, to.Rank + 1}
		from := Square{to.File, to.Rank + 2}
		if occupied.Square(mid) == 0 && from.Rank == Rank7 {
			moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
		}
	}

	// Take SE
	seAttacks := pos.blackPawns.pawnAttacksSE() & enemies
	for se := seAttacks; se != 0; {
		squareIndex := bits.TrailingZeros64(uint64(se))
		se ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File - 1, to.Rank + 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Take SW
	swAttacks := pos.blackPawns.pawnAttacksSW() & enemies
	for sw := swAttacks; sw != 0; {
		squareIndex := bits.TrailingZeros64(uint64(sw))
		sw ^= 1 << squareIndex
		to := indexToSquare(squareIndex)
		from := Square{to.File + 1, to.Rank + 1}
		moves = append(moves, Move{FromSquare: from, ToSquare: to, Promotion: NoPieceType})
	}

	// Promotions
	numMoves := len(moves)
	for i := 0; i < numMoves; i++ {
		if moves[i].ToSquare.Rank == Rank1 {
			base := moves[i]
			moves[i].Promotion = Queen

			for _, promo := range []PieceType{Rook, Knight, Bishop} {
				copy := base
				copy.Promotion = promo
				moves = append(moves, copy)
			}
		}
	}

	return moves
}

func rookMoves(pos *Position) []Move {
	var rooks Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		rooks = pos.whiteRooks
		allies = pos.ColorBitboard(White)
	case Black:
		rooks = pos.blackRooks
		allies = pos.ColorBitboard(Black)
	default:
		return []Move{}
	}
	moves := make([]Move, 0, bits.OnesCount64(uint64(rooks))*10)

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
	var knights Bitboard
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		knights = pos.whiteKnights
		allies = pos.ColorBitboard(White)
	case Black:
		knights = pos.blackKnights
		allies = pos.ColorBitboard(Black)
	default:
		return []Move{}
	}
	moves := make([]Move, 0, bits.OnesCount64(uint64(knights))*8)

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
	var bishops Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		bishops = pos.whiteBishops
		allies = pos.ColorBitboard(White)
	case Black:
		bishops = pos.blackBishops
		allies = pos.ColorBitboard(Black)
	default:
		return []Move{}
	}
	moves := make([]Move, 0, bits.OnesCount64(uint64(bishops))*10)

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
	var queens Bitboard
	var occupied Bitboard = pos.OccupiedBitboard()
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		queens = pos.whiteQueens
		allies = pos.ColorBitboard(White)
	case Black:
		queens = pos.blackQueens
		allies = pos.ColorBitboard(Black)
	default:
		return []Move{}
	}
	moves := make([]Move, 0, bits.OnesCount64(uint64(queens))*15)

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
	moves := make([]Move, 0, 10) // 10 is the most moves a king can ever make
	var kings Bitboard
	var allies Bitboard
	switch pos.SideToMove {
	case White:
		kings = pos.whiteKings
		allies = pos.ColorBitboard(White)
	case Black:
		kings = pos.blackKings
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

	switch pos.SideToMove {
	case White:
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
	case Black:
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
	for i := 0; i < len(pseudoLegalMoves); i++ {
		tempPos := pos.Copy()
		tempPos.Move(pseudoLegalMoves[i])
		tempPos.SideToMove = pos.SideToMove
		if tempPos.IsCheck() {
			pseudoLegalMoves = slices.Delete(pseudoLegalMoves, i, i+1)
			i--
		}
	}
	return pseudoLegalMoves
}
