package chess

import (
	"errors"
	"fmt"
	"math"
	"unicode"
)

type File uint8

const (
	NoFile File = iota
	FileA
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

func (f File) String() string {
	if !isValidFile(f) {
		return "INVALID FILE"
	}
	if f == NoFile {
		return "NO FILE"
	}
	return string('@' + rune(f))
}

func isValidFile(f File) bool {
	return f <= FileH
}

type Rank uint8

const (
	NoRank Rank = iota
	Rank1
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

func (r Rank) String() string {
	if !isValidRank(r) {
		return "INVALID RANK"
	}
	if r == NoRank {
		return "NO RANK"
	}
	return string('0' + rune(r))
}

func isValidRank(r Rank) bool {
	return r <= Rank8
}

type Square struct {
	File File
	Rank Rank
}

var (
	NoSquare = Square{}

	A8 = Square{FileA, Rank8}
	B8 = Square{FileB, Rank8}
	C8 = Square{FileC, Rank8}
	D8 = Square{FileD, Rank8}
	E8 = Square{FileE, Rank8}
	F8 = Square{FileF, Rank8}
	G8 = Square{FileG, Rank8}
	H8 = Square{FileH, Rank8}

	A7 = Square{FileA, Rank7}
	B7 = Square{FileB, Rank7}
	C7 = Square{FileC, Rank7}
	D7 = Square{FileD, Rank7}
	E7 = Square{FileE, Rank7}
	F7 = Square{FileF, Rank7}
	G7 = Square{FileG, Rank7}
	H7 = Square{FileH, Rank7}

	A6 = Square{FileA, Rank6}
	B6 = Square{FileB, Rank6}
	C6 = Square{FileC, Rank6}
	D6 = Square{FileD, Rank6}
	E6 = Square{FileE, Rank6}
	F6 = Square{FileF, Rank6}
	G6 = Square{FileG, Rank6}
	H6 = Square{FileH, Rank6}

	A5 = Square{FileA, Rank5}
	B5 = Square{FileB, Rank5}
	C5 = Square{FileC, Rank5}
	D5 = Square{FileD, Rank5}
	E5 = Square{FileE, Rank5}
	F5 = Square{FileF, Rank5}
	G5 = Square{FileG, Rank5}
	H5 = Square{FileH, Rank5}

	A4 = Square{FileA, Rank4}
	B4 = Square{FileB, Rank4}
	C4 = Square{FileC, Rank4}
	D4 = Square{FileD, Rank4}
	E4 = Square{FileE, Rank4}
	F4 = Square{FileF, Rank4}
	G4 = Square{FileG, Rank4}
	H4 = Square{FileH, Rank4}

	A3 = Square{FileA, Rank3}
	B3 = Square{FileB, Rank3}
	C3 = Square{FileC, Rank3}
	D3 = Square{FileD, Rank3}
	E3 = Square{FileE, Rank3}
	F3 = Square{FileF, Rank3}
	G3 = Square{FileG, Rank3}
	H3 = Square{FileH, Rank3}

	A2 = Square{FileA, Rank2}
	B2 = Square{FileB, Rank2}
	C2 = Square{FileC, Rank2}
	D2 = Square{FileD, Rank2}
	E2 = Square{FileE, Rank2}
	F2 = Square{FileF, Rank2}
	G2 = Square{FileG, Rank2}
	H2 = Square{FileH, Rank2}

	A1 = Square{FileA, Rank1}
	B1 = Square{FileB, Rank1}
	C1 = Square{FileC, Rank1}
	D1 = Square{FileD, Rank1}
	E1 = Square{FileE, Rank1}
	F1 = Square{FileF, Rank1}
	G1 = Square{FileG, Rank1}
	H1 = Square{FileH, Rank1}
)

func (s Square) String() string {
	if !isValidSquare(s) {
		return "INVALID SQUARE"
	}
	if s == NoSquare {
		return "-"
	}
	return s.File.String() + s.Rank.String()
}

func isValidSquare(s Square) bool {
	if s.File == NoFile || s.Rank == NoRank {
		if s.File != NoFile || s.Rank != NoRank {
			return false
		}
	}
	return isValidFile(s.File) && isValidRank(s.Rank)
}

func ParseSquare(s string) (Square, error) {
	if s == "-" {
		return NoSquare, nil
	}
	runes := []rune(s)
	if len(runes) != 2 {
		return NoSquare, errors.New("invalid square - string is not two characters long")
	}
	file, err := parseFile(runes[0])
	if err != nil {
		return NoSquare, fmt.Errorf("invalid square: %w", err)
	}
	rank, err := parseRank(runes[1])
	if err != nil {
		return NoSquare, fmt.Errorf("invalid square: %w", err)
	}
	return Square{file, rank}, nil
}

func parseFile(r rune) (File, error) {
	var f File = File(unicode.ToUpper(r) - '@')
	if !isValidFile(f) || f == NoFile {
		return 0, errors.New("invalid file")
	}
	return f, nil
}

func parseRank(r rune) (Rank, error) {
	var rank Rank = Rank(r - '0')
	if !isValidRank(rank) || rank == NoRank {
		return 0, errors.New("invalid rank")
	}
	return rank, nil
}

func squareToLeft(s Square) Square {
	s.File--
	if !isValidSquare(s) {
		return NoSquare
	}
	return s
}

func squareToRight(s Square) Square {
	s.File++
	if !isValidSquare(s) {
		return NoSquare
	}
	return s
}

func squareAbove(s Square) Square {
	s.Rank++
	if !isValidSquare(s) {
		return NoSquare
	}
	return s
}

func squareBelow(s Square) Square {
	s.Rank--
	if !isValidSquare(s) {
		return NoSquare
	}
	return s
}

func ChebyshevDistance(s1 Square, s2 Square) uint8 {
	if !isValidSquare(s1) || !isValidSquare(s2) {
		return math.MaxUint8
	}
	var fileDistance uint8
	if s1.File-s2.File < 9 {
		fileDistance = uint8(s1.File) - uint8(s2.File)
	} else {
		fileDistance = uint8(s2.File) - uint8(s1.File)
	}
	var rankDistance uint8
	if s1.Rank-s2.Rank < 9 {
		rankDistance = uint8(s1.Rank) - uint8(s2.Rank)
	} else {
		rankDistance = uint8(s2.Rank) - uint8(s1.Rank)
	}
	return min(fileDistance, rankDistance)
}

func ManhattanDistance(s1 Square, s2 Square) uint8 {
	if !isValidSquare(s1) || !isValidSquare(s2) {
		return math.MaxUint8
	}
	var fileDistance uint8
	if s1.File-s2.File < 9 {
		fileDistance = uint8(s1.File) - uint8(s2.File)
	} else {
		fileDistance = uint8(s2.File) - uint8(s1.File)
	}
	var rankDistance uint8
	if s1.Rank-s2.Rank < 9 {
		rankDistance = uint8(s1.Rank) - uint8(s2.Rank)
	} else {
		rankDistance = uint8(s2.Rank) - uint8(s1.Rank)
	}
	return fileDistance + rankDistance
}
