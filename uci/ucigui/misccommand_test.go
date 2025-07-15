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

package ucigui

import (
	"testing"
	"time"

	"github.com/brighamskarda/chess/v2"
)

func TestIdParsing_Name(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("id name \t Stockfish 17.1 author \t stillInName\t \n"))

	parsedCommand := client.commandBuf.Next().(idCommand)

	expected := idCommand{
		idt:   name,
		value: "Stockfish 17.1 author \t stillInName",
	}

	if parsedCommand != expected {
		t.Errorf("IDs do not match: expected %v, got %v", expected, parsedCommand)
	}
}

func TestIdParsing_Author(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("id \t author \t the Stockfish developers (see AUTHORS file) name \t stillInName\t \n"))

	parsedCommand := client.commandBuf.Next().(idCommand)

	expected := idCommand{
		idt:   author,
		value: "the Stockfish developers (see AUTHORS file) name \t stillInName",
	}

	if parsedCommand != expected {
		t.Errorf("IDs do not match: expected %v, got %v", expected, parsedCommand)
	}
}

func TestIdParsing_BadInput(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("id \t jeff  \t the Stockfish developers (see AUTHORS file) \t stillInName\t \n"))
	dummy.stdoutWriter.Write([]byte("id name \t Stockfish 17.1 author \t stillInName\t \n"))

	parsedCommand := client.commandBuf.Next().(idCommand)

	expected := idCommand{
		idt:   name,
		value: "Stockfish 17.1 author \t stillInName",
	}

	if parsedCommand != expected {
		t.Errorf("IDs do not match: expected %v, got %v", expected, parsedCommand)
	}
}

func TestUciokParsing(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("uciok\n"))
	dummy.stdoutWriter.Write([]byte(" \tuciok\t\t fdfj fdk\n"))

	parsedCommand1 := client.commandBuf.Next()
	parsedCommand2 := client.commandBuf.Next()

	if parsedCommand1.commandType() != uciok {
		t.Error("parsedCommand1 is not uciok")
	}
	if parsedCommand2.commandType() != uciok {
		t.Error("parsedCommand2 is not uciok")
	}
}

func TestReadyokParsing(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("readyok\n"))
	dummy.stdoutWriter.Write([]byte(" \treadyok\t\t fdfj fdk\n"))

	parsedCommand1 := client.commandBuf.Next()
	parsedCommand2 := client.commandBuf.Next()

	if parsedCommand1.commandType() != readyok {
		t.Error("parsedCommand1 is not readyok")
	}
	if parsedCommand2.commandType() != readyok {
		t.Error("parsedCommand2 is not readyok")
	}
}

func TestBestMoveParsing(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("bestmove e2e4q\n"))

	parsedCommand1 := client.commandBuf.Next().(bestMove)

	expectedBest := chess.Move{
		FromSquare: chess.E2,
		ToSquare:   chess.E4,
		Promotion:  chess.Queen,
	}

	if parsedCommand1.best != expectedBest {
		t.Errorf("did not get expectedBest: expected %v, got %v", expectedBest, parsedCommand1.best)
	}
	if parsedCommand1.ponder != nil {
		t.Error("ponder was not nil")
	}
}

func TestBestMoveParsing_Ponder(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("bestmove e2e4 ponder a1a3r\n"))

	parsedCommand1 := client.commandBuf.Next().(bestMove)

	expectedPonder := chess.Move{
		FromSquare: chess.A1,
		ToSquare:   chess.A3,
		Promotion:  chess.Rook,
	}

	if *parsedCommand1.ponder != expectedPonder {
		t.Errorf("did not get expectedPonder: expected %v, got %v", expectedPonder, parsedCommand1.ponder)
	}
}

func TestBestMoveParsing_InputOutOfOrder(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("bestmove ponder a1a3r e2e4q\n"))
	dummy.stdoutWriter.Write([]byte("bestmove d8d7\n"))

	parsedCommand1 := client.commandBuf.Next().(bestMove)

	expectedBest := chess.Move{
		FromSquare: chess.D8,
		ToSquare:   chess.D7,
		Promotion:  chess.NoPieceType,
	}

	if parsedCommand1.best != expectedBest {
		t.Errorf("did not get expectedBest: expected %v, got %v", expectedBest, parsedCommand1.best)
	}
}

func TestBestMoveParsing_InvalidMove(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("bestmove z2h3 ponder a1a2\n"))
	dummy.stdoutWriter.Write([]byte("bestmove a1a2 ponder f3\n"))
	dummy.stdoutWriter.Write([]byte("bestmove d8d7\n"))

	parsedCommand1 := client.commandBuf.Next().(bestMove)
	parsedCommand2 := client.commandBuf.Next().(bestMove)

	expectedBest := chess.Move{
		FromSquare: chess.A1,
		ToSquare:   chess.A2,
		Promotion:  chess.NoPieceType,
	}

	if parsedCommand1.best != expectedBest {
		t.Errorf("did not get expectedBest from parsedCommand1: expected %v, got %v", expectedBest, parsedCommand1.best)
	}

	if parsedCommand1.ponder != nil {
		t.Error("expected ponder to be nil")
	}

	expectedBest = chess.Move{
		FromSquare: chess.D8,
		ToSquare:   chess.D7,
		Promotion:  chess.NoPieceType,
	}

	if parsedCommand2.best != expectedBest {
		t.Errorf("did not get expectedBest from parsedCommand2: expected %v, got %v", expectedBest, parsedCommand1.best)
	}
}

func TestBestMoveParsing_RandomWhiteSpace(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("\tbestmove \t d8d7     ponder     \ta1a2\n"))

	parsedCommand1 := client.commandBuf.Next().(bestMove)

	expectedBest := chess.Move{
		FromSquare: chess.D8,
		ToSquare:   chess.D7,
		Promotion:  chess.NoPieceType,
	}
	expectedPonder := chess.Move{
		FromSquare: chess.A1,
		ToSquare:   chess.A2,
		Promotion:  chess.NoPieceType,
	}

	if parsedCommand1.best != expectedBest {
		t.Errorf("did not get expectedBest: expected %v, got %v", expectedBest, parsedCommand1.best)
	}
	if *parsedCommand1.ponder != expectedPonder {
		t.Errorf("did not get expectedPonder: expected %v, got %v", expectedPonder, parsedCommand1.ponder)
	}
}

func TestCopyProtectionParsing(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("copyprotection checking\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpChecking {
		t.Error("status not set to CpChecking")
	}

	dummy.stdoutWriter.Write([]byte("copyprotection ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpOk {
		t.Error("status not set to CpOk")
	}

	dummy.stdoutWriter.Write([]byte("copyprotection error\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpError {
		t.Error("status not set to CpError")
	}
}

func TestCopyProtectionBeforeInit(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if client.CopyrightStatus() != CpUnknown {
		t.Error("wrong cp status before init")
	}
}

func TestCopyProtectionParsing_WithGibberish(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("copyprotection fdkd checking\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpChecking {
		t.Error("status not set to CpChecking")
	}

	dummy.stdoutWriter.Write([]byte("s ddf copyprotection dfdf fd fdf fd ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpOk {
		t.Error("status not set to CpOk")
	}

	dummy.stdoutWriter.Write([]byte("    copyprotection\t   error ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpError {
		t.Error("status not set to CpError")
	}
}

func TestCopyProtectionParsing_BadInput(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("copyprotection ok\n"))
	dummy.stdoutWriter.Write([]byte("copyprotection\n"))
	time.Sleep(10 * time.Millisecond)
	if client.CopyrightStatus() != CpOk {
		t.Error("incomplete copyprotection overwrote CpOk")
	}
}

func TestRegistrationCommandParsing(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("registration checking\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegChecking {
		t.Error("status not set to RegChecking")
	}
	dummy.stdoutWriter.Write([]byte("registration ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegOk {
		t.Error("status not set to RegOk")
	}
	dummy.stdoutWriter.Write([]byte("registration error\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegError {
		t.Error("status not set to RegError")
	}
}

func TestRegistrationBeforeInit(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if client.RegistrationStatus() != RegUnknown {
		t.Error("wrong reg status before init")
	}
}

func TestRegistrationCommandParsing_WithNoise(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("registration gibberish checking\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegChecking {
		t.Error("status not set to RegChecking")
	}
	dummy.stdoutWriter.Write([]byte("random registration text ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegOk {
		t.Error("status not set to RegOk")
	}
	dummy.stdoutWriter.Write([]byte("\t registration  \t error ok\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegError {
		t.Error("status not set to RegError")
	}
}

func TestRegistrationCommandParsing_IncompleteInput(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("registration ok\n"))
	dummy.stdoutWriter.Write([]byte("registration\n"))
	time.Sleep(10 * time.Millisecond)
	if client.RegistrationStatus() != RegOk {
		t.Error("incomplete registration overwrote RegOk")
	}
}
