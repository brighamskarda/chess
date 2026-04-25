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

// mockEngine simply keeps track of how many times each of its functions are called.
//
// Used for testing.
type mockEngine struct {
	output     func(*InfoCmd)
	debugState bool

	initialize int
	name       int
	author     int
	options    int
	debug      int
	setOption  int
	quit       int
}

func (engine *mockEngine) Initialize(o func(*InfoCmd)) {
	engine.output = o
	engine.initialize++
}

func (engine *mockEngine) Name() string {
	engine.name++
	return "mockEngine v0.1"
}

func (engine *mockEngine) Author() string {
	engine.author++
	return "Brigham Skarda"
}

func (engine *mockEngine) Options() []OptionCmd {
	engine.options++
	return []OptionCmd{
		&CheckOptionCmd{
			Name:         "checkOpt",
			DefaultValue: true,
		},
		&SpinOptionCmd{
			Name:         "spinOpt",
			DefaultValue: 3,
			Min:          1,
			Max:          5,
		},
		&ComboOptionCmd{
			Name:         "comboOpt",
			DefaultValue: "one",
			Variants:     []string{"one", "two", "three"},
		},
		&StringOptionCmd{
			Name:         "stringOpt",
			DefaultValue: "sss",
		},
		&ButtonOptionCmd{
			Name: "buttonOpt",
		},
	}
}

func (engine *mockEngine) SetDebug(value bool) {
	engine.debugState = value
	engine.debug++
}

func (engine *mockEngine) SetOption(option SetOptionCmd) {
	engine.setOption++
}

func (engine *mockEngine) Quit() {
	engine.quit++
}
