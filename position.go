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

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const DefaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Position represents all parts of a chess position as specified by FEN notation.
//
// The zero value is usable, though not very useful. You likely will want to use the following instead:
//
//	chess.ParseFEN(DefaultFEN)
type Position struct {
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

// Copy creates a copy of the current position.
func (pos *Position) Copy() *Position {
	newPos := *pos
	return &newPos
}

// Equal returns true if the positions are exactly the same (excluding move counters).
func (pos *Position) Equal(other *Position) bool {
	return pos.whitePawns == other.whitePawns &&
		pos.whiteRooks == other.whiteRooks &&
		pos.whiteKnights == other.whiteKnights &&
		pos.whiteBishops == other.whiteBishops &&
		pos.whiteQueens == other.whiteQueens &&
		pos.whiteKings == other.whiteKings &&
		pos.blackPawns == other.blackPawns &&
		pos.blackRooks == other.blackRooks &&
		pos.blackKnights == other.blackKnights &&
		pos.blackBishops == other.blackBishops &&
		pos.blackQueens == other.blackQueens &&
		pos.blackKings == other.blackKings &&
		pos.sideToMove == other.sideToMove &&
		pos.whiteKsCastle == other.whiteKsCastle &&
		pos.whiteQsCastle == other.whiteQsCastle &&
		pos.blackKsCastle == other.blackKsCastle &&
		pos.blackQsCastle == other.blackQsCastle &&
		pos.enPassant == other.enPassant
}

// ParseFEN returns an error if it could not parse an FEN. It was likely malformed or missing important pieces.
func ParseFEN(fen string) (*Position, error) {
	words := strings.Fields(fen)
	if len(words) != 6 {
		return nil, fmt.Errorf("fen %q could not be parsed: fen should contain 6 distinct sections", fen)
	}
	pos := &Position{}
	err := pos.parseFenBody(words[0])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	err = pos.parseSideToMove(words[1])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	err = pos.parseCastleRights(words[2])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	err = pos.parseEnPassant(words[3])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	err = pos.parseHalfMove(words[4])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	err = pos.parseFullMove(words[5])
	if err != nil {
		return nil, fmt.Errorf("fen %q could not be parsed: %w", fen, err)
	}
	return pos, nil
}

func (pos *Position) parseFenBody(body string) error {
	currentFile := FileA
	currentRank := Rank8
	for _, r := range body {
		if unicode.IsLetter(r) {
			p := parsePiece(string(r))
			if p == NoPiece {
				return fmt.Errorf("could not parse piece %q", r)
			}
			pos.SetPiece(p, Square{currentFile, currentRank})
		} else if unicode.IsNumber(r) {
			currentFile += File(r - '1') // Note this is 1 because file is automatically incremented in loop.
		} else if r == '/' {
			if currentFile != FileH+1 {
				return fmt.Errorf("invalid number of squares on rank %d, found %d", currentRank, currentFile-1)
			}
			currentRank -= 1
			currentFile = NoFile
		} else {
			return fmt.Errorf("encountered unexpected character %q", r)
		}
		currentFile += 1
	}
	if currentRank != Rank1 {
		return fmt.Errorf("invalid number of ranks, ended on rank %d, should be rank 1", currentRank)
	}
	return nil
}

func (pos *Position) parseSideToMove(sideToMove string) error {
	color := parseColor(sideToMove)
	if color == NoColor {
		return fmt.Errorf("could not parse color %q", sideToMove)
	}
	pos.SetSideToMove(color)
	return nil
}

func (pos *Position) parseCastleRights(castleRights string) error {
	if castleRights == "-" {
		return nil
	}
	for _, r := range castleRights {
		switch r {
		case 'K':
			if pos.WhiteKingSideCastle() {
				return fmt.Errorf("white king-side castle set twice")
			}
			pos.SetWhiteKingSideCastle(true)
		case 'Q':
			if pos.WhiteQueenSideCastle() {
				return fmt.Errorf("white queen-side castle set twice")
			}
			pos.SetWhiteQueenSideCastle(true)
		case 'k':
			if pos.BlackKingSideCastle() {
				return fmt.Errorf("black king-side castle set twice")
			}
			pos.SetBlackKingSideCastle(true)
		case 'q':
			if pos.BlackQueenSideCastle() {
				return fmt.Errorf("black queen-side castle set twice")
			}
			pos.SetBlackQueenSideCastle(true)
		default:
			return fmt.Errorf("could not parse castle rights %q", r)
		}
	}
	return nil
}

func (pos *Position) parseEnPassant(enPassant string) error {
	if enPassant == "-" {
		return nil
	}
	square := parseSquare(enPassant)
	if square == NoSquare {
		return fmt.Errorf("could not parse en passant %q", enPassant)
	}
	pos.SetEnPassant(square)
	return nil
}

func (pos *Position) parseHalfMove(halfMove string) error {
	hm, err := strconv.ParseUint(halfMove, 10, 0)
	if err != nil {
		return fmt.Errorf("could not parse half move %q", halfMove)
	}
	pos.SetHalfMove(uint(hm))
	return nil
}

func (pos *Position) parseFullMove(fullMove string) error {
	fm, err := strconv.ParseUint(fullMove, 10, 0)
	if err != nil {
		return fmt.Errorf("could not parse full move %q", fullMove)
	}
	pos.SetFullMove(uint(fm))
	return nil
}

// String generates an FEN string for the current position. See [PrettyString] for getting a board like representation.
func (pos *Position) String() string {
	fen := ""
	fen += pos.boardString() + " "
	fen += pos.sideToMoveString() + " "
	fen += pos.castleRightString() + " "
	fen += pos.EnPassant().String() + " "
	fen += strconv.FormatUint(uint64(pos.HalfMove()), 10) + " "
	fen += strconv.FormatUint(uint64(pos.FullMove()), 10)
	return fen
}

func (pos *Position) boardString() string {
	boardString := ""
	numEmptySquares := 0
	for currentRank := Rank8; currentRank != NoRank; currentRank -= 1 {
		for currentFile := FileA; currentFile <= FileH; currentFile += 1 {
			if piece := pos.Piece(Square{currentFile, currentRank}); piece == NoPiece {
				numEmptySquares += 1
			} else {
				if numEmptySquares > 0 {
					boardString += strconv.Itoa(numEmptySquares)
					numEmptySquares = 0
				}
				boardString += piece.String()
			}
		}
		if numEmptySquares > 0 {
			boardString += strconv.Itoa(numEmptySquares)
			numEmptySquares = 0
		}
		if currentRank != Rank1 {
			boardString += "/"
		}
	}
	return boardString
}

func (pos *Position) castleRightString() string {
	castleRights := ""
	if pos.WhiteKingSideCastle() {
		castleRights += "K"
	}
	if pos.WhiteQueenSideCastle() {
		castleRights += "Q"
	}
	if pos.BlackKingSideCastle() {
		castleRights += "k"
	}
	if pos.BlackQueenSideCastle() {
		castleRights += "q"
	}
	if len(castleRights) == 0 {
		castleRights = "-"
	}
	return castleRights
}

func (pos *Position) sideToMoveString() string {
	if pos.SideToMove() == White {
		return "w"
	}
	if pos.SideToMove() == Black {
		return "b"
	}
	return "-"
}

// PrettyString returns a board like representation of the current position. Uppercase letters are white and lowercase letters are black.
//
// Set whitesPerspective to true to see the board from white's side. Set extraInfo to false to just see the board. Set extraInfo to true to see all the other information stored in an FEN.
func (pos *Position) PrettyString(whitesPerspective bool, extraInfo bool) string {
	s := ""
	if whitesPerspective {
		s += pos.prettyBoardStringWhite()
	} else {
		s += pos.prettyBoardStringBlack()
	}
	if extraInfo {
		s += "\n\n"
		s += pos.extraInfo()
	}
	return s
}

func (pos *Position) prettyBoardStringWhite() string {
	s := ""
	for currentRank := Rank8; currentRank > NoRank; currentRank -= 1 {
		s += currentRank.String()
		for currentFile := FileA; currentFile <= FileH; currentFile += 1 {
			piece := pos.Piece(Square{currentFile, currentRank})
			s += piece.String()
		}
		s += "\n"
	}
	s += " ABCDEFGH"
	return s
}

func (pos *Position) prettyBoardStringBlack() string {
	s := ""
	for currentRank := Rank1; currentRank <= Rank8; currentRank += 1 {
		s += currentRank.String()
		for currentFile := FileH; currentFile > NoFile; currentFile -= 1 {
			piece := pos.Piece(Square{currentFile, currentRank})
			s += piece.String()
		}
		s += "\n"
	}
	s += " HGFEDCBA"
	return s
}

func (pos *Position) extraInfo() string {
	s := ""
	s += "Side To Move: "
	if pos.SideToMove() == White {
		s += "White"
	} else if pos.SideToMove() == Black {
		s += "Black"
	} else {
		s += "-"
	}
	s += "\n"

	s += "Castle Rights: "
	s += pos.castleRightString()
	s += "\n"
	s += "En Passant Square: "
	s += strings.ToUpper(pos.EnPassant().String())
	s += "\n"
	s += "Half Move: "
	s += strconv.FormatUint(uint64(pos.HalfMove()), 10)
	s += "\n"
	s += "Full Move: "
	s += strconv.FormatUint(uint64(pos.FullMove()), 10)
	return s
}

func (pos *Position) SideToMove() Color {
	return pos.sideToMove
}

func (pos *Position) SetSideToMove(c Color) {
	pos.sideToMove = c
}

// WhiteKingSideCastle returns true if white may still castle kingside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (pos *Position) WhiteKingSideCastle() bool {
	return pos.whiteKsCastle
}

// WhiteQueenSideCastle returns true if white may still castle queenside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (pos *Position) WhiteQueenSideCastle() bool {
	return pos.whiteQsCastle
}

// BlackKingSideCastle returns true if black may still castle kingside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (pos *Position) BlackKingSideCastle() bool {
	return pos.blackKsCastle
}

// BlackQueenSideCastle returns true if black may still castle queenside. Note that this does not indicate if the move is currently valid. It is really an indication of if the king or rook have moved yet this game.
func (pos *Position) BlackQueenSideCastle() bool {
	return pos.blackQsCastle
}

func (pos *Position) SetWhiteKingSideCastle(value bool) {
	pos.whiteKsCastle = value
}

func (pos *Position) SetWhiteQueenSideCastle(value bool) {
	pos.whiteQsCastle = value
}

func (pos *Position) SetBlackKingSideCastle(value bool) {
	pos.blackKsCastle = value
}

func (pos *Position) SetBlackQueenSideCastle(value bool) {
	pos.blackQsCastle = value
}

// EnPassant returns the square on to which a pawn may move to perform en-passant. This does not indicate if the move is legal. NoSquare is returned if there is no en passant option.
func (pos *Position) EnPassant() Square {
	return pos.enPassant
}

func (pos *Position) SetEnPassant(s Square) {
	pos.enPassant = s
}

func (pos *Position) HalfMove() uint {
	return pos.halfMove
}

func (pos *Position) SetHalfMove(i uint) {
	pos.halfMove = i
}

func (pos *Position) FullMove() uint {
	return pos.fullMove
}

func (pos *Position) SetFullMove(i uint) {
	pos.fullMove = i
}

// Piece gets the piece on the given square. NoPiece is returned if no piece is present.
func (pos *Position) Piece(s Square) Piece {
	if pos.whitePawns.Square(s) == 1 {
		return WhitePawn
	}
	if pos.whiteRooks.Square(s) == 1 {
		return WhiteRook
	}
	if pos.whiteKnights.Square(s) == 1 {
		return WhiteKnight
	}
	if pos.whiteBishops.Square(s) == 1 {
		return WhiteBishop
	}
	if pos.whiteQueens.Square(s) == 1 {
		return WhiteQueen
	}
	if pos.whiteKings.Square(s) == 1 {
		return WhiteKing
	}

	if pos.blackPawns.Square(s) == 1 {
		return BlackPawn
	}
	if pos.blackRooks.Square(s) == 1 {
		return BlackRook
	}
	if pos.blackKnights.Square(s) == 1 {
		return BlackKnight
	}
	if pos.blackBishops.Square(s) == 1 {
		return BlackBishop
	}
	if pos.blackQueens.Square(s) == 1 {
		return BlackQueen
	}
	if pos.blackKings.Square(s) == 1 {
		return BlackKing
	}

	return NoPiece
}

// SetPiece sets the given piece at the given square.
func (pos *Position) SetPiece(p Piece, s Square) {
	pos.ClearPiece(s)

	switch p {
	case WhitePawn:
		pos.whitePawns = pos.whitePawns.SetSquare(s)
	case WhiteRook:
		pos.whiteRooks = pos.whiteRooks.SetSquare(s)
	case WhiteKnight:
		pos.whiteKnights = pos.whiteKnights.SetSquare(s)
	case WhiteBishop:
		pos.whiteBishops = pos.whiteBishops.SetSquare(s)
	case WhiteQueen:
		pos.whiteQueens = pos.whiteQueens.SetSquare(s)
	case WhiteKing:
		pos.whiteKings = pos.whiteKings.SetSquare(s)

	case BlackPawn:
		pos.blackPawns = pos.blackPawns.SetSquare(s)
	case BlackRook:
		pos.blackRooks = pos.blackRooks.SetSquare(s)
	case BlackKnight:
		pos.blackKnights = pos.blackKnights.SetSquare(s)
	case BlackBishop:
		pos.blackBishops = pos.blackBishops.SetSquare(s)
	case BlackQueen:
		pos.blackQueens = pos.blackQueens.SetSquare(s)
	case BlackKing:
		pos.blackKings = pos.blackKings.SetSquare(s)
	}
}

// ClearPiece removes any piece from the given square.
func (pos *Position) ClearPiece(s Square) {
	pos.whitePawns = pos.whitePawns.ClearSquare(s)
	pos.whiteRooks = pos.whiteRooks.ClearSquare(s)
	pos.whiteKnights = pos.whiteKnights.ClearSquare(s)
	pos.whiteBishops = pos.whiteBishops.ClearSquare(s)
	pos.whiteQueens = pos.whiteQueens.ClearSquare(s)
	pos.whiteKings = pos.whiteKings.ClearSquare(s)

	pos.blackPawns = pos.blackPawns.ClearSquare(s)
	pos.blackRooks = pos.blackRooks.ClearSquare(s)
	pos.blackKnights = pos.blackKnights.ClearSquare(s)
	pos.blackBishops = pos.blackBishops.ClearSquare(s)
	pos.blackQueens = pos.blackQueens.ClearSquare(s)
	pos.blackKings = pos.blackKings.ClearSquare(s)
}

// Bitboard returns a bitboard for the given piece. See also [Position.OccupiedBitboard] and [Position.ColorBitboard].
func (pos *Position) Bitboard(p Piece) Bitboard {
	switch p {
	case WhitePawn:
		return pos.whitePawns
	case WhiteKnight:
		return pos.whiteKnights
	case WhiteBishop:
		return pos.whiteBishops
	case WhiteRook:
		return pos.whiteRooks
	case WhiteQueen:
		return pos.whiteQueens
	case WhiteKing:
		return pos.whiteKings

	case BlackPawn:
		return pos.blackPawns
	case BlackKnight:
		return pos.blackKnights
	case BlackBishop:
		return pos.blackBishops
	case BlackRook:
		return pos.blackRooks
	case BlackQueen:
		return pos.blackQueens
	case BlackKing:
		return pos.blackKings

	default:
		return 0
	}
}

// OccupiedBitboard returns a bitboard indicating all the squares with a piece on them.
func (pos *Position) OccupiedBitboard() Bitboard {
	return pos.whitePawns | pos.whiteKnights | pos.whiteBishops | pos.whiteRooks | pos.whiteQueens | pos.whiteKings |
		pos.blackPawns | pos.blackKnights | pos.blackBishops | pos.blackRooks | pos.blackQueens | pos.blackKings
}

// ColorBitboard returns a bitboard indicating all the squares occupied by pieces of a certain color.
func (pos *Position) ColorBitboard(c Color) Bitboard {
	if c == White {
		return pos.whitePawns | pos.whiteKnights | pos.whiteBishops | pos.whiteRooks | pos.whiteQueens | pos.whiteKings
	} else if c == Black {
		return pos.blackPawns | pos.blackKnights | pos.blackBishops | pos.blackRooks | pos.blackQueens | pos.blackKings
	} else {
		return 0
	}
}

// IsCheck returns true if the side to move has a king under attack from an enemy piece.
func (pos *Position) IsCheck() bool {
	var attackingSide Color
	if pos.SideToMove() == White {
		attackingSide = Black
	} else if pos.SideToMove() == Black {
		attackingSide = White
	} else {
		return false
	}

	attackedSquares := pos.getAttackedSquares(attackingSide)
	kingsInCheck := pos.Bitboard(Piece{pos.SideToMove(), King}) & attackedSquares
	return kingsInCheck > 0
}

// getAttackedSquares returns a bitboard with all the squares the specified color attacks.
func (pos *Position) getAttackedSquares(side Color) Bitboard {
	var attackedSquares Bitboard = 0

	occupied := pos.OccupiedBitboard()
	if side == White {
		attackedSquares |= pos.Bitboard(Piece{side, Pawn}).WhitePawnAttacks()
	} else if side == Black {
		attackedSquares |= pos.Bitboard(Piece{side, Pawn}).BlackPawnAttacks()
	}
	attackedSquares |= pos.Bitboard(Piece{side, Rook}).RookAttacks(occupied)
	attackedSquares |= pos.Bitboard(Piece{side, Knight}).KnightAttacks()
	attackedSquares |= pos.Bitboard(Piece{side, Bishop}).BishopAttacks(occupied)
	attackedSquares |= pos.Bitboard(Piece{side, Queen}).QueenAttacks(occupied)
	attackedSquares |= pos.Bitboard(Piece{side, King}).KingAttacks()
	return attackedSquares
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
func (pos *Position) Move(m Move) {
	pos.SetHalfMove(pos.HalfMove() + 1)

	if pos.isCastle(m) {
		pos.SetEnPassant(NoSquare)
		pos.performCastle(m)
	} else if pos.isPawnMove(m) {
		pos.performPawnMove(m)
	} else {
		pos.SetEnPassant(NoSquare)
		if pos.Piece(m.ToSquare) != NoPiece {
			pos.SetHalfMove(0)
		}
		pos.updateCastleRights(m)
		pos.SetPiece(pos.Piece(m.FromSquare), m.ToSquare)
		pos.ClearPiece(m.FromSquare)
	}

	pos.promotePiece(m.ToSquare, m.Promotion)
	pos.flipSide_incrementFullMove(m)
}

func (pos *Position) updateCastleRights(m Move) {
	if m.FromSquare == E1 || m.FromSquare == A1 {
		pos.SetWhiteQueenSideCastle(false)
	}
	if m.FromSquare == E1 || m.FromSquare == H1 {
		pos.SetWhiteKingSideCastle(false)
	}
	if m.FromSquare == E8 || m.FromSquare == A8 {
		pos.SetBlackQueenSideCastle(false)
	}
	if m.FromSquare == E8 || m.FromSquare == H8 {
		pos.SetBlackKingSideCastle(false)
	}
}

func (pos *Position) promotePiece(s Square, pt PieceType) {
	if pt != NoPieceType {
		piece := pos.Piece(s)
		if piece != NoPiece {
			piece.Type = pt
		}
		pos.SetPiece(piece, s)
	}
}

func (pos *Position) isPawnMove(m Move) bool {
	return pos.Piece(m.FromSquare).Type == Pawn
}

func (pos *Position) performPawnMove(m Move) {
	pos.SetHalfMove(0)
	pos.performPawnMove_takeEnPassant(m)
	pos.performPawnMove_setEnPassant(m)
	piece := pos.Piece(m.FromSquare)
	pos.SetPiece(piece, m.ToSquare)
	pos.ClearPiece(m.FromSquare)
}

func (pos *Position) performPawnMove_takeEnPassant(m Move) {
	if m.ToSquare == pos.EnPassant() {
		if pos.SideToMove() == White {
			pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank - 1})
		} else if pos.SideToMove() == Black {
			pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank + 1})
		} else {
			if pos.Piece(m.FromSquare).Color == White {
				pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank - 1})
			} else if pos.Piece(m.FromSquare).Color == Black {
				pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank + 1})
			}
		}
	}
}

func (pos *Position) performPawnMove_setEnPassant(m Move) {
	pos.SetEnPassant(NoSquare)
	movingPiece := pos.Piece(m.FromSquare)
	if m.FromSquare.File == m.ToSquare.File {
		if m.FromSquare.Rank == 2 && m.ToSquare.Rank == 4 && movingPiece.Color == White {
			pos.SetEnPassant(Square{m.FromSquare.File, m.FromSquare.Rank + 1})
		} else if m.FromSquare.Rank == 7 && m.ToSquare.Rank == 5 && movingPiece.Color == Black {
			pos.SetEnPassant(Square{m.FromSquare.File, m.FromSquare.Rank - 1})
		}
	}
}

func (pos *Position) isCastle(m Move) bool {
	switch m {
	case Move{E1, G1, NoPieceType}: // White King-side castle
		return pos.WhiteKingSideCastle() &&
			pos.Piece(E1) == WhiteKing &&
			pos.Piece(H1) == WhiteRook &&
			pos.Piece(F1) == NoPiece &&
			pos.Piece(G1) == NoPiece
	case Move{E1, C1, NoPieceType}: // White Queen-side castle
		return pos.WhiteQueenSideCastle() &&
			pos.Piece(E1) == WhiteKing &&
			pos.Piece(A1) == WhiteRook &&
			pos.Piece(D1) == NoPiece &&
			pos.Piece(C1) == NoPiece &&
			pos.Piece(B1) == NoPiece
	case Move{E8, G8, NoPieceType}: // Black King-side castle
		return pos.BlackKingSideCastle() &&
			pos.Piece(E8) == BlackKing &&
			pos.Piece(H8) == BlackRook &&
			pos.Piece(F8) == NoPiece &&
			pos.Piece(G8) == NoPiece
	case Move{E8, C8, NoPieceType}: // Black Queen-side castle
		return pos.BlackQueenSideCastle() &&
			pos.Piece(E8) == BlackKing &&
			pos.Piece(A8) == BlackRook &&
			pos.Piece(D8) == NoPiece &&
			pos.Piece(C8) == NoPiece &&
			pos.Piece(B8) == NoPiece
	default:
		return false
	}
}

func (pos *Position) performCastle(m Move) {
	switch m {
	case Move{E1, G1, NoPieceType}:
		pos.SetPiece(WhiteKing, G1)
		pos.ClearPiece(E1)
		pos.SetPiece(WhiteRook, F1)
		pos.ClearPiece(H1)
		pos.SetWhiteKingSideCastle(false)
		pos.SetWhiteQueenSideCastle(false)
	case Move{E1, C1, NoPieceType}:
		pos.SetPiece(WhiteKing, C1)
		pos.ClearPiece(E1)
		pos.SetPiece(WhiteRook, D1)
		pos.ClearPiece(A1)
		pos.SetWhiteKingSideCastle(false)
		pos.SetWhiteQueenSideCastle(false)
	case Move{E8, G8, NoPieceType}:
		pos.SetPiece(BlackKing, G8)
		pos.ClearPiece(E8)
		pos.SetPiece(BlackRook, F8)
		pos.ClearPiece(H8)
		pos.SetBlackKingSideCastle(false)
		pos.SetBlackQueenSideCastle(false)
	case Move{E8, C8, NoPieceType}:
		pos.SetPiece(BlackKing, C8)
		pos.ClearPiece(E8)
		pos.SetPiece(BlackRook, D8)
		pos.ClearPiece(A8)
		pos.SetBlackKingSideCastle(false)
		pos.SetBlackQueenSideCastle(false)
	}
}

func (pos *Position) flipSide_incrementFullMove(m Move) {
	if pos.SideToMove() == Black {
		pos.fullMove++
		pos.SetSideToMove(White)
	} else if pos.SideToMove() == White {
		pos.SetSideToMove(Black)
	} else {
		colorMoved := pos.Piece(m.ToSquare).Color
		if colorMoved == Black {
			pos.fullMove++
			pos.SetSideToMove(White)
		} else if colorMoved == White {
			pos.SetSideToMove(Black)
		}
	}
}
