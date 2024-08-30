package chess

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func moveSetsEqual(m1 []Move, m2 []Move) bool {
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

func TestGeneratePsuedoLegalMoves(t *testing.T) {
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

	pos := getDefaultPosition()
	moves := GeneratePseudoLegalMoves(pos)
	if !moveSetsEqual(defaultMoveSet, moves) {
		t.Error("incorrect result for default board: ", cmp.Diff(defaultMoveSet, moves))
	}

	pos, _ = ParseFen("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14")
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
	moves = GeneratePseudoLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Error("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14 ", cmp.Diff(expectedMoves, moves))
	}

	pos.Turn = Black
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
	moves = GeneratePseudoLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Error("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R b KQkq - 2 14 ", cmp.Diff(expectedMoves, moves))
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

	pos := getDefaultPosition()
	moves := GenerateLegalMoves(pos)
	if !moveSetsEqual(defaultMoveSet, moves) {
		t.Error("incorrect result for default board: ", cmp.Diff(defaultMoveSet, moves))
	}

	pos, _ = ParseFen("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14")
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
	moves = GenerateLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Error("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14 ", cmp.Diff(expectedMoves, moves))
	}

	pos.Turn = Black
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
	moves = GenerateLegalMoves(pos)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Error("incorrect result: fen = r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R b KQkq - 2 14 ", cmp.Diff(expectedMoves, moves))
	}
}

func TestGenerateWhitePawnMovesForward(t *testing.T) {
	pos := &Position{}
	pos.Turn = White

	expectedMoves := []Move{
		{E2, E3, NoPieceType},
		{E2, E4, NoPieceType},
	}
	pos.SetPieceAt(E2, WhitePawn)
	moves := generateWhitePawnMoves(pos, E2)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E3, E4, NoPieceType},
	}
	pos.SetPieceAt(E3, WhitePawn)
	pos.SetPieceAt(E2, NoPiece)
	moves = generateWhitePawnMoves(pos, E3)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when not on starting square: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E7, E8, Queen},
		{E7, E8, Rook},
		{E7, E8, Knight},
		{E7, E8, Bishop},
	}
	pos.SetPieceAt(E3, NoPiece)
	pos.SetPieceAt(E7, WhitePawn)
	moves = generateWhitePawnMoves(pos, E7)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when promoting: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateWhitePawnTakes(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(E6, WhitePawn)
	pos.SetPieceAt(E7, BlackPawn)
	pos.SetPieceAt(D7, BlackPawn)
	expectedMoves := []Move{
		{E6, D7, NoPieceType},
	}
	moves := generateWhitePawnMoves(pos, E6)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPieceAt(F7, BlackPawn)
	expectedMoves = append(expectedMoves, Move{E6, F7, NoPieceType})
	moves = generateWhitePawnMoves(pos, E6)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(E5, WhitePawn)
	pos.SetPieceAt(F7, NoPiece)
	pos.SetPieceAt(D7, NoPiece)
	pos.SetPieceAt(D5, BlackPawn)
	pos.SetPieceAt(F5, BlackPawn)
	pos.EnPassant = D6
	expectedMoves = []Move{
		{E5, D6, NoPieceType},
		{E5, E6, NoPieceType},
	}
	moves = generateWhitePawnMoves(pos, E5)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for EnPassant: expected %v, got %v", expectedMoves, moves)
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(H7, WhitePawn)
	pos.SetPieceAt(H8, BlackQueen)
	pos.SetPieceAt(G8, BlackKnight)
	pos.SetPieceAt(A7, BlackPawn)
	expectedMoves = []Move{
		{H7, G8, Rook},
		{H7, G8, Knight},
		{H7, G8, Bishop},
		{H7, G8, Queen},
	}
	moves = generateWhitePawnMoves(pos, H7)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for promotion: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateRookMovesWhite(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(E4, WhiteRook)
	pos.SetPieceAt(D4, WhitePawn)
	pos.SetPieceAt(E7, BlackRook)
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
	}
	moves := generateRookMoves(pos, E4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(G4, WhiteKnight)
	pos.SetPieceAt(F6, WhiteRook)
	pos.SetPieceAt(E3, BlackKnight)
	expectedMoves := []Move{
		{G4, H6, NoPieceType},
		{G4, H2, NoPieceType},
		{G4, F2, NoPieceType},
		{G4, E3, NoPieceType},
		{G4, E5, NoPieceType},
	}
	moves := generateKnightMoves(pos, G4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(F3, WhiteBishop)
	pos.SetPieceAt(E4, BlackBishop)
	pos.SetPieceAt(D1, WhitePawn)
	expectedMoves := []Move{
		{F3, E4, NoPieceType},
		{F3, G4, NoPieceType},
		{F3, H5, NoPieceType},
		{F3, G2, NoPieceType},
		{F3, H1, NoPieceType},
		{F3, E2, NoPieceType},
	}
	moves := generateBishopMoves(pos, F3)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(F3, BlackQueen)
	pos.SetPieceAt(E3, WhiteQueen)
	pos.SetPieceAt(C6, BlackPawn)
	pos.SetPieceAt(D1, BlackPawn)
	pos.SetPieceAt(F5, BlackPawn)
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
	}
	moves := generateQueenMoves(pos, F3)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateKingMoves(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(A3, WhiteKing)
	pos.SetPieceAt(A4, BlackPawn)
	pos.SetPieceAt(B2, WhitePawn)
	expectedMoves := []Move{
		{A3, A4, NoPieceType},
		{A3, B4, NoPieceType},
		{A3, B3, NoPieceType},
		{A3, A2, NoPieceType},
	}
	moves := generateKingMoves(pos, A3)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateCastleMovesWhite(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.WhiteKingSideCastle = true
	pos.WhiteQueenSideCastle = true
	pos.BlackKingSideCastle = true
	pos.BlackQueenSideCastle = true
	pos.SetPieceAt(E1, WhiteKing)
	pos.SetPieceAt(A1, WhiteRook)
	pos.SetPieceAt(E8, BlackKing)
	pos.SetPieceAt(H8, BlackRook)
	expectedMoves := []Move{
		{E1, C1, NoPieceType},
	}
	moves := generateCastleMoves(pos, E1)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for white queen-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.WhiteQueenSideCastle = false
	expectedMoves = []Move{}
	moves = generateCastleMoves(pos, E1)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for white queen-side castle: expected %v, got %v", expectedMoves, moves)
	}

	// pos.WhiteQueenSideCastle = true
	// pos.SetPieceAt(A2, BlackPawn)
	// expectedMoves = []Move{
	// 	{E1, C1, NoPieceType},
	// }
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: black pawn on A2: expected %v, got %v", expectedMoves, moves)
	// }

	// pos.SetPieceAt(B2, BlackPawn)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: black pawn on B2: expected %v, got %v", expectedMoves, moves)
	// }

	// pos.SetPieceAt(B2, NoPiece)
	// pos.SetPieceAt(E2, BlackPawn)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: black pawn on E2: expected %v, got %v", expectedMoves, moves)
	// }

	// pos.SetPieceAt(E2, NoPiece)
	// pos.SetPieceAt(F2, BlackPawn)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: black pawn on F2: expected %v, got %v", expectedMoves, moves)
	// }

	// pos.SetPieceAt(F2, NoPiece)
	// pos.SetPieceAt(C8, BlackRook)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: black rook on C8: expected %v, got %v", expectedMoves, moves)
	// }

	// pos.SetPieceAt(C8, NoPiece)
	// pos.SetPieceAt(B1, WhiteKnight)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E1)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for white queen-side castle: white knight on B1: expected %v, got %v", expectedMoves, moves)
	// }
}

func TestGenerateBlackPawnMovesForward(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	expectedMoves := []Move{
		{E7, E6, NoPieceType},
		{E7, E5, NoPieceType},
	}
	pos.SetPieceAt(E7, BlackPawn)
	moves := generateBlackPawnMoves(pos, E7)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{}
	pos.SetPieceAt(E6, WhitePawn)
	moves = generateBlackPawnMoves(pos, E7)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when blocked: expected empty list, got %v", moves)
	}

	expectedMoves = []Move{
		{E6, E5, NoPieceType},
	}
	pos.SetPieceAt(E6, BlackPawn)
	pos.SetPieceAt(E7, NoPiece)
	moves = generateBlackPawnMoves(pos, E6)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when not on starting square: expected %v, got %v", expectedMoves, moves)
	}

	expectedMoves = []Move{
		{E2, E1, Queen},
		{E2, E1, Rook},
		{E2, E1, Knight},
		{E2, E1, Bishop},
	}
	pos.SetPieceAt(E6, NoPiece)
	pos.SetPieceAt(E2, BlackPawn)
	moves = generateBlackPawnMoves(pos, E2)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result when promoting: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateBlackPawnTakes(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(E4, BlackPawn)
	pos.SetPieceAt(E3, WhitePawn)
	pos.SetPieceAt(D3, WhitePawn)
	expectedMoves := []Move{
		{E4, D3, NoPieceType},
	}
	moves := generateBlackPawnMoves(pos, E4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPieceAt(F3, WhitePawn)
	expectedMoves = append(expectedMoves, Move{E4, F3, NoPieceType})
	moves = generateBlackPawnMoves(pos, E4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPieceAt(F3, NoPiece)
	pos.SetPieceAt(D3, NoPiece)
	pos.SetPieceAt(D4, WhitePawn)
	pos.SetPieceAt(F4, WhitePawn)
	pos.EnPassant = F3
	expectedMoves = []Move{
		{E4, F3, NoPieceType},
	}
	moves = generateBlackPawnMoves(pos, E4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for EnPassant: expected %v, got %v", expectedMoves, moves)
	}

	pos = &Position{}
	pos.Turn = Black
	pos.SetPieceAt(H2, BlackPawn)
	pos.SetPieceAt(H1, WhiteQueen)
	pos.SetPieceAt(G1, WhiteKnight)
	pos.SetPieceAt(A2, WhitePawn)
	expectedMoves = []Move{
		{H2, G1, Rook},
		{H2, G1, Knight},
		{H2, G1, Bishop},
		{H2, G1, Queen},
	}
	moves = generateBlackPawnMoves(pos, H2)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for promotion: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateRookMovesBlack(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(E4, BlackRook)
	pos.SetPieceAt(D4, BlackPawn)
	pos.SetPieceAt(E7, WhiteRook)
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
	}
	moves := generateRookMoves(pos, E4)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result: expected %v, got %v", expectedMoves, moves)
	}
}

func TestGenerateCastleMovesBlack(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.WhiteKingSideCastle = true
	pos.WhiteQueenSideCastle = true
	pos.BlackKingSideCastle = true
	pos.BlackQueenSideCastle = true
	pos.SetPieceAt(E1, WhiteKing)
	pos.SetPieceAt(A1, WhiteRook)
	pos.SetPieceAt(E8, BlackKing)
	pos.SetPieceAt(H8, BlackRook)
	expectedMoves := []Move{
		{E8, G8, NoPieceType},
	}
	moves := generateCastleMoves(pos, E8)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for black king-side castle: expected %v, got %v", expectedMoves, moves)
	}

	pos.SetPieceAt(F7, WhiteKnight)
	expectedMoves = []Move{
		{E8, G8, NoPieceType},
	}
	moves = generateCastleMoves(pos, E8)
	if !moveSetsEqual(expectedMoves, moves) {
		t.Errorf("incorrect result for black king-side castle: white knight on F7: expected %v, got %v", expectedMoves, moves)
	}

	// pos.SetPieceAt(F7, NoPiece)
	// pos.SetPieceAt(E6, WhiteBishop)
	// expectedMoves = []Move{}
	// moves = generateCastleMoves(pos, E8)
	// if !moveSetsEqual(expectedMoves, moves) {
	// 	t.Errorf("incorrect result for black king-side castle: white bishop on E6: expected %v, got %v", expectedMoves, moves)
	// }
}

func BenchmarkGeneratePseudoLegalMoves(b *testing.B) {
	pos, _ := ParseFen("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GeneratePseudoLegalMoves(pos)
	}
}

func BenchmarkGenerateLegalMoves(b *testing.B) {
	pos, _ := ParseFen("r3kb1r/2p3pp/pp3n2/q4P2/2B1p3/6Q1/PPP2PPP/RNB1K2R w KQkq - 2 14")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateLegalMoves(pos)
	}
}
