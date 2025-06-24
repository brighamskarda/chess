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
	"math/bits"
	"slices"
	"strings"
	"unicode"
)

// Move represents a UCI chess move. UCI chess moves are easily represented by three fields: FromSquare, ToSquare, and an optional Promotion. [UCI Specification]
//
// [UCI Specification]: https://www.shredderchess.com/download/div/uci.zip
type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
}

// MarshalText implements the [encoding.TextMarshaler] interface to encode a move into a UCI compatible format. The form is <FromSquare><ToSquare><OptionalPromotion>, ex. "a7a8q". If fromsquare, tosquare, and promotion are all 0 then "0000" is returned as per the [UCI specification]. An error may be returned if a field is missing or malformed.
//
// [UCI specification]: https://www.shredderchess.com/download/div/uci.zip
func (m Move) MarshalText() (text []byte, err error) {
	if m.FromSquare == NoSquare && m.ToSquare == NoSquare && m.Promotion == NoPieceType {
		return []byte{'0', '0', '0', '0'}, nil
	}
	if !squareOnBoard(m.FromSquare) || !squareOnBoard(m.ToSquare) {
		return nil, fmt.Errorf("could not marshal move %#v, contains squares not on board", m)
	}
	if m.Promotion > King {
		return nil, fmt.Errorf("could not marshal move %#v: invalid promotion", m)
	}
	from, _ := m.FromSquare.MarshalText()
	to, _ := m.ToSquare.MarshalText()
	text = append(text, from...)
	text = append(text, to...)
	if m.Promotion != NoPieceType {
		text = append(text, m.Promotion.String()[0])
	}
	return text, nil
}

// String provides a UCI compatible representation of the square in the form <FromSquare><ToSquare><OptionalPromotion>, ex. "a7a8q". If fromsquare, tosquare, and promotion are all 0 then "0000" is returned as per the [UCI specification]. An error string is returned if any of the fields are invalid.
//
// [UCI specification]: https://www.shredderchess.com/download/div/uci.zip
func (m Move) String() string {
	text, err := m.MarshalText()
	if err != nil {
		return fmt.Sprintf("Invalid Move %#v", m)
	}
	return string(text)
}

// StringSAN provides a move in standard algebraic notation as specified in the [PGN specification].
//
// pos is required to convert a move to SAN. pos should be the position before the move. An error is returned if an SAN string could not be generated with the given move and position. This does not necessarily test for move legality, should this be an illegal move, unexpected results may occur.
//
// [PGN specification]: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.2.3
func (m Move) StringSAN(pos *Position) (string, error) {
	// This is another cry for help, why couldn't the pgn file spec just use UCI notation. This is another set of complex logic that was not necessary.
	if !squareOnBoard(m.FromSquare) || !squareOnBoard(m.ToSquare) {
		return "", fmt.Errorf("could not convert move %v to SAN, contains squares not on board", m)
	}
	if m.Promotion > King {
		return "", fmt.Errorf("could not convert move %v to SAN: invalid promotion", m)
	}
	pos = pos.Copy()
	returnString := ""
	if pos.Piece(m.FromSquare) == NoPiece {
		return "", fmt.Errorf("could not convert move %v to SAN, no piece on FromSquare", m)
	} else if pos.Piece(m.FromSquare).Type == Pawn {
		returnString = m.pawnStringSAN()
	} else {
		returnString = m.normalStringSAN(pos)
	}

	pos.Move(m)
	if pos.IsCheck() {
		if len(LegalMoves(pos)) == 0 {
			returnString += "#"
		} else {
			returnString += "+"
		}
	}
	return returnString, nil
}

func (m Move) pawnStringSAN() string {
	returnString := ""
	if m.FromSquare.File != m.ToSquare.File {
		returnString += m.FromSquare.File.String() + "x"
	}
	toSquare, _ := m.ToSquare.MarshalText()
	returnString += string(toSquare)
	if m.Promotion != NoPieceType {
		returnString += "=" + strings.ToUpper(m.Promotion.String())
	}
	return returnString
}

func (m Move) normalStringSAN(pos *Position) string {
	returnString := ""
	pieceToMove := pos.Piece(m.FromSquare)
	returnString += strings.ToUpper(pieceToMove.String())
	if isFileDisambiguation(m.FromSquare, m.ToSquare, pieceToMove, pos) {
		returnString += m.FromSquare.File.String()
	} else if isRankDisambiguation(m.FromSquare, m.ToSquare, pieceToMove, pos) {
		returnString += m.FromSquare.Rank.String()
	} else if isSquareDisambiguation(m.ToSquare, pieceToMove, pos) {
		fromSquare, _ := m.FromSquare.MarshalText()
		returnString += string(fromSquare)
	}
	if pos.Piece(m.ToSquare) != NoPiece {
		returnString += "x"
	}
	toSquare, _ := m.ToSquare.MarshalText()
	returnString += string(toSquare)
	return returnString
}

func isFileDisambiguation(fromSquare Square, toSquare Square, pieceToMove Piece, pos *Position) bool {
	possibleMoves := []Move{}
	piecesToMove := pos.Bitboard(pieceToMove)
	for piecesToMove != 0 {
		singlePiece := 1 << bits.TrailingZeros(uint(piecesToMove))
		piecesToMove &^= Bitboard(singlePiece)
		attacks := getPieceAttacks(Bitboard(singlePiece), pieceToMove.Type, pos)
		if attacks.Square(toSquare) == 1 {
			fromSquare := indexToSquare(bits.TrailingZeros(uint(singlePiece)))
			possibleMoves = append(possibleMoves, Move{fromSquare, toSquare, NoPieceType})

		}
	}

	legalMoves := []Move{}
	for _, m := range possibleMoves {
		posCopy := pos.Copy()
		posCopy.Move(m)
		posCopy.SideToMove = pos.SideToMove
		if !posCopy.IsCheck() {
			legalMoves = append(legalMoves, m)
		}
	}

	numMovesFromFile := 0

	for _, m := range legalMoves {
		if m.FromSquare.File == fromSquare.File {
			numMovesFromFile++
		}
	}

	return numMovesFromFile == 1 && len(legalMoves) > 1
}

func isRankDisambiguation(fromSquare Square, toSquare Square, pieceToMove Piece, pos *Position) bool {
	possibleMoves := []Move{}
	piecesToMove := pos.Bitboard(pieceToMove)
	for piecesToMove != 0 {
		singlePiece := 1 << bits.TrailingZeros(uint(piecesToMove))
		piecesToMove &^= Bitboard(singlePiece)
		attacks := getPieceAttacks(Bitboard(singlePiece), pieceToMove.Type, pos)
		if attacks.Square(toSquare) == 1 {
			fromSquare := indexToSquare(bits.TrailingZeros(uint(singlePiece)))
			possibleMoves = append(possibleMoves, Move{fromSquare, toSquare, NoPieceType})
		}
	}

	legalMoves := []Move{}
	for _, m := range possibleMoves {
		posCopy := pos.Copy()
		posCopy.Move(m)
		posCopy.SideToMove = pos.SideToMove
		if !posCopy.IsCheck() {
			legalMoves = append(legalMoves, m)
		}
	}

	numMovesFromRank := 0

	for _, m := range legalMoves {
		if m.FromSquare.Rank == fromSquare.Rank {
			numMovesFromRank++
		}
	}

	return numMovesFromRank == 1 && len(legalMoves) > 1

}

func isSquareDisambiguation(toSquare Square, pieceToMove Piece, pos *Position) bool {
	possibleMoves := []Move{}
	piecesToMove := pos.Bitboard(pieceToMove)
	for piecesToMove != 0 {
		singlePiece := 1 << bits.TrailingZeros(uint(piecesToMove))
		piecesToMove &^= Bitboard(singlePiece)
		attacks := getPieceAttacks(Bitboard(singlePiece), pieceToMove.Type, pos)
		if attacks.Square(toSquare) == 1 {
			fromSquare := indexToSquare(bits.TrailingZeros(uint(singlePiece)))
			possibleMoves = append(possibleMoves, Move{fromSquare, toSquare, NoPieceType})
		}
	}

	legalMoves := []Move{}
	for _, m := range possibleMoves {
		posCopy := pos.Copy()
		posCopy.Move(m)
		posCopy.SideToMove = pos.SideToMove
		if !posCopy.IsCheck() {
			legalMoves = append(legalMoves, m)
		}
	}

	return len(legalMoves) > 1
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface to read a move from a UCI compatible format. The format is <FromSquare><ToSquare><OptionalPromotion>, ex. a2c3q.
func (m *Move) UnmarshalText(text []byte) error {
	if len(text) < 4 || len(text) > 5 {
		return fmt.Errorf("could not unmarshal move, expected text to be of len 4 or 5, got len(text) = %d", len(text))
	}
	var fromSquare Square
	err := fromSquare.UnmarshalText(text[0:2])
	if err != nil {
		return fmt.Errorf("could not unmarshal move: %w", err)
	}

	var toSquare Square
	err = toSquare.UnmarshalText(text[2:4])
	if err != nil {
		return fmt.Errorf("could not unmarshal move: %w", err)
	}

	var promotion PieceType
	if len(text) == 5 {
		promotion, err = parsePieceType(text[4])
		if err != nil {
			return fmt.Errorf("could not unmarshal move: %w", err)
		}
	}
	m.FromSquare = fromSquare
	m.ToSquare = toSquare
	m.Promotion = promotion
	return nil
}

type sanMoveType uint8

const (
	unknown sanMoveType = iota
	pawnAdvance
	pawnCapture
	normalMove
	fileDisambiguation
	rankDisambiguation
	squareDisambiguation
	castleMove
)

// ParseSANMove parses a chess move provided in [Standard Algebraic Notation]. Standard Algebraic notation requires position information to know how moves should be parsed. An error is provided if the move could not be parsed.
//
// [Standard Algebraic Notation]: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.2.3
func ParseSANMove(san string, pos *Position) (Move, error) {
	// This comment is a silent cry for help. I only implemented this because the pgn standard uses it. SAN was not designed to be parsed by computers. This was a major effort that could have been easily avoided by the much superior UCI move notation.
	if pos.SideToMove != White && pos.SideToMove != Black {
		return Move{}, fmt.Errorf("could not parse SAN move %q, pos.SideToMove not set", san)
	}

	// These characters don't matter for this parsing.
	san = strings.ReplaceAll(san, "+", "")
	san = strings.ReplaceAll(san, "#", "")
	san = strings.ReplaceAll(san, "!", "")
	san = strings.ReplaceAll(san, "?", "")

	var m Move
	var err error

	switch classifySANMove(san) {
	case fileDisambiguation:
		m, err = parseFileDisambiguation(san, pos)
	case rankDisambiguation:
		m, err = parseRankDisambiguation(san, pos)
	case squareDisambiguation:
		m, err = parseSquareDisambiguation(san)
	case normalMove:
		m, err = parseNormalMove(san, pos)
	case pawnAdvance:
		m, err = parsePawnAdvance(san, pos)
	case pawnCapture:
		m, err = parsePawnCapture(san, pos)
	case castleMove:
		m = parseCastleMove(san, pos)
	case unknown:
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine san move type, the move was likely malformed", san)
	default:
		panic("unexpected chess.sanMoveType")
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: %w", san, err)
	}
	return m, nil
}

func classifySANMove(san string) sanMoveType {
	switch len(san) {
	case 2:
		return pawnAdvance
	case 3:
		if strings.ToLower(san[0:1]) == "o" && san[1] == '-' && strings.ToLower(san[2:3]) == "o" {
			return castleMove
		}
		return normalMove
	case 4:
		return classifySANMove4(san)
	case 5:
		return classifySANMove5(san)
	case 6:
		if san[4] == '=' {
			return pawnCapture
		}
		return squareDisambiguation
	default:
		return unknown
	}
}

func classifySANMove4(san string) sanMoveType {
	if san[0] == 'B' || slices.Contains([]rune{'r', 'n', 'q', 'k'}, unicode.ToLower(rune(san[0]))) {
		if strings.ToLower(san[1:2]) == "x" {
			return normalMove
		} else if slices.Contains([]rune{'1', '2', '3', '4', '5', '6', '7', '8'}, rune(san[1])) {
			return rankDisambiguation
		} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(rune(san[1]))) {
			return fileDisambiguation
		} else {
			return unknown
		}
	} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(rune(san[0]))) {
		if san[2] == '=' {
			return pawnAdvance
		} else if strings.ToLower(san[1:2]) == "x" {
			return pawnCapture
		} else {
			return unknown
		}
	} else {
		return unknown
	}
}

func classifySANMove5(san string) sanMoveType {
	if strings.ToLower(san[0:1]) == "o" && san[1] == '-' && strings.ToLower(san[2:3]) == "o" && san[3] == '-' && strings.ToLower(san[4:5]) == "o" {
		return castleMove
	} else if strings.ToLower(san[2:3]) == "x" {
		if slices.Contains([]rune{'1', '2', '3', '4', '5', '6', '7', '8'}, rune(san[1])) {
			return rankDisambiguation
		} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(rune(san[1]))) {
			return fileDisambiguation
		} else {
			return unknown
		}
	} else {
		return squareDisambiguation
	}
}

func parseCastleMove(san string, pos *Position) Move {
	switch strings.ToLower(san) {
	case "o-o":
		switch pos.SideToMove {
		case White:
			return Move{E1, G1, NoPieceType}
		case Black:
			return Move{E8, G8, NoPieceType}
		}
	case "o-o-o":
		switch pos.SideToMove {
		case White:
			return Move{E1, C1, NoPieceType}
		case Black:
			return Move{E8, C8, NoPieceType}
		}
	}
	panic("unexpected condition: bad san in parseCastleMove")

}

func parsePawnCapture(san string, pos *Position) (Move, error) {
	m := Move{}
	toSquare := &Square{}
	err := toSquare.UnmarshalText([]byte(san[2:4]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse pawn capture: %w", err)
	}
	m.ToSquare = *toSquare

	fromFile, err := parseFile(san[0])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse pawn capture: %w", err)
	}
	switch pos.SideToMove {
	case White:
		m.FromSquare = Square{fromFile, toSquare.Rank - 1}
	case Black:
		m.FromSquare = Square{fromFile, toSquare.Rank + 1}
	}

	if len(san) == 6 {
		// Pawn capture with promotion
		pt, err := parsePieceType(san[5])
		if err != nil {
			return Move{}, fmt.Errorf("could not parse pawn capture: %w", err)
		}
		m.Promotion = pt
	}
	return m, nil
}

func parsePawnAdvance(san string, pos *Position) (Move, error) {
	toSquare := Square{}
	err := toSquare.UnmarshalText([]byte(san[0:2]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse pawn advance: %w", err)
	}

	var fromSquare Square
	switch pos.SideToMove {
	case White:
		fromSquare = Square{toSquare.File, toSquare.Rank - 1}
	case Black:
		fromSquare = Square{toSquare.File, toSquare.Rank + 1}
	}

	// The following take into account the possibility of a double pawn advance.
	if fromSquare.Rank == Rank3 && toSquare.Rank == Rank4 && pos.Piece(fromSquare).Type != Pawn {
		fromSquare.Rank = Rank2
	}
	if fromSquare.Rank == Rank6 && toSquare.Rank == Rank5 && pos.Piece(fromSquare).Type != Pawn {
		fromSquare.Rank = Rank7
	}

	m := Move{fromSquare, toSquare, NoPieceType}

	if len(san) == 4 {
		// Possible promotion
		pt, err := parsePieceType(san[3])
		if err != nil {
			return Move{}, fmt.Errorf("could not parse pawn advance: %w", err)
		}
		m.Promotion = pt
	}

	return m, nil
}

func parseNormalMove(san string, pos *Position) (Move, error) {
	pt, err := parsePieceType(san[0])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse normal move: %w", err)
	}
	if pt == Pawn || pt == NoPieceType {
		return Move{}, errors.New("could not parse normal move, got Pawn or NoPieceType which should use a different syntax")
	}
	pieceToMove := Piece{pos.SideToMove, pt}
	var toSquare Square
	if len(san) == 3 {
		err = toSquare.UnmarshalText([]byte(san[1:3]))
	} else if len(san) == 4 {
		err = toSquare.UnmarshalText([]byte(san[2:4]))
	} else {
		panic("unexpected condition")
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse normal move: %w", err)
	}
	fromSquare, err := getNormalMoveFromSquare(pieceToMove, toSquare, pos)
	if err != nil {
		return Move{}, fmt.Errorf("could not parse normal move: %w", err)
	}
	return Move{fromSquare, toSquare, NoPieceType}, nil
}

func getNormalMoveFromSquare(pieceToMove Piece, toSquare Square, pos *Position) (Square, error) {
	bb := pos.Bitboard(pieceToMove)
	possibleFromSquares := []Square{}

	for bb != 0 {
		singlePieceBB := Bitboard(1 << bits.TrailingZeros(uint(bb)))
		bb &^= singlePieceBB
		attacks := getPieceAttacks(singlePieceBB, pieceToMove.Type, pos)
		if attacks.Square(toSquare) == 1 {
			possibleFromSquares = append(possibleFromSquares, indexToSquare(bits.TrailingZeros(uint(singlePieceBB))))
		}
	}

	if len(possibleFromSquares) == 0 {
		return NoSquare, errors.New("could not determine from square, no squares seem adequate")
	}
	if len(possibleFromSquares) == 1 {
		return possibleFromSquares[0], nil
	}

	squaresNoCheck := []Square{}
	for _, s := range possibleFromSquares {
		testPosition := pos.Copy()
		testMove := Move{s, toSquare, NoPieceType}
		testPosition.Move(testMove)
		testPosition.SideToMove = pos.SideToMove
		if !testPosition.IsCheck() {
			squaresNoCheck = append(squaresNoCheck, s)
		}
	}

	if len(squaresNoCheck) == 0 {
		return NoSquare, errors.New("could not determine from square, no squares seem adequate")
	} else if len(squaresNoCheck) > 1 {
		return NoSquare, errors.New("could not determine from square, multiple squares seem adequate")
	} else {
		return squaresNoCheck[0], nil
	}
}

func getPieceAttacks(bb Bitboard, t PieceType, pos *Position) Bitboard {
	switch t {
	case Bishop:
		return bb.BishopAttacks(pos.OccupiedBitboard())
	case King:
		return bb.KingAttacks()
	case Knight:
		return bb.KnightAttacks()
	case Queen:
		return bb.QueenAttacks(pos.OccupiedBitboard())
	case Rook:
		return bb.RookAttacks(pos.OccupiedBitboard())
	default:
		panic(fmt.Sprintf("unexpected chess.PieceType: %#v", t))
	}
}

func parseSquareDisambiguation(san string) (Move, error) {
	fromSquare := Square{}
	err := fromSquare.UnmarshalText([]byte(san[1:3]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse square disambiguation form: %w", err)
	}

	var toSquare Square
	if strings.Contains(san, "x") && len(san) >= 6 {
		err = toSquare.UnmarshalText([]byte(san[4:6]))
	} else if len(san) >= 5 {
		err = toSquare.UnmarshalText([]byte(san[3:5]))
	} else {
		err = fmt.Errorf("malformed square disambiguation")
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse square disambiguation form: %w", err)
	}

	return Move{fromSquare, toSquare, NoPieceType}, nil
}

func parseFileDisambiguation(san string, pos *Position) (Move, error) {
	// this function should probably be broken down into various sub functions
	var toSquare Square
	var err error
	if strings.Contains(san, "x") || strings.Contains(san, "X") {
		if len(san) < 5 {
			return Move{}, errors.New("could not parse file disambiguation form, san to short")
		}
		err = toSquare.UnmarshalText([]byte(san[3:5]))
	} else {
		if len(san) < 4 {
			return Move{}, errors.New("could not parse file disambiguation form, san to short")
		}
		err = toSquare.UnmarshalText([]byte(san[2:4]))
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse file disambiguation form: %w", err)
	}

	fromFile, err := parseFile(san[1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse file disambiguation form: %w", err)
	}
	pt, err := parsePieceType(san[0])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse file disambiguation form: %w", err)
	}
	if pt == Pawn || pt == NoPieceType {
		return Move{}, errors.New("could not parse file disambiguation form, got Pawn or NoPieceType which should use a different syntax")
	}
	pieceToMove := Piece{pos.SideToMove, pt}
	possibleFromRanks := []Rank{}
	for r := Rank1; r <= Rank8; r++ {
		possibleFromSquare := Square{fromFile, r}
		if pos.Piece(possibleFromSquare) != pieceToMove {
			continue
		}
		attacks := getPieceAttacks(1<<squareToIndex(possibleFromSquare), pt, pos)
		if attacks.Square(toSquare) != 1 {
			continue
		}
		possibleFromRanks = append(possibleFromRanks, r)
	}

	if lenRanks := len(possibleFromRanks); lenRanks == 0 {
		return Move{}, errors.New("could not parse file disambiguation form, no ranks seem adequate")
	} else if lenRanks == 1 {
		return Move{Square{fromFile, possibleFromRanks[0]}, toSquare, NoPieceType}, nil
	}

	squaresNoCheck := []Square{}
	for _, r := range possibleFromRanks {
		testPosition := pos.Copy()
		testMove := Move{Square{fromFile, r}, toSquare, NoPieceType}
		testPosition.Move(testMove)
		testPosition.SideToMove = pos.SideToMove
		if !testPosition.IsCheck() {
			squaresNoCheck = append(squaresNoCheck, Square{fromFile, r})
		}
	}

	if lenSquares := len(squaresNoCheck); lenSquares == 0 {
		return Move{}, errors.New("could not parse file disambiguation form, no ranks seem adequate")
	} else if lenSquares > 1 {
		return Move{}, errors.New("could not parse file disambiguation form, multiple ranks seem adequate")
	} else {
		return Move{squaresNoCheck[0], toSquare, NoPieceType}, nil
	}

}

func parseRankDisambiguation(san string, pos *Position) (Move, error) {
	// this function should probably be broken down into various sub functions
	var toSquare Square
	var err error
	if strings.Contains(san, "x") || strings.Contains(san, "X") {
		if len(san) < 5 {
			return Move{}, errors.New("could not parse rank disambiguation form, san to short")
		}
		err = toSquare.UnmarshalText([]byte(san[3:5]))
	} else {
		if len(san) < 4 {
			return Move{}, errors.New("could not parse rank disambiguation form, san to short")
		}
		err = toSquare.UnmarshalText([]byte(san[2:4]))
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse rank disambiguation form: %w", err)
	}

	fromRank, err := parseRank(san[1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse rank disambiguation form: %w", err)
	}
	pt, err := parsePieceType(san[0])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse rank disambiguation form: %w", err)
	}
	if pt == Pawn || pt == NoPieceType {
		return Move{}, errors.New("could not parse rank disambiguation form, got Pawn or NoPieceType which should use a different syntax")
	}
	pieceToMove := Piece{pos.SideToMove, pt}
	possibleFromFiles := []File{}
	for f := FileA; f <= FileH; f++ {
		possibleFromSquare := Square{f, fromRank}
		if pos.Piece(possibleFromSquare) != pieceToMove {
			continue
		}
		attacks := getPieceAttacks(1<<squareToIndex(possibleFromSquare), pt, pos)
		if attacks.Square(toSquare) != 1 {
			continue
		}
		possibleFromFiles = append(possibleFromFiles, f)
	}

	if lenFiles := len(possibleFromFiles); lenFiles == 0 {
		return Move{}, errors.New("could not parse rank disambiguation form, no files seem adequate")
	} else if lenFiles == 1 {
		return Move{Square{possibleFromFiles[0], fromRank}, toSquare, NoPieceType}, nil
	}

	squaresNoCheck := []Square{}
	for _, f := range possibleFromFiles {
		testPosition := pos.Copy()
		testMove := Move{Square{f, fromRank}, toSquare, NoPieceType}
		testPosition.Move(testMove)
		testPosition.SideToMove = pos.SideToMove
		if !testPosition.IsCheck() {
			squaresNoCheck = append(squaresNoCheck, Square{f, fromRank})
		}
	}

	if lenSquares := len(squaresNoCheck); lenSquares == 0 {
		return Move{}, errors.New("could not parse rank disambiguation form, no files seem adequate")
	} else if lenSquares > 1 {
		return Move{}, errors.New("could not parse rank disambiguation form, multiple files seem adequate")
	} else {
		return Move{squaresNoCheck[0], toSquare, NoPieceType}, nil
	}
}
