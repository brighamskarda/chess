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
	"testing"

	"github.com/brighamskarda/chess/v2"
)

func TestIdCommandMarshal(t *testing.T) {
	var cmd engineToClientCmd = &idCmd{
		isAuthor: false,
		id:       "Brigham Skarda",
	}

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "id name Brigham Skarda\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	cmd = &idCmd{
		isAuthor: true,
		id:       "Brigham Skarda",
	}

	text, err = cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "id author Brigham Skarda\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestUciokCommandMarshal(t *testing.T) {
	var cmd engineToClientCmd = &uciokCmd{}

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "uciok\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestReadyokCommandMarshal(t *testing.T) {
	var cmd engineToClientCmd = &readyokCmd{}

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "readyok\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestBestMoveCommandMarshal(t *testing.T) {
	var cmd engineToClientCmd = &bestMoveCmd{
		move:       chess.Move{FromSquare: chess.A1, ToSquare: chess.A2},
		ponderMove: &chess.Move{FromSquare: chess.A1, ToSquare: chess.A2, Promotion: chess.Knight},
	}

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "bestmove a1a2 ponder a1a2n\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	cmd = &bestMoveCmd{
		move: chess.Move{FromSquare: chess.A1, ToSquare: chess.A2},
	}

	text, err = cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "bestmove a1a2\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestCopyprotectionCommandMarshal(t *testing.T) {
	copyprotectCmd := copyprotectChecking
	var cmd engineToClientCmd = &copyprotectCmd

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "copyprotection checking\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	copyprotectCmd = copyprotectOk
	text, err = copyprotectCmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "copyprotection ok\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	copyprotectCmd = copyprotectError
	text, err = copyprotectCmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "copyprotection error\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestRegistrationCommandMarshal(t *testing.T) {
	registrationCmd := registerChecking
	var cmd engineToClientCmd = &registrationCmd

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "registration checking\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	registrationCmd = registerOk
	text, err = registrationCmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "registration ok\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	registrationCmd = registerError
	text, err = registrationCmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "registration error\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}
}

func TestInfoCommandMarshal_AllFields(t *testing.T) {
	cmd := &infoCmd{
		depth:    OptionalOf(20),
		seldepth: OptionalOf(25),
		time:     OptionalOf(1200),
		nodes:    OptionalOf(5000000),
		pv:       []chess.Move{{FromSquare: chess.E2, ToSquare: chess.E4}, {FromSquare: chess.E7, ToSquare: chess.E5}},
		multipv:  OptionalOf(1),
		score: OptionalOf(infoScore{score: 35,
			isMate:       true,
			isLowerbound: true}),
		currmove:       OptionalOf(chess.Move{FromSquare: chess.G1, ToSquare: chess.F3}),
		currmovenumber: OptionalOf(1),
		hashfull:       OptionalOf(500),
		nps:            OptionalOf(450000),
		tbhits:         OptionalOf(100),
		sbhits:         OptionalOf(50),
		cpuload:        OptionalOf(200),
		refutation:     []chess.Move{{FromSquare: chess.D1, ToSquare: chess.H5}, {FromSquare: chess.G6, ToSquare: chess.H5}},
		currline: OptionalOf(currentLine{
			cpunr: OptionalOf(1),
			moves: []chess.Move{{FromSquare: chess.A2, ToSquare: chess.A4}},
		}),
		stringMsg: OptionalOf("engine is thinking deeply"),
	}

	text, err := cmd.marshalText()
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Note: Field order in UCI is flexible, but "string" MUST be last.
	// This expectation assumes a standard logical order.
	expected := "info depth 20 seldepth 25 time 1200 nodes 5000000 pv e2e4 e7e5 multipv 1 " +
		"score mate 35 lowerbound currmove g1f3 currmovenumber 1 hashfull 500 " +
		"nps 450000 tbhits 100 sbhits 50 cpuload 200 refutation d1h5 g6h5 " +
		"currline 1 a2a4 string engine is thinking deeply\n"

	if string(text) != expected {
		t.Errorf("\ngot:      %q\nexpected: %q", string(text), expected)
	}
}

func TestInfoCommandMarshal(t *testing.T) {
	tests := []struct {
		name     string
		cmd      infoCmd
		expected string
	}{
		{
			name:     "Empty Info",
			cmd:      infoCmd{},
			expected: "info\n",
		},
		{
			name: "Basic Search Stats",
			cmd: infoCmd{
				depth: OptionalOf(12),
				nodes: OptionalOf(123456),
				nps:   OptionalOf(100000),
			},
			expected: "info depth 12 nodes 123456 nps 100000\n",
		},
		{
			name: "Score Centipawns",
			cmd: infoCmd{
				score: OptionalOf(infoScore{
					score: -45,
				}),
			},
			expected: "info score cp -45\n",
		},
		{
			name: "Score Mate and Upperbound",
			cmd: infoCmd{
				score: OptionalOf(infoScore{
					score:        5,
					isMate:       true,
					isUpperbound: true,
				}),
			},
			expected: "info score mate 5 upperbound\n",
		},
		{
			name: "PV with Moves",
			cmd: infoCmd{
				depth: OptionalOf(2),
				pv: []chess.Move{
					{FromSquare: chess.E2, ToSquare: chess.E4},
					{FromSquare: chess.E7, ToSquare: chess.E5},
				},
			},
			expected: "info depth 2 pv e2e4 e7e5\n",
		},
		{
			name: "Current Move and Progress",
			cmd: infoCmd{
				currmove:       OptionalOf(chess.Move{FromSquare: chess.G1, ToSquare: chess.F3}),
				currmovenumber: OptionalOf(1),
				hashfull:       OptionalOf(400),
			},
			expected: "info currmove g1f3 currmovenumber 1 hashfull 400\n",
		},
		{
			name: "Current Line with CPU",
			cmd: infoCmd{
				currline: OptionalOf(currentLine{
					cpunr: OptionalOf(1),
					moves: []chess.Move{
						{FromSquare: chess.D2, ToSquare: chess.D4},
						{FromSquare: chess.D7, ToSquare: chess.D5},
					},
				}),
			},
			expected: "info currline 1 d2d4 d7d5\n",
		},
		{
			name: "Refutation Line",
			cmd: infoCmd{
				refutation: []chess.Move{
					{FromSquare: chess.D1, ToSquare: chess.H5},
					{FromSquare: chess.G6, ToSquare: chess.H5},
				},
			},
			expected: "info refutation d1h5 g6h5\n",
		},
		{
			name: "Arbitrary String",
			cmd: infoCmd{
				stringMsg: OptionalOf("checking for early draw"),
			},
			expected: "info string checking for early draw\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using pointer to cmd to satisfy potential engineToClientCmd interface
			var cmd engineToClientCmd = &tt.cmd

			text, err := cmd.marshalText()
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}

			if string(text) != tt.expected {
				t.Errorf("got %q, expected %q", string(text), tt.expected)
			}
		})
	}
}

func TestOptionCommands_UCICompliance(t *testing.T) {
	tests := []struct {
		name     string
		cmd      engineToClientCmd
		expected string
	}{
		{
			name: "Type Check - Pondering ability",
			cmd: &CheckOptionCmd{
				Name:         "Ponder",
				DefaultValue: true,
			},
			expected: "option name Ponder type check default true\n",
		},
		{
			name: "Type Spin - Hash size with range",
			cmd: &SpinOptionCmd{
				Name:         "Hash",
				DefaultValue: 1,
				Min:          1,
				Max:          128,
			},
			expected: "option name Hash type spin default 1 min 1 max 128\n",
		},
		{
			name: "Type Combo - Playing styles",
			cmd: &ComboOptionCmd{
				Name:         "Style",
				DefaultValue: "Normal",
				Variants:     []string{"Solid", "Normal", "Risky"},
			},
			expected: "option name Style type combo default Normal var Solid var Normal var Risky\n",
		},
		{
			name: "Type Button - Clear internal state",
			cmd: &ButtonOptionCmd{
				Name: "Clear Hash",
			},
			expected: "option name Clear Hash type button\n",
		},
		{
			name: "Type String - Empty path handling",
			cmd: &StringOptionCmd{
				Name:         "NalimovPath",
				DefaultValue: "",
			},
			expected: "option name NalimovPath type string default <empty>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := tt.cmd.marshalText()
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}
			if string(text) != tt.expected {
				t.Errorf("\ngot:      %q\nexpected: %q", string(text), tt.expected)
			}
		})
	}
}
