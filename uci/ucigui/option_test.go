// Copyright (C) 2025 Brigham Skarda

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

// MOST OF THE CODE IN THIS FILE WAS WRITTEN BY MICROSOFT COPILOT.

package ucigui

import (
	"slices"
	"testing"
)

func (o *Option) equals(other *Option) bool {
	if o == nil || other == nil {
		return o == other // both must be nil
	}

	if o.Name != other.Name || o.OType != other.OType {
		return false
	}

	if (o.Default == nil) != (other.Default == nil) ||
		(o.Default != nil && other.Default != nil && *o.Default != *other.Default) {
		return false
	}

	if (o.Min == nil) != (other.Min == nil) ||
		(o.Min != nil && other.Min != nil && *o.Min != *other.Min) {
		return false
	}

	if (o.Max == nil) != (other.Max == nil) ||
		(o.Max != nil && other.Max != nil && *o.Max != *other.Max) {
		return false
	}

	return slices.Equal(o.Var, other.Var)
}

func TestOptionParsing_CheckType(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name Nullmove type check default true\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "true"
	expected := &Option{
		Name:    "Nullmove",
		OType:   Check,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_SpinType(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name SkillLevel type spin min -1 max 99 default 0\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	min := -1
	max := 99
	def := "0"
	expected := &Option{
		Name:    "SkillLevel",
		OType:   Spin,
		Min:     &min,
		Max:     &max,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_ComboType(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name Style type combo var Opt1 var Opt2 default Opt1\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "Opt1"
	expected := &Option{
		Name:    "Style",
		OType:   Combo,
		Default: &def,
		Var:     []string{"Opt1", "Opt2"},
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_ButtonType(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name ClearHash type button \n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	expected := &Option{
		Name:  "ClearHash",
		OType: Button,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_StringType(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name BookPath type string default My Favorite Book\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "My Favorite Book"
	expected := &Option{
		Name:    "BookPath",
		OType:   String,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_VarsWithSpaces(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name BookPath type combo default Default 1 var Default 1 var Default 2\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	if *parsedCommand.Default != "Default 1" {
		t.Errorf("options do not match: expected %v, got %v", "Default 1", parsedCommand.Default)
	}

	expected := []string{"Default 1", "Default 2"}
	if !slices.Equal(parsedCommand.Var, expected) {
		t.Errorf("options do not match: expected %v, got %v", expected, parsedCommand.Var)
	}
}

func TestOptionParsing_BadCheck(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name Nullmove type check default notABool\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	expected := &Option{
		Name:  "Nullmove",
		OType: Check,
	}

	if expected.Name != parsedCommand.Name || expected.OType != parsedCommand.OType {
		t.Errorf("option Names or OType do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_BadSpin(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name SkillLevel type spin min hi max lol default 0\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "0"
	expected := &Option{
		Name:    "SkillLevel",
		OType:   Spin,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_StringWithKeywords(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name BookPath type string default My min Favorite default Book max ra\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "My min Favorite default Book max ra"
	expected := &Option{
		Name:    "BookPath",
		OType:   String,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_InvalidOptions(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := []string{
		"option name Hash type combo\n",
		"option name NalimovPath type combo\n",
		"option name NalimovCache type combo\n",
		"option name Ponder type combo\n",
		"option name OwnBook type combo\n",
		"option name MultiPV type combo\n",
		"option name UCI_ShowCurrLine type combo\n",
		"option name UCI_ShowRefutations type combo\n",
		"option name UCI_LimitStrength type combo\n",
		"option name UCI_Elo type combo\n",
		"option name UCI_AnalyseMode type combo\n",
		"option name UCI_Opponent type combo\n",
		"option name UCI_EngineAbout type combo\n",
		"option name UCI_ShredderbasesPath type combo\n",
		"option name UCI_SetPositionValue type combo\n",
	}
	for _, o := range options {
		dummy.stdoutWriter.Write([]byte(o))
	}
	dummy.stdoutWriter.Write([]byte("option name Nullmove type check default true\n"))

	parsedCommand := (<-client.commandBuf).(*Option)
	if parsedCommand.Name != "Nullmove" {
		t.Errorf("parsed invalid command %v", parsedCommand.Name)
	}
}

func TestOptionParsing_ValidOptions(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := []string{
		"option name Hash type spin\n",
		"option name NalimovPath type string\n",
		"option name NalimovCache type spin\n",
		"option name Ponder type check\n",
		"option name OwnBook type check\n",
		"option name MultiPV type spin\n",
		"option name UCI_ShowCurrLine type check\n",
		"option name UCI_ShowRefutations type check\n",
		"option name UCI_LimitStrength type check\n",
		"option name UCI_Elo type spin\n",
		"option name UCI_AnalyseMode type check\n",
		"option name UCI_Opponent type string\n",
		"option name UCI_EngineAbout type string\n",
		"option name UCI_ShredderbasesPath type string\n",
		"option name UCI_SetPositionValue type string\n",
	}

	for _, o := range options {
		dummy.stdoutWriter.Write([]byte(o))
	}

	for range options {
		_ = (<-client.commandBuf).(*Option)
	}

}

func TestOptionParsing_ArbitraryWhiteSpace(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option     name\t   \tBook Path type combo default \tDefault 1 var \t    Default \t2   var Default \t2\t\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "Default \t2"
	expected := &Option{
		Name:    "Book Path",
		OType:   Combo,
		Default: &def,
		Var:     []string{"Default 1", def},
	}

	if parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_BadTokenBefore(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("joe mama option name Nullmove type check default true\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	def := "true"
	expected := &Option{
		Name:    "Nullmove",
		OType:   Check,
		Default: &def,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}

func TestOptionParsing_BadTokenAtEnd(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("option name Nullmove type check my mama\n"))

	parsedCommand := (<-client.commandBuf).(*Option)

	expected := &Option{
		Name:  "Nullmove",
		OType: Check,
	}

	if !parsedCommand.equals(expected) {
		t.Errorf("options do not match: expected %v, got %v", *expected, *parsedCommand)
	}
}
