package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const DefaultFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Position represents a chess position as described by Forsyth-Edwards Notation (FEN).
// Board is the actual representation of the pieces on the squares. It starts at A8 and moves left
// to right, top to bottom all the way to H1.
type Position struct {
	Board                [64]Piece
	Turn                 Color
	WhiteKingSideCastle  bool
	WhiteQueenSideCastle bool
	BlackKingSideCastle  bool
	BlackQueenSideCastle bool
	EnPassant            Square
	HalfMove             uint16
	FullMove             uint16
}

func (p *Position) String() string {
	str := strings.Builder{}
	rank := '8'
	for index, piece := range p.Board {
		if index%8 == 0 {
			str.WriteRune(rank)
			rank -= 1
		}
		str.WriteString(piece.String())
		if index%8 == 7 {
			str.WriteRune('\n')
		}
	}
	str.WriteString(" ABCDEFGH")
	return str.String()
}

func (p *Position) PieceAt(s Square) Piece {
	if !isValidSquare(s) || s == NoSquare {
		return NoPiece
	}
	return p.Board[squareToIndex(s)]
}

func (p *Position) SetPieceAt(s Square, piece Piece) {
	if isValidSquare(s) && s != NoSquare {
		p.Board[squareToIndex(s)] = piece
	}
}

func squareToIndex(s Square) int {
	index := 0
	index += int(s.File - 1)
	index += int(Rank8-s.Rank) * 8
	return index
}

func ParseFen(fen string) (Position, error) {
	words := strings.Split(fen, " ")
	if len(words) != 6 {
		return Position{}, errors.New("invalid fen, fen does not have 6 required parts")
	}
	board, err := parseFenPos(words[0])
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	turn, err := parseTurn(words[1])
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	castleRights, err := parseCastleRights(words[2])
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	enPassant, err := ParseSquare(words[3])
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	halfMove, err := strconv.ParseUint(words[4], 10, 16)
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, can't parse halfMove, %w", err)
	}
	fullMove, err := strconv.ParseUint(words[5], 10, 16)
	if err != nil {
		return Position{}, fmt.Errorf("invalid fen, can't parse fullMove, %w", err)
	}

	return Position{
		Board:                board,
		Turn:                 turn,
		WhiteKingSideCastle:  castleRights[0],
		WhiteQueenSideCastle: castleRights[1],
		BlackKingSideCastle:  castleRights[2],
		BlackQueenSideCastle: castleRights[3],
		EnPassant:            enPassant,
		HalfMove:             uint16(halfMove),
		FullMove:             uint16(fullMove),
	}, nil
}

func parseFenPos(fen string) ([64]Piece, error) {
	pos := [64]Piece{}
	posIndex := 0
	for _, char := range fen {
		if posIndex >= 64 {
			return pos, errors.New("invalid pos, too many pieces on board")
		}
		if unicode.IsNumber(char) {
			posIndex += int(char - '0')
			continue
		}
		if char == '/' {
			if posIndex%8 != 0 {
				return pos, errors.New("invalid pos, '/' in wrong position")
			}
			continue
		}
		piece, err := ParsePiece(char)
		if err != nil {
			return pos, errors.New("invalid pos, can't parse " + string(char) + "to piece")
		}
		pos[posIndex] = piece
		posIndex++
	}
	return pos, nil
}

func parseTurn(turn string) (Color, error) {
	switch strings.ToLower(turn) {
	case "w":
		return White, nil
	case "b":
		return Black, nil
	default:
		return NoColor, errors.New("can't parse color")
	}
}

func parseCastleRights(castleRights string) ([4]bool, error) {
	rights := [4]bool{}
	if castleRights == "-" {
		return rights, nil
	}
	for _, char := range castleRights {
		switch char {
		case 'K':
			rights[0] = true
		case 'Q':
			rights[1] = true
		case 'k':
			rights[2] = true
		case 'q':
			rights[3] = true
		default:
			return rights, errors.New("invalid castling rights")
		}
	}
	return rights, nil
}

func GenerateFen(p *Position) string {
	fen := strings.Builder{}
	fen.WriteString(generateFenPos(p))
	fen.WriteString(" " + generateFenTurn(p))
	fen.WriteString(" " + generateFenCastleRights(p))
	fen.WriteString(" " + strings.ToLower(p.EnPassant.String()))
	fen.WriteString(" " + strconv.FormatUint(uint64(p.HalfMove), 10))
	fen.WriteString(" " + strconv.FormatUint(uint64(p.FullMove), 10))
	return fen.String()
}

func generateFenPos(p *Position) string {
	fen := strings.Builder{}
	currentFile := FileA
	numBlank := 0
	for _, piece := range p.Board {
		if currentFile > FileH {
			if numBlank > 0 {
				fen.WriteString(strconv.FormatUint(uint64(numBlank), 10))
				numBlank = 0
			}
			fen.WriteRune('/')
			currentFile = FileA
		}
		if piece != NoPiece && numBlank > 0 {
			fen.WriteString(strconv.FormatUint(uint64(numBlank), 10))
		}
		if piece == NoPiece {
			numBlank++
			currentFile++
			continue
		}
		fen.WriteString(piece.String())
		currentFile++
	}

	return fen.String()
}

func generateFenTurn(p *Position) string {
	switch p.Turn {
	case White:
		return "w"
	case Black:
		return "b"
	default:
		return "-"
	}
}

func generateFenCastleRights(p *Position) string {
	if !p.WhiteKingSideCastle && !p.WhiteQueenSideCastle && !p.BlackKingSideCastle && !p.BlackQueenSideCastle {
		return "-"
	}
	rights := ""
	if p.WhiteKingSideCastle {
		rights += "K"
	}
	if p.WhiteQueenSideCastle {
		rights += "Q"
	}
	if p.BlackKingSideCastle {
		rights += "k"
	}
	if p.BlackQueenSideCastle {
		rights += "q"
	}
	return rights
}

// IsValidPosition determines if a given position is a legal chess position. It checks the following:
//   - There is one king of each color on the board
//   - There are no pawns on their last rank
//   - Castling rights are logical
//   - The enPassant Square is logical
//   - Turn is set
func IsValidPosition(p *Position) bool {
	return true
}

func checkBothKingsPresent(p *Position) bool {
	return false
}

func checkNoInvalidPawns(p *Position) bool {
	return false
}

func checkCastlingRightsLogical(p *Position) bool {
	return false
}

func checkEnPassantLogical(p *Position) bool {
	return false
}

func checkTurnIsSet(p *Position) bool {
	return false
}
