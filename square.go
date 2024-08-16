package chess

import (
	"errors"
	"fmt"
)

type File uint8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

func (f File) String() string {
	if !isFileValid(f) {
		return "INVALID FILE"
	}

	return string('A' + rune(f))
}

func isFileValid(f File) bool {
	return f <= FileH
}

type Rank uint8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

func (r Rank) String() string {
	if !isRankValid(r) {
		return "INVALID RANK"
	}
	return string('1' + rune(r))
}

func isRankValid(r Rank) bool {
	return r <= Rank8
}

type Square struct {
	File File
	Rank Rank
}

var (
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

func (s *Square) String() string {
	if s == nil {
		return "-"
	}
	if !isSquareValid(*s) {
		return "INVALID SQUARE"
	}
	return s.File.String() + s.Rank.String()
}

func isSquareValid(s Square) bool {
	return isFileValid(s.File) && isRankValid(s.Rank)
}

func ParseSquare(s string) (*Square, error) {
	if s == "-" {
		return nil, nil
	}

	runes := []rune(s)
	if len(runes) != 2 {
		return nil, errors.New("invalid square - string is not two characters long")
	}
	file, err := parseFile(runes[0])
	if err != nil {
		return nil, fmt.Errorf("invalid square: %w", err)
	}
	rank, err := parseRank(runes[1])
	if err != nil {
		return nil, fmt.Errorf("invalid square: %w", err)
	}
	return &Square{file, rank}, nil
}

func parseFile(r rune) (File, error) {
	var f File = File(r - 'A')
	if !isFileValid(f) {
		return 0, errors.New("invalid file")
	}
	return f, nil
}

func parseRank(r rune) (Rank, error) {
	var rank Rank = Rank(r - '1')
	if !isRankValid(rank) {
		return 0, errors.New("invalid rank")
	}
	return rank, nil
}
