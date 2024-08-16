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
	if !fileIsValid(f) {
		return "INVALID FILE"
	}
	return string('A' + rune(f))
}

func fileIsValid(f File) bool {
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
	if !rankIsValid(r) {
		return "INVALID RANK"
	}
	return string('1' + rune(r))
}

func rankIsValid(r Rank) bool {
	return r <= Rank8
}

type Square struct {
	File File
	Rank Rank
}

func (s Square) String() string {
	if !squareIsValid(s) {
		return "INVALID SQUARE"
	}
	return s.File.String() + s.Rank.String()
}

func squareIsValid(s Square) bool {
	return fileIsValid(s.File) && rankIsValid(s.Rank)
}

func ParseSquare(s string) (Square, error) {
	runes := []rune(s)
	if len(runes) != 2 {
		return Square{}, errors.New("invalid square - string is not two characters long")
	}
	file, err := parseFile(runes[0])
	if err != nil {
		return Square{}, fmt.Errorf("invalid square: %w", err)
	}
	rank, err := parseRank(runes[1])
	if err != nil {
		return Square{}, fmt.Errorf("invalid square: %w", err)
	}
	return Square{file, rank}, nil
}

func parseFile(r rune) (File, error) {
	var f File = File(r - 'A')
	if !fileIsValid(f) {
		return 0, errors.New("invalid file")
	}
	return f, nil
}

func parseRank(r rune) (Rank, error) {
	var rank Rank = Rank(r - '1')
	if !rankIsValid(rank) {
		return 0, errors.New("invalid rank")
	}
	return rank, nil
}
