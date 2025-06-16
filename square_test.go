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
	"testing"
)

func TestSquareMarshalText(t *testing.T) {
	expected := "a1"
	actual, err := Square{FileA, Rank1}.MarshalText()

	if err != nil {
		t.Errorf("got an error")
	}
	if expected != string(actual) {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	expected = "-"
	actual, err = Square{NoFile, NoRank}.MarshalText()

	if err != nil {
		t.Errorf("got an error")
	}
	if expected != string(actual) {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	_, err = Square{128, 255}.MarshalText()

	if err == nil {
		t.Errorf("did not get an error")
	}
}

func TestSquareUnmarshal(t *testing.T) {
	s := &Square{}
	err := s.UnmarshalText([]byte("a1"))
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if *s != A1 {
		t.Errorf("unmarshal provided incorrect results")
	}

	err = s.UnmarshalText([]byte("H8"))
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if *s != H8 {
		t.Errorf("unmarshal provided incorrect results")
	}

	err = s.UnmarshalText([]byte("-"))
	if err != nil {
		t.Errorf("got unexpected error")
	}
	if *s != NoSquare {
		t.Errorf("unmarshal provided incorrect results")
	}
}

func TestSquareUnmarshalError(t *testing.T) {
	s := &Square{FileC, Rank5}
	err := s.UnmarshalText([]byte(""))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}

	err = s.UnmarshalText([]byte("  "))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}

	err = s.UnmarshalText([]byte("a1-"))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}

	err = s.UnmarshalText([]byte("b"))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}

	err = s.UnmarshalText([]byte("c9"))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}

	err = s.UnmarshalText([]byte("i2"))
	if err == nil {
		t.Errorf("did not get error")
	}
	if *s != C5 {
		t.Errorf("unmarshal changed on error")
	}
}
