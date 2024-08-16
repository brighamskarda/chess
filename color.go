package chess

type Color uint8

const (
	NoColor Color = iota
	White
	Black
)

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

func colorIsValid(c Color) bool {
	return c <= 2
}
