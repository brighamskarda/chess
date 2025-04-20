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
	"strings"
)

// Move represents a UCI chess move.
type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
}

// ParseUCI parses a move string of the form <FromSquare><ToSquare><OptionalPromotion>. (e.g. a2c3 or H2H1q. Returns an error if it could not parse.
func ParseUCIMove(uci string) (Move, error) {
	uci = strings.ToLower(uci)
	if len(uci) < 4 || len(uci) > 5 {
		return Move{}, errors.New("uci move string not 4 or 5 characters long")
	}
	fromSquare := parseSquare(uci[0:2])
	toSquare := parseSquare(uci[2:4])
	promotion := NoPieceType
	if len(uci) == 5 {
		promotion = parsePieceType(uci[4:5])
	}
	if fromSquare == NoSquare || toSquare == NoSquare {
		return Move{}, fmt.Errorf("could not parse move square, %q", uci)
	}

	return Move{fromSquare, toSquare, promotion}, nil
}

// String provides a UCI compatible string of the move in the form <FromSquare><ToSquare><OptionalPromotion>
func (m Move) String() string {
	promotion := m.Promotion.String()
	if promotion == "-" {
		promotion = ""
	}
	return m.FromSquare.String() + m.ToSquare.String() + promotion
}
