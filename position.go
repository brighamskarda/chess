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

// String returns a representation of the board from white's perspective. No other information from the position is printed. It is recommended to use [Position.FormatString] to print from black's perspective.
func (p *Position) String() string {
	return p.FormatString(false)
}

// FormatString returns a representation of the board. No other information from the position is printed.
// blacksPerspective should be true if you which to print the board as if you were on black's side.
func (p *Position) FormatString(blacksPerspective bool) string {
	if blacksPerspective {
		str := strings.Builder{}
		rank := '1'
		for index := len(p.Board) - 1; index >= 0; index-- {
			piece := p.Board[index]
			if index%8 == 7 {
				str.WriteRune(rank)
				rank += 1
			}
			str.WriteString(piece.String())
			if index%8 == 0 {
				str.WriteRune('\n')
			}
		}
		str.WriteString(" HGFEDCBA")
		return str.String()
	} else {
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

// Position.Move does no checking of move legality. For checked moves use [Game.Move], or check that your move is in the
// list provided by [GenerateLegalMoves]. All parts of the position are updated including en passant and castling rights based
// how the the move interacts with the board.
func (p *Position) Move(m Move) {
	if !isValidMove(m) {
		return
	}
	p.updateMoveCounts(m)
	p.movePiece(m)
	p.updateTurn()
	p.updateCastleRights(m)
	p.updateEnPassant(m)
}

func (p *Position) updateMoveCounts(m Move) {
	if p.PieceAt(m.FromSquare).Type == Pawn || p.PieceAt(m.ToSquare) != NoPiece || isCastleMove(p, m) {
		p.HalfMove = 0
	} else {
		p.HalfMove++
	}
	if p.Turn == Black {
		p.FullMove++
	}
}

func (p *Position) movePiece(m Move) {
	if isCastleMove(p, m) {
		p.performCastleMove(m)
		return
	}
	pieceToMove := p.PieceAt(m.FromSquare)
	if m.Promotion != NoPieceType {
		pieceToMove.Type = m.Promotion
	}
	p.SetPieceAt(m.FromSquare, NoPiece)
	p.SetPieceAt(m.ToSquare, pieceToMove)
	if m.ToSquare == p.EnPassant && p.PieceAt(m.ToSquare).Type == Pawn {
		if p.PieceAt(m.ToSquare).Color == White {
			p.SetPieceAt(Square{File: m.ToSquare.File, Rank: m.ToSquare.Rank - 1}, NoPiece)
		}
		if p.PieceAt(m.ToSquare).Color == Black {
			p.SetPieceAt(Square{File: m.ToSquare.File, Rank: m.ToSquare.Rank + 1}, NoPiece)
		}
	}
}

func (p *Position) performCastleMove(m Move) {
	p.SetPieceAt(m.ToSquare, p.PieceAt(m.FromSquare))
	p.SetPieceAt(m.FromSquare, NoPiece)
	if m.ToSquare == G1 {
		p.SetPieceAt(H1, NoPiece)
		p.SetPieceAt(F1, WhiteRook)
	}
	if m.ToSquare == C1 {
		p.SetPieceAt(A1, NoPiece)
		p.SetPieceAt(D1, WhiteRook)
	}
	if m.ToSquare == G8 {
		p.SetPieceAt(H8, NoPiece)
		p.SetPieceAt(F8, BlackRook)
	}
	if m.ToSquare == C8 {
		p.SetPieceAt(A8, NoPiece)
		p.SetPieceAt(D8, BlackRook)
	}
}

func (p *Position) updateTurn() {
	if p.Turn == White {
		p.Turn = Black
	} else if p.Turn == Black {
		p.Turn = White
	}
}

func (p *Position) updateCastleRights(m Move) {
	switch m.FromSquare {
	case E1:
		p.WhiteKingSideCastle = false
		p.WhiteQueenSideCastle = false
	case E8:
		p.BlackKingSideCastle = false
		p.BlackQueenSideCastle = false
	case A1:
		p.WhiteQueenSideCastle = false
	case H1:
		p.WhiteKingSideCastle = false
	case A8:
		p.BlackQueenSideCastle = false
	case H8:
		p.BlackKingSideCastle = false
	}
}

func (p *Position) updateEnPassant(m Move) {
	if p.PieceAt(m.ToSquare).Type == Pawn &&
		((m.ToSquare.Rank == Rank4 && m.FromSquare.Rank == Rank2) ||
			(m.ToSquare.Rank == Rank5 && m.FromSquare.Rank == Rank7)) {
		if m.FromSquare.Rank == Rank2 {
			p.EnPassant = Square{File: m.ToSquare.File, Rank: Rank3}
		} else if m.FromSquare.Rank == Rank7 {
			p.EnPassant = Square{File: m.ToSquare.File, Rank: Rank6}

		}
	} else {
		p.EnPassant = NoSquare
	}
}

func isCastleMove(p *Position, m Move) bool {
	return (m.FromSquare == E1 && (m.ToSquare == G1 || m.ToSquare == C1) && p.PieceAt(m.FromSquare) == WhiteKing) ||
		(m.FromSquare == E8 && (m.ToSquare == G8 || m.ToSquare == C8) && p.PieceAt(E8) == BlackKing)
}

func squareToIndex(s Square) int {
	index := 0
	index += int(s.File - 1)
	index += int(Rank8-s.Rank) * 8
	return index
}

func indexToSquare(index int) Square {
	file := File(index%8 + 1)
	rank := Rank(8 - (index / 8))
	square := Square{file, rank}
	if !isValidSquare(square) {
		return NoSquare
	}
	return square
}

// ParseFen take only fully formed valid FEN strings. All parts of the FEN must be present, though the position need not necessarily be valid.
func ParseFen(fen string) (*Position, error) {
	words := strings.Split(fen, " ")
	if len(words) != 6 {
		return &Position{}, errors.New("invalid fen, fen does not have 6 required parts")
	}
	board, err := parseFenPos(words[0])
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	turn, err := parseTurn(words[1])
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	castleRights, err := parseCastleRights(words[2])
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	enPassant, err := ParseSquare(words[3])
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, %w", err)
	}
	halfMove, err := strconv.ParseUint(words[4], 10, 16)
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, can't parse halfMove, %w", err)
	}
	fullMove, err := strconv.ParseUint(words[5], 10, 16)
	if err != nil {
		return &Position{}, fmt.Errorf("invalid fen, can't parse fullMove, %w", err)
	}

	return &Position{
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

func (p *Position) GenerateFen() string {
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
			numBlank = 0
		}
		if piece == NoPiece {
			numBlank++
			currentFile++
			continue
		}
		fen.WriteString(piece.String())
		currentFile++
	}
	if numBlank > 0 {
		fen.WriteString(strconv.FormatUint(uint64(numBlank), 10))
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
//   - All pieces are valid chess pieces
func IsValidPosition(p *Position) bool {
	return checkKings(p) &&
		checkNoInvalidPawns(p) &&
		checkCastlingRightsLogical(p) &&
		checkEnPassantLogical(p) &&
		checkTurnIsSet(p) &&
		checkAllPiecesValid(p)
}

func checkKings(p *Position) bool {
	numWhiteKings := 0
	numBlackKings := 0
	for _, piece := range p.Board {
		if piece == WhiteKing {
			numWhiteKings++
		}
		if piece == BlackKing {
			numBlackKings++
		}
	}
	if numWhiteKings != 1 || numBlackKings != 1 {
		return false
	}
	return true
}

func checkNoInvalidPawns(p *Position) bool {
	return checkNoInvalidWhitePawns(p) && checkNoInvalidBlackPawns(p)
}

func checkNoInvalidWhitePawns(p *Position) bool {
	for i := 0; i < 8; i++ {
		if p.Board[i] == WhitePawn {
			return false
		}
	}
	return true
}

func checkNoInvalidBlackPawns(p *Position) bool {
	for i := 56; i < 64; i++ {
		if p.Board[i] == BlackPawn {
			return false
		}
	}
	return true
}

func checkCastlingRightsLogical(p *Position) bool {
	if p.WhiteKingSideCastle {
		if p.PieceAt(E1) != WhiteKing || p.PieceAt(H1) != WhiteRook {
			return false
		}
	}
	if p.WhiteQueenSideCastle {
		if p.PieceAt(E1) != WhiteKing || p.PieceAt(A1) != WhiteRook {
			return false
		}
	}
	if p.BlackKingSideCastle {
		if p.PieceAt(E8) != BlackKing || p.PieceAt(H8) != BlackRook {
			return false
		}
	}
	if p.BlackQueenSideCastle {
		if p.PieceAt(E8) != BlackKing || p.PieceAt(A8) != BlackRook {
			return false
		}
	}
	return true
}

func checkEnPassantLogical(p *Position) bool {
	if p.EnPassant == NoSquare {
		return true
	}
	if !isValidColor(p.Turn) || p.Turn == NoColor {
		return false
	}
	if !isValidSquare(p.EnPassant) {
		return false
	}

	if p.Turn == White {
		return checkValidBlackEnPassant(p)
	} else {
		return checkValidWhiteEnPassant(p)
	}
}

func checkValidWhiteEnPassant(p *Position) bool {
	if p.EnPassant.Rank != Rank3 {
		return false
	}
	if p.PieceAt(p.EnPassant) != NoPiece {
		return false
	}
	expectedSquare := Square{p.EnPassant.File, p.EnPassant.Rank + 1}
	return p.PieceAt(expectedSquare) == WhitePawn
}

func checkValidBlackEnPassant(p *Position) bool {
	if p.EnPassant.Rank != Rank6 {
		return false
	}
	if p.PieceAt(p.EnPassant) != NoPiece {
		return false
	}
	expectedSquare := Square{p.EnPassant.File, p.EnPassant.Rank - 1}
	return p.PieceAt(expectedSquare) == BlackPawn
}

func checkTurnIsSet(p *Position) bool {
	return p.Turn == White || p.Turn == Black
}

func checkAllPiecesValid(p *Position) bool {
	for _, piece := range p.Board {
		if !isValidPiece(piece) {
			return false
		}
	}
	return true
}

func findKing(p *Position, c Color) Square {
	for index, piece := range p.Board {
		if piece.Type == King && piece.Color == c {
			return indexToSquare(index)
		}
	}
	return NoSquare
}
