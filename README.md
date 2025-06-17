# brighamskarda/chess

**chess** is a go module with useful utilities for playing and manipulating the game of chess. It was created to expand the selection of chess libraries available in golang. As of now, one of the only fleshed out libraries is [CorentinGS/chess](https://github.com/corentings/chess) (who has done some admittedly great work along with [notnil](https://github.com/notnil/chess)). Some of the functionality provided in this library includes:

- Pseudo-legal move generation
- Legal move generation
- FEN position parsing
- Extensive PGN support
- Bitboard utilities for move generation

## Performance

This module is designed with performance in mind. It aims to be performant enough for engine development, while still being easy to use.

It performs similarly to [CorentinGS/chess](https://github.com/corentings/chess), with one notable exception. Legal move generation in this module is up to **40%** faster (sometimes even more).

| Postition | brighamskarda | CorentinGS |
| --------- | ------------- | ---------- |
| Starting  | 2820ns        | 4448ns     |
| Midgame   | 6026ns        | 9987ns     |
| Endgame   | 2167ns        | 4468ns     |

Of course, benchmarking is a fickle art, and a deeper review will show that CorentinGS has done a great job of minimizing heap allocations in his library.

## Future Development

There is still a lot of work to do to make this library as great as possible. Currently there are two future releases planned.

- v2.1 - UCI support to aid engine development
- v2.2 - Chess 960 support

## Contributing

I'm open to suggestions and contributions. Feel free to post them on github issues and pull requests. You can also email me directly at <brighamskarda@gmail.com>.

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
