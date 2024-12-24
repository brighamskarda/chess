package chess

// IsCheck returns true if the side to move is currently in check.
func (p *Position) IsCheck() bool {
	if !isValidColor(p.Turn) || p.Turn == NoColor {
		return false
	}
	kingSquare := findKing(p, p.Turn)
	if kingSquare == NoSquare {
		return false
	}

	return isCheckPawn(p, kingSquare) ||
		isCheckRookQueen(p, kingSquare) ||
		isCheckKnight(p, kingSquare) ||
		isCheckBishopQueen(p, kingSquare) ||
		isCheckKing(p, kingSquare)
}

func isCheckPawn(p *Position, kingSquare Square) bool {
	if p.Turn == Black {
		return isCheckWhitePawn(p, kingSquare)
	}
	return isCheckBlackPawn(p, kingSquare)
}

func isCheckWhitePawn(p *Position, kingSquare Square) bool {
	squareToCheck := kingSquare
	squareToCheck.File--
	squareToCheck.Rank--
	if p.PieceAt(squareToCheck) == WhitePawn {
		return true
	}
	squareToCheck.File += 2
	if p.PieceAt(squareToCheck) == WhitePawn {
		return true
	}
	return false
}

func isCheckBlackPawn(p *Position, kingSquare Square) bool {
	squareToCheck := kingSquare
	squareToCheck.File--
	squareToCheck.Rank++
	if p.PieceAt(squareToCheck) == BlackPawn {
		return true
	}
	squareToCheck.File += 2
	if p.PieceAt(squareToCheck) == BlackPawn {
		return true
	}
	return false
}

func isCheckRookQueen(p *Position, kingSquare Square) bool {
	for squareToTest := squareToLeft(kingSquare); squareToTest != NoSquare; squareToTest = squareToLeft(squareToTest) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Rook) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareToRight(kingSquare); squareToTest != NoSquare; squareToTest = squareToRight(squareToTest) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Rook) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareAbove(kingSquare); squareToTest != NoSquare; squareToTest = squareAbove(squareToTest) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Rook) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareBelow(kingSquare); squareToTest != NoSquare; squareToTest = squareBelow(squareToTest) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Rook) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}
	return false
}

func isCheckKnight(p *Position, kingSquare Square) bool {
	squareToTest := squareAbove(squareAbove(squareToLeft(kingSquare)))
	piece := p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareAbove(squareAbove(squareToRight(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareAbove(squareToRight(squareToRight(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(squareToRight(squareToRight(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(squareBelow(squareToRight(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(squareBelow(squareToLeft(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(squareToLeft(squareToLeft(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareAbove(squareToLeft(squareToLeft(kingSquare)))
	piece = p.PieceAt(squareToTest)
	if piece.Type == Knight && piece.Color != p.Turn {
		return true
	}

	return false
}

func isCheckBishopQueen(p *Position, kingSquare Square) bool {
	for squareToTest := squareAbove(squareToLeft(kingSquare)); squareToTest != NoSquare; squareToTest = squareAbove(squareToLeft(squareToTest)) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Bishop) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareAbove(squareToRight(kingSquare)); squareToTest != NoSquare; squareToTest = squareAbove(squareToRight(squareToTest)) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Bishop) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareBelow(squareToRight(kingSquare)); squareToTest != NoSquare; squareToTest = squareBelow(squareToRight(squareToTest)) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Bishop) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	for squareToTest := squareBelow(squareToLeft(kingSquare)); squareToTest != NoSquare; squareToTest = squareBelow(squareToLeft(squareToTest)) {
		piece := p.PieceAt(squareToTest)
		if (piece.Type == Queen || piece.Type == Bishop) && piece.Color != p.Turn {
			return true
		}
		if piece.Type != NoPieceType {
			break
		}
	}

	return false
}

func isCheckKing(p *Position, kingSquare Square) bool {
	squareToTest := squareAbove(kingSquare)
	piece := p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareToRight(squareAbove(kingSquare))
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareToRight(kingSquare)
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(squareToRight(kingSquare))
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareBelow(kingSquare)
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareToLeft(squareBelow(kingSquare))
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareToLeft(kingSquare)
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}

	squareToTest = squareAbove(squareToLeft(kingSquare))
	piece = p.PieceAt(squareToTest)
	if piece.Type == King && piece.Color != p.Turn {
		return true
	}
	return false
}

// IsCheckMate returns true is the side to move is in check and has no legal moves.
func IsCheckMate(p *Position) bool {
	return p.IsCheck() && len(GenerateLegalMoves(p)) == 0
}

// IsStaleMate does not check the fifty move rule. It only checks if a player is not able to move, and is not in check.
func IsStaleMate(p *Position) bool {
	return !p.IsCheck() && len(GenerateLegalMoves(p)) == 0
}
