package chess

// GeneratePseudoLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development. Pseudo-legal moves are moves that would be
// legal if one's own king were not in check. Castling moves that are directly blocked by a piece are not included,
// but castling moves that would be stopped because the kings path is in check are included.
func (p *Position) GeneratePseudoLegalMoves() []Move {
	pseudoLegalMoves := []Move{}
	for index, piece := range p.Board {
		if piece.Color == p.Turn {
			switch piece.Type {
			case Pawn:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generatePawnMoves(indexToSquare(index))...)
			case Rook:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateRookMoves(indexToSquare(index))...)
			case Knight:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateKnightMoves(indexToSquare(index))...)
			case Bishop:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateBishopMoves(indexToSquare(index))...)
			case Queen:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateQueenMoves(indexToSquare(index))...)
			case King:
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateKingMoves(indexToSquare(index))...)
				pseudoLegalMoves = append(pseudoLegalMoves, p.generateCastleMoves(indexToSquare(index))...)
			}
		}
	}
	return pseudoLegalMoves
}

func (p *Position) generatePawnMoves(s Square) []Move {
	if p.Turn == White {
		return p.generateWhitePawnMoves(s)
	} else if p.Turn == Black {
		return p.generateBlackPawnMoves(s)
	} else {
		return []Move{}
	}
}

func (p *Position) generateWhitePawnMoves(s Square) []Move {
	pawnMoves := []Move{}
	pawnMoves = append(pawnMoves, p.generateWhitePawnSingleMoveForward(s)...)
	pawnMoves = append(pawnMoves, p.generateWhitePawnDoubleMoveForward(s)...)
	pawnMoves = append(pawnMoves, p.generateWhitePawnTakesPiece(s)...)
	pawnMoves = append(pawnMoves, p.generateWhitePawnTakesEnPassant(s)...)
	pawnMoves = append(pawnMoves, p.generateWhitePawnPromotion(s)...)
	return pawnMoves
}

func (p *Position) generateWhitePawnSingleMoveForward(s Square) []Move {
	if s.Rank < Rank7 && p.PieceAt(Square{s.File, s.Rank + 1}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank + 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateWhitePawnDoubleMoveForward(s Square) []Move {
	if s.Rank == Rank2 && p.PieceAt(Square{s.File, s.Rank + 1}) == NoPiece && p.PieceAt(Square{s.File, s.Rank + 2}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank + 2}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateWhitePawnTakesPiece(s Square) []Move {
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

func (p *Position) generateWhitePawnTakesEnPassant(s Square) []Move {
	if p.EnPassant == (Square{s.File - 1, s.Rank + 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank + 1}, Promotion: NoPieceType}}
	}
	if p.EnPassant == (Square{s.File + 1, s.Rank + 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank + 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateWhitePawnPromotion(s Square) []Move {
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

func (p *Position) generateBlackPawnMoves(s Square) []Move {
	pawnMoves := []Move{}
	pawnMoves = append(pawnMoves, p.generateBlackPawnSingleMoveForward(s)...)
	pawnMoves = append(pawnMoves, p.generateBlackPawnDoubleMoveForward(s)...)
	pawnMoves = append(pawnMoves, p.generateBlackPawnTakesPiece(s)...)
	pawnMoves = append(pawnMoves, p.generateBlackPawnTakesEnPassant(s)...)
	pawnMoves = append(pawnMoves, p.generateBlackPawnPromotion(s)...)
	return pawnMoves
}

func (p *Position) generateBlackPawnSingleMoveForward(s Square) []Move {
	if s.Rank > Rank2 && p.PieceAt(Square{s.File, s.Rank - 1}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank - 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateBlackPawnDoubleMoveForward(s Square) []Move {
	if s.Rank == Rank7 && p.PieceAt(Square{s.File, s.Rank - 1}) == NoPiece && p.PieceAt(Square{s.File, s.Rank - 2}) == NoPiece {
		return []Move{{FromSquare: s, ToSquare: Square{s.File, s.Rank - 2}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateBlackPawnTakesPiece(s Square) []Move {
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

func (p *Position) generateBlackPawnTakesEnPassant(s Square) []Move {
	if p.EnPassant == (Square{s.File - 1, s.Rank - 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File - 1, s.Rank - 1}, Promotion: NoPieceType}}
	}
	if p.EnPassant == (Square{s.File + 1, s.Rank - 1}) {
		return []Move{{FromSquare: s, ToSquare: Square{s.File + 1, s.Rank - 1}, Promotion: NoPieceType}}
	}
	return []Move{}
}

func (p *Position) generateBlackPawnPromotion(s Square) []Move {
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

func (p *Position) generateRookMoves(s Square) []Move {
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

func (p *Position) generateKnightMoves(s Square) []Move {
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

func (p *Position) generateBishopMoves(s Square) []Move {
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

func (p *Position) generateQueenMoves(s Square) []Move {
	moves := []Move{}
	moves = append(moves, p.generateRookMoves(s)...)
	moves = append(moves, p.generateBishopMoves(s)...)
	return moves
}

func (p *Position) generateKingMoves(s Square) []Move {
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

func (p *Position) generateCastleMoves(s Square) []Move {
	if p.Turn == White && s == E1 {
		return p.generateWhiteCastleMoves()
	}
	if p.Turn == Black && s == E8 {
		return p.generateBlackCastleMoves()
	}
	return []Move{}
}

func (p *Position) generateWhiteCastleMoves() []Move {
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

func (p *Position) generateBlackCastleMoves() []Move {
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
func (p *Position) GenerateLegalMoves() []Move {
	pseudoLegalMoves := p.GeneratePseudoLegalMoves()
	isCurrentPositionCheck := p.IsCheck()
	legalMoves := []Move{}
	for _, move := range pseudoLegalMoves {
		var tempPosition Position = *p
		tempPosition.Move(move)
		tempPosition.Turn = p.Turn
		castleMove := isCastleMove(p, move)
		if !tempPosition.IsCheck() && ((castleMove && !isCurrentPositionCheck) || !castleMove) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}
