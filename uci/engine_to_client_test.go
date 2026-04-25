// Copyright (C) 2026 Brigham Skarda
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
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
	var cmd engineToClientCmd = &readyOkCmd{}

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
	var cmd engineToClientCmd = &BestMove{
		Move:       chess.Move{FromSquare: chess.A1, ToSquare: chess.A2},
		PonderMove: OptionalOf(chess.Move{FromSquare: chess.A1, ToSquare: chess.A2, Promotion: chess.Knight}),
	}

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "bestmove a1a2 ponder a1a2n\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	cmd = &BestMove{
		Move: chess.Move{FromSquare: chess.A1, ToSquare: chess.A2},
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
	copyprotectCmd := copyProtectChecking
	var cmd engineToClientCmd = &copyprotectCmd

	text, err := cmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected := "copyprotection checking\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	copyprotectCmd = copyProtectOk
	text, err = copyprotectCmd.marshalText()
	if err != nil {
		t.Error("got unexpected error")
	}
	expected = "copyprotection ok\n"
	if string(text) != expected {
		t.Errorf("got %q, expected %q", text, expected)
	}

	copyprotectCmd = copyProtectError
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
	cmd := &InfoCmd{
		Depth:    OptionalOf(20),
		SelDepth: OptionalOf(25),
		Time:     OptionalOf(1200),
		Nodes:    OptionalOf(5000000),
		Pv:       []chess.Move{{FromSquare: chess.E2, ToSquare: chess.E4}, {FromSquare: chess.E7, ToSquare: chess.E5}},
		MultiPv:  OptionalOf(1),
		Score: OptionalOf(InfoScore{Score: 35,
			IsMate:       true,
			IsLowerBound: true}),
		CurrMove:       OptionalOf(chess.Move{FromSquare: chess.G1, ToSquare: chess.F3}),
		CurrMoveNumber: OptionalOf(1),
		HashFull:       OptionalOf(500),
		Nps:            OptionalOf(450000),
		TbHits:         OptionalOf(100),
		SbHits:         OptionalOf(50),
		CpuLoad:        OptionalOf(200),
		Refutation:     []chess.Move{{FromSquare: chess.D1, ToSquare: chess.H5}, {FromSquare: chess.G6, ToSquare: chess.H5}},
		CurrLine: OptionalOf(CurrentLine{
			CpuNr: OptionalOf(1),
			Moves: []chess.Move{{FromSquare: chess.A2, ToSquare: chess.A4}},
		}),
		StringMsg: OptionalOf("engine is thinking deeply"),
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
		cmd      InfoCmd
		expected string
	}{
		{
			name:     "Empty Info",
			cmd:      InfoCmd{},
			expected: "info\n",
		},
		{
			name: "Basic Search Stats",
			cmd: InfoCmd{
				Depth: OptionalOf(12),
				Nodes: OptionalOf(123456),
				Nps:   OptionalOf(100000),
			},
			expected: "info depth 12 nodes 123456 nps 100000\n",
		},
		{
			name: "Score Centipawns",
			cmd: InfoCmd{
				Score: OptionalOf(InfoScore{
					Score: -45,
				}),
			},
			expected: "info score cp -45\n",
		},
		{
			name: "Score Mate and Upperbound",
			cmd: InfoCmd{
				Score: OptionalOf(InfoScore{
					Score:        5,
					IsMate:       true,
					IsUpperBound: true,
				}),
			},
			expected: "info score mate 5 upperbound\n",
		},
		{
			name: "PV with Moves",
			cmd: InfoCmd{
				Depth: OptionalOf(2),
				Pv: []chess.Move{
					{FromSquare: chess.E2, ToSquare: chess.E4},
					{FromSquare: chess.E7, ToSquare: chess.E5},
				},
			},
			expected: "info depth 2 pv e2e4 e7e5\n",
		},
		{
			name: "Current Move and Progress",
			cmd: InfoCmd{
				CurrMove:       OptionalOf(chess.Move{FromSquare: chess.G1, ToSquare: chess.F3}),
				CurrMoveNumber: OptionalOf(1),
				HashFull:       OptionalOf(400),
			},
			expected: "info currmove g1f3 currmovenumber 1 hashfull 400\n",
		},
		{
			name: "Current Line with CPU",
			cmd: InfoCmd{
				CurrLine: OptionalOf(CurrentLine{
					CpuNr: OptionalOf(1),
					Moves: []chess.Move{
						{FromSquare: chess.D2, ToSquare: chess.D4},
						{FromSquare: chess.D7, ToSquare: chess.D5},
					},
				}),
			},
			expected: "info currline 1 d2d4 d7d5\n",
		},
		{
			name: "Refutation Line",
			cmd: InfoCmd{
				Refutation: []chess.Move{
					{FromSquare: chess.D1, ToSquare: chess.H5},
					{FromSquare: chess.G6, ToSquare: chess.H5},
				},
			},
			expected: "info refutation d1h5 g6h5\n",
		},
		{
			name: "Arbitrary String",
			cmd: InfoCmd{
				StringMsg: OptionalOf("checking for early draw"),
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
			cmd: &CheckOption{
				Name:         "Ponder",
				DefaultValue: true,
			},
			expected: "option name Ponder type check default true\n",
		},
		{
			name: "Type Spin - Hash size with range",
			cmd: &SpinOption{
				Name:         "Hash",
				DefaultValue: 1,
				Min:          1,
				Max:          128,
			},
			expected: "option name Hash type spin default 1 min 1 max 128\n",
		},
		{
			name: "Type Combo - Playing styles",
			cmd: &ComboOption{
				Name:         "Style",
				DefaultValue: "Normal",
				Variants:     []string{"Solid", "Normal", "Risky"},
			},
			expected: "option name Style type combo default Normal var Solid var Normal var Risky\n",
		},
		{
			name: "Type Button - Clear internal state",
			cmd: &ButtonOption{
				Name: "Clear Hash",
			},
			expected: "option name Clear Hash type button\n",
		},
		{
			name: "Type String - Empty path handling",
			cmd: &StringOption{
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
