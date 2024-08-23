package chess

// GeneratePseudoLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development.
func GeneratePseudoLegalMoves(p *Position) []Move {
	pseudoLegalMoves := []Move{}
	pseudoLegalMoves = append(pseudoLegalMoves, generatePawnMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateRookMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateKnightMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateBishopMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateQueenMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateKingMoves(p)...)
	pseudoLegalMoves = append(pseudoLegalMoves, generateCastleMoves(p)...)
	return pseudoLegalMoves
}

func generatePawnMoves(p *Position) []Move {
	if p.Turn == White {
		return generateWhitePawnMoves(p)
	} else if p.Turn == Black {
		return generateBlackPawnMoves(p)
	} else {
		return []Move{}
	}
}

func generateWhitePawnMoves(p *Position) []Move {
	return []Move{}
}

func generateBlackPawnMoves(p *Position) []Move {
	return []Move{}
}

func generateRookMoves(p *Position) []Move {
	return []Move{}
}

func generateKnightMoves(p *Position) []Move {
	return []Move{}
}

func generateBishopMoves(p *Position) []Move {
	return []Move{}
}

func generateQueenMoves(p *Position) []Move {
	return []Move{}
}

func generateKingMoves(p *Position) []Move {
	return []Move{}
}

func generateCastleMoves(p *Position) []Move {
	return []Move{}
}

// GenerateLegalMoves expects a valid position. Behavior is undefined for invalid positions. This is to improve
// performance since move generation is a vital part to engine development.
func GenerateLegalMoves(p *Position) []Move {
	return []Move{}
}
