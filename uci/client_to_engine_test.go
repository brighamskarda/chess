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
	"fmt"
	"slices"
	"testing"

	"github.com/brighamskarda/chess/v2"
)

func TestUciCommandUnmarshal(t *testing.T) {
	if (&uciCmd{}).unmarshalText([]byte("uci\n")) != nil {
		t.Errorf("incorrect result for valid uci command, got error")
	}

	if (&uciCmd{}).unmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid uci command, did not get error")
	}
}

func TestDebugCommandUnmarshal(t *testing.T) {
	if (&debugCmd{}).unmarshalText([]byte("debug on\n")) != nil {
		t.Errorf("incorrect result for valid debug command, got error")
	}

	if (&debugCmd{}).unmarshalText([]byte("debug off\n")) != nil {
		t.Errorf("incorrect result for valid debug command, got error")
	}

	if (&debugCmd{}).unmarshalText([]byte("debug of\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}

	if (&debugCmd{}).unmarshalText([]byte("debu on\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}

	if (&debugCmd{}).unmarshalText([]byte("on\n")) == nil {
		t.Errorf("incorrect result for valid debug command, did not get error")
	}
}

func TestCmdGetMessage(t *testing.T) {
	cmd := &uciCmd{}
	cmd.unmarshalText([]byte("uci\n"))
	if string(cmd.getCmdText()) != "uci\n" {
		t.Errorf("command text was not stored properly")
	}

	cmd2 := &debugCmd{}
	cmd2.unmarshalText([]byte("debug on\n"))
	if string(cmd2.getCmdText()) != "debug on\n" {
		t.Errorf("command text was not stored properly")
	}
}

func TestIsReadyCommandUnmarshal(t *testing.T) {
	if (&isReadyCmd{}).unmarshalText([]byte("isready\n")) != nil {
		t.Errorf("incorrect result for valid isReady command, got error")
	}

	if (&isReadyCmd{}).unmarshalText([]byte("isReadyy\n")) == nil {
		t.Errorf("incorrect result for valid isReady command, did not get error")
	}
}

func TestCheckOptionCommand(t *testing.T) {
	checkOption := SetCheckOptionCmd{}
	if checkOption.unmarshalText([]byte("setoption name check  box value true\n")) != nil {
		t.Errorf("incorrect result for valid setOption command, got error")
	}
	if checkOption.Checkbox != true {
		t.Errorf("incorrect result for valid setOption command, checkbox should be true")
	}
	if checkOption.name != "check  box" {
		t.Errorf("incorrect result for valid setOption command, optionId should be \"checkbox\"")
	}

	if checkOption.unmarshalText([]byte("setoption name the best option value false\n")) != nil {
		t.Errorf("incorrect result for valid setOption command, got error")
	}
	if checkOption.Checkbox != false {
		t.Errorf("incorrect result for valid setOption command, checkbox should be true")
	}
	if checkOption.name != "the best option" {
		t.Errorf("incorrect result for valid setOption command, optionId should be \"the best option\"")
	}

	if checkOption.unmarshalText([]byte("setoption name badOpt false\n")) == nil {
		t.Errorf("incorrect result for valid setOption command, expected error")
	}
	if checkOption.unmarshalText([]byte("setoption name badOpt value alse\n")) == nil {
		t.Errorf("incorrect result for valid setOption command, expected error")
	}
}

func TestSpinOptionCommand(t *testing.T) {
	spinOption := SetSpinOptionCmd{}

	// Valid case
	if err := spinOption.unmarshalText([]byte("setoption name Selectivity value 2\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if spinOption.Value != 2 || spinOption.name != "Selectivity" {
		t.Errorf("incorrect parse: got value %d, id %s", spinOption.Value, spinOption.name)
	}

	// Valid case with spaces in name
	if err := spinOption.unmarshalText([]byte("setoption name Multi  PV value 4\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if spinOption.name != "Multi  PV" {
		t.Errorf("expected two spaces in name")
	}

	// Invalid: non-numeric value
	if err := spinOption.unmarshalText([]byte("setoption name Selectivity value fast\n")); err == nil {
		t.Errorf("expected error for non-numeric spin value")
	}
}

func TestStringOptionCommand(t *testing.T) {
	strOption := SetStringOptionCmd{}

	// Valid case
	if err := strOption.unmarshalText([]byte("setoption name Nalimov  Path value C:\\che  ss\\tablebases\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if strOption.Value != "C:\\che  ss\\tablebases" {
		t.Errorf("expected path string, got %s", strOption.Value)
	}
	if strOption.name != "Nalimov  Path" {
		t.Errorf("expected two spaces in string, got %s", strOption.Value)
	}

	// Valid case: empty string
	if err := strOption.unmarshalText([]byte("setoption name Author value <empty>\n")); err != nil {
		t.Errorf("expected no error for empty string value")
	}

	// Invalid case: empty string
	if err := strOption.unmarshalText([]byte("setoption name Author value \n")); err == nil {
		t.Errorf("expected error for empty string value")
	}
}

func TestButtonOptionCommand(t *testing.T) {
	btnOption := SetButtonOptionCmd{}

	// Valid case: Button options do not have a "value" suffix
	if err := btnOption.unmarshalText([]byte("setoption name Clear  Hash\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if btnOption.name != "Clear  Hash" {
		t.Errorf("expected optionId 'Clear  Hash', got '%s'", btnOption.name)
	}
}

func TestRegisterCommand(t *testing.T) {
	registerCmd := RegisterCmd{}

	if err := registerCmd.unmarshalText([]byte("register later\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !registerCmd.Later {
		t.Errorf("expected Later to be %v, got %v", true, registerCmd.Later)
	}

	if err := registerCmd.unmarshalText([]byte("register name cornelius the  third\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.Later {
		t.Errorf("expected Later to be %v, got %v", false, registerCmd.Later)
	}
	if registerCmd.Name.Value() != "cornelius the  third" {
		t.Errorf("expected Name to be %q, got %v", "cornelius the  third", registerCmd.Name)
	}
	if registerCmd.Code.HasValue() {
		t.Errorf("expected Code to be empty, got %v", registerCmd.Code)
	}

	if err := registerCmd.unmarshalText([]byte("register code cornelius the  third\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.Later {
		t.Errorf("expected Later to be %v, got %v", false, registerCmd.Later)
	}
	if registerCmd.Code.Value() != "cornelius the  third" {
		t.Errorf("expected Code to be %q, got %v", "cornelius the  third", registerCmd.Code)
	}
	if registerCmd.Name.HasValue() {
		t.Errorf("expected Name to be empty, got %v", registerCmd.Name)
	}

	if err := registerCmd.unmarshalText([]byte("register name cornelius the  third code 433\n")); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if registerCmd.Later {
		t.Errorf("expected Later to be %v, got %v", false, registerCmd.Later)
	}
	if registerCmd.Name.Value() != "cornelius the  third" {
		t.Errorf("expected Name to be %q, got %v", "cornelius the  third", registerCmd.Name)
	}
	if registerCmd.Code.Value() != "433" {
		t.Errorf("expected Code to be %q, got %v", "433", registerCmd.Code)
	}

	if err := registerCmd.unmarshalText([]byte("register code\n")); err == nil {
		t.Errorf("code without value should produce an error")
	}

	if err := registerCmd.unmarshalText([]byte("register name\n")); err == nil {
		t.Errorf("name without value should produce an error")
	}
}

func TestUciNewGameCommandUnmarshal(t *testing.T) {
	if (&uciNewGameCmd{}).unmarshalText([]byte("ucinewgame\n")) != nil {
		t.Errorf("incorrect result for valid ucinewgame command, got error")
	}

	if (&uciNewGameCmd{}).unmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid ucinewgame command, did not get error")
	}
}

func TestPositionCommandUnmarshalStartPos(t *testing.T) {
	defaultPos := &chess.Position{}
	defaultPos.UnmarshalText([]byte(chess.DefaultFEN))

	posCmd := &positionCmd{}
	if posCmd.unmarshalText([]byte("position startpos\n")) != nil {
		t.Errorf("got error for startpos keyword")
	}
	if !posCmd.position.Equal(defaultPos) {
		t.Errorf("startpos is not default position")
	}
	if posCmd.moves != nil {
		t.Errorf("expected moves to be nil")
	}

	if posCmd.unmarshalText([]byte("position startpos moves\n")) != nil {
		t.Errorf("got error for empty moves")
	}
	if posCmd.moves == nil || len(posCmd.moves) != 0 {
		t.Errorf("expected empty move slice")
	}

	if posCmd.unmarshalText([]byte("position startpos moves e2e4 d7d5\n")) != nil {
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
	if posCmd.unmarshalText([]byte("position fen "+otherFen+"\n")) != nil {
		t.Errorf("got error for other fen")
	}
	if !posCmd.position.Equal(otherPos) {
		t.Errorf("startpos is not default position")
	}
	if posCmd.moves != nil {
		t.Errorf("expected moves to be nil")
	}

	if posCmd.unmarshalText([]byte("position fen "+otherFen+" moves\n")) != nil {
		t.Errorf("got error for empty moves")
	}
	if posCmd.moves == nil || len(posCmd.moves) != 0 {
		t.Errorf("expected empty move slice")
	}

	if posCmd.unmarshalText([]byte("position fen "+otherFen+" moves e2e4 d7d5\n")) != nil {
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
	if posCmd.unmarshalText([]byte("position fen "+otherFen+"\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
	if posCmd.unmarshalText([]byte("position fen "+otherFen+" moves\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
	if posCmd.unmarshalText([]byte("position fen "+otherFen+" moves e2e4 d7d5\n")) == nil {
		t.Errorf("did not get error for invalid FEN")
	}
}

func TestStopCommandUnmarshal(t *testing.T) {
	if (&stopCmd{}).unmarshalText([]byte("stop\n")) != nil {
		t.Errorf("incorrect result for valid stop command, got error")
	}

	if (&stopCmd{}).unmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid stop command, did not get error")
	}
}

func TestPonderHitCommandUnmarshal(t *testing.T) {
	if (&ponderhitCmd{}).unmarshalText([]byte("ponderhit\n")) != nil {
		t.Errorf("incorrect result for valid ponderhit command, got error")
	}

	if (&ponderhitCmd{}).unmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid ponderhit command, did not get error")
	}
}

func TestQuitCommandUnmarshal(t *testing.T) {
	if (&quitCmd{}).unmarshalText([]byte("quit\n")) != nil {
		t.Errorf("incorrect result for valid quit command, got error")
	}

	if (&quitCmd{}).unmarshalText([]byte("ui\n")) == nil {
		t.Errorf("incorrect result for valid quit command, did not get error")
	}
}

func TestGoCmdUnmarshal(t *testing.T) {
	goCmd := &EvaluateCmd{}
	if err := goCmd.unmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime 999 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err != nil {
		t.Errorf("incorrect result for valid go command, got error: %v", err)
	}
	validateGoCmd(goCmd, t)

	if goCmd.unmarshalText([]byte("go\n")) != nil {
		t.Errorf("incorrect result for valid go command, got error")
	}

	if goCmd.SearchMoves != nil ||
		goCmd.Ponder ||
		goCmd.Wtime.HasValue() ||
		goCmd.Btime.HasValue() ||
		goCmd.Winc.HasValue() ||
		goCmd.Binc.HasValue() ||
		goCmd.MovesToGo.HasValue() ||
		goCmd.Depth.HasValue() ||
		goCmd.Nodes.HasValue() ||
		goCmd.Mate.HasValue() ||
		goCmd.MoveTime.HasValue() ||
		goCmd.Infinite {
		t.Error("variables were set for empty go command")
	}

	var genericCommand clientToEngineCmd = goCmd
	if string(genericCommand.getCmdText()) != "go\n" {
		t.Error("getCmdText returned wrong value")
	}
}

func validateGoCmd(goCmd *EvaluateCmd, t *testing.T) {
	if !slices.Equal(goCmd.SearchMoves, []chess.Move{
		{FromSquare: chess.E2, ToSquare: chess.E4},
		{FromSquare: chess.D7, ToSquare: chess.D5}}) {
		t.Error("search moves do not match")
	}
	if !goCmd.Ponder {
		t.Error("ponder not set")
	}
	if !goCmd.Wtime.HasValue() || goCmd.Wtime.Value() != 999 {
		t.Error("wtime not set correctly")
	}
	if !goCmd.Btime.HasValue() || goCmd.Btime.Value() != 888 {
		t.Error("btime not set correctly")
	}
	if !goCmd.Winc.HasValue() || goCmd.Winc.Value() != 777 {
		t.Error("winc not set correctly")
	}
	if !goCmd.Binc.HasValue() || goCmd.Binc.Value() != 666 {
		t.Error("binc not set correctly")
	}
	if !goCmd.MovesToGo.HasValue() || goCmd.MovesToGo.Value() != 555 {
		t.Error("movestogo not set correctly")
	}
	if !goCmd.Depth.HasValue() || goCmd.Depth.Value() != 444 {
		t.Error("depth not set correctly")
	}
	if !goCmd.Nodes.HasValue() || goCmd.Nodes.Value() != 333 {
		t.Error("nodes not set correctly")
	}
	if !goCmd.Mate.HasValue() || goCmd.Mate.Value() != 222 {
		t.Error("mate not set correctly")
	}
	if !goCmd.MoveTime.HasValue() || goCmd.MoveTime.Value() != 111 {
		t.Error("movetime not set correctly")
	}
	if !goCmd.Infinite {
		t.Error("infinite not set")
	}
}

func TestGoCmdUnmarshalError(t *testing.T) {
	goCmd := &EvaluateCmd{}
	if err := goCmd.unmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: missing value")
	}

	if err := goCmd.unmarshalText([]byte("go searchmoves ponder wtime 999 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err != nil {
		t.Errorf("incorrect result for invalid go command, expected no error: no search moves should be parsed anyways")
	}

	if err := goCmd.unmarshalText([]byte("go searchmoves e2e4 d7d5 ponder wtime 999 888 btime 888 winc " +
		"777 binc 666 movestogo 555 depth 444 nodes 333 mate 222 movetime 111 infinite\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: extra value")
	}

	if err := goCmd.unmarshalText([]byte("go searchmove\n")); err == nil {
		t.Errorf("incorrect result for invalid go command, expected error: invalid command")
	}
}

func TestParseClientToEngineCmd_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // prefix or type check
	}{
		{"uci", "uci \n", "*uci.uciCmd"},
		{"debug", "debug on\n", "*uci.debugCmd"},
		{"isready", "isready \n", "*uci.isReadyCmd"},
		{"setoption check", "setoption name MyCheck value true\n", "*uci.SetCheckOptionCmd"},
		{"setoption spin", "setoption name MySpin value 10\n", "*uci.SetSpinOptionCmd"},
		{"setoption string", "setoption name MyStr value hello\n", "*uci.SetStringOptionCmd"},
		{"setoption button", "setoption name MyButton\n", "*uci.SetButtonOptionCmd"},
		{"register", "register later\n", "*uci.RegisterCmd"},
		{"ucinewgame", "ucinewgame \n", "*uci.uciNewGameCmd"},
		{"position", "position startpos moves e2e4\n", "*uci.positionCmd"},
		{"go", "go infinite\n", "*uci.EvaluateCmd"},
		{"stop", "stop \n", "*uci.stopCmd"},
		{"ponderhit", "ponderhit \n", "*uci.ponderhitCmd"},
		{"quit", "quit \n", "*uci.quitCmd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parseClientToEngineCmd([]byte(tt.input))
			if err != nil {
				t.Fatalf("failed to parse %q: %v", tt.name, err)
			}
			actualType := fmt.Sprintf("%T", cmd)
			if actualType != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, actualType)
			}
		})
	}
}

func TestParseClientToEngineCmd_LeadingInvalidWords(t *testing.T) {
	// The implementation uses strings.Index to find the first valid keyword.
	// This ensures that "garbage uci" still parses as a uciCmd.
	input := "garbage words uci \n"
	cmd, err := parseClientToEngineCmd([]byte(input))
	if err != nil {
		t.Fatalf("expected successful parse despite leading garbage, got error: %v", err)
	}

	if _, ok := cmd.(*uciCmd); !ok {
		t.Errorf("expected *uci.uciCmd, got %T", cmd)
	}

	// Verify it captured the text starting from the valid command
	if string(cmd.getCmdText()) != "uci \n" {
		t.Errorf("expected command text to be 'uci \\n', got %q", string(cmd.getCmdText()))
	}
}
