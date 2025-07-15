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

import "testing"

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
		t.Errorf("parsedCommand1 is not uciok")
	}
	if parsedCommand2.commandType() != uciok {
		t.Errorf("parsedCommand2 is not uciok")
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
		t.Errorf("parsedCommand1 is not readyok")
	}
	if parsedCommand2.commandType() != readyok {
		t.Errorf("parsedCommand2 is not readyok")
	}
}
