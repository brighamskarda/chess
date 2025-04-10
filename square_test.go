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

func TestSquareString(t *testing.T) {
	expected := "a1"
	actual := Square{FileA, Rank1}.String()

	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}

	expected = "--"
	actual = Square{NoFile, NoRank}.String()

	if expected != actual {
		t.Errorf("incorrect result: expected %q, got %q", expected, actual)
	}
}
