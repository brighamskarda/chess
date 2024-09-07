package chess

import (
	"testing"
)

func TestParseUCIMove(t *testing.T) {
	_, err := ParseUCIMove("")
	if err == nil {
		t.Error("incorrect result: input : expected error, got nil")
	}

	_, err = ParseUCIMove("E7E8  ")
	if err == nil {
		t.Error("incorrect result: input E7E8  : expected error, got nil")
	}

	_, err = ParseUCIMove("E7E")
	if err == nil {
		t.Error("incorrect result: input E7E: expected error, got nil")
	}

	move, err := ParseUCIMove("E7E8Q")
	if err != nil {
		t.Error("incorrect result: input E7E8Q: expected nil, got error")
	}
	expectedMove := Move{E7, E8, Queen}
	if move != expectedMove {
		t.Errorf("incorrect result: input E7E8Q: expected %v, got %v", expectedMove, move)
	}

	move, err = ParseUCIMove("e7e8")
	if err != nil {
		t.Error("incorrect result: input e7e8: expected nil, got error")
	}
	expectedMove = Move{E7, E8, NoPieceType}
	if move != expectedMove {
		t.Errorf("incorrect result: input e7e8: expected %v, got %v", expectedMove, move)
	}
}

func TestParseSANMovePawn(t *testing.T) {
	pos := getDefaultPosition()
	moveString := "e4"
	expectedMove := Move{E2, E4, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "e3"
	expectedMove = Move{E2, E3, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	pos.Turn = Black
	moveString = "e5"
	expectedMove = Move{E7, E5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "e6"
	expectedMove = Move{E7, E6, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
}

func TestParseSANMoveOtherPiecesWhite(t *testing.T) {
	pos := getDefaultPosition()
	moveString := "Nc3"
	expectedMove := Move{B1, C3, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	pos.SetPieceAt(D2, NoPiece)

	moveString = "Bh6"
	expectedMove = Move{C1, H6, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "Qd5"
	expectedMove = Move{D1, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "Kd2"
	expectedMove = Move{E1, D2, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	pos.SetPieceAt(A2, NoPiece)

	moveString = "Ra4"
	expectedMove = Move{A1, A4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
}

func TestParseSANMoveOtherPiecesBlack(t *testing.T) {
	pos := getDefaultPosition()
	pos.Turn = Black
	moveString := "Nc6"
	expectedMove := Move{B8, C6, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	pos.SetPieceAt(D7, NoPiece)

	moveString = "Bh3"
	expectedMove = Move{C8, H3, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "Qd5"
	expectedMove = Move{D8, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	moveString = "Kd7"
	expectedMove = Move{E8, D7, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}

	pos.SetPieceAt(A7, NoPiece)

	moveString = "Ra4"
	expectedMove = Move{A8, A4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
}

func TestParseSANMoveInvalid(t *testing.T) {
	pos := getDefaultPosition()
	moveString := "e5"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Ra2"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Ra2"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Nb3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Ba3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Qd2"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Kf1"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Qd7"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	pos.SetPieceAt(D2, NoPiece)
	moveString = "Qd8"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	pos = getDefaultPosition()
	pos.Turn = Black
	moveString = "d8"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Nb6"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Nc3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Qd6"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(G3, WhitePawn)
	pos.SetPieceAt(G4, BlackKnight)
	moveString = "gxg4"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}
}

func TestParseSANMoveAmbiguousRook(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(C3, BlackRook)
	pos.SetPieceAt(E3, BlackRook)
	moveString := "Rd3"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Rcd3"
	expectedMove := Move{C3, D3, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Rc3d3"
	expectedMove = Move{C3, D3, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Red3"
	expectedMove = Move{E3, D3, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveAmbiguousKnight(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(C3, WhiteKnight)
	pos.SetPieceAt(E3, WhiteKnight)
	moveString := "Nd5"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Ncd5"
	expectedMove := Move{C3, D5, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Nc3d5"
	expectedMove = Move{C3, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Ned5"
	expectedMove = Move{E3, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveAmbiguousBishop(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(C3, BlackBishop)
	pos.SetPieceAt(E3, BlackBishop)
	moveString := "Bd4"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Bcd4"
	expectedMove := Move{C3, D4, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Bc3d4"
	expectedMove = Move{C3, D4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Bed4"
	expectedMove = Move{E3, D4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveAmbiguousQueen(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(D3, WhiteQueen)
	pos.SetPieceAt(F3, WhiteQueen)
	moveString := "Qd5"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Qdd5"
	expectedMove := Move{D3, D5, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Qd3d5"
	expectedMove = Move{D3, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Qfd5"
	expectedMove = Move{F3, D5, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	pos = &Position{}
	pos.Turn = White
	pos.SetPieceAt(D3, WhiteQueen)
	pos.SetPieceAt(D5, WhiteQueen)
	moveString = "Qe4"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Q5e4"
	expectedMove = Move{D5, E4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Qd5e4"
	expectedMove = Move{D5, E4, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveCapturesWhitePawn(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(F3, WhitePawn)
	pos.SetPieceAt(G4, BlackKnight)
	moveString := "g4"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "fg4"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "fxg4"
	expectedMove := Move{F3, G4, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveCapturesWhitePawnEnPassant(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.EnPassant = G6
	pos.SetPieceAt(H5, WhitePawn)
	pos.SetPieceAt(G5, BlackPawn)
	moveString := "hg6"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "hxg6"
	expectedMove := Move{H5, G6, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	pos.EnPassant = NoSquare
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}
}

func TestParseSANMoveCapturesBlackPawn(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.SetPieceAt(F3, WhiteQueen)
	pos.SetPieceAt(G4, BlackPawn)
	moveString := "f3"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "gf3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "gxf3"
	expectedMove := Move{G4, F3, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveCapturesBlackPawnEnPassant(t *testing.T) {
	pos := &Position{}
	pos.Turn = Black
	pos.EnPassant = H3
	pos.SetPieceAt(H4, WhitePawn)
	pos.SetPieceAt(G4, BlackPawn)
	moveString := "gh3"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "gxh3"
	expectedMove := Move{G4, H3, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveOtherCaptures(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(D5, WhiteQueen)
	pos.SetPieceAt(D3, WhiteQueen)
	pos.SetPieceAt(F3, BlackKnight)
	moveString := "Qf3"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Qxf3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Q3f3"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "Q3xf3"
	expectedMove := Move{D3, F3, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "Qd3xf3"
	expectedMove = Move{D3, F3, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveWhitePromotions(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(C7, WhitePawn)
	moveString := "c8"
	_, err := ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "c8Q"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "c8=Q"
	expectedMove := Move{C7, C8, Queen}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	pos.SetPieceAt(B8, BlackKnight)

	moveString = "b8"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "cxb8"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "cxb8=K"
	_, err = ParseSANMove(pos, moveString)
	if err == nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, "error", nil)
	}

	moveString = "cxb8=R"
	expectedMove = Move{C7, B8, Rook}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveCastling(t *testing.T) {
	pos := getDefaultPosition()
	moveString := "O-O"
	expectedMove := Move{E1, G1, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "O-O-O"
	expectedMove = Move{E1, C1, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	pos = getDefaultPosition()
	pos.Turn = Black
	moveString = "O-O"
	expectedMove = Move{E8, G8, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	moveString = "O-O-O"
	expectedMove = Move{E8, C8, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}

func TestParseSANMoveIgnoresOtherSymbols(t *testing.T) {
	pos := &Position{}
	pos.Turn = White
	pos.SetPieceAt(A8, BlackKing)
	pos.SetPieceAt(B1, WhiteRook)
	moveString := "Ra1+"
	expectedMove := Move{B1, A1, NoPieceType}
	move, err := ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}

	pos.SetPieceAt(B2, WhiteRook)
	moveString = "Ra2#"
	expectedMove = Move{B1, A1, NoPieceType}
	move, err = ParseSANMove(pos, moveString)
	if err != nil {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, nil, err)
	}
	if move != expectedMove {
		t.Errorf("incorrect result: input %s: expected %v, got %v", moveString, expectedMove, move)
	}
}
