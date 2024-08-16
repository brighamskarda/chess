package chess

type Color uint8

const (
	NONE Color = iota
	WHITE
	BLACK
)

func (c Color) String() string {
	switch c {
	case NONE:
		return "NONE"
	case WHITE:
		return "WHITE"
	case BLACK:
		return "BLACK"
	default:
		return "INVALID COLOR"
	}
}
