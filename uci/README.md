# brighamskarda/chess/uci

**uci** is a one of a kind go package that makes chess engine development easier than ever before!

It solves a major hurdle many new (and experienced) developers face when making a chess engine.
And that is making it UCI compatible.

The [Universal Chess Interface](https://en.wikipedia.org/wiki/Universal_Chess_Interface) 
is the most widespread communication protocol for chess engines.
If you want your chess engine to be compatible with other pieces of chess software,
you need to implement UCI.

And this package makes it easier than ever.
All you need to do is make a struct that satisfies the 
[ChessEngine interface](https://pkg.go.dev/github.com/brighamskarda/chess/v2/uci#ChessEngine).
Once you've done that pass it into a 
[UciEngineBroker](https://pkg.go.dev/github.com/brighamskarda/chess/v2/uci#UciEngineBroker),
call start, and you are done.
It is really that easy.

Now don't get too worried about about the size and rather lengthy documentation
of the ChessEngine interface.
Most of it is optional.
To prove it to you,
here is all the code you need to write for a UCI compatible chess engine.

```go
package main

import (
	"context"
	"os"

	"github.com/brighamskarda/chess/v2"
	"github.com/brighamskarda/chess/v2/uci"
)

// MinimalEngine will implement the simplest chess engine possible.
type MinimalEngine struct {
	// We will store the position to evaluate here.
	position *chess.Position
}

func (engine *MinimalEngine) Initialize(ignore func(*uci.InfoCmd)) {
	// We don't need to initialize anything
	// and for simple engines it is okay to ignore the info command function.
}

func (engine *MinimalEngine) CopyProtection() bool {
	// Most engines won't use this, just return true.
	return true
}

func (engine *MinimalEngine) Register(ignore *uci.RegisterCmd) bool {
	// Most engines won't use this either, just return true.
	return true
}

func (engine *MinimalEngine) Name() string {
	// Give your engine a cool name.
	return "GALACTIC CRUSHER!!!"
}

func (engine *MinimalEngine) Author() string {
	// Don't forget to take credit for your hard work.
	return "John Smith"
}

func (engine *MinimalEngine) Options() []uci.Option {
	// No options are required by the UCI specification.
	return nil
}

func (engine *MinimalEngine) SetDebug(ignore bool) {
	// Out engine doesn't have a debug mode.
}

func (engine *MinimalEngine) SetOption(ignore uci.SetOption) {
	// Our engine doesn't support any options
}

func (engine *MinimalEngine) NewGame() {
	// New game isn't required for simple engines.
}

func (engine *MinimalEngine) SetPosition(pos *chess.Position, moves []chess.Move) {
	// We need to save the position we are being sent.
	engine.position = pos

	// We then need to perform the moves we were given on that position.
	for _, m := range moves {
		engine.position.Move(m)
	}

	// Some engines will store the move history
	// to see if they are entering a stalemate situation,
	// but this engine doesn't care.
}

func (engine *MinimalEngine) Evaluate(ignore *uci.EvaluateCmd) *uci.BestMove {
	// This is the brains of our engine.
	// For a super simple engine like this we don't need to worry about the options given in EvaluateCmd.

	// We will just return the first valid Move
	thisIsTheBest := chess.LegalMoves(engine.position)[0]
	return &uci.BestMove{
		Move: thisIsTheBest,
	}
}

func (engine *MinimalEngine) Stop() {
	// Our evaluation logic is so simple we never need to worry about ending it early.
}

func (engine *MinimalEngine) PonderHit() {
	// Our engine doesn't support pondering.
}

func (engine *MinimalEngine) Quit() {
	// We don't have any open files or long running searches we need to stop.
}

func main() {
	// We make our engine.
	myEngine := &MinimalEngine{}

	// We pass it into the broker.
	broker := uci.UciEngineBroker{
		Engine: myEngine,
		Input:  os.Stdin,
		Output: os.Stdout,
	}

	// broker.Start handles the rest.
	// context.Background(), is the default context for the go program.
	broker.Start(context.Background())

	// Voilà!!! You are now running a compliant UCI chess engine. Go ahead and try it.
	// I ran this engine in the Arena Chess, Lucas Chess, and En Croissant GUIs
	// and it worked flawlessly in all of them.
}
```

I've tested engines developed with this library in multiple chess GUIs 
including Arena Chess, Lucas Chess, and En Croissant.
By following the guidelines below 
it is quite easy to develop a chess engine that
passes the UCI compliance checks provided by the
[fastchess](https://github.com/Disservin/fastchess/tree/master) developers.

I hope you enjoy developing chess engines with this library.
Happy Coding!

## Guidelines For Creating Chess Engines

While the example above shows just how easy it is to make a UCI compliant chess engine,
here are a few other things you will probably want to keep in mind.

### Info Commands

Use the info command function passed into Initialize.

If you've ever used chess software to analyze a position,
the scores and lines of play that you see are made possible by the info command.

Just be sure to use it with a buffered channel
so you don't slow down your engine.

**At a minimum it is recommended to at least send a final score info when finishing an evaluation.**

### Copy Protection and Registration

Unless you are creating commercial software
(please see [LICENSE](../LICENSE)), 
you really don't need these.

### Options

While the UCI specification doesn't have any required options,
here are some good ones to implement.

#### Threads

Threads is not part of the UCI specification, 
but if you are making a multithreaded engine
use this option to specify how many threads it can use.

Default to 1 thread.

#### Hash

If your engine uses a hash table, it should also implement this option
so UCI clients can set the size of the table.
(Critical for ensuring equal competition with other chess engines.)

#### MultiPV

This option is super useful for UCI clients
as it allows your engine to output more than one of the lines it found.

#### Ponder

This tells the UCI client if your engine is able to ponder 
(think on the opponents time).

### Quitting

While the Quit function can be empty most of the time without adverse side effects,
it is good practice to ensure that the all go routines are finished
and open files are closed.

Try not to send anything to the info command callback function after Quit is finished either.

Using [contexts](https://pkg.go.dev/context), 
[wait groups](https://pkg.go.dev/sync#WaitGroup), 
and/or [atomic values](https://pkg.go.dev/sync/atomic#Bool)
to help you ensure that all tasks are finished.

### Example

For more examples of how to write chess engines using this library see https://github.com/brighamskarda/Chess-Engines-In-Go

## Contributing

I'm open to suggestions and contributions. 
Feel free to post them on github issues and pull requests. 
You can also email me directly at [brighamskarda@gmail.com](mailto:brighamskarda@gmail.com). 
Also take a look at [CONTRIBUTING.md](../CONTRIBUTING.md).