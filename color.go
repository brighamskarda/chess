package chess

type Color uint8

const (
	NoColor Color = iota
	White
	Black
)

// String returns one of NO-COLOR, WHITE, BLACK, or INVALID COLOR
func (c Color) String() string {
	switch c {
	case NoColor:
		return "NO-COLOR"
	case White:
		return "WHITE"
	case Black:
		return "BLACK"
	default:
		return "INVALID COLOR"
	}
}

func isValidColor(c Color) bool {
	return c <= 2
}
