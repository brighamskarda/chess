package chess

// Bitboard is a partial representation of a chess board where each bit signifies the presence or absence of a piece.
// Square A1 is represented by the most significant bit (on little-endian systems), and h8 is the least significant bit.
type BitBoard uint64

// Board is a representation of all the pieces on chess board using a [BitBoard] for each kind of piece.
// It is recommended to use [Position] unless you have a specific reason for using bitboards.
type Board struct {
	WhiteKings   BitBoard
	BlackKings   BitBoard
	WhiteQueens  BitBoard
	BlackQueens  BitBoard
	WhiteBishops BitBoard
	BlackBishops BitBoard
	WhiteKnights BitBoard
	BlackKnights BitBoard
	WhiteRooks   BitBoard
	BlackRooks   BitBoard
	WhitePawns   BitBoard
	BlackPawns   BitBoard
}

// GetBitBoard returns a bitBoard for the piece specified.
func (pos *Position) GetBitBoard(piece Piece) BitBoard {
	var bb BitBoard
	var bitPosition int8 = 63
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			if pos.PieceAt(Square{Rank: rank, File: file}) == piece {
				bb |= (1 << bitPosition)
			}
			bitPosition--
		}
	}
	return bb
}

// GetBoard returns a struct with bitboards for each of the piece types.
func (pos *Position) GetBoard() Board {
	board := Board{}
	var bitPosition int8 = 63
	// TODO small speed improvement here by just using a number instead of generating square structs
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			switch pos.PieceAt(Square{Rank: rank, File: file}) {
			case WhiteKing:
				board.WhiteKings |= (1 << bitPosition)
			case WhiteQueen:
				board.WhiteQueens |= (1 << bitPosition)
			case WhiteRook:
				board.WhiteRooks |= (1 << bitPosition)
			case WhiteBishop:
				board.WhiteBishops |= (1 << bitPosition)
			case WhiteKnight:
				board.WhiteKnights |= (1 << bitPosition)
			case WhitePawn:
				board.WhitePawns |= (1 << bitPosition)
			case BlackKing:
				board.BlackKings |= (1 << bitPosition)
			case BlackQueen:
				board.BlackQueens |= (1 << bitPosition)
			case BlackRook:
				board.BlackRooks |= (1 << bitPosition)
			case BlackBishop:
				board.BlackBishops |= (1 << bitPosition)
			case BlackKnight:
				board.BlackKnights |= (1 << bitPosition)
			case BlackPawn:
				board.BlackPawns |= (1 << bitPosition)
			}

			bitPosition--
		}
	}
	return board
}
