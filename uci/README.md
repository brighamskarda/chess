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
// Simple example here.
```

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

Here is a bigger example that uses a few more features of the interface.

```go
// Insert bigger example here.
```

## Contributing

I'm open to suggestions and contributions. 
Feel free to post them on github issues and pull requests. 
You can also email me directly at [brighamskarda@gmail.com](mailto:brighamskarda@gmail.com). 
Also take a look at [CONTRIBUTING.md](../CONTRIBUTING.md).