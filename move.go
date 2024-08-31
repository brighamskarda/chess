package chess

import (
	"errors"
	"fmt"
	"strings"
)

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
}

func (m Move) String() string {
	returnString := m.FromSquare.String() + m.ToSquare.String()
	if m.Promotion != NoPieceType {
		returnString += m.Promotion.String()
	}
	return returnString
}

// ParseUCIMove expects a UCI compatible move string. Format should be Square1Square2Promotion, where promotion is optional.
func ParseUCIMove(s string) (Move, error) {
	if len(s) != 4 && len(s) != 5 {
		return Move{}, fmt.Errorf("invalid move string: string not 4 or 5 characters long: %s", s)
	}
	fromSquare, err := ParseSquare(s[0:2])
	if err != nil {
		return Move{}, fmt.Errorf("invalid move string: %w", err)
	}
	toSquare, err := ParseSquare(s[2:4])
	if err != nil {
		return Move{}, fmt.Errorf("invalid move string: %w", err)
	}
	promotion := NoPieceType
	if len(s) == 5 {
		promotion, err = parsePieceType(rune(s[4]))
		if err != nil {
			return Move{}, fmt.Errorf("invalid move string: %w", err)
		}
	}

	return Move{fromSquare, toSquare, promotion}, nil
}

func ParseSANMove(p *Position, s string) (Move, error) {
	cleanedString := strings.ReplaceAll(s, "+", "")
	cleanedString = strings.ReplaceAll(cleanedString, "#", "")

	if p.Turn != White && p.Turn != Black {
		return Move{}, errors.New("could not parse SAN move: position turn is not set to white or black")
	}

	if len(cleanedString) == 2 {
		return parseSANBasicPawnMove(p, cleanedString)
	}

	return Move{}, nil
}

func parseSANBasicPawnMove(p *Position, s string) (Move, error) {
	square, err := ParseSquare(s)
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN basic pawn move: input %s: %w", s, err)
	}
	if p.Turn == White {
		return parseSANBasicPawnMoveWhite(p, square)
	}
	return parseSANBasicPawnMoveBlack(p, square)
}

func parseSANBasicPawnMoveWhite(p *Position, s Square) (Move, error) {
	if p.PieceAt(squareBelow(s)) == WhitePawn {
		return Move{FromSquare: squareBelow(s), ToSquare: s, Promotion: NoPieceType}, nil
	}
	if s.Rank == 4 && p.PieceAt(squareBelow(s)) == NoPiece && p.PieceAt(squareBelow(squareBelow(s))) == WhitePawn {
		return Move{FromSquare: squareBelow(squareBelow(s)), ToSquare: s, Promotion: NoPieceType}, nil
	}
	return Move{}, errors.New("could not parse SAN basic pawn move")
}

func parseSANBasicPawnMoveBlack(p *Position, s Square) (Move, error) {
	if p.PieceAt(squareAbove(s)) == BlackPawn {
		return Move{FromSquare: squareAbove(s), ToSquare: s, Promotion: NoPieceType}, nil
	}
	if s.Rank == 5 && p.PieceAt(squareAbove(s)) == NoPiece && p.PieceAt(squareAbove(squareAbove(s))) == BlackPawn {
		return Move{FromSquare: squareAbove(squareAbove(s)), ToSquare: s, Promotion: NoPieceType}, nil
	}
	return Move{}, errors.New("could not parse SAN basic pawn move")
}

// IsValidMove makes sure each of the elements in Move m are logical. Namely that the squares can be found on a chess board.
func IsValidMove(m Move) bool {
	return isValidSquare(m.FromSquare) && m.FromSquare != NoSquare &&
		isValidSquare(m.ToSquare) && m.ToSquare != NoSquare &&
		isValidPieceType(m.Promotion)
}
