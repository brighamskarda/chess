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
	moves := []Move{}
	moveFunctions := [4]func(Square) Square{squareToLeft, squareToRight, squareAbove, squareBelow}
	for _, moveFunc := range moveFunctions {
		for toSquare := moveFunc(s); toSquare != NoSquare; toSquare = moveFunc(toSquare) {
			pieceAtSquare := p.PieceAt(toSquare)
			if pieceAtSquare.Color == p.Turn {
				break
			}
			moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
			if pieceAtSquare.Color != NoColor {
				break
			}
		}
	}
	return moves
}

func generateKnightMoves(p *Position, s Square) []Move {
	moves := []Move{}
	toSquare := squareAbove(squareAbove(squareToRight(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareAbove(squareToRight(squareToRight(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(squareToRight(squareToRight(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(squareBelow(squareToRight(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(squareBelow(squareToLeft(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(squareToLeft(squareToLeft(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareAbove(squareToLeft(squareToLeft(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareAbove(squareAbove(squareToLeft(s)))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	return moves
}

func generateBishopMoves(p *Position, s Square) []Move {
	moves := []Move{}
	verticalFunctions := [2]func(Square) Square{squareAbove, squareBelow}
	horizontalFunctions := [2]func(Square) Square{squareToLeft, squareToRight}
	for _, verticalFunc := range verticalFunctions {
		for _, horizontalFunc := range horizontalFunctions {
			for toSquare := verticalFunc(horizontalFunc(s)); toSquare != NoSquare; toSquare = verticalFunc(horizontalFunc(toSquare)) {
				pieceAtSquare := p.PieceAt(toSquare)
				if pieceAtSquare.Color == p.Turn {
					break
				}
				moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
				if pieceAtSquare.Color != NoColor {
					break
				}
			}
		}
	}
	return moves
}

func generateQueenMoves(p *Position, s Square) []Move {
	moves := []Move{}
	moves = append(moves, generateRookMoves(p, s)...)
	moves = append(moves, generateBishopMoves(p, s)...)
	return moves
}

func generateKingMoves(p *Position, s Square) []Move {
	moves := []Move{}
	toSquare := squareAbove(s)
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareToRight(squareAbove(s))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareToRight(s)
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(squareToRight(s))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareBelow(s)
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareToLeft(squareBelow(s))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareToLeft(s)
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	toSquare = squareAbove(squareToLeft(s))
	if toSquare != NoSquare && p.PieceAt(toSquare).Color != p.Turn {
		moves = append(moves, Move{FromSquare: s, ToSquare: toSquare, Promotion: NoPieceType})
	}
	return moves
}

func generateCastleMoves(p *Position, s Square) []Move {
	if p.Turn == White && s == E1 {
		return generateWhiteCastleMoves(p)
	}
	if p.Turn == Black && s == E8 {
		return generateBlackCastleMoves(p)
	}
	return []Move{}
}

func generateWhiteCastleMoves(p *Position) []Move {
	moves := []Move{}
	if p.WhiteKingSideCastle &&
		p.PieceAt(E1) == WhiteKing &&
		p.PieceAt(F1) == NoPiece &&
		p.PieceAt(G1) == NoPiece &&
		p.PieceAt(H1) == WhiteRook {
		moves = append(moves, Move{FromSquare: E1, ToSquare: G1, Promotion: NoPieceType})
	}
	if p.WhiteQueenSideCastle &&
		p.PieceAt(E1) == WhiteKing &&
		p.PieceAt(D1) == NoPiece &&
		p.PieceAt(C1) == NoPiece &&
		p.PieceAt(B1) == NoPiece &&
		p.PieceAt(A1) == WhiteRook {
		moves = append(moves, Move{FromSquare: E1, ToSquare: C1, Promotion: NoPieceType})
	}
	return moves
}

func generateBlackCastleMoves(p *Position) []Move {
	moves := []Move{}
	if p.BlackKingSideCastle &&
		p.PieceAt(E8) == BlackKing &&
		p.PieceAt(F8) == NoPiece &&
		p.PieceAt(G8) == NoPiece &&
		p.PieceAt(H8) == BlackRook {
		moves = append(moves, Move{FromSquare: E8, ToSquare: G8, Promotion: NoPieceType})
	}
	if p.BlackQueenSideCastle &&
		p.PieceAt(E8) == BlackKing &&
		p.PieceAt(D8) == NoPiece &&
		p.PieceAt(C8) == NoPiece &&
		p.PieceAt(B8) == NoPiece &&
		p.PieceAt(A8) == BlackRook {
		moves = append(moves, Move{FromSquare: E8, ToSquare: C8, Promotion: NoPieceType})
	}
	return moves
}

// GenerateLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development.
func GenerateLegalMoves(p *Position) []Move {
	return []Move{}
}
