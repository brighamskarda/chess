package chess

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
}

func (m Move) String() string {
	return ""
}

// ParseMove expects a UCI compatible move string. Format should be Square1Square2Promotion, where promotion is optional.
func ParseMove(s string) (Move, error) {
	return Move{}, nil
}

// IsValidMove makes sure each of its elements are logical. Namely that the squares can be found on a chess board.
func IsValidMove(m Move) bool {
	return isValidSquare(m.FromSquare) && m.FromSquare != NoSquare &&
		isValidSquare(m.ToSquare) && m.ToSquare != NoSquare &&
		isValidPieceType(m.Promotion)
}
