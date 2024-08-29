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
	return []Move{}
}

func generateBlackPawnMoves(p *Position, s Square) []Move {
	return []Move{}
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
