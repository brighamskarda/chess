// Copyright (C) 2025 Brigham Skarda

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package chess

const DefaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Board represents all parts of a chess board as specified by FEN notation.
//
// The zero value is usable, though not very useful. You likely will want to use the following instead:
// 		chess.ParseFEN(DefaultFEN)
type Board struct {
	whitePawns   Bitboard
	whiteRooks   Bitboard
	whiteKnights Bitboard
	whiteBishops Bitboard
	whiteQueens  Bitboard
	whiteKings   Bitboard

	blackPawns   Bitboard
	blackRooks   Bitboard
	blackKnights Bitboard
	blackBishops Bitboard
	blackQueens  Bitboard
	blackKings   Bitboard

	sideToMove Color

	whiteKsCastle bool
	whiteQsCastle bool
	blackKsCastle bool
	blackQsCastle bool

	enPassant Square

	halfMove uint
	fullMove uint
}

// ParseFEN returns an error if it could not parse an FEN. It was likely malformed or missing important pieces.
func ParseFEN(fen string) (*Board, error) {
	return &Board{}, nil
}

// String generates an FEN string for the current board. See [PrettyString] for getting a board like representation.
func (b *Board) String() string {
	return ""
}

// PrettyString returns a board like representation of the current board. Uppercase letters are white and lowercase letters are black.
//
// Set whitesPerspective to true to see the board from white's side. Set extraInfo to false to just see the board. Set extraInfo to true to see all the other information stored in an FEN.
func (b *Board) PrettyString(whitesPerspective bool, extraInfo bool) string {
	return ""
}

func (b *Board) SideToMove() Color {
	return b.sideToMove
}

func (b *Board) SetSideToMove(c Color) {
	b.sideToMove = c
}

// WhiteKingSideCastle returns true if white may still castle kingside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (b *Board) WhiteKingSideCastle() bool {
	return b.whiteKsCastle
}

// WhiteQueenSideCastle returns true if white may still castle queenside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (b *Board) WhiteQueenSideCastle() bool {
	return b.whiteQsCastle
}

// BlackKingSideCastle returns true if black may still castle kingside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (b *Board) BlackKingSideCastle() bool {
	return b.blackKsCastle
}

// BlackQueenSideCastle returns true if black may still castle queenside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (b *Board) BlackQueenSideCastle() bool {
	return b.blackQsCastle
}

func (b *Board) SetWhiteKingSideCastle(value bool) {
	b.whiteKsCastle = value
}

func (b *Board) SetWhiteQueenSideCastle(value bool) {
	b.whiteQsCastle = value
}

func (b *Board) SetBlackKingSideCastle(value bool) {
	b.blackKsCastle = value
}

func (b *Board) SetBlackQueenSideCastle(value bool) {
	b.blackQsCastle = value
}

// EnPassant returns the square on to which a pawn may move to perform en-passant. This does not indicate if the move is legal. NoSquare is returned if there is no en passant option.
func (b *Board) EnPassant() Square {
	return b.enPassant
}

func (b *Board) SetEnPassant(s Square) {
	b.enPassant = s
}

func (b *Board) HalfMove() uint {
	return b.halfMove
}

func (b *Board) SetHalfMove(i uint) {
	b.halfMove = i
}

func (b *Board) FullMove() uint {
	return b.fullMove
}

func (b *Board) SetFullMove(i uint) {
	b.fullMove = i
}

// Piece gets the piece on the given square. NoPiece is returned if no piece is present.
func (b *Board) Piece(s Square) Piece {
	if b.whitePawns.Square(s) == 1 {
		return WhitePawn
	}
	if b.whiteRooks.Square(s) == 1 {
		return WhiteRook
	}
	if b.whiteKnights.Square(s) == 1 {
		return WhiteKnight
	}
	if b.whiteBishops.Square(s) == 1 {
		return WhiteBishop
	}
	if b.whiteQueens.Square(s) == 1 {
		return WhiteQueen
	}
	if b.whiteKings.Square(s) == 1 {
		return WhiteKing
	}

	if b.blackPawns.Square(s) == 1 {
		return BlackPawn
	}
	if b.blackRooks.Square(s) == 1 {
		return BlackRook
	}
	if b.blackKnights.Square(s) == 1 {
		return BlackKnight
	}
	if b.blackBishops.Square(s) == 1 {
		return BlackBishop
	}
	if b.blackQueens.Square(s) == 1 {
		return BlackQueen
	}
	if b.blackKings.Square(s) == 1 {
		return BlackKing
	}

	return NoPiece
}

// SetPiece sets the given piece at the given square.
func (b *Board) SetPiece(p Piece, s Square) {
	b.ClearPiece(s)

	switch p {
	case WhitePawn:
		b.whitePawns = b.whitePawns.SetSquare(s)
	case WhiteRook:
		b.whiteRooks = b.whiteRooks.SetSquare(s)
	case WhiteKnight:
		b.whiteKnights = b.whiteKnights.SetSquare(s)
	case WhiteBishop:
		b.whiteBishops = b.whiteBishops.SetSquare(s)
	case WhiteQueen:
		b.whiteQueens = b.whiteQueens.SetSquare(s)
	case WhiteKing:
		b.whiteKings = b.whiteKings.SetSquare(s)

	case BlackPawn:
		b.blackPawns = b.blackPawns.SetSquare(s)
	case BlackRook:
		b.blackRooks = b.blackRooks.SetSquare(s)
	case BlackKnight:
		b.blackKnights = b.blackKnights.SetSquare(s)
	case BlackBishop:
		b.blackBishops = b.blackBishops.SetSquare(s)
	case BlackQueen:
		b.blackQueens = b.blackQueens.SetSquare(s)
	case BlackKing:
		b.blackKings = b.blackKings.SetSquare(s)
	}
}

// ClearPiece removes any piece from the given square.
func (b *Board) ClearPiece(s Square) {
	b.whitePawns = b.whitePawns.ClearSquare(s)
	b.whiteRooks = b.whiteRooks.ClearSquare(s)
	b.whiteKnights = b.whiteKnights.ClearSquare(s)
	b.whiteBishops = b.whiteBishops.ClearSquare(s)
	b.whiteQueens = b.whiteQueens.ClearSquare(s)
	b.whiteKings = b.whiteKings.ClearSquare(s)

	b.blackPawns = b.blackPawns.ClearSquare(s)
	b.blackRooks = b.blackRooks.ClearSquare(s)
	b.blackKnights = b.blackKnights.ClearSquare(s)
	b.blackBishops = b.blackBishops.ClearSquare(s)
	b.blackQueens = b.blackQueens.ClearSquare(s)
	b.blackKings = b.blackKings.ClearSquare(s)
}

// IsCheck returns true if the side to move has a king under attack from an enemy piece.
func (b *Board) IsCheck() bool {
	return true
}

// Move performs chess moves in such a way that if all moves are legal, the FEN will always be properly updated. The rules it follows are listed below.
//
// 1. By default the following happens:
//
// 1a. The piece at the from square is moved to the to square and promoted. (This also includes moving NoPiece, in which case the promotion is not applied.)
//
// 1b. The half move counter is incremented.
//
// 1c. The side to move is flipped (or set to the opposite of the piece moved if not previously set [stays on NoColor if not set and side to move is not set]).
//
// 1d. If the side to move flips from black to white then the full move counter is incremented.
//
// 1e. En-passant is set to NoSquare.
//
// 2. If a pawn advances, or a piece is taken the half move counter is reset.
//
// 3. If a pawn advances two spaces forward from its starting rank en-passant is set to the square right behind its current position.
//
// 4. If a king or rook moves from their starting square (in standard chess, 960 is not supported) then the corresponding castle rights are set to false.
//
// 5. If one of the four possible castle moves if executed and the castle rights still exist, and there are no pieces in the way, then the appropriate castle move will be applied. (Check will not block a castle move)
func (b *Board) Move(m Move) {

}

// TODO implement Move, IsCheck, PrettyString, String, ParseFEN.
