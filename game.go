// Copyright (C) 2025 Brigham Skarda

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package chess

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Result represents the result of a chess [Game].
type Result uint8

const (
	NoResult Result = iota
	WhiteWins
	BlackWins
	Draw
)

// MarshalText returns "1-0" for [WhiteWins], "0-1" for [BlackWins], "1/2-1/2" for [Draw], and * otherwise. err is always nil.
func (r Result) MarshalText() (text []byte, err error) {
	switch r {
	case NoResult:
		return []byte("*"), nil
	case WhiteWins:
		return []byte("1-0"), nil
	case BlackWins:
		return []byte("0-1"), nil
	case Draw:
		return []byte("1/2-1/2"), nil
	default:
		return nil, fmt.Errorf("invalid Result %#v", r)
	}
}

// UnmarshalText sets [WhiteWins] for "1-0", [BlackWins] for "0-1", [Draw] for "1/2-1/2", and [NoResult] for "*". Error is returned if text is not one of these four values.
func (r *Result) UnmarshalText(text []byte) error {
	switch string(text) {
	case "1-0":
		*r = WhiteWins
		return nil
	case "0-1":
		*r = BlackWins
		return nil
	case "1/2-1/2":
		*r = Draw
		return nil
	case "*":
		*r = NoResult
		return nil
	default:
		return fmt.Errorf("could not parse result %q", text)
	}
}

// PgnMove is an expanded move struct used in [Game]. It provides fields for Numeric Annotation Glyphs, commentary and Recursive Annotation Variation (RAV - move variations).
type PgnMove struct {
	Move              Move
	NumericAnnotation uint8
	PreCommentary     []string
	PostCommentary    []string
	// Variations supports multiple variations. Hence it is a 2d slice. The first move in a variation should replace the current move.
	Variations [][]PgnMove
}

// Copy provides a copy of the PgnMove. This is a deep copy so the variations are separate.
func (m PgnMove) Copy() PgnMove {
	newPgnMove := PgnMove{
		Move:              m.Move,
		NumericAnnotation: m.NumericAnnotation,
		PreCommentary:     slices.Clone(m.PreCommentary),
		PostCommentary:    slices.Clone(m.PostCommentary),
		Variations:        make([][]PgnMove, 0, len(m.Variations)),
	}
	for _, variation := range m.Variations {
		newVariation := make([]PgnMove, 0, len(variation))
		for _, move := range variation {
			newVariation = append(newVariation, move.Copy())
		}
		newPgnMove.Variations = append(newPgnMove.Variations, newVariation)
	}
	return newPgnMove
}

// Game represents all parts of the PGN game notation standard found here. https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm.
//
// It ensures that only legal moves are performed and keeps track of move history. It also provides utilities for determining if a draw can be claimed.
//
// This library does not support chess 960, largely because special castling rights are not implemented. Otherwise starting games from arbitrary positions is supported.
//
// The field descriptions given are pulled from the PGN specification.
type Game struct {
	pos         *Position
	moveHistory []PgnMove
	// moves are the current legal moves, nil if they haven't been generated. If it is an empty slice then there are no legal moves.
	moves []Move

	// Event should be reasonably descriptive. Abbreviations are to be avoided unless absolutely necessary. A consistent event naming should be used to help facilitate database scanning. If the name of the event is unknown, a single question mark should appear as the tag value.
	Event string
	// Site should include city and region names along with a standard name for the country. The use of the IOC (International Olympic Committee) three letter names is suggested for those countries where such codes are available. If the site of the event is unknown, a single question mark should appear as the tag value. A comma may be used to separate a city from a region. No comma is needed to separate a city or region from the IOC country code. A later section of this document gives a list of three letter nation codes along with a few additions for "locations" not covered by the IOC.
	Site string
	// Date gives the starting date for the game. (Note: this is not necessarily the same as the starting date for the event.) The date is given with respect to the local time of the site given in the Event tag. The Date tag value field always uses a standard ten character format: "YYYY.MM.DD". The first four characters are digits that give the year, the next character is a period, the next two characters are digits that give the month, the next character is a period, and the final two characters are digits that give the day of the month. If the any of the digit fields are not known, then question marks are used in place of the digits.
	Date string
	// Round gives the playing round for the game. In a match competition, this value is the number of the game played. If the use of a round number is inappropriate, then the field should be a single hyphen character. If the round is unknown, a single question mark should appear as the tag value.
	//
	// Some organizers employ unusual round designations and have multipart playing rounds and sometimes even have conditional rounds. In these cases, a multipart round identifier can be made from a sequence of integer round numbers separated by periods. The leftmost integer represents the most significant round and succeeding integers represent round numbers in descending hierarchical order.
	Round string
	// 	White is the name of the player or players of the white pieces. The names are given as they would appear in a telephone directory. The family or last name appears first. If a first name or first initial is available, it is separated from the family name by a comma and a space. Finally, one or more middle initials may appear. (Wherever a comma appears, the very next character should be a space. Wherever an initial appears, the very next character should be a period.) If the name is unknown, a single question mark should appear as the tag value.
	//
	// The intent is to allow meaningful ASCII sorting of the tag value that is independent of regional name formation customs. If more than one person is playing the white pieces, the names are listed in alphabetical order and are separated by the colon character between adjacent entries. A player who is also a computer program should have appropriate version information listed after the name of the program.
	//
	// The format used in the FIDE Rating Lists is appropriate for use for player name tags.
	White string
	// Black is the name of the player or players of the black pieces. The names are given here as they are for the White tag value.
	Black string
	// Result field is the result of the game. It is always exactly the same as the game termination marker that concludes the associated movetext. It is always one of four possible values: "1-0" (White wins), "0-1" (Black wins), "1/2-1/2" (drawn game), and "*" (game still in progress, game abandoned, or result otherwise unknown). Note that the digit zero is used in both of the first two cases; not the letter "O".
	//
	// Note that this field is not just a string, a type was provided to make it slightly easier to use for computers.
	Result Result
	// OtherTags is intended for custom PGN game tags. Some examples are provided here: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c9
	OtherTags map[string]string
}

// NewGame returns a fresh game of chess with the starting position initialized. Tags are set as follows:
//
// * Event - ?
//
// * Site - https://github.com/brighamskarda/chess
//
// * Date - <CurrentDate>
//
// * Round - 1
//
// * White - ?
//
// * Black - ?
//
// * Result - NoResult
func NewGame() *Game {
	pos := &Position{}
	pos.UnmarshalText([]byte(DefaultFEN))
	date := time.Now()
	return &Game{
		pos:         pos,
		moveHistory: []PgnMove{},
		moves:       nil,
		Event:       "?",
		Site:        "https://github.com/brighamskarda/chess",
		Date:        date.Format("2006.01.02"),
		Round:       "1",
		White:       "?",
		Black:       "?",
		Result:      NoResult,
		OtherTags:   map[string]string{},
	}
}

// NewGameFromFEN starts a game specified from the provided fen string. Returns an error if [Position.UnmarshalText] could not parse the FEN. Values are set the same as [NewGame]. The SetUp and FEN tags are also filled into OtherTags.
func NewGameFromFEN(fen string) (*Game, error) {
	pos := &Position{}
	err := pos.UnmarshalText([]byte(fen))
	if err != nil {
		return nil, fmt.Errorf("could not make game: %w", err)
	}
	date := time.Now()
	return &Game{
		pos:         pos,
		moveHistory: []PgnMove{},
		moves:       nil,
		Event:       "?",
		Site:        "https://github.com/brighamskarda/chess",
		Date:        date.Format("2006.01.02"),
		Round:       "1",
		White:       "?",
		Black:       "?",
		Result:      NoResult,
		OtherTags: map[string]string{
			"SetUp": "1",
			"FEN":   fen,
		},
	}, nil
}

// UnmarshalText is capable of unmarshaling a single game in pgn format.
//
// Games can start from any position, but all the moves must be legal.
//
// See also [ParsePGN]
func (g *Game) UnmarshalText(text []byte) error {
	// Be sure to not read lines beginning with %. These are comments.
	// Semicolons are commentary and go to the end of the line.
	pgn := string(text)
	pgn = strings.TrimSpace(pgn)
	lines := strings.Split(strings.ReplaceAll(pgn, "\r\n", "\n"), "\n")
	lines = removeCommentLines(lines)
	tags, movetext := separateTagsAndMovetext(lines)
	if len(movetext) == 0 {
		return fmt.Errorf("expected an empty newline after tags followed by movetext")
	}

	newG := NewGame()
	err := newG.parseTags(tags)
	if err != nil {
		return fmt.Errorf("could not parse pgn tags: %w", err)
	}

	err = newG.parseMovetext(movetext)
	if err != nil {
		return fmt.Errorf("could not parse movetext: %w", err)
	}

	*g = *newG
	return nil
}

func removeCommentLines(lines []string) []string {
	var filteredLines []string
	for _, line := range lines {
		if len(line) == 0 || line[0] != '%' {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

func separateTagsAndMovetext(lines []string) (tags []string, movetext []string) {
	emptyLineIndex := slices.Index(lines, "")
	if emptyLineIndex == -1 {
		return nil, nil
	}
	return lines[0:emptyLineIndex], lines[emptyLineIndex+1:]

}

func (g *Game) parseTags(lines []string) error {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] != '[' || line[len(line)-1] != ']' {
			return fmt.Errorf("tag %q missing square braces", line)
		}
		if len(line) <= 2 {
			// empty tag
			continue
		}
		if err := g.parseSingleTag(line[1 : len(line)-1]); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) parseSingleTag(tag string) error {
	i := 0
	for ; i < len(tag) && tag[i] != ' '; i++ {
	}
	name := tag[0:i]
	i++
	if i >= len(tag) {
		return fmt.Errorf("tag missing value after key")
	}
	if tag[i] != '"' {
		return fmt.Errorf("expected a single space in between tag name and opening quote, %q", tag)
	}
	i++
	bodyStart := i
	for ; i < len(tag)-1 && tag[i] != '"'; i++ {
	}
	if i >= len(tag) || tag[i] != '"' {
		return fmt.Errorf("missing closing quote for tag body, %q", tag)
	}
	body := tag[bodyStart:i]
	return g.setTag(name, body)
}

// setTag automatically set the 7 tag roster, and the Setup and FEN tags. When setting the FEN tag it will set the position as well.
func (g *Game) setTag(name string, body string) error {
	switch name {
	case "Event":
		g.Event = body
	case "Site":
		g.Site = body
	case "Date":
		g.Date = body
	case "Round":
		g.Round = body
	case "White":
		g.White = body
	case "Black":
		g.Black = body
	case "Result":
		err := g.Result.UnmarshalText([]byte(body))
		if err != nil {
			return fmt.Errorf("could not parse Result")
		}
	case "FEN":
		g.OtherTags["FEN"] = body
		err := g.pos.UnmarshalText([]byte(body))
		if err != nil {
			return fmt.Errorf("could not parse fen: %w", err)
		}
	default:
		g.OtherTags[name] = body
	}
	return nil
}

func (g *Game) parseMovetext(lines []string) error {
	text := ""
	for _, line := range lines {
		text += line
		text += "\n"
	}
	tokens, err := tokenizeMovetext(text)
	if err != nil {
		return err
	}
	if len(tokens) == 0 || tokens[len(tokens)-1].tokenType != result {
		return fmt.Errorf("there is no result at end of pgn")
	}

	moveHis, err := createMoveHistory(tokens, g.Position())
	if err != nil {
		return err
	}
	for _, m := range moveHis {
		err := g.Move(m.Move)
		if err != nil {
			return fmt.Errorf("found illegal move %s in pgn", m.Move)
		}
	}
	g.moveHistory = moveHis

	err = g.Result.UnmarshalText([]byte(tokens[len(tokens)-1].body))
	if err != nil {
		return fmt.Errorf("could not parse result %q", tokens[len(tokens)-1].body)
	}
	return nil
}

func createMoveHistory(tokens []pgnToken, pos *Position) ([]PgnMove, error) {
	moveHis := []PgnMove{}
	prevPos := pos.Copy()
	preCommentary := []string{}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.tokenType {
		case commentary:
			if len(moveHis) == 0 {
				preCommentary = append(preCommentary, t.body)
			} else {
				moveHis[len(moveHis)-1].PostCommentary = append(moveHis[len(moveHis)-1].PostCommentary, t.body)
			}
		case move:
			prevPos = pos.Copy()
			m, err := ParseSANMove(t.body, pos)
			if err != nil {
				return nil, err
			}
			moveHis = append(moveHis, PgnMove{
				Move:              m,
				NumericAnnotation: 0,
				PreCommentary:     preCommentary,
				PostCommentary:    []string{},
				Variations:        [][]PgnMove{},
			})
			if !slices.Contains(LegalMoves(pos), m) {
				return nil, fmt.Errorf("pgn contains illegal move: %s", m)
			}
			pos.Move(m)
			preCommentary = []string{}
		case moveNum:
			// no action needed
		case numericAnnotation:
			if len(moveHis) == 0 {
				return nil, fmt.Errorf("can't apply numeric annotation without a move")
			}
			nag, err := strconv.ParseUint(t.body[1:], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("numeric annotation glyph is not between 0 and 255")
			}
			moveHis[len(moveHis)-1].NumericAnnotation = uint8(nag)
		case numericSuffixAnnotation:
			if len(moveHis) == 0 {
				return nil, fmt.Errorf("can't apply numeric suffix annotation without a move")
			}
			nag, err := parseNumericSuffixAnnotation(t.body)
			if err != nil {
				return nil, err
			}
			moveHis[len(moveHis)-1].NumericAnnotation = nag
		case ravOpen:
			closeToken := i + findRavClose(tokens[i:])
			rav, err := createMoveHistory(tokens[i+1:closeToken], prevPos.Copy())
			if err != nil {
				return nil, err
			}
			if len(moveHis) == 0 {
				return nil, fmt.Errorf("can't have variation without a move first")
			}
			moveHis[len(moveHis)-1].Variations = append(moveHis[len(moveHis)-1].Variations, rav)
			i = closeToken
		case ravClose:
			// no action needed
		case result:
			if i != len(tokens)-1 {
				return nil, fmt.Errorf("encountered a result before the end of the pgn")
			}
		default:
			panic(fmt.Sprintf("unexpected chess.pgnTokenType: %#v", t.tokenType))
		}
	}
	return moveHis, nil
}

// findRavClose finds the accompanying closing token index from a slice of tokens given the first token is an open token. Returns len(tokens) if it was not found.
func findRavClose(tokens []pgnToken) int {
	numOpenTokens := 0
	for i, t := range tokens {
		if t.tokenType == ravOpen {
			numOpenTokens++
		}
		if t.tokenType == ravClose {
			if numOpenTokens == 1 {
				return i
			} else {
				numOpenTokens--
			}
		}
	}
	return len(tokens)
}

func parseNumericSuffixAnnotation(nag string) (uint8, error) {
	switch nag {
	case "!":
		return 1, nil
	case "?":
		return 2, nil
	case "!!":
		return 3, nil
	case "??":
		return 4, nil
	case "!?":
		return 5, nil
	case "?!":
		return 6, nil
	default:
		return 0, fmt.Errorf("unrecognized numeric suffix annotation %q", nag)
	}
}

type pgnTokenType uint8

const (
	moveNum pgnTokenType = iota
	move
	commentary
	numericAnnotation
	numericSuffixAnnotation
	ravOpen
	ravClose
	result
)

type pgnToken struct {
	tokenType pgnTokenType
	body      string
}

func tokenizeMovetext(text string) ([]pgnToken, error) {
	// This could perhaps use a nicer tokenization system. But it seems to work well for my current needs.
	tokens := []pgnToken{}
	words := splitWordsPreserveWhitespace(text)
	for i := 0; i < len(words); i++ {
		if len(words[i]) == 0 {
			continue
		}
		if words[i] == "1-0" ||
			words[i] == "0-1" ||
			words[i] == "1/2-1/2" ||
			words[i] == "*" {
			// Result
			tokens = append(tokens, pgnToken{result, words[i]})
		} else if unicode.IsDigit(rune(words[i][0])) {
			// Move number
			nonPeriodIndex := firstNonPeriod(words[i][1:]) + 1
			tokens = append(tokens, pgnToken{moveNum, words[i][0:nonPeriodIndex]})
			words[i] = words[i][nonPeriodIndex:]
			i--
		} else if words[i][0] == '(' {
			// Begin recursive annotation variation
			tokens = append(tokens, pgnToken{ravOpen, "("})
			words[i] = words[i][1:]
			i--
		} else if words[i][0] == ')' {
			// End recursive annotation variation
			tokens = append(tokens, pgnToken{ravClose, ")"})
			words[i] = words[i][1:]
			i--
		} else if words[i][0] == '$' {
			// Numeric annotation
			nonDigitIndex := firstNonDigit(words[i][1:]) + 1
			tokens = append(tokens, pgnToken{numericAnnotation, words[i][0:nonDigitIndex]})
			words[i] = words[i][nonDigitIndex:]
			i--
		} else if words[i][0] == '!' || words[i][0] == '?' {
			// Numeric suffix annotation
			if len(words[i]) > 1 && (words[i][1] == '!' || words[i][1] == '?') {
				tokens = append(tokens, pgnToken{numericSuffixAnnotation, words[i][0:2]})
				words[i] = words[i][2:]
				i--
			} else {
				tokens = append(tokens, pgnToken{numericSuffixAnnotation, words[i][0:1]})
				words[i] = words[i][1:]
				i--
			}
		} else if words[i][0] == ';' {
			// Line comment
			words[i] = words[i][1:]
			comment := ""
			for words[i] != "\n" {
				comment += words[i]
				i++
			}
			tokens = append(tokens, pgnToken{commentary, strings.TrimSpace(comment)})
		} else if words[i][0] == '{' {
			// Curly brace comment
			words[i] = words[i][1:]
			comment := ""
			for !strings.Contains(words[i], "}") {
				comment += words[i]
				i++
				if i >= len(words) {
					return nil, fmt.Errorf("unmatched { in movetext")
				}
			}
			braceIndex := strings.Index(words[i], "}")
			comment += words[i][0:braceIndex]
			tokens = append(tokens, pgnToken{commentary, strings.TrimSpace(comment)})
			words[i] = words[i][braceIndex+1:]
			i--
		} else if words[i][0] == '}' {
			return nil, fmt.Errorf("unmatched } in movetext")
		} else if !unicode.IsSpace([]rune(words[i])[0]) {
			// Move
			endMoveIndex := strings.IndexAny(words[i], "!?{}()$;")
			if endMoveIndex == -1 {
				endMoveIndex = len(words[i])
			}
			movetext := words[i][0:endMoveIndex]
			words[i] = words[i][endMoveIndex:]
			i--
			tokens = append(tokens, pgnToken{move, movetext})
		}
	}
	return tokens, nil
}

func firstNonDigit(s string) int {
	for i, r := range s {
		if !unicode.IsDigit(r) {
			return i
		}
	}
	return len(s)
}

func firstNonPeriod(s string) int {
	for i, r := range s {
		if r != '.' {
			return i
		}
	}
	return len(s)
}

func splitWordsPreserveWhitespace(s string) []string {
	words := []string{}
	for s != "" {
		whitespaceIndex := 0
		for i, r := range s {
			if unicode.IsSpace(r) {
				whitespaceIndex = i
				break
			}
		}
		if whitespaceIndex == 0 {
			words = append(words, s[0:1])
			s = s[1:]
			continue
		}
		words = append(words, s[0:whitespaceIndex])
		words = append(words, s[whitespaceIndex:whitespaceIndex+1])
		s = s[whitespaceIndex+1:]
	}
	return words
}

// ParsePGN reads until the provided reader returns [io.EOF] and provides a list of the games parsed from the PGN.
// If an error is encountered before reaching EOF, all the games that could be parsed will be returned with the error.
// For large pgn files this function may take a few seconds.
// See also [Game.UnmarshalText].
func ParsePGN(rd io.Reader) ([]*Game, error) {
	bufReader := bufio.NewReader(rd)
	games := []*Game{}
	gameParseErrors := []error{}
	pgn, err := extractSinglePgn(bufReader)
	for ; err == nil || (err == io.EOF && len(pgn) > 0); pgn, err = extractSinglePgn(bufReader) {
		newG := &Game{}
		gameErr := newG.UnmarshalText(pgn)
		if gameErr != nil {
			gameParseErrors = append(gameParseErrors, gameErr)
		} else {
			games = append(games, newG)
		}
		if err == io.EOF {
			break
		}
	}
	if err != io.EOF {
		return games, fmt.Errorf("error parsing entire file: %w", err)
	}
	if len(gameParseErrors) > 0 {
		errorString := strings.Builder{}
		errorString.WriteRune('[')
		for _, e := range gameParseErrors {
			errorString.WriteString("\"" + e.Error() + "\",\n")
		}
		errorString.WriteRune(']')
		return games, fmt.Errorf("error parsing games, here is a list of all encountered errors: %s", errorString.String())
	}
	return games, nil
}

func extractSinglePgn(bufrd *bufio.Reader) ([]byte, error) {
	pgn := []byte{}
	// get tag section
	for {
		buf, err := bufrd.ReadBytes('\n')
		if err != nil {
			if len(pgn) == 0 && err == io.EOF {
				return nil, err
			} else {
				return nil, fmt.Errorf("could not parse a complete pgn: %w", err)
			}
		}

		pgn = append(pgn, buf...)
		if isEmptyLine(buf) {
			break
		}
	}

	// get movetext section
	reachedEOF := false
	for {
		buf, err := bufrd.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading pgn file: %w", err)
		}

		pgn = append(pgn, buf...)
		if err == io.EOF {
			reachedEOF = true
			break
		}
		if isEmptyLine(buf) {
			break
		}
	}

	if reachedEOF {
		return pgn, io.EOF
	}
	return pgn, nil
}

func isEmptyLine(buf []byte) bool {
	return (len(buf) == 1 && buf[0] == '\n') ||
		(len(buf) == 2 && buf[0] == '\r' && buf[1] == '\n')
}

// Copy returns a copy of the game.
func (g *Game) Copy() *Game {
	return &Game{
		pos:         g.pos.Copy(),
		moveHistory: g.MoveHistory(),
		moves:       slices.Clone(g.moves),
		Event:       g.Event,
		Site:        g.Site,
		Date:        g.Date,
		Round:       g.Round,
		White:       g.White,
		Black:       g.Black,
		Result:      g.Result,
		OtherTags:   maps.Clone(g.OtherTags),
	}
}

// LegalMoves returns a copy the current legal moves, cached for performance.
func (g *Game) LegalMoves() []Move {
	return slices.Clone(g.legalMoves())
}

// legalMoves returns the struct with the current legal moves. It is not a copy, don't modify it. Cached for performance.
func (g *Game) legalMoves() []Move {
	if g.moves == nil {
		g.moves = LegalMoves(g.pos)
	}
	return g.moves
}

// IsCheckmate returns true if the side to move is in check and there are no legal moves.
func (g *Game) IsCheckmate() bool {
	if len(g.legalMoves()) == 0 && g.pos.IsCheck() {
		return true
	}
	return false
}

// IsStalemate returns true if the side to move is not in check and has no legal moves.
func (g *Game) IsStalemate() bool {
	if len(g.legalMoves()) == 0 && !g.pos.IsCheck() {
		return true
	}
	return false
}

// CanClaimDraw returns true if the side to move can claim a draw either due to the 50 move rule, or three fold repetition.
func (g *Game) CanClaimDraw() bool {
	return true
}

// CanClaimDrawThreeFold returns true if the side to move can claim a draw due to three fold repetition.
func (g *Game) CanClaimDrawThreeFold() bool {
	return true
}

// Move performs the given move m only if it is legal. Otherwise an error is produced.
//
// If Result is set then it will be set to [NoResult]. If the game ends (in stalemate or checkmate) then Result will also be set appropriately.
func (g *Game) Move(m Move) error {
	if !slices.Contains(g.legalMoves(), m) {
		return errors.New("illegal move")
	}
	g.pos.Move(m)
	g.moves = nil
	g.moveHistory = append(g.moveHistory, PgnMove{
		Move:              m,
		NumericAnnotation: 0,
		PreCommentary:     []string{},
		PostCommentary:    []string{},
		Variations:        [][]PgnMove{},
	})
	g.setResult()
	return nil
}

// setResult sets the game result, defaults to NoResult.
func (g *Game) setResult() {
	if g.IsStalemate() {
		g.Result = Draw
	} else if g.IsCheckmate() {
		if g.pos.SideToMove == White {
			g.Result = BlackWins
		} else if g.pos.SideToMove == Black {
			g.Result = WhiteWins
		}
	} else {
		g.Result = NoResult
	}
}

// MoveUCI parses and performs a UCI chess move. https://www.wbec-ridderkerk.nl/html/UCIProtocol.html
//
// Errors are returned if m could not be parsed or the move was illegal.
func (g *Game) MoveUCI(m string) error {
	var move Move
	err := move.UnmarshalText([]byte(m))
	if err != nil {
		return err
	}
	return g.Move(move)
}

// MoveSAN parses and performs a SAN (Standard Algebraic Notation) chess move. https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.2.3
//
// Errors are returned if m could not be parsed or the move was illegal.
func (g *Game) MoveSAN(m string) error {
	move, err := ParseSANMove(m, g.pos)
	if err != nil {
		return err
	}
	return g.Move(move)
}

// Position returns a copy of the current position.
func (g *Game) Position() *Position {
	return g.pos.Copy()
}

// PositionPly returns a copy of the position at a certain ply (half move). 0 returns the initial game position.
//
// If a negative number is provided, or ply goes beyond the number of moves played nil is returned. Ply always starts at 0 and increments by 1, even if the game starts at move 16 for example.
func (g *Game) PositionPly(ply int) *Position {
	pos := &Position{}
	if g.OtherTags["SetUp"] == "1" {

		err := pos.UnmarshalText([]byte(g.OtherTags["FEN"]))
		if err != nil {
			panic("game somehow got invalid FEN starting position, can get position at ply")
		}
	} else {
		pos.UnmarshalText([]byte(DefaultFEN))
	}
	for _, m := range g.moveHistory[:ply] {
		pos.Move(m.Move)
	}
	return pos
}

// MoveHistory returns a copy of all the moves played this game with their annotations, commentary and variations. Will not return nil.
func (g *Game) MoveHistory() []PgnMove {
	moveHistoryCopy := make([]PgnMove, 0, len(g.moveHistory))
	for _, move := range g.moveHistory {
		moveHistoryCopy = append(moveHistoryCopy, move.Copy())
	}
	return moveHistoryCopy
}

// AnnotateMove applies a numeric annotation glyph (NAG) to the specified move number. NAG's can be found here: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c10
//
// plyNum starts at 0 for the first move. Any previous nag is overwritten.
func (g *Game) AnnotateMove(plyNum int, nag uint8) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	g.moveHistory[plyNum].NumericAnnotation = nag
	return nil
}

// CommentAfterMove appends a comment to the specified move. Returns and error if plyNum is out of range.
//
// plyNum starts at 0 for the first move.
func (g *Game) CommentAfterMove(plyNum int, comment string) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	g.moveHistory[plyNum].PostCommentary = append(g.moveHistory[plyNum].PostCommentary, comment)
	return nil
}

// CommentBeforeMove appends appends a comment to be displayed before a move.
// Commentary is not well defined in the pgn specification, thus in most situations it is impossible to tell if a comment should be associated with the move right before it, or right after it. By default comments will be associated with the move right before them, but in some cases (such as the start of a game, or start of a variation) it is possible to have a comment that must precede a move. This is to say that if you marshal and then unmarshal a game, most comments will be parsed as being after a move.
//
// plyNum starts at 0 for the first move.
func (g *Game) CommentBeforeMove(plyNum int, comment string) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	g.moveHistory[plyNum].PreCommentary = append(g.moveHistory[plyNum].PreCommentary, comment)
	return nil
}

// DeleteCommentAfter deletes a comment after the specified move.
//
// plyNum and commentNum start at 0 for the first move.
func (g *Game) DeleteCommentAfter(plyNum int, commentNum int) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	if commentNum < 0 || commentNum >= len(g.moveHistory[plyNum].PostCommentary) {
		return fmt.Errorf("commentNum is too large or to small: len(PostCommentary) = %d, commentNum = %d", len(g.moveHistory[plyNum].PostCommentary), commentNum)
	}
	g.moveHistory[plyNum].PostCommentary = slices.Delete(g.moveHistory[plyNum].PostCommentary, commentNum, commentNum+1)
	return nil
}

func (g *Game) DeleteCommentBefore(plyNum int, commentNum int) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	if commentNum < 0 || commentNum >= len(g.moveHistory[plyNum].PreCommentary) {
		return fmt.Errorf("commentNum is too large or to small: len(PreCommentary) = %d, commentNum = %d", len(g.moveHistory[plyNum].PreCommentary), commentNum)
	}
	g.moveHistory[plyNum].PreCommentary = slices.Delete(g.moveHistory[plyNum].PreCommentary, commentNum, commentNum+1)
	return nil
}

// MakeVariation adds a set of variation moves to the specified move. The variation should begin with a move that replaces the current move. Variation moves must be legal.
//
// plyNum starts at 0 for the first move. moves is not copied, it is simply added to the variations list for the appropriate pgn move.
func (g *Game) MakeVariation(plyNum int, moves []PgnMove) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	pos := g.PositionPly(plyNum)
	if !isLegalVariation(pos, moves) {
		return fmt.Errorf("variation contains illegal moves")
	}

	g.moveHistory[plyNum].Variations = append(g.moveHistory[plyNum].Variations, moves)
	return nil
}

func isLegalVariation(p *Position, variation []PgnMove) bool {
	for _, m := range variation {
		for _, v := range m.Variations {
			if !isLegalVariation(p.Copy(), v) {
				return false
			}
		}
		if !isLegalMove(p, m.Move) {
			return false
		}
		p.Move(m.Move)
	}
	return true
}

func isLegalMove(p *Position, m Move) bool {
	return slices.Contains(LegalMoves(p), m)
}

// DeleteVariation deletes the variation at the plyNum. Give the index you want to delete (starting at 0).
func (g *Game) DeleteVariation(plyNum int, variationNum int) error {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d", len(g.moveHistory), plyNum)
	}
	if variationNum < 0 || variationNum >= len(g.moveHistory[plyNum].Variations) {
		return fmt.Errorf("variationNum is too large or to small: len(Variations) = %d, variationNum = %d", len(g.moveHistory[plyNum].Variations), variationNum)
	}
	g.moveHistory[plyNum].Variations = slices.Delete(g.moveHistory[plyNum].Variations, variationNum, variationNum+1)
	return nil
}

// GetVariation returns a new game where the specified variation is followed. All other specified variations are preserved, and the main line is also preserved as a variation.
func (g *Game) GetVariation(plyNum int, variationNum int) (*Game, error) {
	if plyNum < 0 || plyNum >= len(g.moveHistory) {
		return nil, fmt.Errorf("plyNum is too large or to small: len(moveHistory) = %d, plyNum = %d",
			len(g.moveHistory), plyNum)
	}
	if variationNum < 0 || variationNum >= len(g.moveHistory[plyNum].Variations) {
		return nil, fmt.Errorf("variationNum is too large or to small: len(Variations) = %d, variationNum = %d",
			len(g.moveHistory[plyNum].Variations), variationNum)
	}

	newMoveHistory := g.MoveHistory()[0:plyNum]
	// Make the specified variation the main move history
	for _, m := range g.moveHistory[plyNum].Variations[variationNum] {
		newMoveHistory = append(newMoveHistory, m.Copy())
	}

	// Make the main move history a variation
	variation := []PgnMove{}
	for _, m := range g.moveHistory[plyNum:] {
		variation = append(variation, m.Copy())
	}
	newMoveHistory[plyNum].Variations = append(newMoveHistory[plyNum].Variations, variation)

	// Add other move variations
	for varNum, variation := range g.moveHistory[plyNum].Variations {
		if varNum == variationNum {
			continue
		}
		variationCopy := []PgnMove{}
		for _, m := range variation {
			variationCopy = append(variationCopy, m.Copy())
		}
		newMoveHistory[plyNum].Variations = append(newMoveHistory[plyNum].Variations, variationCopy)
	}

	newGame := g.Copy()
	newGame.pos = newGame.PositionPly(0)
	newGame.moveHistory = []PgnMove{}
	newGame.moves = nil
	for _, m := range newMoveHistory {
		newGame.Move(m.Move)
	}

	newGame.moveHistory = newMoveHistory

	return newGame, nil
}

// MarshalText implements [encoding.TextMarshaler]. It provides the game as a valid PGN that can be written to a file. Multiple PGNs can be written to the same file. Just be sure to separate them with a new line.
//
// The seven tag roster will appear in order, then all other tags will appear in alphabetical order for consistency. err is always nil.
func (g *Game) MarshalText() (text []byte, err error) {
	lines := make([]string, 0, 10)
	err = g.addTags(&lines)
	if err != nil {
		return nil, fmt.Errorf("could not marshal tags: %w", err)
	}
	lines = append(lines, "")
	err = g.addMoveText(&lines)
	if err != nil {
		return nil, fmt.Errorf("could not marshal movetext: %w", err)
	}
	pgn := strings.Builder{}
	for i, l := range lines {
		pgn.WriteString(l)
		if i != len(lines)-1 {
			pgn.WriteString("\n")
		}
	}
	return []byte(pgn.String()), nil
}

func (g *Game) addTags(lines *[]string) error {
	*lines = append(*lines, fmt.Sprintf("[Event %q]", g.Event))
	*lines = append(*lines, fmt.Sprintf("[Site %q]", g.Site))
	*lines = append(*lines, fmt.Sprintf("[Date %q]", g.Date))
	*lines = append(*lines, fmt.Sprintf("[Round %q]", g.Round))
	*lines = append(*lines, fmt.Sprintf("[White %q]", g.White))
	*lines = append(*lines, fmt.Sprintf("[Black %q]", g.Black))
	rstStr, err := g.Result.MarshalText()
	if err != nil {
		return err
	}
	*lines = append(*lines, fmt.Sprintf("[Result %q]", rstStr))
	g.addOtherTags(lines)
	return nil
}

func (g *Game) addReducedTags(lines *[]string) error {
	*lines = append(*lines, fmt.Sprintf("[Event %q]", g.Event))
	*lines = append(*lines, fmt.Sprintf("[Site %q]", g.Site))
	*lines = append(*lines, fmt.Sprintf("[Date %q]", g.Date))
	*lines = append(*lines, fmt.Sprintf("[Round %q]", g.Round))
	*lines = append(*lines, fmt.Sprintf("[White %q]", g.White))
	*lines = append(*lines, fmt.Sprintf("[Black %q]", g.Black))
	rstStr, err := g.Result.MarshalText()
	if err != nil {
		return err
	}
	*lines = append(*lines, fmt.Sprintf("[Result %q]", rstStr))
	if g.OtherTags["SetUp"] == "1" {
		*lines = append(*lines, fmt.Sprintf("[FEN %q]", g.OtherTags["FEN"]))
		*lines = append(*lines, fmt.Sprintf("[SetUp %q]", g.OtherTags["SetUp"]))
	}
	return nil
}

func (g *Game) addOtherTags(lines *[]string) {
	keys := make([]string, 0, len(g.OtherTags))
	for k := range g.OtherTags {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		*lines = append(*lines, fmt.Sprintf("[%s %q]", k, g.OtherTags[k]))
	}
}

func (g *Game) addMoveText(lines *[]string) error {
	currPos := g.PositionPly(0)
	currentLine := strings.Builder{}
	currentLine.Grow(80)

	includeBlackMoveNum := currPos.SideToMove == Black
	for _, m := range g.moveHistory {
		for _, comment := range m.PreCommentary {
			cmtStr := fmt.Sprintf(" {%s}", comment)
			appendToPgnLine(cmtStr, &currentLine, lines)
		}
		if currPos.SideToMove == White {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "."
			appendToPgnLine(moveNum, &currentLine, lines)
		}
		if includeBlackMoveNum {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "..."
			appendToPgnLine(moveNum, &currentLine, lines)
		}

		temp, err := m.Move.StringSAN(currPos)
		if err != nil {
			return err
		}
		sanMove := " " + temp
		nag := nagString(m.NumericAnnotation)
		if m.NumericAnnotation <= 6 {
			sanMove += nag
			appendToPgnLine(sanMove, &currentLine, lines)
		} else {
			appendToPgnLine(sanMove, &currentLine, lines)
			appendToPgnLine(" "+nag, &currentLine, lines)
		}

		for _, comment := range m.PostCommentary {
			cmtStr := fmt.Sprintf(" {%s}", comment)
			appendToPgnLine(cmtStr, &currentLine, lines)
		}

		for _, variation := range m.Variations {
			err := appendVariation(currPos.Copy(), variation, &currentLine, lines)
			if err != nil {
				return fmt.Errorf("could not add variations: %w", err)
			}
		}

		currPos.Move(m.Move)
		if (len(m.PostCommentary) > 0 || len(m.Variations) > 0) && currPos.SideToMove == Black {
			includeBlackMoveNum = true
		} else {
			includeBlackMoveNum = false
		}
	}
	result, err := g.Result.MarshalText()
	if err != nil {
		return err
	}
	appendToPgnLine(" "+string(result), &currentLine, lines)
	*lines = append(*lines, currentLine.String())
	return nil
}

func (g *Game) addReducedMoveText(lines *[]string) error {
	currPos := g.PositionPly(0)
	currentLine := strings.Builder{}
	currentLine.Grow(80)

	includeBlackMoveNum := currPos.SideToMove == Black
	for _, m := range g.moveHistory {
		if currPos.SideToMove == White {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "."
			currentLine.WriteString(moveNum)
		}
		if includeBlackMoveNum {
			moveNum := " " + strconv.FormatUint(uint64(currPos.FullMove), 10) + "..."
			currentLine.WriteString(moveNum)
			includeBlackMoveNum = false
		}
		temp, err := m.Move.StringSAN(currPos)
		if err != nil {
			return err
		}
		sanMove := " " + temp
		currentLine.WriteString(sanMove)
		currPos.Move(m.Move)
	}
	result, err := g.Result.MarshalText()
	if err != nil {
		return err
	}
	currentLine.WriteString(" " + string(result))
	*lines = append(*lines, strings.TrimSpace(currentLine.String()))
	return nil
}

// appendToPgnLine appends string s to currentLine. If currentLine would be longer than 80, then it is appended to lines and currentLine is reset to the value of s.
func appendToPgnLine(s string, currentLine *strings.Builder, lines *[]string) {
	if len(s)+currentLine.Len() > 80 {
		*lines = append(*lines, currentLine.String())
		currentLine.Reset()
	}
	if currentLine.Len() == 0 {
		s = strings.TrimSpace(s)
	}
	currentLine.WriteString(s)
}

func nagString(nag uint8) string {
	switch nag {
	case 0:
		return ""
	case 1:
		return "!"
	case 2:
		return "?"
	case 3:
		return "!!"
	case 4:
		return "??"
	case 5:
		return "!?"
	case 6:
		return "?!"
	default:
		return "$" + strconv.FormatUint(uint64(nag), 10)
	}
}

func appendVariation(currPos *Position, moves []PgnMove, currentLine *strings.Builder, lines *[]string) error {
	includeBlackMoveNum := currPos.SideToMove == Black
	for i, m := range moves {
		if currPos.SideToMove == White {
			moveNum := strconv.FormatUint(uint64(currPos.FullMove), 10) + "."
			if i == 0 {
				moveNum = "(" + moveNum
			}
			moveNum = " " + moveNum
			appendToPgnLine(moveNum, currentLine, lines)
		}
		if includeBlackMoveNum {
			moveNum := strconv.FormatUint(uint64(currPos.FullMove), 10) + "..."
			if i == 0 {
				moveNum = "(" + moveNum
			}
			moveNum = " " + moveNum
			appendToPgnLine(moveNum, currentLine, lines)
		}
		temp, err := m.Move.StringSAN(currPos)
		if err != nil {
			return err
		}
		sanMove := " " + temp
		nag := nagString(m.NumericAnnotation)
		if m.NumericAnnotation <= 6 {
			sanMove += nag
			if i == len(moves)-1 && len(m.PostCommentary) == 0 && len(m.Variations) == 0 {
				sanMove += ")"
			}
			appendToPgnLine(sanMove, currentLine, lines)
		} else {
			appendToPgnLine(sanMove, currentLine, lines)
			if i == len(moves)-1 && len(m.PostCommentary) == 0 && len(m.Variations) == 0 {
				nag += ")"
			}
			appendToPgnLine(" "+nag, currentLine, lines)
		}

		for j, comment := range m.PostCommentary {
			cmtStr := fmt.Sprintf(" {%s}", comment)
			if i == len(moves)-1 && j == len(m.PostCommentary)-1 && len(m.Variations) == 0 {
				cmtStr += ")"
			}
			appendToPgnLine(cmtStr, currentLine, lines)
		}

		for j, variation := range m.Variations {
			err := appendVariation(currPos.Copy(), variation, currentLine, lines)
			if err != nil {
				return err
			}
			if i == len(moves)-1 && j == len(m.Variations)-1 {
				appendToPgnLine(")", currentLine, lines)
			}
		}

		currPos.Move(m.Move)
		if (m.NumericAnnotation > 6 || len(m.PostCommentary) > 0 || len(m.Variations) > 0) && currPos.SideToMove == Black {
			includeBlackMoveNum = true
		} else {
			includeBlackMoveNum = false
		}
	}
	return nil
}

// MarshalTextReduced provides the game as a valid PGN following these rules: https://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c3.2.4
//
// It essentially removes all unnecessary information from the PGN, useful for archival purposes. err is always nil.
func (g *Game) MarshalTextReduced() (text []byte, err error) {
	lines := make([]string, 0, 10)
	err = g.addReducedTags(&lines)
	if err != nil {
		return nil, fmt.Errorf("could not marshal tags: %w", err)
	}
	lines = append(lines, "")
	err = g.addReducedMoveText(&lines)
	if err != nil {
		return nil, fmt.Errorf("could not marshal movetext: %w", err)
	}
	pgn := strings.Builder{}
	for i, l := range lines {
		pgn.WriteString(l)
		if i != len(lines)-1 {
			pgn.WriteString("\n")
		}
	}
	return []byte(pgn.String()), nil
}
