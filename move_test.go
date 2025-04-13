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

func TestMoveString(t *testing.T) {
	expected := "a1b2"
	actual := Move{A1, B2, NoPieceType}.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
	expected = "h2c1q"
	actual = Move{H2, C1, Queen}.String()
	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}

func TestParseUCIMove(t *testing.T) {
	expected := Move{A1, A2, NoPieceType}
	actual, err := ParseUCIMove("a1a2")
	if expected != actual {
		t.Errorf("incorrect result: expected %v, got %v", expected, actual)
	}
	if err != nil {
		t.Errorf("incorrect result for \"a1a2\": expected err to be nil")
	}

	expected = Move{H2, C1, Queen}
	actual, err = ParseUCIMove("h2c1q")
	if expected != actual {
		t.Errorf("incorrect result: expected %v, got %v", expected, actual)
	}
	if err != nil {
		t.Errorf("incorrect result for \"h2c1q\": expected err to be nil")
	}
}

func TestParseUCIMoveErr(t *testing.T) {
	_, err := ParseUCIMove("a1c")
	if err == nil {
		t.Error("Expected err to be nil.")
	}
}
