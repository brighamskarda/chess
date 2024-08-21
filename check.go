package chess

func IsCheck(p *Position) bool {
	return true
}

func IsCheckMate(p *Position) bool {
	return true
}

// IsStaleMate does not check the fifty rule move. It only checks if a player is not able to move, and is not in check.
func IsStaleMate(p *Position) bool {
	return true
}
