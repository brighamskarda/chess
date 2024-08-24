package chess

import (
	"testing"
)

func TestParseMove(t *testing.T) {
	_, err := ParseMove("")
	if err == nil {
		t.Error("incorrect result: input : expected error, got nil")
	}

	_, err = ParseMove("E7E8  ")
	if err == nil {
		t.Error("incorrect result: input E7E8  : expected error, got nil")
	}

	_, err = ParseMove("E7E")
	if err == nil {
		t.Error("incorrect result: input E7E: expected error, got nil")
	}

	move, err := ParseMove("E7E8Q")
	if err != nil {
		t.Error("incorrect result: input E7E8Q: expected nil, got error")
	}
	expectedMove := Move{E7, E8, Queen}
	if move != expectedMove {
		t.Errorf("incorrect result: input E7E8Q: expected %v, got %v", expectedMove, move)
	}

	move, err = ParseMove("e7e8")
	if err != nil {
		t.Error("incorrect result: input e7e8: expected nil, got error")
	}
	expectedMove = Move{E7, E8, NoPieceType}
	if move != expectedMove {
		t.Errorf("incorrect result: input e7e8: expected %v, got %v", expectedMove, move)
	}
}
