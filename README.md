# brighamskarda/chess

**chess** is a go module with useful utilities for playing and manipulating the game of chess. It was created to expand the selection of chess libraries available in golang. As of now, there are few fleshed out libraries that are documented and performant. Some of the functionality provided in this library includes:

- Pseudo-legal move generation
- Legal move generation
- FEN position parsing
- Extensive PGN support
- Bitboard utilities for move generation

## Performance

This module is designed with performance in mind. It aims to be performant enough for engine development, while still being easy to use. Performance testing (Perft) is one way to test the performance and correctness of a chess library. It involves generating all legal moves for a position up to a certain depth. Perft(6) indicates the time for a library to generate all legal moves 6 plys deep.

### Popular Go Chess Library Performance Comparison

| Repository            | Starting Position (Perft 6) | KiwiPete Position (Perft 5) | End Position (Perft 6) |
| --------------------- | --------------------------- | --------------------------- | ---------------------- |
| brighamskarda/chess   | 7.58                        | 12.04                       | 0.69                   |
| CorentinGS/chess      | 32.49                       | 52.94                       | 4.95                   |
| dylhunn/dragontoothmg | 1.44                        | 1.46                        | 0.22                   |
| keelus/chess          | 435                         | \*                          | 37.86                  |
| malbrecht/chess       | 77.39                       | 147.97                      | 5.44                   |
| eightsquared/chess    | 7.075                       | 22.35                       | 1.01                   |

Starting Position - `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`

Kiwi Pete Position - `r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1` (a popular midgame position for testing)

End Position - `8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1`

Of course there are other hyper-specialized go programs and libraries that are able to achieve even faster times (such as [this one](https://github.com/bluescreen10/chester) which claims 300ms for perft 6 from the starting position). For this list I am looking at libraries that show up in the first couple pages of search results, and that provide a usable API with at least some documentation. If you have a library you think should be here let me know.

But as the results show this library is one of the quickest. Notably it is about 4x as fast as CorentingGS/chess (a fork of notnil/chess that focuses on reducing memory allocations). dylhunn/dragontoothmg is an amazing library for pure speed if that is what you are looking for. It lacks PGN utilities though.

## Future Development

There is still a lot of work to do to make this library as great as possible. Currently there are two future releases planned.

- v2.1 - UCI support to aid engine development
- v2.2 - Chess 960 support

## Contributing

I'm open to suggestions and contributions. Feel free to post them on github issues and pull requests. You can also email me directly at <mailto:brighamskarda@gmail.com>.

## Usage

You can find the full documentation at <https://pkg.go.dev/github.com/brighamskarda/chess/v2>. Here are a few useful tips though:

- Use the **Game** type to play, and keep track of your games. This struct is how you can read and manipulate PGN's.
- Use the **Position** type for high performance applications. Using Position.Move is much faster than Game.Move since it doesn't validate the move, or keep a history.
- If your developing an engine, use **Position.Bitboard** to get bitboards for pieces and colors. Bit operations can accelerate your engine.

### Example Game Against Dumb Computer

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brighamskarda/chess/v2"
)

func main() {
	// Initialize a new game of chess
	myGame := chess.NewGame()
	// Get input reader
	input := bufio.NewReader(os.Stdin)

	// Game Loop
	for !myGame.IsStalemate() && !myGame.IsCheckmate() {
		currentPosition := myGame.Position()
		if currentPosition.SideToMove == chess.White {
			// Player plays white
			fmt.Println(currentPosition.String(true, false))
			fmt.Print("Enter Move (<square1><square2><promotion>): ")
			playerInput, _ := input.ReadString('\n')
			playerInput = strings.TrimSpace(playerInput)
			// MoveUCI automatically parses the players move
			err := myGame.MoveUCI(playerInput)
			if err != nil {
				fmt.Println("Invalid Move")
			}
		} else if currentPosition.SideToMove == chess.Black {
			// Computer plays the first legal move it sees
			legalMoves := myGame.LegalMoves()
			myGame.Move(legalMoves[0])
			fmt.Printf("\nBlack performed move %v\n\n", legalMoves[0])
		}
	}

	// Game detected stalemate or checkmate, so we print the result.
	switch myGame.Result {
	case chess.Draw:
		fmt.Println("Draw")
	case chess.WhiteWins:
		fmt.Println("White Wins")
	case chess.BlackWins:
		fmt.Println("Black Wins")
	default:
		panic("game ended without result.")
	}
}
```
