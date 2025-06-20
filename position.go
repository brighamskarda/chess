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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const DefaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Position represents all parts of a chess position as specified by [Forsyth-Edwards Notation].
//
// The zero value is usable, though not very useful. See example for how to initialize the starting position.
//
// [Forsyth-Edwards Notation]: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c16.1
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

	SideToMove Color

	// WhiteKsCastle should be true if it is still possible to castle.
	WhiteKsCastle bool
	// WhiteQsCastle should be true if it is still possible to castle.
	WhiteQsCastle bool
	// BlackKsCastle should be true if it is still possible to castle.
	BlackKsCastle bool
	// BlackQsCastle should be true if it is still possible to castle.
	BlackQsCastle bool

	// EnPassant should be [NoSquare] if no en passant options are available.
	EnPassant Square

	HalfMove uint
	FullMove uint
}

// Copy creates a copy of the current position.
func (pos *Position) Copy() *Position {
	newPos := *pos
	return &newPos
}

// Equal returns true if the positions are the same, excluding move counters.
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
		pos.SideToMove == other.SideToMove &&
		pos.WhiteKsCastle == other.WhiteKsCastle &&
		pos.WhiteQsCastle == other.WhiteQsCastle &&
		pos.BlackKsCastle == other.BlackKsCastle &&
		pos.BlackQsCastle == other.BlackQsCastle &&
		pos.EnPassant == other.EnPassant
}

// UnmarshalText is an implementation of the [encoding.TextUnmarshaler] interface. It expects text in [Forsyth-Edwards Notation]. It returns an error if it could not parse fen. It was likely malformed or missing important pieces.
//
// [Forsyth-Edwards Notation]: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c16.1
func (pos *Position) UnmarshalText(fen []byte) error {
	words := strings.Fields(string(fen))
	if len(words) != 6 {
		return fmt.Errorf("pos %q could not be unmarshaled: fen should contain 6 distinct sections", fen)
	}
	p := &Position{}
	err := p.parseFenBody(words[0])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	err = p.parseSideToMove(words[1])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	err = p.parseCastleRights(words[2])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	err = p.parseEnPassant(words[3])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	err = p.parseHalfMove(words[4])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	err = p.parseFullMove(words[5])
	if err != nil {
		return fmt.Errorf("pos %q could not be unmarshaled: %w", fen, err)
	}
	*pos = *p
	return nil
}

func (pos *Position) parseFenBody(body string) error {
	currentFile := FileA
	currentRank := Rank8
	for _, r := range body {
		if unicode.IsLetter(r) {
			p, err := parsePiece(string(r))
			if err != nil {
				return fmt.Errorf("could not parse fen body: %w", err)
			}
			pos.SetPiece(p, Square{currentFile, currentRank})
		} else if unicode.IsNumber(r) {
			currentFile += File(r - '1') // Note this is 1 because file is automatically incremented in loop.
		} else if r == '/' {
			if currentFile != FileH+1 {
				return fmt.Errorf("could not parse fen body, invalid number of squares on rank %d", currentRank)
			}
			currentRank -= 1
			currentFile = NoFile
		} else {
			return fmt.Errorf("could not parse fen body, encountered unexpected character %q", r)
		}
		currentFile += 1
	}
	if currentRank != Rank1 {
		return fmt.Errorf("could not parse fen body, ended on rank %v, should be Rank1", currentRank)
	}
	return nil
}

func (pos *Position) parseSideToMove(sideToMove string) error {
	color := parseColor(sideToMove)
	if color == NoColor {
		return fmt.Errorf("could not parse side to move %q", sideToMove)
	}
	pos.SideToMove = color
	return nil
}

func (pos *Position) parseCastleRights(castleRights string) error {
	if castleRights == "-" {
		return nil
	}
	for _, r := range castleRights {
		switch r {
		case 'K':
			if pos.WhiteKsCastle {
				return errors.New("could not parse castle rights, white king-side castle set twice")
			}
			pos.WhiteKsCastle = true
		case 'Q':
			if pos.WhiteQsCastle {
				return errors.New("could not parse castle rights, white queen-side castle set twice")
			}
			pos.WhiteQsCastle = true
		case 'k':
			if pos.BlackKsCastle {
				return errors.New("could not parse castle rights, black king-side castle set twice")
			}
			pos.BlackKsCastle = true
		case 'q':
			if pos.BlackQsCastle {
				return errors.New("could not parse castle rights, black queen-side castle set twice")
			}
			pos.BlackQsCastle = true
		default:
			return fmt.Errorf("could not parse castle rights, invalid character %q", r)
		}
	}
	return nil
}

func (pos *Position) parseEnPassant(enPassant string) error {
	if enPassant == "-" {
		return nil
	}
	square := Square{}
	err := square.UnmarshalText([]byte(enPassant))
	if err != nil {
		return fmt.Errorf("could not parse en passant: %w", err)
	}
	pos.EnPassant = square
	return nil
}

func (pos *Position) parseHalfMove(halfMove string) error {
	hm, err := strconv.ParseUint(halfMove, 10, 0)
	if err != nil {
		return fmt.Errorf("could not parse half move: %w", err)
	}
	pos.HalfMove = uint(hm)
	return nil
}

func (pos *Position) parseFullMove(fullMove string) error {
	fm, err := strconv.ParseUint(fullMove, 10, 0)
	if err != nil {
		return fmt.Errorf("could not parse full move %w", err)
	}
	pos.FullMove = uint(fm)
	return nil
}

// MarshalText is an implementation of the [encoding.TextMarshaler] interface. It provides the [FEN] representation of the board and returns an error if position contains invalid fields. See also [Position.String] for a more human readable form of the position.
//
// [FEN]: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c16.1
func (pos *Position) MarshalText() (text []byte, err error) {
	fen := ""
	board := pos.boardString()
	fen += board + " "
	stm, err := pos.sideToMoveString()
	if err != nil {
		return nil, fmt.Errorf("could not marshal position: %w", err)
	}
	fen += stm + " "
	fen += pos.castleRightString() + " "
	enPassant, err := pos.EnPassant.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("could not marshal position, could not marshal en passant: %w", err)
	}
	fen += string(enPassant) + " "
	fen += strconv.FormatUint(uint64(pos.HalfMove), 10) + " "
	fen += strconv.FormatUint(uint64(pos.FullMove), 10)
	return []byte(fen), nil
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
	if pos.WhiteKsCastle {
		castleRights += "K"
	}
	if pos.WhiteQsCastle {
		castleRights += "Q"
	}
	if pos.BlackKsCastle {
		castleRights += "k"
	}
	if pos.BlackQsCastle {
		castleRights += "q"
	}
	if len(castleRights) == 0 {
		castleRights = "-"
	}
	return castleRights
}

func (pos *Position) sideToMoveString() (string, error) {
	if pos.SideToMove == White {
		return "w", nil
	}
	if pos.SideToMove == Black {
		return "b", nil
	}
	return "", errors.New("side to move not set")
}

// String returns a board like representation of the current position. Uppercase letters are white and lowercase letters are black.
//
// Set whitesPerspective to true to see the board from white's side. Set extraInfo to false to just see the board. Set extraInfo to true to see all the other information stored in an FEN.
func (pos *Position) String(whitesPerspective bool, extraInfo bool) string {
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
	switch pos.SideToMove {
	case White:
		s += "White"
	case Black:
		s += "Black"
	default:
		s += "-"
	}
	s += "\n"

	s += "Castle Rights: "
	s += pos.castleRightString()
	s += "\n"
	s += "En Passant Square: "
	enPassant, err := pos.EnPassant.MarshalText()
	if err != nil {
		enPassant = []byte{'-'}
	}
	s += string(enPassant)
	s += "\n"
	s += "Half Move: "
	s += strconv.FormatUint(uint64(pos.HalfMove), 10)
	s += "\n"
	s += "Full Move: "
	s += strconv.FormatUint(uint64(pos.FullMove), 10)
	return s
}

// Piece gets the piece on the given square. [NoPiece] is returned if no piece is present, or square is invalid.
func (pos *Position) Piece(s Square) Piece {
	if !squareOnBoard(s) {
		return NoPiece
	}
	squareIndex := squareToIndex(s)
	if pos.whitePawns.Bit(squareIndex) == 1 {
		return WhitePawn
	}
	if pos.blackPawns.Bit(squareIndex) == 1 {
		return BlackPawn
	}
	if pos.whiteRooks.Bit(squareIndex) == 1 {
		return WhiteRook
	}
	if pos.whiteKnights.Bit(squareIndex) == 1 {
		return WhiteKnight
	}
	if pos.whiteBishops.Bit(squareIndex) == 1 {
		return WhiteBishop
	}
	if pos.blackRooks.Bit(squareIndex) == 1 {
		return BlackRook
	}
	if pos.blackKnights.Bit(squareIndex) == 1 {
		return BlackKnight
	}
	if pos.blackBishops.Bit(squareIndex) == 1 {
		return BlackBishop
	}
	if pos.whiteQueens.Bit(squareIndex) == 1 {
		return WhiteQueen
	}
	if pos.blackQueens.Bit(squareIndex) == 1 {
		return BlackQueen
	}
	if pos.whiteKings.Bit(squareIndex) == 1 {
		return WhiteKing
	}
	if pos.blackKings.Bit(squareIndex) == 1 {
		return BlackKing
	}

	return NoPiece
}

// SetPiece sets p on square s. If p or s are invalid nothings happens.
func (pos *Position) SetPiece(p Piece, s Square) {
	if !squareOnBoard(s) {
		return
	}
	squareIndex := squareToIndex(s)

	pos.ClearPiece(s)

	switch p {
	case WhitePawn:
		pos.whitePawns = pos.whitePawns.SetBit(squareIndex)
	case WhiteRook:
		pos.whiteRooks = pos.whiteRooks.SetBit(squareIndex)
	case WhiteKnight:
		pos.whiteKnights = pos.whiteKnights.SetBit(squareIndex)
	case WhiteBishop:
		pos.whiteBishops = pos.whiteBishops.SetBit(squareIndex)
	case WhiteQueen:
		pos.whiteQueens = pos.whiteQueens.SetBit(squareIndex)
	case WhiteKing:
		pos.whiteKings = pos.whiteKings.SetBit(squareIndex)

	case BlackPawn:
		pos.blackPawns = pos.blackPawns.SetBit(squareIndex)
	case BlackRook:
		pos.blackRooks = pos.blackRooks.SetBit(squareIndex)
	case BlackKnight:
		pos.blackKnights = pos.blackKnights.SetBit(squareIndex)
	case BlackBishop:
		pos.blackBishops = pos.blackBishops.SetBit(squareIndex)
	case BlackQueen:
		pos.blackQueens = pos.blackQueens.SetBit(squareIndex)
	case BlackKing:
		pos.blackKings = pos.blackKings.SetBit(squareIndex)
	}
}

// ClearPiece removes any piece from the given square. Nothing happens if s is invalid.
func (pos *Position) ClearPiece(s Square) {
	if !squareOnBoard(s) {
		return
	}
	squareIndex := squareToIndex(s)

	pos.whitePawns = pos.whitePawns.ClearBit(squareIndex)
	pos.whiteRooks = pos.whiteRooks.ClearBit(squareIndex)
	pos.whiteKnights = pos.whiteKnights.ClearBit(squareIndex)
	pos.whiteBishops = pos.whiteBishops.ClearBit(squareIndex)
	pos.whiteQueens = pos.whiteQueens.ClearBit(squareIndex)
	pos.whiteKings = pos.whiteKings.ClearBit(squareIndex)

	pos.blackPawns = pos.blackPawns.ClearBit(squareIndex)
	pos.blackRooks = pos.blackRooks.ClearBit(squareIndex)
	pos.blackKnights = pos.blackKnights.ClearBit(squareIndex)
	pos.blackBishops = pos.blackBishops.ClearBit(squareIndex)
	pos.blackQueens = pos.blackQueens.ClearBit(squareIndex)
	pos.blackKings = pos.blackKings.ClearBit(squareIndex)
}

// Bitboard returns a bitboard for the given piece. If p is [NoPiece] then a bitboard with all the unoccupied squares is returned. If p is invalid 0 is returned. See also [Position.OccupiedBitboard] and [Position.ColorBitboard].
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

	case NoPiece:
		return ^(pos.whitePawns | pos.whiteKnights | pos.whiteBishops | pos.whiteRooks |
			pos.whiteQueens | pos.whiteKings | pos.blackPawns | pos.blackKnights |
			pos.blackBishops | pos.blackRooks | pos.blackQueens | pos.blackKings)

	default:
		return 0
	}
}

// OccupiedBitboard returns a bitboard indicating all the squares with a piece on them.
func (pos *Position) OccupiedBitboard() Bitboard {
	return pos.whitePawns | pos.whiteKnights | pos.whiteBishops | pos.whiteRooks | pos.whiteQueens | pos.whiteKings |
		pos.blackPawns | pos.blackKnights | pos.blackBishops | pos.blackRooks | pos.blackQueens | pos.blackKings
}

// ColorBitboard returns a bitboard indicating all the squares occupied by pieces of a certain color. Returns 0 if NoColor or invalid color.
func (pos *Position) ColorBitboard(c Color) Bitboard {
	switch c {
	case White:
		return pos.whitePawns | pos.whiteKnights | pos.whiteBishops | pos.whiteRooks | pos.whiteQueens | pos.whiteKings
	case Black:
		return pos.blackPawns | pos.blackKnights | pos.blackBishops | pos.blackRooks | pos.blackQueens | pos.blackKings
	default:
		return 0
	}
}

// IsCheck returns true if the side to move has a king under attack from an enemy piece. If side to move is not set false is returned.
func (pos *Position) IsCheck() bool {
	occupied := pos.OccupiedBitboard()

	switch pos.SideToMove {
	case White:
		return pos.whiteKings&pos.blackPawns.BlackPawnAttacks() > 0 ||
			pos.whiteKings&pos.blackKnights.KnightAttacks() > 0 ||
			pos.whiteKings&(pos.blackBishops|pos.blackQueens).BishopAttacks(occupied) > 0 ||
			pos.whiteKings&(pos.blackRooks|pos.blackQueens).RookAttacks(occupied) > 0 ||
			pos.whiteKings&pos.blackKings.KingAttacks() > 0
	case Black:
		return pos.blackKings&pos.whitePawns.WhitePawnAttacks() > 0 ||
			pos.blackKings&pos.whiteKnights.KnightAttacks() > 0 ||
			pos.blackKings&(pos.whiteBishops|pos.whiteQueens).BishopAttacks(occupied) > 0 ||
			pos.blackKings&(pos.whiteRooks|pos.whiteQueens).RookAttacks(occupied) > 0 ||
			pos.blackKings&pos.whiteKings.KingAttacks() > 0
	default:
		return false
	}
}

// getAttackedSquares returns a bitboard with all the squares the specified color attacks.
func (pos *Position) getAttackedSquares(side Color) Bitboard {
	var attackedSquares Bitboard = 0

	occupied := pos.OccupiedBitboard()
	switch side {
	case White:
		attackedSquares |= pos.whitePawns.WhitePawnAttacks()
		attackedSquares |= (pos.whiteRooks | pos.whiteQueens).RookAttacks(occupied)
		attackedSquares |= pos.whiteKnights.KnightAttacks()
		attackedSquares |= (pos.whiteBishops | pos.whiteQueens).BishopAttacks(occupied)
		attackedSquares |= pos.whiteKings.KingAttacks()

	case Black:
		attackedSquares |= pos.blackPawns.BlackPawnAttacks()
		attackedSquares |= (pos.blackRooks | pos.blackQueens).RookAttacks(occupied)
		attackedSquares |= pos.blackKnights.KnightAttacks()
		attackedSquares |= (pos.blackBishops | pos.blackQueens).BishopAttacks(occupied)
		attackedSquares |= pos.blackKings.KingAttacks()
	}
	return attackedSquares
}

// Move performs chess moves in such a way that if all moves are legal, the FEN will always be properly updated. If m contains invalid squares or promotions the results are undefined. The rules it follows are listed below.
//
//  1. By default the following happens:
//
//     a. The piece at the from square is moved to the to square and promoted. (This also includes moving NoPiece, in which case the promotion is not applied.)
//
//     b. The half move counter is incremented.
//
//     c. The side to move is flipped (or set to the opposite of the piece moved if not previously set [stays on NoColor if not set and side to move is not set]).
//
//     d. If the side to move flips from black to white then the full move counter is incremented.
//
//     e. En-passant is set to NoSquare.
//
//  2. If a pawn advances, or a piece is taken the half move counter is reset.
//
//  3. If a pawn advances two spaces forward from its starting rank en-passant is set to the square right behind its current position.
//
//  4. If a king or rook moves from their starting square (in standard chess, 960 is not supported) then the corresponding castle rights are set to false.
//
//  5. If one of the four possible castle moves if executed and the castle rights still exist, and there are no pieces in the way, then the appropriate castle move will be applied. (Check will not block a castle move)
func (pos *Position) Move(m Move) {
	pos.HalfMove = pos.HalfMove + 1

	if pos.isCastle(m) {
		pos.EnPassant = NoSquare
		pos.performCastle(m)
	} else if pos.isPawnMove(m) {
		pos.performPawnMove(m)
	} else {
		pos.EnPassant = NoSquare
		if pos.Piece(m.ToSquare) != NoPiece {
			pos.HalfMove = 0
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
		pos.WhiteQsCastle = false
	}
	if m.FromSquare == E1 || m.FromSquare == H1 {
		pos.WhiteKsCastle = false
	}
	if m.FromSquare == E8 || m.FromSquare == A8 {
		pos.BlackQsCastle = false
	}
	if m.FromSquare == E8 || m.FromSquare == H8 {
		pos.BlackKsCastle = false
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
	squareIndex := squareToIndex(m.FromSquare)
	return pos.whitePawns.Bit(squareIndex) == 1 || pos.blackPawns.Bit(squareIndex) == 1
}

func (pos *Position) performPawnMove(m Move) {
	pos.HalfMove = 0
	pos.performPawnMove_takeEnPassant(m)
	pos.performPawnMove_setEnPassant(m)
	if squareIndex := squareToIndex(m.FromSquare); pos.whitePawns.Bit(squareIndex) == 1 {
		pos.SetPiece(WhitePawn, m.ToSquare)
		pos.whitePawns = pos.whitePawns.ClearBit(squareIndex)
	} else {
		pos.SetPiece(BlackPawn, m.ToSquare)
		pos.blackPawns = pos.blackPawns.ClearBit(squareIndex)
	}
}

func (pos *Position) performPawnMove_takeEnPassant(m Move) {
	if m.ToSquare == pos.EnPassant {
		switch pos.SideToMove {
		case White:
			pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank - 1})
		case Black:
			pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank + 1})
		default:
			if pos.Piece(m.FromSquare).Color == White {
				pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank - 1})
			} else if pos.Piece(m.FromSquare).Color == Black {
				pos.ClearPiece(Square{m.ToSquare.File, m.ToSquare.Rank + 1})
			}
		}
	}
}

func (pos *Position) performPawnMove_setEnPassant(m Move) {
	pos.EnPassant = NoSquare
	movingPiece := pos.Piece(m.FromSquare)
	if m.FromSquare.File == m.ToSquare.File {
		if m.FromSquare.Rank == 2 && m.ToSquare.Rank == 4 && movingPiece.Color == White {
			pos.EnPassant = Square{m.FromSquare.File, m.FromSquare.Rank + 1}
		} else if m.FromSquare.Rank == 7 && m.ToSquare.Rank == 5 && movingPiece.Color == Black {
			pos.EnPassant = Square{m.FromSquare.File, m.FromSquare.Rank - 1}
		}
	}
}

func (pos *Position) isCastle(m Move) bool {
	switch m {
	case Move{E1, G1, NoPieceType}: // White King-side castle
		return pos.WhiteKsCastle &&
			pos.whiteKings.Bit(4) == 1 &&
			pos.whiteRooks.Bit(7) == 1 &&
			pos.Piece(F1) == NoPiece &&
			pos.Piece(G1) == NoPiece
	case Move{E1, C1, NoPieceType}: // White Queen-side castle
		return pos.WhiteQsCastle &&
			pos.whiteKings.Bit(4) == 1 &&
			pos.whiteRooks.Bit(0) == 1 &&
			pos.Piece(D1) == NoPiece &&
			pos.Piece(C1) == NoPiece &&
			pos.Piece(B1) == NoPiece
	case Move{E8, G8, NoPieceType}: // Black King-side castle
		return pos.BlackKsCastle &&
			pos.blackKings.Bit(60) == 1 &&
			pos.blackRooks.Bit(63) == 1 &&
			pos.Piece(F8) == NoPiece &&
			pos.Piece(G8) == NoPiece
	case Move{E8, C8, NoPieceType}: // Black Queen-side castle
		return pos.BlackQsCastle &&
			pos.blackKings.Bit(60) == 1 &&
			pos.blackRooks.Bit(56) == 1 &&
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
		pos.whiteKings = pos.whiteKings.SetBit(6).ClearBit(4)
		pos.whiteRooks = pos.whiteRooks.SetBit(5).ClearBit(7)
		pos.WhiteKsCastle = false
		pos.WhiteQsCastle = false
	case Move{E1, C1, NoPieceType}:
		pos.whiteKings = pos.whiteKings.SetBit(2).ClearBit(4)
		pos.whiteRooks = pos.whiteRooks.SetBit(3).ClearBit(0)
		pos.WhiteKsCastle = false
		pos.WhiteQsCastle = false
	case Move{E8, G8, NoPieceType}:
		pos.blackKings = pos.blackKings.SetBit(62).ClearBit(60)
		pos.blackRooks = pos.blackRooks.SetBit(61).ClearBit(63)
		pos.BlackKsCastle = false
		pos.BlackQsCastle = false
	case Move{E8, C8, NoPieceType}:
		pos.blackKings = pos.blackKings.SetBit(58).ClearBit(60)
		pos.blackRooks = pos.blackRooks.SetBit(59).ClearBit(56)
		pos.BlackKsCastle = false
		pos.BlackQsCastle = false
	}
}

func (pos *Position) flipSide_incrementFullMove(m Move) {
	switch pos.SideToMove {
	case Black:
		pos.FullMove++
		pos.SideToMove = White
	case White:
		pos.SideToMove = Black
	default:
		colorMoved := pos.Piece(m.ToSquare).Color
		switch colorMoved {
		case Black:
			pos.FullMove++
			pos.SideToMove = White
		case White:
			pos.SideToMove = Black
		}
	}
}
