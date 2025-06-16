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
)

// File is a vertical column of squares as seen on a chess board. The zero value is [NoFile], and the files A-H can be represented.
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

// String returns a single lowercase letter representation of f if valid, else an error string. NoFile returns "-".
func (f File) String() string {
	switch f {
	case NoFile:
		return "-"
	case FileA:
		return "a"
	case FileB:
		return "b"
	case FileC:
		return "c"
	case FileD:
		return "d"
	case FileE:
		return "e"
	case FileF:
		return "f"
	case FileG:
		return "g"
	case FileH:
		return "h"
	default:
		return fmt.Sprintf("Unknown File %d", f)
	}
}

// Rank is a horizontal row of squares as seen on a chess board. The zero value is [NoRank], and the ranks 1-8 can be represented.
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

// String returns a single number representation of r if valid, else an error string. NoRank returns "-".
func (r Rank) String() string {
	switch r {
	case NoRank:
		return "-"
	case Rank1:
		return "1"
	case Rank2:
		return "2"
	case Rank3:
		return "3"
	case Rank4:
		return "4"
	case Rank5:
		return "5"
	case Rank6:
		return "6"
	case Rank7:
		return "7"
	case Rank8:
		return "8"
	default:
		return fmt.Sprintf("Unknown Rank %d", r)
	}
}

// Square represents one of 64 squares on a chess board. The zero value represents [NoSquare].
type Square struct {
	File File
	Rank Rank
}

var (
	NoSquare = Square{File: NoFile, Rank: NoRank}

	A1 = Square{File: FileA, Rank: Rank1}
	A2 = Square{File: FileA, Rank: Rank2}
	A3 = Square{File: FileA, Rank: Rank3}
	A4 = Square{File: FileA, Rank: Rank4}
	A5 = Square{File: FileA, Rank: Rank5}
	A6 = Square{File: FileA, Rank: Rank6}
	A7 = Square{File: FileA, Rank: Rank7}
	A8 = Square{File: FileA, Rank: Rank8}

	B1 = Square{File: FileB, Rank: Rank1}
	B2 = Square{File: FileB, Rank: Rank2}
	B3 = Square{File: FileB, Rank: Rank3}
	B4 = Square{File: FileB, Rank: Rank4}
	B5 = Square{File: FileB, Rank: Rank5}
	B6 = Square{File: FileB, Rank: Rank6}
	B7 = Square{File: FileB, Rank: Rank7}
	B8 = Square{File: FileB, Rank: Rank8}

	C1 = Square{File: FileC, Rank: Rank1}
	C2 = Square{File: FileC, Rank: Rank2}
	C3 = Square{File: FileC, Rank: Rank3}
	C4 = Square{File: FileC, Rank: Rank4}
	C5 = Square{File: FileC, Rank: Rank5}
	C6 = Square{File: FileC, Rank: Rank6}
	C7 = Square{File: FileC, Rank: Rank7}
	C8 = Square{File: FileC, Rank: Rank8}

	D1 = Square{File: FileD, Rank: Rank1}
	D2 = Square{File: FileD, Rank: Rank2}
	D3 = Square{File: FileD, Rank: Rank3}
	D4 = Square{File: FileD, Rank: Rank4}
	D5 = Square{File: FileD, Rank: Rank5}
	D6 = Square{File: FileD, Rank: Rank6}
	D7 = Square{File: FileD, Rank: Rank7}
	D8 = Square{File: FileD, Rank: Rank8}

	E1 = Square{File: FileE, Rank: Rank1}
	E2 = Square{File: FileE, Rank: Rank2}
	E3 = Square{File: FileE, Rank: Rank3}
	E4 = Square{File: FileE, Rank: Rank4}
	E5 = Square{File: FileE, Rank: Rank5}
	E6 = Square{File: FileE, Rank: Rank6}
	E7 = Square{File: FileE, Rank: Rank7}
	E8 = Square{File: FileE, Rank: Rank8}

	F1 = Square{File: FileF, Rank: Rank1}
	F2 = Square{File: FileF, Rank: Rank2}
	F3 = Square{File: FileF, Rank: Rank3}
	F4 = Square{File: FileF, Rank: Rank4}
	F5 = Square{File: FileF, Rank: Rank5}
	F6 = Square{File: FileF, Rank: Rank6}
	F7 = Square{File: FileF, Rank: Rank7}
	F8 = Square{File: FileF, Rank: Rank8}

	G1 = Square{File: FileG, Rank: Rank1}
	G2 = Square{File: FileG, Rank: Rank2}
	G3 = Square{File: FileG, Rank: Rank3}
	G4 = Square{File: FileG, Rank: Rank4}
	G5 = Square{File: FileG, Rank: Rank5}
	G6 = Square{File: FileG, Rank: Rank6}
	G7 = Square{File: FileG, Rank: Rank7}
	G8 = Square{File: FileG, Rank: Rank8}

	H1 = Square{File: FileH, Rank: Rank1}
	H2 = Square{File: FileH, Rank: Rank2}
	H3 = Square{File: FileH, Rank: Rank3}
	H4 = Square{File: FileH, Rank: Rank4}
	H5 = Square{File: FileH, Rank: Rank5}
	H6 = Square{File: FileH, Rank: Rank6}
	H7 = Square{File: FileH, Rank: Rank7}
	H8 = Square{File: FileH, Rank: Rank8}
)

// AllSquares is an array of all the squares on a chess board for convenience.
var AllSquares = [64]Square{
	A1, A2, A3, A4, A5, A6, A7, A8,
	B1, B2, B3, B4, B5, B6, B7, B8,
	C1, C2, C3, C4, C5, C6, C7, C8,
	D1, D2, D3, D4, D5, D6, D7, D8,
	E1, E2, E3, E4, E5, E6, E7, E8,
	F1, F2, F3, F4, F5, F6, F7, F8,
	G1, G2, G3, G4, G5, G6, G7, G8,
	H1, H2, H3, H4, H5, H6, H7, H8,
}

// MarshalText is an implementation of the [encoding.TextMarshaler] interface. It provides the square in the form "a1". An error is returned if the square is not valid. [NoSquare] produces "-". See also [Square.String]
func (s Square) MarshalText() (text []byte, err error) {
	if s == NoSquare {
		return []byte{'-'}, nil
	}
	if !squareOnBoard(s) {
		return nil, fmt.Errorf("cannot marshal invalid square %#v", s)
	}
	return []byte{s.File.String()[0], s.Rank.String()[0]}, nil
}

// String provides a two letter textual representation of s in the form "a1". [NoSquare] produces "-". If s is invalid and error string is returned.
func (s Square) String() string {
	text, err := s.MarshalText()
	if err != nil {
		return fmt.Sprintf("Unknown Square %#v", s)
	}
	return string(text)
}

// UnmarshalText is an implementation of the [encoding.TextUnmarshaler] interface. It expects text in the form of a two character square like "a1" or "A1". A "-" will provide [NoSquare]. All other cases result in an error, and s is unmodified.
func (s *Square) UnmarshalText(text []byte) error {
	if string(text) == "-" {
		*s = NoSquare
		return nil
	}

	if len(text) != 2 {
		return fmt.Errorf("could not unmarshal square %q, text should have length 2", text)
	}
	f, err := parseFile(text[0])
	if err != nil {
		return fmt.Errorf("could not parse square %q: %w", text, err)
	}
	r, err := parseRank(text[1])
	if err != nil {
		return fmt.Errorf("could not parse square %q, %w", text, err)
	}
	s.File = f
	s.Rank = r
	return nil
}

func parseFile(f byte) (File, error) {
	switch f {
	case 'a', 'A':
		return FileA, nil
	case 'b', 'B':
		return FileB, nil
	case 'c', 'C':
		return FileC, nil
	case 'd', 'D':
		return FileD, nil
	case 'e', 'E':
		return FileE, nil
	case 'f', 'F':
		return FileF, nil
	case 'g', 'G':
		return FileG, nil
	case 'h', 'H':
		return FileH, nil
	default:
		return NoFile, fmt.Errorf("could not parse file %q", f)
	}
}

func parseRank(r byte) (Rank, error) {
	switch r {
	case '1':
		return Rank1, nil
	case '2':
		return Rank2, nil
	case '3':
		return Rank3, nil
	case '4':
		return Rank4, nil
	case '5':
		return Rank5, nil
	case '6':
		return Rank6, nil
	case '7':
		return Rank7, nil
	case '8':
		return Rank8, nil
	default:
		return NoRank, fmt.Errorf("could not parse rank %q", r)
	}
}

func squareOnBoard(s Square) bool {
	return s.File > NoFile && s.File <= FileH &&
		s.Rank > NoRank && s.Rank <= Rank8
}
