package chess

// GeneratePseudoLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development.
func GeneratePseudoLegalMoves(p *Position) []Move {
	pseudoLegalMoves := []Move{}
	for index, piece := range p.Board {
		if piece.Color == p.Turn {
			switch piece.Type {
			case Pawn:
				pseudoLegalMoves = append(pseudoLegalMoves, generatePawnMoves(p, indexToSquare(index))...)
			case Rook:
				pseudoLegalMoves = append(pseudoLegalMoves, generateRookMoves(p, indexToSquare(index))...)
			case Knight:
				pseudoLegalMoves = append(pseudoLegalMoves, generateKnightMoves(p, indexToSquare(index))...)
			case Bishop:
				pseudoLegalMoves = append(pseudoLegalMoves, generateBishopMoves(p, indexToSquare(index))...)
			case Queen:
				pseudoLegalMoves = append(pseudoLegalMoves, generateQueenMoves(p, indexToSquare(index))...)
			case King:
				pseudoLegalMoves = append(pseudoLegalMoves, generateKingMoves(p, indexToSquare(index))...)
				pseudoLegalMoves = append(pseudoLegalMoves, generateCastleMoves(p, indexToSquare(index))...)
			}
		}
	}
	return pseudoLegalMoves
}

func generatePawnMoves(p *Position, s Square) []Move {
	if p.Turn == White {
		return generateWhitePawnMoves(p, s)
	} else if p.Turn == Black {
		return generateBlackPawnMoves(p, s)
	} else {
		return []Move{}
	}
}

func generateWhitePawnMoves(p *Position, s Square) []Move {
	pawnMoves := []Move{}
	pawnMoves = append(pawnMoves, generateWhitePawnSingleMoveForward(p, s)...)
	pawnMoves = append(pawnMoves, generateWhitePawnDoubleMoveForward(p, s)...)
	pawnMoves = append(pawnMoves, generateWhitePawnTakesPiece(p, s)...)
	pawnMoves = append(pawnMoves, generateWhitePawnTakesEnPassant(p, s)...)
	pawnMoves = append(pawnMoves, generateWhitePawnPromotion(p, s)...)
	return pawnMoves
}

func generateWhitePawnSingleMoveForward(p *Position, s Square) []Move {
	if s.Rank < Rank7 && p.PieceAt(Square{s.File, s.Rank + 1}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank + 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateWhitePawnDoubleMoveForward(p *Position, s Square) []Move {
	if s.Rank == Rank2 && p.PieceAt(Square{s.File, s.Rank + 1}) == NoPiece && p.PieceAt(Square{s.File, s.Rank + 2}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank + 2}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateWhitePawnTakesPiece(p *Position, s Square) []Move {
	moves := []Move{}
	pieceUpLeft := p.PieceAt(Square{s.File - 1, s.Rank + 1})
	if s.Rank < Rank7 && pieceUpLeft.Color != NoColor && pieceUpLeft.Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank + 1}, Promotion: NoPieceType})
	}
	pieceUpRight := p.PieceAt(Square{s.File + 1, s.Rank + 1})
	if s.Rank < Rank7 && pieceUpRight.Color != NoColor && pieceUpRight.Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank + 1}, Promotion: NoPieceType})
	}
	return moves
}

func generateWhitePawnTakesEnPassant(p *Position, s Square) []Move {
	if p.EnPassant == (Square{s.File - 1, s.Rank + 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank + 1}, Promotion: NoPieceType}}
	}
	if p.EnPassant == (Square{s.File + 1, s.Rank + 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank + 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateWhitePawnPromotion(p *Position, s Square) []Move {
	movesBeforePromotion := []Move{}
	if s.Rank == Rank7 && p.PieceAt(Square{s.File, s.Rank + 1}) == NoPiece {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File, s.Rank + 1}, Promotion: NoPieceType})
	}
	pieceUpLeft := p.PieceAt(Square{s.File - 1, s.Rank + 1})
	if s.Rank == Rank7 && pieceUpLeft.Color != NoColor && pieceUpLeft.Color != p.Turn {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank + 1}, Promotion: NoPieceType})
	}
	pieceUpRight := p.PieceAt(Square{s.File + 1, s.Rank + 1})
	if s.Rank == Rank7 && pieceUpRight.Color != NoColor && pieceUpRight.Color != p.Turn {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank + 1}, Promotion: NoPieceType})
	}
	movesWithPromotion := []Move{}
	for _, move := range movesBeforePromotion {
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Rook})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Knight})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Bishop})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Queen})
	}
	return movesWithPromotion
}

func generateBlackPawnMoves(p *Position, s Square) []Move {
	pawnMoves := []Move{}
	pawnMoves = append(pawnMoves, generateBlackPawnSingleMoveForward(p, s)...)
	pawnMoves = append(pawnMoves, generateBlackPawnDoubleMoveForward(p, s)...)
	pawnMoves = append(pawnMoves, generateBlackPawnTakesPiece(p, s)...)
	pawnMoves = append(pawnMoves, generateBlackPawnTakesEnPassant(p, s)...)
	pawnMoves = append(pawnMoves, generateBlackPawnPromotion(p, s)...)
	return pawnMoves
}

func generateBlackPawnSingleMoveForward(p *Position, s Square) []Move {
	if s.Rank > Rank2 && p.PieceAt(Square{s.File, s.Rank - 1}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank - 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateBlackPawnDoubleMoveForward(p *Position, s Square) []Move {
	if s.Rank == Rank7 && p.PieceAt(Square{s.File, s.Rank - 1}) == NoPiece && p.PieceAt(Square{s.File, s.Rank - 2}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank - 2}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateBlackPawnTakesPiece(p *Position, s Square) []Move {
	moves := []Move{}
	pieceUpLeft := p.PieceAt(Square{s.File - 1, s.Rank - 1})
	if s.Rank > Rank2 && pieceUpLeft.Color != NoColor && pieceUpLeft.Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank - 1}, Promotion: NoPieceType})
	}
	pieceUpRight := p.PieceAt(Square{s.File + 1, s.Rank - 1})
	if s.Rank > Rank2 && pieceUpRight.Color != NoColor && pieceUpRight.Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank - 1}, Promotion: NoPieceType})
	}
	return moves
}

func generateBlackPawnTakesEnPassant(p *Position, s Square) []Move {
	if p.EnPassant == (Square{s.File - 1, s.Rank - 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank - 1}, Promotion: NoPieceType}}
	}
	if p.EnPassant == (Square{s.File + 1, s.Rank - 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank - 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func generateBlackPawnPromotion(p *Position, s Square) []Move {
	movesBeforePromotion := []Move{}
	if s.Rank == Rank2 && p.PieceAt(Square{s.File, s.Rank - 1}) == NoPiece {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File, s.Rank - 1}, Promotion: NoPieceType})
	}
	pieceUpLeft := p.PieceAt(Square{s.File - 1, s.Rank - 1})
	if s.Rank == Rank2 && pieceUpLeft.Color != NoColor && pieceUpLeft.Color != p.Turn {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank - 1}, Promotion: NoPieceType})
	}
	pieceUpRight := p.PieceAt(Square{s.File + 1, s.Rank - 1})
	if s.Rank == Rank2 && pieceUpRight.Color != NoColor && pieceUpRight.Color != p.Turn {
		movesBeforePromotion = append(movesBeforePromotion, Move{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank - 1}, Promotion: NoPieceType})
	}
	movesWithPromotion := []Move{}
	for _, move := range movesBeforePromotion {
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Rook})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Knight})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Bishop})
		movesWithPromotion = append(movesWithPromotion, Move{move.FromSquare, move.ToSquare, Queen})
	}
	return movesWithPromotion
}

func generateRookMoves(p *Position, s Square) []Move {
	return []Move{}
}

func generateKnightMoves(p *Position, s Square) []Move {
	return []Move{}
}

func generateBishopMoves(p *Position, s Square) []Move {
	return []Move{}
}

func generateQueenMoves(p *Position, s Square) []Move {
	return []Move{}
}

func generateKingMoves(p *Position, s Square) []Move {
	return []Move{}
}

func generateCastleMoves(p *Position, s Square) []Move {
	return []Move{}
}

// GenerateLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development.
func GenerateLegalMoves(p *Position) []Move {
	return []Move{}
}
