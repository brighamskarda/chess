package chess

import (
	"testing"
)

const emptyBoardStr string = `8        
7        
6        
5        
4        
3        
2        
1        
 ABCDEFGH`

const defaultBoardStr string = `8rnbqkbnr
7pppppppp
6        
5        
4        
3        
2PPPPPPPP
1RNBQKBNR
 ABCDEFGH`

func TestPositionString(t *testing.T) {
	pos := &Position{}
	if pos.String() != emptyBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.String(), emptyBoardStr)
	}
	pos = getDefaultPosition()
	if pos.String() != defaultBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.String(), defaultBoardStr)
	}
}

const emptyBoardStrFlipped string = `1        
2        
3        
4        
5        
6        
7        
8        
 HGFEDCBA`

const defaultBoardStrFlipped string = `1RNBKQBNR
2PPPPPPPP
3        
4        
5        
6        
7pppppppp
8rnbkqbnr
 HGFEDCBA`

func TestPositionFormatSTring(t *testing.T) {
	pos := &Position{}
	if pos.FormatString(false) != emptyBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.FormatString(false), emptyBoardStr)
	}
	if pos.FormatString(true) != emptyBoardStrFlipped {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.FormatString(true), emptyBoardStrFlipped)
	}
	pos = getDefaultPosition()
	if pos.FormatString(false) != defaultBoardStr {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.FormatString(false), defaultBoardStr)
	}
	if pos.FormatString(true) != defaultBoardStrFlipped {
		t.Errorf("Default board string incorrect:\nActual:\n%s\nExpected:\n%s", pos.FormatString(true), defaultBoardStrFlipped)
	}
}

func getDefaultPosition() *Position {
	pos := Position{}
	pos.Board[0] = BlackRook
	pos.Board[1] = BlackKnight
	pos.Board[2] = BlackBishop
	pos.Board[3] = BlackQueen
	pos.Board[4] = BlackKing
	pos.Board[5] = BlackBishop
	pos.Board[6] = BlackKnight
	pos.Board[7] = BlackRook
	for i := 8; i < 16; i++ {
		pos.Board[i] = BlackPawn
	}
	for i := 48; i < 56; i++ {
		pos.Board[i] = WhitePawn
	}
	pos.Board[56] = WhiteRook
	pos.Board[57] = WhiteKnight
	pos.Board[58] = WhiteBishop
	pos.Board[59] = WhiteQueen
	pos.Board[60] = WhiteKing
	pos.Board[61] = WhiteBishop
	pos.Board[62] = WhiteKnight
	pos.Board[63] = WhiteRook

	pos.Turn = White

	pos.WhiteKingSideCastle = true
	pos.WhiteQueenSideCastle = true
	pos.BlackKingSideCastle = true
	pos.BlackQueenSideCastle = true

	pos.EnPassant = NoSquare

	pos.HalfMove = 0
	pos.FullMove = 1

	return &pos
}

func TestPositionPieceAt(t *testing.T) {
	pos := getDefaultPosition()
	if pos.PieceAt(C1) != WhiteBishop {
		t.Errorf("pos.At(C1) != WhiteBishop. Actual %s", pos.PieceAt(C1).String())
	}
}

func TestPositionSetPieceAt(t *testing.T) {
	pos := Position{}
	pos.SetPieceAt(D4, BlackPawn)
	if pos.Board[35] != BlackPawn {
		t.Fail()
	}
}

func TestParseFen(t *testing.T) {
	pos, err := ParseFen(DefaultFen)
	if err != nil {
		t.Error("ParseFen set error nil when fen is valid")
	}
	if *pos != *getDefaultPosition() {
		t.Errorf("ParseFen incorrect output. Actual:%+v\nExpected%+v", pos, getDefaultPosition())
	}

	pos, err = ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16")
	if err != nil {
		t.Error("ParseFen set error nil when fen is valid")
	}
	updatedPos := getDefaultPosition()
	updatedPos.Board[0] = WhitePawn
	updatedPos.Turn = Black
	updatedPos.BlackKingSideCastle = false
	updatedPos.EnPassant = E6
	updatedPos.HalfMove = 32
	updatedPos.FullMove = 16
	if *pos != *updatedPos {
		t.Errorf("ParseFen incorrect output. Actual:%+v\nExpected%+v", pos, updatedPos)
	}
}

func TestParseFenInvalid(t *testing.T) {
	_, err := ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3216")
	if err == nil {
		t.Error("ParseFen failed to set error for Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3216")
	}

	_, err = ParseFen("Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3a2 16")
	if err == nil {
		t.Error("ParseFen failed to set error for Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 3a2 16")
	}

	_, err = ParseFen("Pnbqkbnr/pppopppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16")
	if err == nil {
		t.Error("ParseFen failed to set error for Pnbqkbnr/pppopppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16")
	}
}

func BenchmarkParseFen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseFen(DefaultFen)
	}
}

func TestGenerateFen(t *testing.T) {
	pos := getDefaultPosition()
	if pos.GenerateFen() != DefaultFen {
		t.Error("GenerateFen did not generate default fen")
	}
	pos.Board[0] = WhitePawn
	pos.Turn = Black
	pos.BlackKingSideCastle = false
	pos.EnPassant = E6
	pos.HalfMove = 32
	pos.FullMove = 16
	if pos.GenerateFen() != "Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16" {
		t.Errorf(`GenerateFen expected "Pnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQq e6 32 16" but got %s`, pos.GenerateFen())
	}

	expected := "8/7p/3b2p1/2p2p2/3kpn2/7r/8/5K2 b - - 1 46"
	pos, _ = ParseFen(expected)
	actual := pos.GenerateFen()
	if expected != actual {
		t.Errorf("GenerateFen got wrong result: expected %s, got %s", expected, actual)
	}
}

func BenchmarkGenerateFen(b *testing.B) {
	pos := getDefaultPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos.GenerateFen()
	}
}

func TestIsValidPosition(t *testing.T) {
	pos := &Position{}
	if pos.IsValid() {
		t.Error("Zeroed position should not be valid.")
	}

	pos = getDefaultPosition()
	if !pos.IsValid() {
		t.Error("Default position should be valid")
	}
}

func TestIsValidPositionCastleRights(t *testing.T) {
	pos := getDefaultPosition()
	pos.Board[63] = BlackRook
	if pos.IsValid() {
		t.Error("Castle rights not checked properly")
	}
}

func TestIsValidPositionKings(t *testing.T) {
	pos := getDefaultPosition()
	pos.WhiteKingSideCastle = false
	pos.Board[4] = NoPiece
	if pos.IsValid() {
		t.Error("Kings not checked properly")
	}
}

func TestIsValidPositionPawns(t *testing.T) {
	pos := getDefaultPosition()
	pos.Board[0] = WhitePawn
	pos.BlackQueenSideCastle = false
	if pos.IsValid() {
		t.Error("Pawns not checked properly")
	}

	pos = getDefaultPosition()
	pos.Board[63] = BlackPawn
	pos.WhiteKingSideCastle = false
	if pos.IsValid() {
		t.Error("Pawns not checked properly")
	}
}

func TestIsValidPositionEnPassant(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(E4, WhitePawn)
	pos.EnPassant = E3
	pos.Turn = Black
	if !pos.IsValid() {
		t.Error("En Passant was logical")
	}

	pos.EnPassant = E4
	if pos.IsValid() {
		t.Error("En Passant was not logical")
	}

	pos.EnPassant = E6
	if pos.IsValid() {
		t.Error("En Passant was not logical")
	}

	pos.SetPieceAt(E5, WhitePawn)
	if pos.IsValid() {
		t.Error("En Passant was not logical")
	}

	pos.SetPieceAt(E5, BlackPawn)
	if pos.IsValid() {
		t.Error("En Passant was not logical")
	}

	pos.Turn = White
	if !pos.IsValid() {
		t.Error("En Passant was logical")
	}
}

func TestIsValidPositionEnPassantPawnOnSquare(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(E4, WhitePawn)
	pos.SetPieceAt(E2, NoPiece)
	pos.SetPieceAt(E3, WhitePawn)
	pos.EnPassant = E3
	pos.Turn = Black
	if pos.IsValid() {
		t.Error("incorrect result for pawns on e4 and e3: expected false, got true")
	}
}

func TestIsValidPositionTurn(t *testing.T) {
	pos := getDefaultPosition()
	pos.Turn = NoColor
	if pos.IsValid() {
		t.Error("Turn was not set, position not valid")
	}
}

func TestIsValidPositionInvalidPieces(t *testing.T) {
	pos := getDefaultPosition()
	pos.Board[9].Color = NoColor
	if pos.IsValid() {
		t.Error("Invalid Piece on board, position not valid")
	}

	pos.Board[9].Color = Black
	pos.Board[9].Type = NoPieceType
	if pos.IsValid() {
		t.Error("Invalid Piece on board, position not valid")
	}
}

func BenchmarkIsValidPosition(b *testing.B) {
	pos := getDefaultPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos.IsValid()
	}
}

func TestIndexToSquare(t *testing.T) {
	square := indexToSquare(0)
	if square != A8 {
		t.Errorf("incorrect result: input 0: expected A8, got %v", square)
	}

	square = indexToSquare(33)
	if square != B4 {
		t.Errorf("incorrect result: input 33: expected B4, got %v", square)
	}
}

func TestMove(t *testing.T) {
	pos := getDefaultPosition()
	pos.Move(Move{FromSquare: E2, ToSquare: E4, Promotion: NoPieceType})
	if pos.PieceAt(E4) != WhitePawn ||
		pos.PieceAt(E2) != NoPiece ||
		pos.EnPassant != E3 ||
		pos.Turn != Black {
		t.Errorf("incorrect result: move E2-E4: result %#v", pos)
	}

	pos.Move(Move{FromSquare: E7, ToSquare: E6, Promotion: NoPieceType})
	if pos.EnPassant != NoSquare ||
		pos.HalfMove != 0 ||
		pos.FullMove != 2 {
		t.Errorf("incorrect result: move E7-E6: result %#v", pos)
	}

	pos.Move(Move{FromSquare: B1, ToSquare: C3, Promotion: NoPieceType})
	if pos.HalfMove != 1 ||
		pos.FullMove != 2 {
		t.Errorf("incorrect result: move B1-C3: result %#v", pos)
	}
}

func TestMoveCastle(t *testing.T) {
	pos := getDefaultPosition()
	pos.SetPieceAt(F1, NoPiece)
	pos.SetPieceAt(G1, NoPiece)
	pos.Move(Move{FromSquare: E1, ToSquare: G1, Promotion: NoPieceType})
	if pos.PieceAt(G1) != WhiteKing ||
		pos.PieceAt(F1) != WhiteRook ||
		pos.WhiteKingSideCastle ||
		pos.WhiteQueenSideCastle ||
		!pos.BlackKingSideCastle ||
		!pos.BlackQueenSideCastle {
		t.Errorf("incorrect result: move E1-G1: result %#v", pos)
	}

	pos = getDefaultPosition()
	pos.Move(Move{FromSquare: E1, ToSquare: C1, Promotion: NoPieceType})
	if pos.PieceAt(C1) != WhiteKing ||
		pos.PieceAt(D1) != WhiteRook ||
		pos.WhiteKingSideCastle ||
		pos.WhiteQueenSideCastle ||
		!pos.BlackKingSideCastle ||
		!pos.BlackQueenSideCastle {
		t.Errorf("incorrect result: move E1-C1: result %#v", pos)
	}

	pos = getDefaultPosition()
	pos.Turn = Black
	pos.Move(Move{FromSquare: E8, ToSquare: C8, Promotion: NoPieceType})
	if pos.PieceAt(C8) != BlackKing ||
		pos.PieceAt(D8) != BlackRook ||
		!pos.WhiteKingSideCastle ||
		!pos.WhiteQueenSideCastle ||
		pos.BlackKingSideCastle ||
		pos.BlackQueenSideCastle {
		t.Errorf("incorrect result: move E8-C8: result %#v", pos)
	}

	pos = getDefaultPosition()
	pos.Turn = Black
	pos.Move(Move{FromSquare: E8, ToSquare: G8, Promotion: NoPieceType})
	if pos.PieceAt(G8) != BlackKing ||
		pos.PieceAt(F8) != BlackRook ||
		!pos.WhiteKingSideCastle ||
		!pos.WhiteQueenSideCastle ||
		pos.BlackKingSideCastle ||
		pos.BlackQueenSideCastle {
		t.Errorf("incorrect result: move E8-G8: result %#v", pos)
	}
}

func TestMoveEnPassant(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(E5, WhitePawn)
	pos.SetPieceAt(D5, BlackPawn)
	pos.EnPassant = D6
	pos.Move(Move{FromSquare: E5, ToSquare: D6, Promotion: NoPieceType})
	if pos.PieceAt(D5) != NoPiece ||
		pos.PieceAt(D6) != WhitePawn ||
		pos.PieceAt(E5) != NoPiece ||
		pos.EnPassant != NoSquare {
		t.Errorf("incorrect result: move E5-D6: result %#v", pos)
	}

	pos = &Position{}
	pos.Turn = Black
	pos.SetPieceAt(E4, WhitePawn)
	pos.SetPieceAt(D4, BlackPawn)
	pos.EnPassant = E3
	pos.Move(Move{FromSquare: D4, ToSquare: E3, Promotion: NoPieceType})
	if pos.PieceAt(D4) != NoPiece ||
		pos.PieceAt(E3) != BlackPawn ||
		pos.PieceAt(E4) != NoPiece ||
		pos.EnPassant != NoSquare {
		t.Errorf("incorrect result: move D4-E3: result %#v", pos)
	}
}

func TestMovePromotion(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(D6, WhitePawn)
	pos.Move(Move{FromSquare: D6, ToSquare: D8, Promotion: Queen})
	if pos.PieceAt(D6) != NoPiece ||
		pos.PieceAt(D8) != WhiteQueen {
		t.Errorf("incorrect result: move D6-D8Q: result %#v", pos)
	}
}
