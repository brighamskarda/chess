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
	"strings"
)

// Color can be [NoColor], [White], or [Black].
type Color uint8

const (
	NoColor Color = iota
	White
	Black
)

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case NoColor:
		return "NoColor"
	case White:
		return "White"
	default:
		return "Unknown Color"
	}
}

func parseColor(s string) Color {
	switch strings.ToLower(s) {
	case "w":
		return White
	case "b":
		return Black
	default:
		return NoColor
	}
}
