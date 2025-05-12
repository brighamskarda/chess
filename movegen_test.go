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

// Many of these tests were converted from V1 of the library, as such their naming and organization may not be the tidiest.

import (
	"slices"
	"testing"
)

func moveSetsEqual(m1 []Move, m2 []Move) bool {
	slices.SortFunc(m1, moveSortFunc)
	slices.SortFunc(m2, moveSortFunc)
	if len(m1) != len(m2) {
		return false
	}
	for _, move := range m1 {
		if !slices.Contains(m2, move) {
			return false
		}
	}
	return true
}

func moveSortFunc(a, b Move) int {
	if squareToIndex(a.FromSquare) < squareToIndex(b.FromSquare) {
		return -1
	}
	if squareToIndex(a.FromSquare) > squareToIndex(b.FromSquare) {
		return 1
	}
	if squareToIndex(a.ToSquare) < squareToIndex(b.ToSquare) {
		return -1
	}
	if squareToIndex(a.ToSquare) > squareToIndex(b.ToSquare) {
		return 1
	}
	return 0
}

func TestGeneratePseudoLegalMoves(t *testing.T) {
	defaultMoveSet := []Move{
		{A2, A3, NoPieceType},
		{A2, A4, NoPieceType},
		{B2, B3, NoPieceType},
		{B2, B4, NoPieceType},
		{C2, C3, NoPieceType},
		{C2, C4, NoPieceType},
		{D2, D3, NoPieceType},
		{D2, D4, NoPieceType},
		{E2, E3, NoPieceType},
		{E2, E4, NoPieceType},
		{F2, F3, NoPieceType},
		{F2, F4, NoPieceType},
		{G2, G3, NoPieceType},
		{G2, G4, NoPieceType},
		{H2, H3, NoPieceType},
		{H2, H4, NoPieceType},
		{B1, A3, NoPieceType},
		{B1, C3, NoPieceType},
		{G1, F3, NoPieceType},
		{G1, H3, NoPieceType},
	}

	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	moves := PseudoLegalMoves(pos)
	if !moveSetsEqual(defaultMoveSet, moves) {
		t.Errorf("incorrect result for default board: expected %v, got %v", defaultMoveSet, moves)
	}

	pos = &Position{}
	pos.UnmarshalText([]byte("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14"))
	expectedMoves := []Move{
		{A2, A3, NoPieceType},
		{A2, A4, NoPieceType},
		{B2, B3, NoPieceType},
		{B2, B4, NoPieceType},
		{C2, C3, NoPieceType},
		{F2, F3, NoPieceType},
		{F2, F4, NoPieceType},
		{H2, H3, NoPieceType},
		{H2, H4, NoPieceType},
		{C4, B5, NoPieceType},
		{C4, A6, NoPieceType},
		{C4, D5, NoPieceType},
		{C4, E6, NoPieceType},
		{C4, F7, NoPieceType},
		{C4, G8, NoPieceType},
		{C4, D3, NoPieceType},
		{C4, E2, NoPieceType},
		{C4, F1, NoPieceType},
		{C4, B3, NoPieceType},
		{G3, H3, NoPieceType},
		{G3, H4, NoPieceType},
		{G3, G4, NoPieceType},
		{G3, G5, NoPieceType},
		{G3, G6, NoPieceType},
		{G3, G7, NoPieceType},
		{G3, F4, NoPieceType},
		{G3, E5, NoPieceType},
		{G3, D6, NoPieceType},
		{G3, C7, NoPieceType},
		{G3, F3, NoPieceType},
		{G3, E3, NoPieceType},
		{G3, D3, NoPieceType},
		{G3, C3, NoPieceType},
		{G3, B3, NoPieceType},
		{G3, A3, NoPieceType},
		{B1, A3, NoPieceType},
		{B1, C3, NoPieceType},
		{B1, D2, NoPieceType},
		{C1, D2, NoPieceType},
		{C1, E3, NoPieceType},
		{C1, F4, NoPieceType},
		{C1, G5, NoPieceType},
		{C1, H6, NoPieceType},
		{E1, D1, NoPieceType},
		{E1, D2, NoPieceType},
		{E1, E2, NoPieceType},
		{E1, F1, NoPieceType},
		{H1, G1, NoPieceType},
		{H1, F1, NoPieceType},
	}
	moves = PseudoLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14: expected %v, got %v", expectedMoves, moves)
	}

	pos.SideToMove = Black
	expectedMoves = []Move{
		{A8, A7, NoPieceType},
		{A8, B8, NoPieceType},
		{A8, C8, NoPieceType},
		{A8, D8, NoPieceType},
		{E8, D8, NoPieceType},
		{E8, D7, NoPieceType},
		{E8, E7, NoPieceType},
		{E8, F7, NoPieceType},
		{E8, C8, NoPieceType},
		{F8, E7, NoPieceType},
		{F8, D6, NoPieceType},
		{F8, C5, NoPieceType},
		{F8, B4, NoPieceType},
		{F8, A3, NoPieceType},
		{H8, G8, NoPieceType},
		{C7, C6, NoPieceType},
		{C7, C5, NoPieceType},
		{G7, G6, NoPieceType},
		{G7, G5, NoPieceType},
		{H7, H6, NoPieceType},
		{H7, H5, NoPieceType},
		{B6, B5, NoPieceType},
		{F6, D7, NoPieceType},
		{F6, D5, NoPieceType},
		{F6, G4, NoPieceType},
		{F6, H5, NoPieceType},
		{F6, G8, NoPieceType},
		{A5, A4, NoPieceType},
		{A5, A3, NoPieceType},
		{A5, A2, NoPieceType},
		{A5, B5, NoPieceType},
		{A5, C5, NoPieceType},
		{A5, D5, NoPieceType},
		{A5, E5, NoPieceType},
		{A5, F5, NoPieceType},
		{A5, B4, NoPieceType},
		{A5, C3, NoPieceType},
		{A5, D2, NoPieceType},
		{A5, E1, NoPieceType},
		{E4, E3, NoPieceType},
	}
	moves = PseudoLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R b KQkq - 2 14: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateLegalMoves(t *testing.T) {
	defaultMoveSet := []Move{
		{A2, A3, NoPieceType},
		{A2, A4, NoPieceType},
		{B2, B3, NoPieceType},
		{B2, B4, NoPieceType},
		{C2, C3, NoPieceType},
		{C2, C4, NoPieceType},
		{D2, D3, NoPieceType},
		{D2, D4, NoPieceType},
		{E2, E3, NoPieceType},
		{E2, E4, NoPieceType},
		{F2, F3, NoPieceType},
		{F2, F4, NoPieceType},
		{G2, G3, NoPieceType},
		{G2, G4, NoPieceType},
		{H2, H3, NoPieceType},
		{H2, H4, NoPieceType},
		{B1, A3, NoPieceType},
		{B1, C3, NoPieceType},
		{G1, F3, NoPieceType},
		{G1, H3, NoPieceType},
	}

	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	moves := LegalMoves(pos)
	if !moveSetsEqual(defaultMoveSet, moves) {
		t.Errorf("incorrect result for default board: expected %v, got %v", defaultMoveSet, moves)
	}

	pos = &Position{}
	pos.UnmarshalText([]byte("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14"))
	expectedMoves := []Move{
		{B2, B4, NoPieceType},
		{C2, C3, NoPieceType},
		{G3, C3, NoPieceType},
		{B1, C3, NoPieceType},
		{B1, D2, NoPieceType},
		{C1, D2, NoPieceType},
		{E1, D1, NoPieceType},
		{E1, E2, NoPieceType},
		{E1, F1, NoPieceType},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14: expected %v, got %v", expectedMoves, moves)
	}

	pos.SideToMove = Black
	expectedMoves = []Move{
		{A8, A7, NoPieceType},
		{A8, B8, NoPieceType},
		{A8, C8, NoPieceType},
		{A8, D8, NoPieceType},
		{E8, D8, NoPieceType},
		{E8, D7, NoPieceType},
		{E8, E7, NoPieceType},
		{E8, C8, NoPieceType},
		{F8, E7, NoPieceType},
		{F8, D6, NoPieceType},
		{F8, C5, NoPieceType},
		{F8, B4, NoPieceType},
		{F8, A3, NoPieceType},
		{H8, G8, NoPieceType},
		{C7, C6, NoPieceType},
		{C7, C5, NoPieceType},
		{G7, G6, NoPieceType},
		{G7, G5, NoPieceType},
		{H7, H6, NoPieceType},
		{H7, H5, NoPieceType},
		{B6, B5, NoPieceType},
		{F6, D7, NoPieceType},
		{F6, D5, NoPieceType},
		{F6, G4, NoPieceType},
		{F6, H5, NoPieceType},
		{F6, G8, NoPieceType},
		{A5, A4, NoPieceType},
		{A5, A3, NoPieceType},
		{A5, A2, NoPieceType},
		{A5, B5, NoPieceType},
		{A5, C5, NoPieceType},
		{A5, D5, NoPieceType},
		{A5, E5, NoPieceType},
		{A5, F5, NoPieceType},
		{A5, B4, NoPieceType},
		{A5, C3, NoPieceType},
		{A5, D2, NoPieceType},
		{A5, E1, NoPieceType},
		{E4, E3, NoPieceType},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R b KQkq - 2 14: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateWhitePawnMovesForward(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White

	expectedMoves := []Move{
		{E2, E3, NoPieceType},
		{E2, E4, NoPieceType},
	}
	pos.SetPiece(WhitePawn, E2)
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E3, E4, NoPieceType},
	}
	pos.SetPiece(WhitePawn, E3)
	pos.ClearPiece(E2)
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when not on starting square: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E7, E8, Queen},
		{E7, E8, Rook},
		{E7, E8, Knight},
		{E7, E8, Bishop},
	}
	pos.ClearPiece(E3)
	pos.SetPiece(WhitePawn, E7)
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when promoting: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateWhitePawnTakes(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhitePawn, E6)
	pos.SetPiece(BlackPawn, E7)
	pos.SetPiece(BlackPawn, D7)
	expectedMoves := []Move{
		{E6, D7, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPiece(BlackPawn, F7)
	expectedMoves = append(expectedMoves, Move{E6, F7, NoPieceType})
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos = &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhitePawn, E5)
	pos.SetPiece(BlackPawn, D5)
	pos.SetPiece(BlackPawn, F5)
	pos.EnPassant = D6
	expectedMoves = []Move{
		{E5, D6, NoPieceType},
		{E5, E6, NoPieceType},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for EnPassant: expected %v, got %v", expectedMoves, moves)
	}

	pos = &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhitePawn, H7)
	pos.SetPiece(BlackQueen, H8)
	pos.SetPiece(BlackKnight, G8)
	pos.SetPiece(BlackPawn, A7)
	expectedMoves = []Move{
		{H7, G8, Rook},
		{H7, G8, Knight},
		{H7, G8, Bishop},
		{H7, G8, Queen},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for promotion: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateRookMovesWhite(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteRook, E4)
	pos.SetPiece(WhitePawn, D4)
	pos.SetPiece(BlackRook, E7)
	expectedMoves := []Move{
		{E4, E5, NoPieceType},
		{E4, E6, NoPieceType},
		{E4, E7, NoPieceType},
		{E4, F4, NoPieceType},
		{E4, G4, NoPieceType},
		{E4, H4, NoPieceType},
		{E4, E3, NoPieceType},
		{E4, E2, NoPieceType},
		{E4, E1, NoPieceType},
		{D4, D5, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKnight, G4)
	pos.SetPiece(WhitePawn, F6)
	pos.SetPiece(BlackKnight, E3)
	expectedMoves := []Move{
		{G4, H6, NoPieceType},
		{G4, H2, NoPieceType},
		{G4, F2, NoPieceType},
		{G4, E3, NoPieceType},
		{G4, E5, NoPieceType},
		{F6, F7, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteBishop, F3)
	pos.SetPiece(BlackBishop, E4)
	pos.SetPiece(WhitePawn, D1)
	expectedMoves := []Move{
		{F3, E4, NoPieceType},
		{F3, G4, NoPieceType},
		{F3, H5, NoPieceType},
		{F3, G2, NoPieceType},
		{F3, H1, NoPieceType},
		{F3, E2, NoPieceType},
		{D1, D2, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black
	pos.SetPiece(BlackQueen, F3)
	pos.SetPiece(WhiteQueen, E3)
	pos.SetPiece(BlackPawn, C6)
	pos.SetPiece(BlackPawn, D1)
	pos.SetPiece(BlackPawn, F5)
	expectedMoves := []Move{
		{F3, E3, NoPieceType},
		{F3, E4, NoPieceType},
		{F3, D5, NoPieceType},
		{F3, F4, NoPieceType},
		{F3, G4, NoPieceType},
		{F3, H5, NoPieceType},
		{F3, G3, NoPieceType},
		{F3, H3, NoPieceType},
		{F3, G2, NoPieceType},
		{F3, H1, NoPieceType},
		{F3, F2, NoPieceType},
		{F3, F1, NoPieceType},
		{F3, E2, NoPieceType},
		{F5, F4, NoPieceType},
		{C6, C5, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateKingMoves(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.SetPiece(WhiteKing, A3)
	pos.SetPiece(BlackPawn, A4)
	pos.SetPiece(WhitePawn, B2)
	expectedMoves := []Move{
		{A3, A4, NoPieceType},
		{A3, B4, NoPieceType},
		{A3, A2, NoPieceType},
		{B2, B3, NoPieceType},
		{B2, B4, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateCastleMovesWhite(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = White
	pos.WhiteKsCastle = true
	pos.WhiteQsCastle = true
	pos.BlackKsCastle = true
	pos.BlackQsCastle = true
	pos.SetPiece(WhiteKing, E1)
	pos.SetPiece(WhiteRook, A1)
	pos.SetPiece(BlackKing, E8)
	pos.SetPiece(BlackRook, H8)
	expectedMoves := []Move{
		{E1, C1, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !slices.Contains(moves, expectedMoves[0]) {
		t.Errorf("incorrect result for white queen-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.WhiteQsCastle = false
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.WhiteQsCastle = true
	pos.SetPiece(BlackPawn, A2)
	expectedMoves = []Move{
		{E1, C1, NoPieceType},
	}
	moves = LegalMoves(pos)
	if !slices.Contains(moves, expectedMoves[0]) {
		t.Errorf("incorrect result for white queen-side castle: black pawn on A2: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPiece(BlackPawn, B2)
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: black pawn on B2: expected %v, got %v", expectedMoves, moves)
	}

	pos.ClearPiece(B2)
	pos.SetPiece(BlackPawn, E2)
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: black pawn on E2: expected %v, got %v", expectedMoves, moves)
	}

	pos.ClearPiece(E2)
	pos.SetPiece(BlackPawn, F2)
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: black pawn on F2: expected %v, got %v", expectedMoves, moves)
	}

	pos.ClearPiece(F2)
	pos.SetPiece(BlackRook, C8)
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: black rook on C8: expected %v, got %v", expectedMoves, moves)
	}

	pos.ClearPiece(C8)
	pos.SetPiece(WhiteKnight, B1)
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E1, C1, NoPieceType}) {
		t.Errorf("incorrect result for white queen-side castle: white knight on B1: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateBlackPawnMovesForward(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black

	expectedMoves := []Move{
		{E7, E6, NoPieceType},
		{E7, E5, NoPieceType},
	}
	pos.SetPiece(BlackPawn, E7)
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{}
	pos.SetPiece(WhitePawn, E6)
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when blocked: expected empty list, got %v", moves)
	}

	expectedMoves = []Move{
		{E6, E5, NoPieceType},
	}
	pos.SetPiece(BlackPawn, E6)
	pos.ClearPiece(E7)
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when not on starting square: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E2, E1, Queen},
		{E2, E1, Rook},
		{E2, E1, Knight},
		{E2, E1, Bishop},
	}
	pos.ClearPiece(E6)
	pos.SetPiece(BlackPawn, E2)
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when promoting: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateBlackPawnTakes(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black

	pos.SetPiece(BlackPawn, E4)
	pos.SetPiece(WhitePawn, E3)
	pos.SetPiece(WhitePawn, D3)
	expectedMoves := []Move{
		{E4, D3, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPiece(WhitePawn, F3)
	expectedMoves = append(expectedMoves, Move{E4, F3, NoPieceType})
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.ClearPiece(F3)
	pos.ClearPiece(D3)
	pos.EnPassant = F3
	expectedMoves = []Move{
		{E4, F3, NoPieceType},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for EnPassant: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPiece(BlackPawn, H2)
	pos.SetPiece(WhiteQueen, H1)
	pos.SetPiece(WhiteKnight, G1)
	pos.ClearPiece(E4)
	expectedMoves = []Move{
		{H2, G1, Rook},
		{H2, G1, Knight},
		{H2, G1, Bishop},
		{H2, G1, Queen},
	}
	moves = LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for promotion: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateRookMovesBlack(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black

	pos.SetPiece(BlackRook, E4)
	pos.SetPiece(BlackPawn, D4)
	pos.SetPiece(WhiteRook, E7)
	expectedMoves := []Move{
		{E4, E5, NoPieceType},
		{E4, E6, NoPieceType},
		{E4, E7, NoPieceType},
		{E4, F4, NoPieceType},
		{E4, G4, NoPieceType},
		{E4, H4, NoPieceType},
		{E4, E3, NoPieceType},
		{E4, E2, NoPieceType},
		{E4, E1, NoPieceType},
		{D4, D3, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateCastleMovesBlack(t *testing.T) {
	pos := &Position{}
	pos.SideToMove = Black
	pos.BlackKsCastle = true
	pos.BlackQsCastle = true

	pos.SetPiece(BlackKing, E8)
	pos.SetPiece(BlackRook, H8)
	expectedMoves := []Move{
		{E8, G8, NoPieceType},
	}
	moves := LegalMoves(pos)
	if !slices.Contains(moves, expectedMoves[0]) {
		t.Errorf("incorrect result for black king-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.BlackKsCastle = false
	expectedMoves = []Move{}
	moves = LegalMoves(pos)
	if slices.Contains(moves, Move{E8, G8, NoPieceType}) {
		t.Errorf("incorrect result for black king-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.BlackKsCastle = true
	pos.SetPiece(WhiteKnight, F7)
	expectedMoves = []Move{Move{E8, G8, NoPieceType}}
	moves = LegalMoves(pos)
	if !slices.Contains(moves, Move{E8, G8, NoPieceType}) {
		t.Errorf("incorrect result for black king-side castle: white knight on F7: expected %v, got %v", expectedMoves, moves)
	}
}

func BenchmarkGeneratePseudoLegalMoves(b *testing.B) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PseudoLegalMoves(pos)
	}
}

func BenchmarkGenerateLegalMoves(b *testing.B) {
	pos := &Position{}
	pos.UnmarshalText([]byte("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LegalMoves(pos)
	}
}
