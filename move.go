package chess

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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

	if s == "O-O" || s == "O-O-O" {
		return parseSANCastleMove(p, cleanedString)
	}
	if isSANBasicPawnMove(p, cleanedString) {
		return parseSANBasicPawnMove(p, cleanedString)
	}
	if isSANPawnCapture(p, cleanedString) {
		return parseSANPawnCapture(p, cleanedString)
	}
	if strings.ContainsRune(cleanedString, '=') {
		return parseSANPromotion(p, cleanedString)
	}
	if len(cleanedString) == 3 {
		return parseSANPieceMove(p, cleanedString)
	}

	return Move{}, errors.New("unknown error")
}

func isSANBasicPawnMove(p *Position, s string) bool {
	return len(s) == 2 &&
		!(rune(s[1]) == '8' && p.Turn == White) &&
		!(rune(s[1]) == '1' && p.Turn == Black)
}

func isSANPawnCapture(p *Position, s string) bool {
	return len(s) == 4 &&
		unicode.IsLower(rune(s[0])) &&
		rune(s[1]) == 'x' &&
		!(rune(s[3]) == '8' && p.Turn == White) &&
		!(rune(s[3]) == '1' && p.Turn == Black)
}

func parseSANCastleMove(p *Position, s string) (Move, error) {
	if p.Turn == White && s == "O-O" {
		return Move{E1, G1, NoPieceType}, nil
	}
	if p.Turn == White && s == "O-O-O" {
		return Move{E1, C1, NoPieceType}, nil
	}
	if p.Turn == Black && s == "O-O" {
		return Move{E8, G8, NoPieceType}, nil
	}
	if p.Turn == Black && s == "O-O-O" {
		return Move{E8, C8, NoPieceType}, nil
	}
	return Move{}, fmt.Errorf("could not parseSAN castle move: input %s", s)
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

func parseSANPawnCapture(p *Position, s string) (Move, error) {
	toSquare, err := ParseSquare(s[2:])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN pawn capture move: input %s: %w", s, err)
	}
	file, err := parseFile(rune(s[0]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN pawn capture move: input %s: %w", s, err)
	}
	var fromSquare Square
	if p.Turn == White {
		fromSquare = Square{File: file, Rank: toSquare.Rank - 1}
	}
	if p.Turn == Black {
		fromSquare = Square{File: file, Rank: toSquare.Rank + 1}
	}

	pieceAtFromSquare := p.PieceAt(fromSquare)
	if pieceAtFromSquare == NoPiece || pieceAtFromSquare.Type != Pawn || pieceAtFromSquare.Color != p.Turn {
		return Move{}, fmt.Errorf("invalid SAN pawn capture move: piece at %v is not %v", fromSquare, Piece{p.Turn, Pawn})
	}

	pieceAtToSquare := p.PieceAt(toSquare)
	if pieceAtToSquare.Color == p.Turn || (pieceAtToSquare.Color == NoColor && toSquare != p.EnPassant) {
		return Move{}, fmt.Errorf("invalid SAN pawn capture move: invalid piece to capture: square, %v piece, %v, en passant %v",
			toSquare, pieceAtToSquare, p.EnPassant)
	}

	return Move{FromSquare: fromSquare, ToSquare: toSquare, Promotion: NoPieceType}, nil
}

func parseSANPromotion(p *Position, s string) (Move, error) {
	sNoPromotion := s[:strings.IndexRune(s, '=')]
	move := Move{}
	var err error = nil
	if len(sNoPromotion) == 2 {
		move, err = parseSANBasicPawnMove(p, sNoPromotion)
	} else if len(sNoPromotion) == 4 && unicode.IsLower(rune(sNoPromotion[0])) && rune(sNoPromotion[1]) == 'x' {
		move, err = parseSANPawnCapture(p, sNoPromotion)
	} else {
		return Move{}, fmt.Errorf("could not parse move before promotion: num of chars before '=' is not 2 or 4: input %s", s)
	}
	if err != nil {
		return Move{}, fmt.Errorf("could not parse move before promotion: %w", err)
	}

	promotion, err := parsePieceType(rune(s[len(s)-1]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN promotion: input %s: %w", s, err)
	}
	if promotion == King || promotion == Pawn {
		return Move{}, fmt.Errorf("invalid promotion: can't promote to king or pawn: input %s", s)
	}

	move.Promotion = promotion

	return move, nil
}

func parseSANPieceMove(p *Position, s string) (Move, error) {
	pieceType, err := parsePieceType(rune(s[0]))
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move: invalid piece type: input, %s: %w", s, err)
	}
	square, err := ParseSquare(s[1:])
	if err != nil {
		return Move{}, fmt.Errorf("could not parse SAN move: could not parse destination square: input, %s: %w", s, err)
	}
	switch pieceType {
	case Pawn:
		return Move{}, fmt.Errorf("invalid SAN format: should not specify p for pawn: input %s", s)
	case Rook:
		return parseSANRookMove(p, square)
	case Knight:
		return parseSANKnightMove(p, square)
	case Bishop:
		return parseSANBishopMove(p, square)
	case Queen:
		return parseSANQueenMove(p, square)
	case King:
		return parseSANKingMove(p, square)
	default:
		return Move{}, fmt.Errorf("could not parse SAN move: input, %s", s)
	}
}

// TODO reduce repetition
func parseSANRookMove(p *Position, toSquare Square) (Move, error) {
	isAmbiguous := false
	fromSquare := NoSquare
	for currentSquare := squareToLeft(toSquare); currentSquare != NoSquare; currentSquare = squareToLeft(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Rook && piece.Color == p.Turn {
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareAbove(toSquare); currentSquare != NoSquare; currentSquare = squareAbove(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Rook && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToRight(toSquare); currentSquare != NoSquare; currentSquare = squareToRight(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Rook && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareBelow(toSquare); currentSquare != NoSquare; currentSquare = squareBelow(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Rook && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	if isAmbiguous {
		return Move{}, fmt.Errorf("invalid SAN rook move: ambiguous move (multiple possible pieces)")
	}
	if fromSquare == NoSquare {
		return Move{}, fmt.Errorf("invalid SAN rook move: could not find piece to move")
	}
	return Move{FromSquare: fromSquare, ToSquare: toSquare}, nil
}

func parseSANKnightMove(p *Position, toSquare Square) (Move, error) {
	isAmbiguous := false
	fromSquare := NoSquare
	currentSquare := squareAbove(squareAbove(squareToRight(toSquare)))
	piece := p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		fromSquare = currentSquare
	}
	currentSquare = squareAbove(squareToRight(squareToRight(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareToRight(squareToRight(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareBelow(squareToRight(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareBelow(squareToLeft(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareToLeft(squareToLeft(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareAbove(squareToLeft(squareToLeft(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareAbove(squareAbove(squareToLeft(toSquare)))
	piece = p.PieceAt(currentSquare)
	if piece.Type == Knight && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	if isAmbiguous {
		return Move{}, errors.New("")
	}
	if isAmbiguous {
		return Move{}, fmt.Errorf("invalid SAN knight move: ambiguous move (multiple possible pieces)")
	}
	if fromSquare == NoSquare {
		return Move{}, fmt.Errorf("invalid SAN knight move: could not find piece to move")
	}
	return Move{FromSquare: fromSquare, ToSquare: toSquare}, nil
}

// TODO reduce repetition
func parseSANBishopMove(p *Position, toSquare Square) (Move, error) {
	isAmbiguous := false
	fromSquare := NoSquare
	for currentSquare := squareAbove(squareToLeft(toSquare)); currentSquare != NoSquare; currentSquare = squareAbove(squareToLeft(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Bishop && piece.Color == p.Turn {
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToRight(squareAbove(toSquare)); currentSquare != NoSquare; currentSquare = squareToRight(squareAbove(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Bishop && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareBelow(squareToRight(toSquare)); currentSquare != NoSquare; currentSquare = squareBelow(squareToRight(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Bishop && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToLeft(squareBelow(toSquare)); currentSquare != NoSquare; currentSquare = squareToLeft(squareBelow(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Bishop && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	if isAmbiguous {
		return Move{}, fmt.Errorf("invalid SAN bishop move: ambiguous move (multiple possible pieces)")
	}
	if fromSquare == NoSquare {
		return Move{}, fmt.Errorf("invalid SAN bishop move: could not find piece to move")
	}
	return Move{FromSquare: fromSquare, ToSquare: toSquare}, nil
}

// TODO reduce repetition
func parseSANQueenMove(p *Position, toSquare Square) (Move, error) {
	isAmbiguous := false
	fromSquare := NoSquare
	for currentSquare := squareToLeft(toSquare); currentSquare != NoSquare; currentSquare = squareToLeft(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareAbove(toSquare); currentSquare != NoSquare; currentSquare = squareAbove(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToRight(toSquare); currentSquare != NoSquare; currentSquare = squareToRight(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareBelow(toSquare); currentSquare != NoSquare; currentSquare = squareBelow(currentSquare) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareAbove(squareToLeft(toSquare)); currentSquare != NoSquare; currentSquare = squareAbove(squareToLeft(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToRight(squareAbove(toSquare)); currentSquare != NoSquare; currentSquare = squareToRight(squareAbove(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareBelow(squareToRight(toSquare)); currentSquare != NoSquare; currentSquare = squareBelow(squareToRight(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	for currentSquare := squareToLeft(squareBelow(toSquare)); currentSquare != NoSquare; currentSquare = squareToLeft(squareBelow(currentSquare)) {
		piece := p.PieceAt(currentSquare)
		if piece.Type == Queen && piece.Color == p.Turn {
			if fromSquare != NoSquare {
				isAmbiguous = true
			}
			fromSquare = currentSquare
			break
		}
		if piece != NoPiece {
			break
		}
	}
	if isAmbiguous {
		return Move{}, fmt.Errorf("invalid SAN Queen move: ambiguous move (multiple possible pieces)")
	}
	if fromSquare == NoSquare {
		return Move{}, fmt.Errorf("invalid SAN Queen move: could not find piece to move")
	}
	return Move{FromSquare: fromSquare, ToSquare: toSquare}, nil
}

func parseSANKingMove(p *Position, toSquare Square) (Move, error) {
	isAmbiguous := false
	fromSquare := NoSquare
	currentSquare := squareAbove(toSquare)
	piece := p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		fromSquare = currentSquare
	}
	currentSquare = squareAbove(squareToRight(toSquare))
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareToRight(toSquare)
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareToRight(toSquare))
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(toSquare)
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareBelow(squareToLeft(toSquare))
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareToLeft(toSquare)
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	currentSquare = squareAbove(squareToLeft(toSquare))
	piece = p.PieceAt(currentSquare)
	if piece.Type == King && piece.Color == p.Turn {
		if fromSquare != NoSquare {
			isAmbiguous = true
		}
		fromSquare = currentSquare
	}
	if isAmbiguous {
		return Move{}, errors.New("")
	}
	if isAmbiguous {
		return Move{}, fmt.Errorf("invalid SAN King move: ambiguous move (multiple possible pieces)")
	}
	if fromSquare == NoSquare {
		return Move{}, fmt.Errorf("invalid SAN King move: could not find piece to move")
	}
	return Move{FromSquare: fromSquare, ToSquare: toSquare}, nil
}

// IsValidMove makes sure each of the elements in Move m are logical. Namely that the squares can be found on a chess board.
func IsValidMove(m Move) bool {
	return isValidSquare(m.FromSquare) && m.FromSquare != NoSquare &&
		isValidSquare(m.ToSquare) && m.ToSquare != NoSquare &&
		isValidPieceType(m.Promotion)
}
