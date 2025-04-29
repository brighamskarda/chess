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

// Move represents a UCI chess move.
type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
}

// String provides a UCI compatible string of the move in the form <FromSquare><ToSquare><OptionalPromotion>
func (m Move) String() string {
	promotion := m.Promotion.String()
	if promotion == "-" {
		promotion = ""
	}
	return m.FromSquare.String() + m.ToSquare.String() + promotion
}

// ParseUCI parses a move string of the form <FromSquare><ToSquare><OptionalPromotion>. (e.g. a2c3 or H2H1q. Returns an error if it could not parse.
func ParseUCIMove(uci string) (Move, error) {
	uci = strings.ToLower(uci)
	if len(uci) < 4 || len(uci) > 5 {
		return Move{}, errors.New("uci move string not 4 or 5 characters long")
	}
	fromSquare, _ := ParseSquare(uci[0:2])
	toSquare, _ := ParseSquare(uci[2:4])
	promotion := NoPieceType
	if len(uci) == 5 {
		prom, err := parsePieceType(uci[4:5])
		if err != nil {
			fmt.Errorf("could not parse move promotion, %q", uci)
		}
		promotion = prom
	}
	if fromSquare == NoSquare || toSquare == NoSquare {
		return Move{}, fmt.Errorf("could not parse move square, %q", uci)
	}

	return Move{fromSquare, toSquare, promotion}, nil
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

// ParseSANMove parses a chess move provided in Standard Algebraic Notation. Standard Algebraic notation requires position information to know how moves should be disambiguated.
//
// See the pgn specification to know the exact functionality.
func ParseSANMove(san string, pos *Position) (Move, error) {
	// This comment is a silent cry for help. I only implemented this because the pgn standard uses it. SAN was not designed to be parsed by computers. This was a major effort that could have been easily avoided by the much superior UCI move notation.
	if pos.SideToMove != White && pos.SideToMove != Black {
		return Move{}, errors.New("could not parse SAN move: pos.SideToMove not set")
	}

	// These characters don't matter for this parsing.
	san = strings.ReplaceAll(san, "+", "")
	san = strings.ReplaceAll(san, "#", "")
	switch classifySANMove(san) {
	case fileDisambiguation:
		return parseFileDisambiguation(san, pos)
	case rankDisambiguation:
		return parseRankDisambiguation(san, pos)
	case squareDisambiguation:
		return parseSquareDisambiguation(san, pos)
	case normalMove:
		return parseNormalMove(san, pos)
	case pawnAdvance:
		return parsePawnAdvance(san, pos)
	case pawnCapture:
		return parsePawnCapture(san, pos)
	case castleMove:
		return parseCastleMove(san, pos)
	default:
		panic("unexpected chess.sanMoveType")
	}
}

func classifySANMove(san string) sanMoveType {
	sanRunes := []rune(san)
	switch len(sanRunes) {
	case 2:
		return pawnAdvance
	case 3:
		if unicode.ToLower(sanRunes[0]) == 'o' && sanRunes[1] == '-' && unicode.ToLower(sanRunes[2]) == 'o' {
			return castleMove
		}
		return normalMove
	case 4:
		return classifySANMove4(sanRunes)
	case 5:
		return classifySANMove5(sanRunes)
	case 6:
		if sanRunes[4] == '=' {
			return pawnCapture
		}
		return squareDisambiguation
	default:
		return unknown
	}
}

func classifySANMove4(san []rune) sanMoveType {
	if san[0] == 'B' || slices.Contains([]rune{'r', 'n', 'q', 'k'}, unicode.ToLower(san[0])) {
		if unicode.ToLower(san[1]) == 'x' {
			return normalMove
		} else if slices.Contains([]rune{'1', '2', '3', '4', '5', '6', '7', '8'}, san[1]) {
			return rankDisambiguation
		} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(san[1])) {
			return fileDisambiguation
		} else {
			return unknown
		}
	} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(san[0])) {
		if san[2] == '=' {
			return pawnAdvance
		} else if unicode.ToLower(san[1]) == 'x' {
			return pawnCapture
		} else {
			return unknown
		}
	} else {
		return unknown
	}
}

func classifySANMove5(san []rune) sanMoveType {
	if unicode.ToLower(san[0]) == 'o' && san[1] == '-' && unicode.ToLower(san[2]) == 'o' && san[3] == '-' && unicode.ToLower(san[4]) == 'o' {
		return castleMove
	} else if unicode.ToLower(san[2]) == 'x' {
		if slices.Contains([]rune{'1', '2', '3', '4', '5', '6', '7', '8'}, san[1]) {
			return rankDisambiguation
		} else if slices.Contains([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}, unicode.ToLower(san[1])) {
			return fileDisambiguation
		} else {
			return unknown
		}
	} else {
		return squareDisambiguation
	}
}

func parseCastleMove(san string, pos *Position) (Move, error) {
	switch strings.ToLower(san) {
	case "o-o":
		switch pos.SideToMove {
		case White:
			return Move{E1, G1, NoPieceType}, nil
		case Black:
			return Move{E8, G8, NoPieceType}, nil
		}
	case "o-o-o":
		switch pos.SideToMove {
		case White:
			return Move{E1, C1, NoPieceType}, nil
		case Black:
			return Move{E8, C8, NoPieceType}, nil
		}
	}
	panic("unexpected condition: bad san in parseCastleMove")

}

func parsePawnCapture(san string, pos *Position) (Move, error) {
	m := Move{}

	toSquare, err := ParseSquare(string(san[2:4]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN pawn capture %q: %w", san, err)
	}
	m.ToSquare = toSquare

	fromFile, err := parseFile(san[0:1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN pawn capture %q: %w", san, err)
	}
	switch pos.SideToMove {
	case White:
		m.FromSquare = Square{fromFile, toSquare.Rank - 1}
	case Black:
		m.FromSquare = Square{fromFile, toSquare.Rank + 1}
	}

	if len(san) == 6 {
		// Pawn capture with promotion
		pt, err := parsePieceType(san[5:6])
		if err != nil {
			return Move{}, fmt.Errorf("could not parse SAN pawn capture %q: could not parse promotion", san)
		}
		m.Promotion = pt
	}
	return m, nil
}

func parsePawnAdvance(san string, pos *Position) (Move, error) {
	toSquare, err := ParseSquare(san[0:2])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN pawn move %q", san)
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
		pt, err := parsePieceType(san[3:4])
		if err != nil {
			return Move{}, fmt.Errorf("could not parse SAN pawn move %q: could not parse promotion", san)
		}
		m.Promotion = pt
	}

	return m, nil
}

func parseNormalMove(san string, pos *Position) (Move, error) {
	pt, err := parsePieceType(san[0:1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: %w", san, err)
	}
	pieceToMove := Piece{pos.SideToMove, pt}
	var toSquare Square
	if len(san) == 3 {
		toSquare, err = ParseSquare(san[1:3])
	} else if len(san) == 4 {
		toSquare, err = ParseSquare(san[2:4])
	} else {
		panic("unexpected condition")
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not parse to square", san)
	}
	fromSquare, err := getNormalMoveFromSquare(pieceToMove, toSquare, pos)
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine from square", san)
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
		return NoSquare, errors.New("could not determine from square")
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

	if len(squaresNoCheck) != 1 {
		return NoSquare, errors.New("could not determine from square")
	} else {
		return squaresNoCheck[0], nil
	}
}

func getPieceAttacks(bb Bitboard, t PieceType, pos *Position) Bitboard {
	switch t {
	case Bishop:
		return bb.bishopAttacks(pos.OccupiedBitboard())
	case King:
		return bb.kingAttacks()
	case Knight:
		return bb.knightAttacks()
	case Queen:
		return bb.queenAttacks(pos.OccupiedBitboard())
	case Rook:
		return bb.rookAttacks(pos.OccupiedBitboard())
	default:
		panic(fmt.Sprintf("unexpected chess.PieceType: %#v", t))
	}
}

func parseSquareDisambiguation(san string, pos *Position) (Move, error) {
	fromSquare, err := ParseSquare(san[1:3])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not parse from square", san)
	}

	var toSquare Square
	if strings.Contains(san, "x") {
		toSquare, err = ParseSquare(san[4:6])
	} else {
		toSquare, err = ParseSquare(san[3:5])
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not parse to square", san)
	}

	return Move{fromSquare, toSquare, NoPieceType}, nil
}

func parseFileDisambiguation(san string, pos *Position) (Move, error) {
	// this function should probably be broken down into various sub functions
	var toSquare Square
	var err error
	if strings.Contains(san, "x") || strings.Contains(san, "X") {
		toSquare, err = ParseSquare(san[3:5])
	} else {
		toSquare, err = ParseSquare(san[2:4])
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not parse to square", san)
	}

	fromFile, err := parseFile(san[1:2])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine from square", san)
	}
	pt, err := parsePieceType(san[0:1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine from square", san)
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

	if len(possibleFromRanks) == 0 {
		return Move{}, errors.New("could not determine from square")
	}
	if len(possibleFromRanks) == 1 {
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

	if len(squaresNoCheck) != 1 {
		return Move{}, errors.New("could not determine from square")
	} else {
		return Move{squaresNoCheck[0], toSquare, NoPieceType}, nil
	}

}

func parseRankDisambiguation(san string, pos *Position) (Move, error) {
	// this function should probably be broken down into various sub functions
	var toSquare Square
	var err error
	if strings.Contains(san, "x") || strings.Contains(san, "X") {
		toSquare, err = ParseSquare(san[3:5])
	} else {
		toSquare, err = ParseSquare(san[2:4])
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not parse to square", san)
	}

	fromRank, err := parseRank(san[1:2])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine from square", san)
	}
	pt, err := parsePieceType(san[0:1])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move %q: could not determine from square", san)
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

	if len(possibleFromFiles) == 0 {
		return Move{}, errors.New("could not determine from square")
	}
	if len(possibleFromFiles) == 1 {
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

	if len(squaresNoCheck) != 1 {
		return Move{}, errors.New("could not determine from square")
	} else {
		return Move{squaresNoCheck[0], toSquare, NoPieceType}, nil
	}

}
