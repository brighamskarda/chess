package chess

import (
	"testing"
)

func TestSquareToString(t *testing.T) {
	testSquare := Square{
		FileA, Rank1,
	}
	if testSquare.String() != "A1" {
		t.Error("Square A1 String() did not equal \"A1\"")
	}

	testSquare.File = FileH
	testSquare.Rank = Rank8
	if testSquare.String() != "H8" {
		t.Error("Square H8 String() did not equal \"H8\"")
	}
}

func TestParseSquare(t *testing.T) {
	square, err := ParseSquare("A1")
	if square.File != FileA || square.Rank != Rank1 {
		t.Error("ParseSquare(A1) returned wrong square")
	}
	if err != nil {
		t.Error("ParseSquare(A1) returned error")
	}

	square, err = ParseSquare("H8")
	if square.File != FileH || square.Rank != Rank8 {
		t.Error("ParseSquare(H8) returned wrong square")
	}
	if err != nil {
		t.Error("ParseSquare(H8) returned error")
	}
}

func TestParseSquareErrors(t *testing.T) {
	_, err := ParseSquare("I8")
	if err == nil {
		t.Error("ParseSquare(I8) did not return error")
	}

	_, err = ParseSquare("@8")
	if err == nil {
		t.Error("ParseSquare(@8) did not return error")
	}

	_, err = ParseSquare("H9")
	if err == nil {
		t.Error("ParseSquare(I9) did not return error")
	}

	_, err = ParseSquare("H0")
	if err == nil {
		t.Error("ParseSquare(H0) did not return error")
	}

	_, err = ParseSquare("H")
	if err == nil {
		t.Error("ParseSquare(H) did not return error")
	}

	_, err = ParseSquare("H8H")
	if err == nil {
		t.Error("ParseSquare(H8H) did not return error")
	}
}
