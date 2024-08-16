package chess

import "strings"

// Position represents a chess position as described by Forsyth-Edwards Notation (FEN).
// Board is the actual representation of the pieces on the squares. It starts at A8 and moves left
// to right, top to bottom all the way to H1.
type Position struct {
	Board                [64]Piece
	Turn                 Color
	WhiteKingSideCastle  bool
	WhiteQueenSideCastle bool
	BlackKingSideCastle  bool
	BlackQueenSideCastle bool
	EnPassant            *Square
	HalfMove             uint16
	FullMove             uint16
}

func (p *Position) String() string {
	str := strings.Builder{}
	rank := '8'
	for index, piece := range p.Board {
		if index%8 == 0 {
			str.WriteRune(rank)
			rank -= 1
		}
		str.WriteString(piece.String())
		if index%8 == 7 {
			str.WriteRune('\n')
		}
	}
	str.WriteString(" ABCDEFGH")
	return str.String()
}

func (p *Position) PieceAt(s Square) Piece {

	return p.Board[squareToIndex(s)]
}

func (p *Position) SetPieceAt(s Square, piece Piece) {
	p.Board[squareToIndex(s)] = piece
}

func squareToIndex(s Square) int {
	index := 0
	index += int(s.File)
	index += int(Rank8-s.Rank) * 8
	return index
}
