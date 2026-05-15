// Copyright (C) 2026 Brigham Skarda

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

package uci

import (
	"time"

	"github.com/brighamskarda/chess/v2"
)

const mockEngineEvaluateTime = time.Millisecond * 500

// mockEngine simply keeps track of how many times each of its functions are called.
//
// Used for testing.
type mockEngine struct {
	output      func(*InfoCmd)
	debugState  bool
	position    *chess.Position
	moveHistory []chess.Move
	shouldStop  chan struct{}
	stopPonder  chan struct{}

	initialize     int
	copyProtection int
	register       int
	name           int
	author         int
	options        int
	debug          int
	setOption      int
	newGame        int
	setPosition    int
	evaluate       int
	stop           int
	ponderHit      int
	quit           int
}

func (engine *mockEngine) Initialize(o func(*InfoCmd)) {
	engine.output = o
	engine.shouldStop = make(chan struct{})
	engine.stopPonder = make(chan struct{})
	engine.initialize++
}

func (engine *mockEngine) CopyProtection() bool {
	engine.copyProtection++
	return true
}

func (engine *mockEngine) Register(cmd *RegisterCmd) bool {
	engine.register++
	return true
}

func (engine *mockEngine) Name() string {
	engine.name++
	return "mockEngine v0.1"
}

func (engine *mockEngine) Author() string {
	engine.author++
	return "Brigham Skarda"
}

func (engine *mockEngine) Options() []Option {
	engine.options++
	return []Option{
		&CheckOption{
			Name:         "checkOpt",
			DefaultValue: true,
		},
		&SpinOption{
			Name:         "spinOpt",
			DefaultValue: 3,
			Min:          1,
			Max:          5,
		},
		&ComboOption{
			Name:         "comboOpt",
			DefaultValue: "one",
			Variants:     []string{"one", "two", "three"},
		},
		&StringOption{
			Name:         "stringOpt",
			DefaultValue: "sss",
		},
		&ButtonOption{
			Name: "buttonOpt",
		},
	}
}

func (engine *mockEngine) SetDebug(value bool) {
	engine.debugState = value
	engine.debug++
}

func (engine *mockEngine) SetOption(option SetOption) {
	engine.setOption++
}

func (engine *mockEngine) NewGame() {
	engine.newGame++
}

func (engine *mockEngine) SetPosition(pos *chess.Position, his []chess.Move) {
	engine.position = pos
	engine.moveHistory = his
	engine.setPosition++
}

func (engine *mockEngine) Evaluate(cmd *EvaluateCmd) BestMove {
	engine.evaluate++
	timer := time.After(mockEngineEvaluateTime)

	if cmd.Ponder {
		select {
		case <-engine.shouldStop:
			return BestMove{
				Move: chess.Move{
					FromSquare: chess.E2,
					ToSquare:   chess.E4,
				},
				PonderMove: OptionalOf(chess.Move{
					FromSquare: chess.D7,
					ToSquare:   chess.D5,
				},
				),
			}
		case <-engine.stopPonder:
		}
	}

	select {
	case <-engine.shouldStop:
	case <-timer:
	}

	return BestMove{
		Move: chess.Move{
			FromSquare: chess.E2,
			ToSquare:   chess.E4,
		},
		PonderMove: OptionalOf(chess.Move{
			FromSquare: chess.D7,
			ToSquare:   chess.D5,
		},
		),
	}
}

func (engine *mockEngine) Stop() {
	engine.stop++
	engine.shouldStop <- struct{}{}
}

func (engine *mockEngine) PonderHit() {
	engine.ponderHit++
	engine.stopPonder <- struct{}{}
}

func (engine *mockEngine) Quit() {
	engine.quit++
}
