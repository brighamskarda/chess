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
	"slices"
	"testing"

	"github.com/brighamskarda/chess/v2"
)

func TestUciCommandUnmarshal(t *testing.T) {
	if (&uciCmd{}).UnmarshalText([]byte("uci\n")) != nil {
		t.Errorf("incorrect result for valid uci command, got error")
	}

	if (&uciCmd{}).UnmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid uci command, did not get error")
	}
}

func TestDebugCommandUnmarshal(t *testing.T) {
	if (&debugCmd{}).UnmarshalText([]byte("debug on\n")) != nil {
		t.Errorf("incorrect result for valid debug command, got error")
	}

	if (&debugCmd{}).UnmarshalText([]byte("debug off\n")) != nil {
		t.Errorf("incorrect result for valid debug command, got error")
	}

	if (&debugCmd{}).UnmarshalText([]byte("debug of\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}

	if (&debugCmd{}).UnmarshalText([]byte("debu on\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}

	if (&debugCmd{}).UnmarshalText([]byte("on\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}
}

func TestCmdGetMessage(t *testing.T) {
	cmd := &uciCmd{}
	cmd.UnmarshalText([]byte("uci\n"))
	if string(cmd.getCmdText()) != "uci\n" {
		t.Errorf("command text was not stored properly")
	}

	cmd2 := &debugCmd{}
	cmd2.UnmarshalText([]byte("debug on\n"))
	if string(cmd2.getCmdText()) != "debug on\n" {
		t.Errorf("command text was not stored properly")
	}
}

func TestIsReadyCommandUnmarshal(t *testing.T) {
	if (&isReadyCmd{}).UnmarshalText([]byte("isready\n")) != nil {
		t.Errorf("incorrect result for valid isReady command, got error")
	}

	if (&isReadyCmd{}).UnmarshalText([]byte("isReadyy\n")) == nil {
		t.Errorf("incorrect result for valid isReady command, did not get error")
	}
}

func TestCheckOptionCommand(t *testing.T) {
	checkOption := setCheckOptionCmd{}
	if checkOption.UnmarshalText([]byte("setoption name check  box value true\n")) != nil {
		t.Errorf("incorrect result for valid setOption command, got error")
	}
	if checkOption.checkbox != true {
		t.Errorf("incorrect result for valid setOption command, checkbox should be true")
	}
	if checkOption.name != "check  box" {
		t.Errorf("incorrect result for valid setOption command, optionId should be \"checkbox\"")
	}

	if checkOption.UnmarshalText([]byte("setoption name the best option value false\n")) != nil {
		t.Errorf("incorrect result for valid setOption command, got error")
	}
	if checkOption.checkbox != false {
		t.Errorf("incorrect result for valid setOption command, checkbox should be true")
	}
	if checkOption.name != "the best option" {
		t.Errorf("incorrect result for valid setOption command, optionId should be \"the best option\"")
	}

	if checkOption.UnmarshalText([]byte("setoption name badOpt false\n")) == nil {
		t.Errorf("incorrect result for valid setOption command, expected error")
	}
	if checkOption.UnmarshalText([]byte("setoption name badOpt value alse\n")) == nil {
		t.Errorf("incorrect result for valid setOption command, expected error")
	}
}

func TestSpinOptionCommand(t *testing.T) {
	spinOption := setSpinOptionCmd{}

	// Valid case
	if err := spinOption.UnmarshalText([]byte("setoption name Selectivity value 2\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if spinOption.value != 2 || spinOption.name != "Selectivity" {
		t.Errorf("incorrect parse: got value %d, id %s", spinOption.value, spinOption.name)
	}

	// Valid case with spaces in name
	if err := spinOption.UnmarshalText([]byte("setoption name Multi  PV value 4\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if spinOption.name != "Multi  PV" {
		t.Errorf("expected two spaces in name")
	}

	// Invalid: non-numeric value
	if err := spinOption.UnmarshalText([]byte("setoption name Selectivity value fast\n")); err == nil {
		t.Errorf("expected error for non-numeric spin value")
	}
}

func TestStringOptionCommand(t *testing.T) {
	strOption := setStringOptionCmd{}

	// Valid case
	if err := strOption.UnmarshalText([]byte("setoption name Nalimov  Path value C:\\che  ss\\tablebases\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if strOption.value != "C:\\che  ss\\tablebases" {
		t.Errorf("expected path string, got %s", strOption.value)
	}
	if strOption.name != "Nalimov  Path" {
		t.Errorf("expected two spaces in string, got %s", strOption.value)
	}

	// Valid case: empty string
	if err := strOption.UnmarshalText([]byte("setoption name Author value <empty>\n")); err != nil {
		t.Errorf("expected no error for empty string value")
	}

	// Invalid case: empty string
	if err := strOption.UnmarshalText([]byte("setoption name Author value \n")); err == nil {
		t.Errorf("expected error for empty string value")
	}
}

func TestComboOptionCommand(t *testing.T) {
	comboOption := setComboOptionCmd{}

	// Valid case
	if err := comboOption.UnmarshalText([]byte("setoption name Style  Boy value Aggressive  Man\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if comboOption.value != "Aggressive  Man" || comboOption.name != "Style  Boy" {
		t.Errorf("incorrect combo parse: got %s", comboOption.value)
	}

	// Invalid: Missing "value" keyword
	if err := comboOption.UnmarshalText([]byte("setoption name Style Aggressive\n")); err == nil {
		t.Errorf("expected error for missing 'value' keyword")
	}
}

func TestButtonOptionCommand(t *testing.T) {
	btnOption := setButtonOptionCmd{}

	// Valid case: Button options do not have a "value" suffix
	if err := btnOption.UnmarshalText([]byte("setoption name Clear  Hash\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if btnOption.name != "Clear  Hash" {
		t.Errorf("expected optionId 'Clear  Hash', got '%s'", btnOption.name)
	}
}

func TestRegisterCommand(t *testing.T) {
	registerCmd := registerCmd{}

	if err := registerCmd.UnmarshalText([]byte("register later\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.regType != later {
		t.Errorf("expected regType to be %d, got %d", later, registerCmd.regType)
	}

	if err := registerCmd.UnmarshalText([]byte("register name cornelius the  third\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.regType != name {
		t.Errorf("expected regType to be %d, got %d", name, registerCmd.regType)
	}
	if registerCmd.value != "cornelius the  third" {
		t.Errorf("expected regType to be %q, got %q", "cornelius the  third", registerCmd.value)
	}

	if err := registerCmd.UnmarshalText([]byte("register code cornelius the  third\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.regType != code {
		t.Errorf("expected regType to be %d, got %d", name, registerCmd.regType)
	}
	if registerCmd.value != "cornelius the  third" {
		t.Errorf("expected regType to be %q, got %q", "cornelius the  third", registerCmd.value)
	}

	if err := registerCmd.UnmarshalText([]byte("register code\n")); err == nil {
		t.Errorf("code without value should produce an error")
	}

	if err := registerCmd.UnmarshalText([]byte("register name\n")); err == nil {
		t.Errorf("name without value should produce an error")
	}
}

func TestUciNewGameCommandUnmarshal(t *testing.T) {
	if (&uciNewGameCmd{}).UnmarshalText([]byte("ucinewgame\n")) != nil {
		t.Errorf("incorrect result for valid ucinewgame command, got error")
	}

	if (&uciNewGameCmd{}).UnmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid ucinewgame command, did not get error")
	}
}

func TestPositionCommandUnmarshalStartPos(t *testing.T) {
	defaultPos := &chess.Position{}
	defaultPos.UnmarshalText([]byte(chess.DefaultFEN))

	posCmd := &positionCmd{}
	if posCmd.UnmarshalText([]byte("position startpos\n")) != nil {
		t.Errorf("got error for startpos keyword")
	}
	if !posCmd.position.Equal(defaultPos) {
		t.Errorf("startpos is not default position")
	}
	if posCmd.moves != nil {
		t.Errorf("expected moves to be nil")
	}

	if posCmd.UnmarshalText([]byte("position startpos moves\n")) != nil {
		t.Errorf("got error for empty moves")
	}
	if posCmd.moves == nil || len(posCmd.moves) != 0 {
		t.Errorf("expected empty move slice")
	}

	if posCmd.UnmarshalText([]byte("position startpos moves e2e4 d7d5\n")) != nil {
		t.Errorf("got error for filled moves")
	}
	if len(posCmd.moves) != 2 {
		t.Errorf("expected filled move slice")
	}
	if posCmd.moves[0] != (chess.Move{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType}) ||
		posCmd.moves[1] != (chess.Move{FromSquare: chess.D7, ToSquare: chess.D5, Promotion: chess.NoPieceType}) {
		t.Errorf("moves where incorrect")
	}
}

func TestPositionCommandUnmarshalOtherPos(t *testing.T) {
	otherFen := "4k2r/6r1/8/8/8/8/3R4/R3K3 w Qk - 0 1"
	otherPos := &chess.Position{}
	otherPos.UnmarshalText([]byte(otherFen))

	posCmd := &positionCmd{}
	if posCmd.UnmarshalText([]byte("position fen "+otherFen+"\n")) != nil {
		t.Errorf("got error for other fen")
	}
	if !posCmd.position.Equal(otherPos) {
		t.Errorf("startpos is not default position")
	}
	if posCmd.moves != nil {
		t.Errorf("expected moves to be nil")
	}

	if posCmd.UnmarshalText([]byte("position fen "+otherFen+" moves\n")) != nil {
		t.Errorf("got error for empty moves")
	}
	if posCmd.moves == nil || len(posCmd.moves) != 0 {
		t.Errorf("expected empty move slice")
	}

	if posCmd.UnmarshalText([]byte("position fen "+otherFen+" moves e2e4 d7d5\n")) != nil {
		t.Errorf("got error for filled moves")
	}
	if len(posCmd.moves) != 2 {
		t.Errorf("expected filled move slice")
	}
	if posCmd.moves[0] != (chess.Move{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType}) ||
		posCmd.moves[1] != (chess.Move{FromSquare: chess.D7, ToSquare: chess.D5, Promotion: chess.NoPieceType}) {
		t.Errorf("moves where incorrect")
	}
}

func TestPositionCommandUnmarshalInvalidPos(t *testing.T) {
	otherFen := "4k2r/6r1/8/8/7/8/3R4/R3K3 w Qk - 0 1"
	otherPos := &chess.Position{}
	otherPos.UnmarshalText([]byte(otherFen))

	posCmd := &positionCmd{}
	if posCmd.UnmarshalText([]byte("position fen "+otherFen+"\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
	if posCmd.UnmarshalText([]byte("position fen "+otherFen+" moves\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
	if posCmd.UnmarshalText([]byte("position fen "+otherFen+" moves e2e4 d7d5\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
}

func TestStopCommandUnmarshal(t *testing.T) {
	if (&stopCmd{}).UnmarshalText([]byte("stop\n")) != nil {
		t.Errorf("incorrect result for valid stop command, got error")
	}

	if (&stopCmd{}).UnmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid stop command, did not get error")
	}
}

func TestPonderHitCommandUnmarshal(t *testing.T) {
	if (&ponderhitCmd{}).UnmarshalText([]byte("ponderhit\n")) != nil {
		t.Errorf("incorrect result for valid ponderhit command, got error")
	}

	if (&ponderhitCmd{}).UnmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid ponderhit command, did not get error")
	}
}

func TestQuitCommandUnmarshal(t *testing.T) {
	if (&quitCmd{}).UnmarshalText([]byte("quit\n")) != nil {
		t.Errorf("incorrect result for valid quit command, got error")
	}

	if (&quitCmd{}).UnmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid quit command, did not get error")
	}
}

func TestGoCmdUnmarshal(t *testing.T) {
	goCmd := &goCmd{}
	if err := goCmd.UnmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime 999 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err != nil {
		t.Errorf("incorrect result for valid go command, got error: %v", err)
	}
	validateGoCmd(goCmd, t)

	if goCmd.UnmarshalText([]byte("go\n")) != nil {
		t.Errorf("incorrect result for valid go command, got error")
	}

	if goCmd.searchMoves != nil ||
		goCmd.ponder ||
		goCmd.wtime.HasValue() ||
		goCmd.btime.HasValue() ||
		goCmd.winc.HasValue() ||
		goCmd.binc.HasValue() ||
		goCmd.movestogo.HasValue() ||
		goCmd.depth.HasValue() ||
		goCmd.nodes.HasValue() ||
		goCmd.mate.HasValue() ||
		goCmd.movetime.HasValue() ||
		goCmd.infinite {
		t.Error("variables were set for empty go command")
	}

	var genericCommand clientToEngineCmd = goCmd
	if string(genericCommand.getCmdText()) != "go\n" {
		t.Error("getCmdText returned wrong value")
	}
}

func validateGoCmd(goCmd *goCmd, t *testing.T) {
	if !slices.Equal(goCmd.searchMoves, []chess.Move{
		{FromSquare: chess.E2, ToSquare: chess.E4},
		{FromSquare: chess.D7, ToSquare: chess.D5}}) {
		t.Error("search moves do not match")
	}
	if !goCmd.ponder {
		t.Error("ponder not set")
	}
	if !goCmd.wtime.HasValue() || goCmd.wtime.Value() != 999 {
		t.Error("wtime not set correctly")
	}
	if !goCmd.btime.HasValue() || goCmd.btime.Value() != 888 {
		t.Error("btime not set correctly")
	}
	if !goCmd.winc.HasValue() || goCmd.winc.Value() != 777 {
		t.Error("winc not set correctly")
	}
	if !goCmd.binc.HasValue() || goCmd.binc.Value() != 666 {
		t.Error("binc not set correctly")
	}
	if !goCmd.movestogo.HasValue() || goCmd.movestogo.Value() != 555 {
		t.Error("movestogo not set correctly")
	}
	if !goCmd.depth.HasValue() || goCmd.depth.Value() != 444 {
		t.Error("depth not set correctly")
	}
	if !goCmd.nodes.HasValue() || goCmd.nodes.Value() != 333 {
		t.Error("nodes not set correctly")
	}
	if !goCmd.mate.HasValue() || goCmd.mate.Value() != 222 {
		t.Error("mate not set correctly")
	}
	if !goCmd.movetime.HasValue() || goCmd.movetime.Value() != 111 {
		t.Error("movetime not set correctly")
	}
	if !goCmd.infinite {
		t.Error("infinite not set")
	}
}

func TestGoCmdUnmarshalError(t *testing.T) {
	goCmd := &goCmd{}
	if err := goCmd.UnmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: missing value")
	}

	if err := goCmd.UnmarshalText([]byte("go searchmoves ponder wtime 999 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err != nil {
		t.Errorf("incorrect result for invalid go command, expected no error: no search moves should be parsed anyways")
	}

	if err := goCmd.UnmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime 999 888 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: extra value")
	}

	if err := goCmd.UnmarshalText([]byte("go searchmove\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: invalid command")
	}
}
