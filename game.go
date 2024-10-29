package chess

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"time"
)

// Game is guaranteed to always represent a valid game of chess. Invalid positions are not allowed,
// but the move history may not represent an entire game (undoing all moves may not lead to the starting
// chess position). Game can be used to parse and generate PGNs. Additionally Game provides many wrapper methods for
// functions that normally take position as an argument, increasing convenience and decreasing the need to copy Game's
// position excessively.
//
// The zero value for Game is not valid and should not be used. Game should always be duplicated using the
// [Game.Copy] function.
type Game struct {
	position    *Position
	moveHistory []Move
	tags        map[string]string
}

type Result byte

const (
	NoResult Result = iota
	WhiteWins
	BlackWins
	Draw
)

func isValidResult(r Result) bool {
	return r <= Draw
}

func parseResult(s string) Result {
	switch s {
	case "1-0":
		return WhiteWins
	case "0-1":
		return BlackWins
	case "1/2-1/2":
		return Draw
	default:
		return NoResult
	}
}

// String returns
//   - "1-0" if [WhiteWins]
//   - "0-1" if [BlackWins]
//   - "1/2-1/2" if [Draw]
//   - "*" otherwise
func (r Result) String() string {
	switch r {
	case WhiteWins:
		return "1-0"
	case BlackWins:
		return "0-1"
	case Draw:
		return "1/2-1/2"
	default:
		return "*"
	}
}

// NewGame returns a [*Game] representing the starting position for a game of chess.
func NewGame() *Game {
	position, _ := ParseFen(DefaultFen)
	currentDate := time.Now()
	tags := map[string]string{
		"Event":  "Golang chess match",
		"Site":   "github.com/brighamskarda/chess",
		"Date":   fmt.Sprintf("%4d.%2d.%2d", currentDate.Year(), currentDate.Month(), currentDate.Day()),
		"Round":  "1",
		"White":  "White Player",
		"Black":  "Black Player",
		"Result": "*",
	}
	return &Game{
		position:    position,
		moveHistory: []Move{},
		tags:        tags,
	}
}

// Move performs the given move. If move m is not legal g remains unchanged and an error is returned.
// If the move is legal the result tag is set to * (NoResult). If the position ends in checkmate
// or stalemate the result tag is updated accordingly.
func (g *Game) Move(m Move) error {
	legalMoves := GenerateLegalMoves(g.position)
	if !slices.Contains(legalMoves, m) {
		return fmt.Errorf("%s is not a legal move", m)
	}
	g.position.Move(m)
	g.moveHistory = append(g.moveHistory, m)
	if IsCheckMate(g.position) {
		if g.position.Turn == Black {
			g.SetResult(WhiteWins)
		}
		if g.position.Turn == White {
			g.SetResult(BlackWins)
		}
	} else if IsStaleMate(g.position) {
		g.SetResult(Draw)
	} else {
		g.SetResult(NoResult)
	}
	return nil
}

// MoveSan is a helper function that automatically performs an SAN formatted move. SAN format is specified here: http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.2.3
func (g *Game) MoveSan(s string) error {
	move, err := ParseSANMove(g.position, s)
	if err != nil {
		return err
	}
	return g.Move(move)
}

// Returns a copy of current game.
func (g *Game) Copy() *Game {
	positionCopy := *g.position
	gameCopy := &Game{
		position:    &positionCopy,
		moveHistory: slices.Clone(g.moveHistory),
		tags:        maps.Clone(g.tags),
	}
	return gameCopy
}

// String returns the games current position as a string..
func (g *Game) String() string {
	return g.position.String()
}

// Position returns a copy of the game's position.
func (g *Game) Position() *Position {
	var pos Position = *g.position
	return &pos
}

func (g *Game) Turn() Color {
	return g.position.Turn
}

func (g *Game) WhiteKingSideCastle() bool {
	return g.position.WhiteKingSideCastle
}

func (g *Game) WhiteQueenSideCastle() bool {
	return g.position.WhiteQueenSideCastle
}

func (g *Game) BlackKingSideCastle() bool {
	return g.position.BlackKingSideCastle
}

func (g *Game) BlackQueenSideCastle() bool {
	return g.position.BlackQueenSideCastle
}

func (g *Game) HalfMove() uint16 {
	return g.position.HalfMove
}

func (g *Game) FullMove() uint16 {
	return g.position.FullMove
}

func (g *Game) EnPassant() Square {
	return g.position.EnPassant
}

func (g *Game) IsCheckMate() bool {
	return IsCheckMate(g.position)
}

func (g *Game) IsStaleMate() bool {
	return IsStaleMate(g.position)
}

func (g *Game) GetResult() Result {
	return parseResult(g.tags["Result"])
}

func (g *Game) SetResult(r Result) {
	g.tags["Result"] = r.String()
}

func (g *Game) LegalMoves() []Move {
	return GenerateLegalMoves(g.position)
}

func (g *Game) GetTag(t string) (string, error) {
	s, ok := g.tags[t]
	if !ok {
		return s, fmt.Errorf("game does not contain tag \"%s\"", t)
	}
	return s, nil
}

// SetTag sets any tag for the game so that it will show up in the pgn file. The Result, SetUp, and FEN tags cannot be set with this function. Please use the [Game.SetResult] function to set the result, and the [Game.SetPosition] function to set the other two tags.
func (g *Game) SetTag(tag string, value string) {
	if tag == "Result" || tag == "SetUp" || tag == "FEN" {
		return
	}
	g.tags[tag] = value
}

// Remove tag will remove any pgn tag except the 7 required tags specified [here], and the SetUp and FEN
// tags specified [here]. If you wish to remove the SetUp and FEN tags it is best to simply make a new game.
//
// [here]: http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.1.1
//
// [here.]: http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c9.7
func (g *Game) RemoveTag(tag string) {
	requiredTags := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result", "SetUp", "FEN"}
	if slices.Contains(requiredTags, tag) {
		return
	}
	delete(g.tags, tag)
}

func (g *Game) GetAllTags() map[string]string {
	return maps.Clone(g.tags)
}

// SetPosition sets the games position to the given position only if the position is a valid chess
// position. The result tag is updated to match if the game is in mate, or could still be going.
// Move history is cleared, and the pgn tags "SetUp" and "FEN" are set accordingly.
func (g *Game) SetPosition(p *Position) error {
	if !IsValidPosition(p) {
		return errors.New("invalid position: game can only set valid chess positions")
	}
	*g.position = *p
	g.moveHistory = []Move{}
	g.tags["SetUp"] = "1"
	g.tags["FEN"] = GenerateFen(p)
	if IsCheckMate(p) {
		if p.Turn == White {
			g.SetResult(BlackWins)
		} else {
			g.SetResult(WhiteWins)
		}
	} else if IsStaleMate(p) {
		g.SetResult(Draw)
	} else {
		g.SetResult(NoResult)
	}
	return nil
}

// HasThreeFoldRepetition returns true if the game has been in the exact same position (including castling rights)
// At least three times at any point during the entire game.
func (g *Game) HasThreeFoldRepetition() bool {
	// TODO Use a map instead for more efficiency
	allPositions := generateAllGamePositions(g)
	for index, pos1 := range allPositions[:len(allPositions)-1] {
		numEquivalentPositions := 1
		for _, pos2 := range allPositions[index+1:] {
			if positionsEqualNoMoveCounter(&pos1, &pos2) {
				numEquivalentPositions++
				if numEquivalentPositions >= 3 {
					return true
				}
			}
		}
	}
	return false
}

func generateAllGamePositions(g *Game) []Position {
	pos, _ := ParseFen(DefaultFen)
	if fen, err := g.GetTag("FEN"); err == nil {
		pos, _ = ParseFen(fen)
	}
	allPositions := make([]Position, 0, g.FullMove()*2+1)
	allPositions = append(allPositions, *pos)
	for _, move := range g.moveHistory {
		pos.Move(move)
		allPositions = append(allPositions, *pos)
	}
	return allPositions
}

func positionsEqualNoMoveCounter(pos1 *Position, pos2 *Position) bool {
	return pos1.Board == pos2.Board &&
		pos1.Turn == pos2.Turn &&
		pos1.WhiteKingSideCastle == pos2.WhiteKingSideCastle &&
		pos1.WhiteQueenSideCastle == pos2.WhiteQueenSideCastle &&
		pos1.BlackKingSideCastle == pos2.BlackKingSideCastle &&
		pos1.BlackQueenSideCastle == pos2.BlackQueenSideCastle &&
		pos1.EnPassant == pos2.EnPassant
}

// CanClaimDraw returns true if one of the following conditions is true and the game is not in checkmate
//   - The half move counter is >= 100 (indicating that no piece has been taken, nor pawn moved forward for 50 moves)
//   - The game contains a three fold repetition (the exact same position has occurred three times in the game)
//   - The game is in stalemate
func (g *Game) CanClaimDraw() bool {
	return (g.position.HalfMove >= 100 ||
		g.HasThreeFoldRepetition() ||
		g.IsStaleMate()) && !g.IsCheckMate()
}

// WritePgn writes a pgn representation of g to w.
func WritePgn(g *Game, w io.Writer) error {
	sevenTags := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}
	for _, tag := range sevenTags {
		_, err := fmt.Fprintf(w, "[%s \"%s\"]\n", tag, g.tags[tag])
		if err != nil {
			return fmt.Errorf("unable to write pgn: %w", err)
		}
	}
	for tag, val := range g.tags {
		if !slices.Contains(sevenTags, tag) {
			_, err := fmt.Fprintf(w, "[%s \"%s\"]\n", tag, val)
			if err != nil {
				return fmt.Errorf("unable to write pgn: %w", err)
			}
		}
	}
	_, err := fmt.Fprint(w, "\n")
	if err != nil {
		return fmt.Errorf("unable to write pgn: %w", err)
	}

	newGame := NewGame()
	if fen, keyExists := g.tags["FEN"]; keyExists {
		new_position, err := ParseFen(fen)
		if err != nil {
			return fmt.Errorf("unable to write pgn: game contains invalid FEN tag: %w", err)
		}
		newGame.SetPosition(new_position)
	}
	counter := 1
	for _, move := range g.moveHistory {
		sanMoveString := move.SanString(newGame.position)
		if newGame.Turn() == White {
			_, err := fmt.Fprintf(w, "%d. %s ", counter, sanMoveString)
			if err != nil {
				return fmt.Errorf("unable to write pgn: %w", err)
			}
			counter++
		}
		if newGame.Turn() == Black {
			_, err := fmt.Fprintf(w, "%s ", sanMoveString)
			if err != nil {
				return fmt.Errorf("unable to write pgn: %w", err)
			}
		}
		err := newGame.Move(move)
		if err != nil {
			return fmt.Errorf("unable to write pgn: %w", err)
		}
	}
	_, err = fmt.Fprintf(w, "%s", g.GetResult().String())
	if err != nil {
		return fmt.Errorf("unable to write pgn: %w", err)
	}
	return nil
}

// ReadPgn attempts to create a [Game] from r. Parsing should be improved in the future, but for now only well formatted
// pgns containing a single game are accepted. Refer to http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm for
// specific details of how a pgn should be formatted.
func ReadPgn(r io.Reader) (*Game, error) {
	pgn_bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read pgn failed: %w", err)
	}

	pgn_lines := strings.Split(strings.Replace(string(pgn_bytes), "\r", "", -1), "\n")

	game := NewGame()
	var result Result
	tagsComplete := false
	for _, line := range pgn_lines {
		if tagsComplete {
			err := parsePgnMoves(game, line)
			if err != nil {
				return nil, fmt.Errorf("read pgn failed: %w", err)
			}
		} else if line == "" {
			tagsComplete = true
			result = game.GetResult()
		} else {
			err := parsePgnTag(game, line)
			if err != nil {
				return nil, fmt.Errorf("read pgn failed: %w", err)
			}
		}
	}

	game.SetResult(result)

	return game, nil
}

func parsePgnMoves(g *Game, moves string) error {
	splitMoves := strings.Split(moves, " ")
	possibleResults := []string{"1-0", "0-1", "1/2-1/2", "*"}

	for _, move := range splitMoves {
		if strings.Contains(move, ".") || slices.Contains(possibleResults, move) {
			continue
		}
		err := g.MoveSan(move)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsePgnTag(g *Game, tag string) error {
	splitTag := strings.SplitN(tag[1:len(tag)-1], " ", 2)
	if len(splitTag) != 2 {
		return fmt.Errorf("invalid pgn tag: %s", tag)
	}
	if splitTag[0] == "Result" {
		g.SetResult(parseResult(strings.ReplaceAll(splitTag[1], "\"", "")))
	}
	g.SetTag(splitTag[0], strings.ReplaceAll(splitTag[1], "\"", ""))
	return nil
}
