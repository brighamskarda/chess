package chess_test

import (
	"fmt"
	"slices"

	"github.com/brighamskarda/chess/v2"
)

func ExampleBitboard() {
	var myBitboard chess.Bitboard = 0
	myBitboard = myBitboard.SetBit(0)
	myBitboard = myBitboard.SetSquare(chess.G1)
	fmt.Printf("myBitboard:\n%v\n\n", myBitboard)

	// Say this bitboard represents some bishops. Lets see the squares they attack.
	var occupied chess.Bitboard = myBitboard // For this example the bishops are the only pieces on the board.
	var bishopAttacks chess.Bitboard = myBitboard.BishopAttacks(occupied)
	fmt.Printf("bishopAttacks:\n%v", bishopAttacks)
	// Output:
	// myBitboard:
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 10000010
	//
	// bishopAttacks:
	// 00000001
	// 10000010
	// 01000100
	// 00101000
	// 00010000
	// 00101000
	// 01000101
	// 00000000
}

func ExampleBitboard_String() {
	var myBitboard chess.Bitboard = 0
	myBitboard = myBitboard.SetSquare(chess.C2)
	fmt.Printf("myBitboard:\n%s", myBitboard.String())
	// Output:
	// myBitboard:
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00000000
	// 00100000
	// 00000000
}

func ExampleResult() {
	// Results are usually part of a game
	g := chess.NewGame()
	gameResult := g.Result

	// You can convert results to a textual representation. A * means no result (game still going).
	str := gameResult.String()
	fmt.Printf("gameResult: %v", str)
	// Output:
	// gameResult: *
}

func ExamplePgnMove() {
	// PgnMove is usually retrieved as a list of moves from a game.
	game := &chess.Game{}
	err := game.UnmarshalText([]byte(`[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

{PreComments are useful for commentary on the game as a whole.} 1. e4! Nc6 $2
{I personally don't like this move.} {Others may disagree though} (1... c5 
{Many people like the Sicilian Defense more.}) *`))
	if err != nil {
		panic(fmt.Sprintf("issue unmarshaling example game: %s", err))
	}

	// Here are our moves!
	myPgnMoves := game.MoveHistory()

	fmt.Printf("PreComments Example: %s\n", myPgnMoves[0].PreCommentary[0])
	fmt.Printf("PostComment One: %s\n", myPgnMoves[1].PostCommentary[0])
	fmt.Printf("PostComment Two: %s\n", myPgnMoves[1].PostCommentary[1])
	// The ! after white's move is an inline numeric annotation that translates to 1. This indicates a good move.
	fmt.Printf("Inline Numeric Annotation: %d\n", myPgnMoves[0].NumericAnnotation)
	// The $2 after black's move is a more traditional numeric annotation that means this was a bad move.
	fmt.Printf("Traditional Numeric Annotation: %d\n", myPgnMoves[1].NumericAnnotation)
	// This code gets the first move, of the first variation for black's move.
	// This variation suggests c5 instead of Nc6.
	fmt.Printf("Variation Move: %v\n", myPgnMoves[1].Variations[0][0].Move)
	fmt.Printf("Even Variations Support Comments: %s", myPgnMoves[1].Variations[0][0].PostCommentary[0])
	// Output:
	// PreComments Example: PreComments are useful for commentary on the game as a whole.
	// PostComment One: I personally don't like this move.
	// PostComment Two: Others may disagree though
	// Inline Numeric Annotation: 1
	// Traditional Numeric Annotation: 2
	// Variation Move: c7c5
	// Even Variations Support Comments: Many people like the Sicilian Defense more.
}

func ExampleGame() {
	game := chess.NewGame()

	game.White = "Gopher 1"
	game.Black = "Gopher 2"
	game.Date = "2000.01.01"

	legalMoves := game.LegalMoves()
	myMove := chess.Move{chess.E2, chess.E4, chess.NoPieceType}

	if slices.Contains(legalMoves, myMove) {
		game.Move(myMove)
	} else {
		panic("illegal move")
	}

	if game.IsStalemate() {
		println("Stalemate")
		return
	}

	game.Result = chess.BlackWins
	pgn, _ := game.MarshalText()
	fmt.Printf("%s", pgn)
	// Output:
	// [Event "?"]
	// [Site "https://github.com/brighamskarda/chess"]
	// [Date "2000.01.01"]
	// [Round "1"]
	// [White "Gopher 1"]
	// [Black "Gopher 2"]
	// [Result "0-1"]
	//
	// 1. e4 0-1
}

func ExampleGame_UnmarshalText() {
	game := &chess.Game{}

	examplePgn := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2000.01.01"]
[Round "1"]
[White "Gopher 1"]
[Black "Gopher 2"]
[Result "0-1"]

1. e4 0-1`

	err := game.UnmarshalText([]byte(examplePgn))
	if err != nil {
		panic("unmarshaling failed")
	}

	fmt.Printf("White: %v\n", game.White)
	fmt.Printf("Black: %v\n", game.Black)
	moveHist := game.MoveHistory()
	fmt.Printf("First Move: %v", moveHist[0].Move)
	// Output:
	// White: Gopher 1
	// Black: Gopher 2
	// First Move: e2e4
}

func ExampleGame_AnnotateMove() {
	game := chess.NewGame()
	game.Move(chess.Move{chess.E2, chess.E4, chess.NoPieceType})

	// I want to indicate that I thought this was a great first move.
	err := game.AnnotateMove(0, 3)
	if err != nil {
		fmt.Println("got error")
	}

	// To see my annotation I can look at the move history.
	fmt.Printf("%v", game.MoveHistory()[0].NumericAnnotation)
	// Output:
	// 3
}

func ExampleGame_CommentAfterMove() {
	game := chess.NewGame()
	game.Move(chess.Move{chess.E2, chess.E4, chess.NoPieceType})

	// I can make multiple comments after a move.
	game.CommentAfterMove(0, "comment 1")
	game.CommentAfterMove(0, "comment 2")

	// To see my comments I can look at the move history.
	moveHist := game.MoveHistory()
	fmt.Println(moveHist[0].PostCommentary[0])
	fmt.Println(moveHist[0].PostCommentary[1])

	// I can also delete comments.
	game.DeleteCommentAfter(0, 1)
	// Output:
	// comment 1
	// comment 2
}

func ExampleGame_CommentBeforeMove() {
	game := chess.NewGame()
	game.Move(chess.Move{chess.E2, chess.E4, chess.NoPieceType})

	// I can make multiple comments before a move.
	game.CommentBeforeMove(0, "comment 1")
	game.CommentBeforeMove(0, "comment 2")

	// To see my comments I can look at the move history.
	moveHist := game.MoveHistory()
	fmt.Println(moveHist[0].PreCommentary[0])
	fmt.Println(moveHist[0].PreCommentary[1])

	// I can also delete comments.
	game.DeleteCommentBefore(0, 1)
	// Output:
	// comment 1
	// comment 2
}

func ExampleGame_MakeVariation() {
	game := chess.NewGame()
	game.Move(chess.Move{chess.E2, chess.E4, chess.NoPieceType})

	// I want to see what might happen if I played a different move first.
	myVariationSequence := []chess.PgnMove{{
		Move:              chess.Move{chess.D2, chess.D4, chess.NoPieceType},
		NumericAnnotation: 0,
		PreCommentary:     []string{"I can include comments in variations too."},
		PostCommentary:    []string{},
		Variations:        [][]chess.PgnMove{},
	}}

	err := game.MakeVariation(0, myVariationSequence)
	if err != nil {
		panic("i probably got an error because myVariationSequence was invalid")
	}

	// There are two ways I can see the variation I made.
	// 1. I can get it in the move history
	moveHist := game.MoveHistory()
	fmt.Println(moveHist[0].Variations[0][0].Move)

	// 2. I can also make a new game that follows my variation instead of the main line.
	myNewGame, _ := game.GetVariation(0, 0)
	fmt.Println(myNewGame.MoveHistory()[0].Move)
	// Output:
	// d2d4
	// d2d4
}

func ExampleGame_MarshalText() {
	// Here I parse (unmarshal) a pgn to get a game
	game := &chess.Game{}
	examplePgn := `[Event "?"]
[Site "https://github.com/brighamskarda/chess"]
[Date "2000.01.01"]
[Round "1"]
[White "Gopher 1"]
[Black "Gopher 2"]
[Result "0-1"]

1. e4 0-1`
	err := game.UnmarshalText([]byte(examplePgn))
	if err != nil {
		panic("unmarshaling failed")
	}

	// Now when I marshal the game I should get a functionally identical pgn string.
	regeneratedPgn, err := game.MarshalText()
	if err != nil {
		panic("there was an issue marshaling the game")
	}

	fmt.Println(string(regeneratedPgn) == examplePgn)
	// Output:
	// true
}

func ExamplePosition() {
	pos := &chess.Position{}
	// Initialize the default position
	pos.UnmarshalText([]byte(chess.DefaultFEN))
	fmt.Println(pos.String(true, false))
	// Output:
	// 8rnbqkbnr
	// 7pppppppp
	// 6--------
	// 5--------
	// 4--------
	// 3--------
	// 2PPPPPPPP
	// 1RNBQKBNR
	//  ABCDEFGH
}
