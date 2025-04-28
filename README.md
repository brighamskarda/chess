# brighamskarda/chess

I used this project to learn GoLang. This library supports all of the following features:

- FEN parsing and generation
- PGN parsing and generation
- Legal move generation
- Pseudo-legal move generation
- UCI move parsing and generation
- SAN move parsing and generation
- Checks for Checkmate and Stalemate (including three-fold repetition)

All functionality has been thoroughly tested (including parsing and writing over 20,000 PGNs which hits nearly every part of the code base including move generation).

## Usage

I recommend taking at look at the docs at <https://pkg.go.dev/> as all non-obvious functions should be documented there.

If you are looking to create a chess application I recommend using the Game struct as it keeps track of move history and always ensures that the game is in a valid state.

For engine development I recommend using the Position struct as it allows for quick and easy access and modification of the game state. Using Move on a position does not check for move validity which increases engine speed.

## Future Development

I am currently working on version 2.0 of this library. It features a fully bitboard based position representation.
Currently moves generationa for version 2 is about twice as fast as the most popular chess library in go. https://github.com/corentings/chess

## Contact info

I'm very happy to consider changes to the API to make things *feel better*. Feel free to email me, or post an issue. I'm even open to pull requests should you feel the urge to contribute.
