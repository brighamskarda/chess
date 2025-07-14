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
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/brighamskarda/chess/v2"
)

func (info *Info) equals(other *Info) bool {
	if info == nil || other == nil {
		return info == other
	}

	compareUint := func(x, y *uint) bool {
		if x == nil || y == nil {
			return x == y
		}
		return *x == *y
	}

	compareInt := func(x, y *int) bool {
		if x == nil || y == nil {
			return x == y
		}
		return *x == *y
	}

	compareString := func(x, y *string) bool {
		if x == nil || y == nil {
			return x == y
		}
		return *x == *y
	}

	compareMove := func(x, y *chess.Move) bool {
		if x == nil || y == nil {
			return x == y
		}
		return *x == *y
	}

	compareScore := func(x, y *Score) bool {
		if x == nil || y == nil {
			return x == y
		}
		return compareInt(x.Cp, y.Cp) &&
			compareInt(x.Mate, y.Mate) &&
			x.Lowerbound == y.Lowerbound &&
			x.Upperbound == y.Upperbound
	}

	compareCurrline := func(x, y *Currline) bool {
		if x == nil || y == nil {
			return x == y
		}
		if !compareUint(x.Cpunr, y.Cpunr) {
			return false
		}
		if len(x.Moves) != len(y.Moves) {
			return false
		}
		for i := range x.Moves {
			if x.Moves[i] != y.Moves[i] {
				return false
			}
		}
		return true
	}

	compareMoves := func(x, y []chess.Move) bool {
		if len(x) != len(y) {
			return false
		}
		for i := range x {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	}

	return compareUint(info.Depth, other.Depth) &&
		compareUint(info.Seldepth, other.Seldepth) &&
		compareUint(info.Time, other.Time) &&
		compareUint(info.Nodes, other.Nodes) &&
		compareMoves(info.Pv, other.Pv) &&
		compareUint(info.Multipv, other.Multipv) &&
		compareScore(info.Score, other.Score) &&
		compareMove(info.Currmove, other.Currmove) &&
		compareUint(info.Currmovenumber, other.Currmovenumber) &&
		compareUint(info.Hashfull, other.Hashfull) &&
		compareUint(info.Nps, other.Nps) &&
		compareUint(info.Tbhits, other.Tbhits) &&
		compareUint(info.CpuLoad, other.CpuLoad) &&
		compareString(info.String, other.String) &&
		compareMoves(info.Refutation, other.Refutation) &&
		compareCurrline(info.Currline, other.Currline)
}

func (info *Info) string() string {
	if info == nil {
		return "<nil Info>"
	}

	derefUint := func(p *uint) string {
		if p == nil {
			return "nil"
		}
		return fmt.Sprintf("%d", *p)
	}

	derefInt := func(p *int) string {
		if p == nil {
			return "nil"
		}
		return fmt.Sprintf("%d", *p)
	}

	derefString := func(p *string) string {
		if p == nil {
			return "nil"
		}
		return *p
	}

	derefMove := func(p *chess.Move) string {
		if p == nil {
			return "nil"
		}
		return p.String()
	}

	formatMoves := func(moves []chess.Move) string {
		s := make([]string, len(moves))
		for i, m := range moves {
			s[i] = m.String()
		}
		return "[" + strings.Join(s, " ") + "]"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Info{Depth: %s, Seldepth: %s, Time: %s, Nodes: %s, ",
		derefUint(info.Depth), derefUint(info.Seldepth), derefUint(info.Time), derefUint(info.Nodes))

	fmt.Fprintf(&b, "Pv: %s, Multipv: %s, Score: {",
		formatMoves(info.Pv), derefUint(info.Multipv))
	if info.Score != nil {
		fmt.Fprintf(&b, "Cp: %s, Mate: %s, Lowerbound: %v, Upperbound: %v",
			derefInt(info.Score.Cp), derefInt(info.Score.Mate), info.Score.Lowerbound, info.Score.Upperbound)
	} else {
		b.WriteString("<nil>")
	}
	b.WriteString("}, ")

	fmt.Fprintf(&b, "Currmove: %s, Currmovenumber: %s, Hashfull: %s, Nps: %s, ",
		derefMove(info.Currmove), derefUint(info.Currmovenumber), derefUint(info.Hashfull), derefUint(info.Nps))

	fmt.Fprintf(&b, "Tbhits: %s, Cpuload: %s, String: %q, ",
		derefUint(info.Tbhits), derefUint(info.CpuLoad), derefString(info.String))

	fmt.Fprintf(&b, "Refutation: %s, Currline: ", formatMoves(info.Refutation))
	if info.Currline != nil {
		fmt.Fprintf(&b, "{Cpunr: %s, Moves: %s}", derefUint(info.Currline.Cpunr), formatMoves(info.Currline.Moves))
	} else {
		b.WriteString("<nil>")
	}
	b.WriteString("}")

	return b.String()
}

func TestInfoParsing_Empty(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()
	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("info \n"))

	info := client.ReadInfo()

	if !info.equals(&Info{}) {
		t.Errorf("info not empty: %v", info.string())
	}
}

func TestInfoParsing_Blocks(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()
	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	done := make(chan struct{})

	go func() {
		client.ReadInfo()
		done <- struct{}{}
	}()

	select {
	case <-done:
		t.Fatal("did not block")
	case <-ctx.Done():
	}

	dummy.stdoutWriter.Write([]byte("info \n"))
	<-done
}

func TestInfoParsing_OverwriteOldInfos(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()
	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := range 129 {
		dummy.stdoutWriter.Write([]byte(fmt.Sprintf("info depth %v\n", i)))
	}

	time.Sleep(10 * time.Millisecond)

	info := client.ReadInfo()
	if info.Depth == nil || *info.Depth != 1 {
		if info.Depth != nil {
			t.Errorf("expected depth 1, got depth %v", *info.Depth)
		} else {
			t.Error("got nil depth")
		}
	}
}

func TestInfoParsing_Depth(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()
	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	dummy.stdoutWriter.Write([]byte("info depth 69\n"))

	info := client.ReadInfo()
	if info.Depth == nil {
		t.Fatal("got nil depth")
	}
	if *info.Depth != 69 {
		t.Errorf("expected depth 69, got depth %v", *info.Depth)
	}
}

func TestInfoParsing_Seldepth(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info seldepth 42\n"))
	info := client.ReadInfo()

	if info.Seldepth == nil {
		t.Fatal("expected Seldepth to be non-nil, got nil")
	}
	if *info.Seldepth != 42 {
		t.Errorf("expected Seldepth to be 42, got %d", *info.Seldepth)
	}
}

func TestInfoParsing_Time(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info time 1234\n"))
	info := client.ReadInfo()

	if info.Time == nil {
		t.Fatal("expected Time to be non-nil, got nil")
	}
	if *info.Time != 1234 {
		t.Errorf("expected Time to be 1234, got %d", *info.Time)
	}
}

func TestInfoParsing_Currmovenumber(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info currmovenumber 7\n"))
	info := client.ReadInfo()

	if info.Currmovenumber == nil {
		t.Fatal("expected Currmovenumber to be non-nil, got nil")
	}
	if *info.Currmovenumber != 7 {
		t.Errorf("expected Currmovenumber to be 7, got %d", *info.Currmovenumber)
	}
}

func TestInfoParsing_Hashfull(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info hashfull 888\n"))
	info := client.ReadInfo()

	if info.Hashfull == nil {
		t.Fatal("expected Hashfull to be non-nil, got nil")
	}
	if *info.Hashfull != 888 {
		t.Errorf("expected Hashfull to be 888, got %d", *info.Hashfull)
	}
}

func TestInfoParsing_CpuLoad(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info cpuload 960\n"))
	info := client.ReadInfo()

	if info.CpuLoad == nil {
		t.Fatal("expected CpuLoad to be non-nil, got nil")
	}
	if *info.CpuLoad != 960 {
		t.Errorf("expected CpuLoad to be 960, got %d", *info.CpuLoad)
	}
}

func TestInfoParsing_String(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info string thinking hard\n"))
	info := client.ReadInfo()

	if info.String == nil {
		t.Fatal("expected String to be non-nil, got nil")
	}
	if *info.String != "thinking hard" {
		t.Errorf("expected String to be 'thinking hard', got %q", *info.String)
	}
}

func TestInfoParsing_Nodes(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info nodes 424242\n"))
	info := client.ReadInfo()

	if info.Nodes == nil {
		t.Fatal("expected Nodes to be non-nil, got nil")
	}
	if *info.Nodes != 424242 {
		t.Errorf("expected Nodes to be 424242, got %d", *info.Nodes)
	}
}

func TestInfoParsing_Multipv(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info multipv 3\n"))
	info := client.ReadInfo()

	if info.Multipv == nil {
		t.Fatal("expected Multipv to be non-nil, got nil")
	}
	if *info.Multipv != 3 {
		t.Errorf("expected Multipv to be 3, got %d", *info.Multipv)
	}
}

func TestInfoParsing_ScoreCp(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info score cp 23\n"))
	info := client.ReadInfo()

	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	if info.Score.Cp == nil {
		t.Fatal("expected Score.Cp to be non-nil, got nil")
	}
	if *info.Score.Cp != 23 {
		t.Errorf("expected Score.Cp to be 23, got %d", *info.Score.Cp)
	}
}

func TestInfoParsing_ScoreMate(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info score mate -2\n"))
	info := client.ReadInfo()

	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	if info.Score.Mate == nil {
		t.Fatal("expected Score.Mate to be non-nil, got nil")
	}
	if *info.Score.Mate != -2 {
		t.Errorf("expected Score.Mate to be -2, got %d", *info.Score.Mate)
	}
}

func TestInfoParsing_ScoreLowerbound(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info score cp 30 lowerbound\n"))
	info := client.ReadInfo()

	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	if info.Score.Cp == nil || *info.Score.Cp != 30 {
		t.Errorf("expected Score.Cp to be 30, got %v", info.Score.Cp)
	}
	if info.Score.Lowerbound != true {
		t.Errorf("expected Score.Lowerbound to be true, got %v", info.Score.Lowerbound)
	}
	if info.Score.Upperbound != false {
		t.Errorf("expected Score.Upperbound to be false, got %v", info.Score.Upperbound)
	}
}

func TestInfoParsing_ScoreUpperbound(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info score cp 30 upperbound\n"))
	info := client.ReadInfo()

	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	if info.Score.Cp == nil || *info.Score.Cp != 30 {
		t.Errorf("expected Score.Cp to be 30, got %v", info.Score.Cp)
	}
	if info.Score.Upperbound != true {
		t.Errorf("expected Score.Upperbound to be true, got %v", info.Score.Upperbound)
	}
	if info.Score.Lowerbound != false {
		t.Errorf("expected Score.Lowerbound to be false, got %v", info.Score.Lowerbound)
	}
}

func TestInfoParsing_Currmove(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info currmove g1f3\n"))
	info := client.ReadInfo()

	if info.Currmove == nil {
		t.Fatal("expected Currmove to be non-nil, got nil")
	}
	expected := chess.Move{
		FromSquare: chess.G1,
		ToSquare:   chess.F3,
		Promotion:  chess.NoPieceType,
	}
	if *info.Currmove != expected {
		t.Errorf("expected Currmove to be %v, got %v", expected, *info.Currmove)
	}
}

func TestInfoParsing_Nps(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info nps 314159\n"))
	info := client.ReadInfo()

	if info.Nps == nil {
		t.Fatal("expected Nps to be non-nil, got nil")
	}
	if *info.Nps != 314159 {
		t.Errorf("expected Nps to be 314159, got %d", *info.Nps)
	}
}

func TestInfoParsing_Tbhits(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info tbhits 12\n"))
	info := client.ReadInfo()

	if info.Tbhits == nil {
		t.Fatal("expected Tbhits to be non-nil, got nil")
	}
	if *info.Tbhits != 12 {
		t.Errorf("expected Tbhits to be 12, got %d", *info.Tbhits)
	}
}

func TestInfoParsing_Refutation(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info refutation d1h5 g6h5\n"))
	info := client.ReadInfo()

	expected := []chess.Move{
		{FromSquare: chess.D1, ToSquare: chess.H5, Promotion: chess.NoPieceType},
		{FromSquare: chess.G6, ToSquare: chess.H5, Promotion: chess.NoPieceType},
	}

	if info.Refutation == nil {
		t.Fatal("expected Refutation to be non-nil, got nil")
	}
	if !slices.Equal(info.Refutation, expected) {
		t.Errorf("expected Refutation %v, got %v", expected, info.Refutation)
	}
}

func TestInfoParsing_Currline(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info currline 2 e2e4 e7e5\n"))
	info := client.ReadInfo()

	if info.Currline == nil {
		t.Fatal("expected Currline to be non-nil, got nil")
	}

	if info.Currline.Cpunr == nil {
		t.Fatal("expected Currline.Cpunr to be non-nil, got nil")
	}
	if *info.Currline.Cpunr != 2 {
		t.Errorf("expected Currline.Cpunr to be 2, got %d", *info.Currline.Cpunr)
	}

	expectedMoves := []chess.Move{
		{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType},
		{FromSquare: chess.E7, ToSquare: chess.E5, Promotion: chess.NoPieceType},
	}
	if !slices.Equal(info.Currline.Moves, expectedMoves) {
		t.Errorf("expected Currline.Moves %v, got %v", expectedMoves, info.Currline.Moves)
	}
}

func TestInfoParsing_AllFields(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info depth 20 seldepth 25 time 987 nodes 123456 pv e2e4 e7e5 multipv 2 score cp 50 mate -1 lowerbound upperbound currmove d2d4 currmovenumber 3 hashfull 999 nps 3000000 tbhits 15 cpuload 850 refutation f2f4 e5f4 currline 4 b1c3 g8f6 string evaluating king safety\n"))

	info := client.ReadInfo()
	assertAllInfoFields(t, info)
}

func assertAllInfoFields(t *testing.T, info *Info) {
	// Scalars
	assertUint(t, "Depth", info.Depth, 20)
	assertUint(t, "Seldepth", info.Seldepth, 25)
	assertUint(t, "Time", info.Time, 987)
	assertUint(t, "Nodes", info.Nodes, 123456)
	assertUint(t, "Multipv", info.Multipv, 2)
	assertUint(t, "Currmovenumber", info.Currmovenumber, 3)
	assertUint(t, "Hashfull", info.Hashfull, 999)
	assertUint(t, "Nps", info.Nps, 3000000)
	assertUint(t, "Tbhits", info.Tbhits, 15)
	assertUint(t, "Cpuload", info.CpuLoad, 850)

	// Score
	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	assertInt(t, "Score.Cp", info.Score.Cp, 50)
	assertInt(t, "Score.Mate", info.Score.Mate, -1)
	if !info.Score.Lowerbound {
		t.Fatal("expected Lowerbound to be true")
	}
	if !info.Score.Upperbound {
		t.Fatal("expected Upperbound to be true")
	}

	// Pv
	expectedPv := []chess.Move{
		{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType},
		{FromSquare: chess.E7, ToSquare: chess.E5, Promotion: chess.NoPieceType},
	}
	if !slices.Equal(info.Pv, expectedPv) {
		t.Errorf("expected Pv %v, got %v", expectedPv, info.Pv)
	}

	// Currmove
	expectedCurr := chess.Move{FromSquare: chess.D2, ToSquare: chess.D4, Promotion: chess.NoPieceType}
	if info.Currmove == nil || *info.Currmove != expectedCurr {
		t.Errorf("expected Currmove %v, got %v", expectedCurr, info.Currmove)
	}

	// String
	if info.String == nil || *info.String != "evaluating king safety" {
		t.Errorf("expected String to be 'evaluating king safety', got %v", info.String)
	}

	// Refutation
	expectedRef := []chess.Move{
		{FromSquare: chess.F2, ToSquare: chess.F4, Promotion: chess.NoPieceType},
		{FromSquare: chess.E5, ToSquare: chess.F4, Promotion: chess.NoPieceType},
	}
	if !slices.Equal(info.Refutation, expectedRef) {
		t.Errorf("expected Refutation %v, got %v", expectedRef, info.Refutation)
	}

	// Currline
	if info.Currline == nil {
		t.Fatal("expected Currline to be non-nil, got nil")
	}
	assertUint(t, "Currline.Cpunr", info.Currline.Cpunr, 4)
	expectedLine := []chess.Move{
		{FromSquare: chess.B1, ToSquare: chess.C3, Promotion: chess.NoPieceType},
		{FromSquare: chess.G8, ToSquare: chess.F6, Promotion: chess.NoPieceType},
	}
	if !slices.Equal(info.Currline.Moves, expectedLine) {
		t.Errorf("expected Currline.Moves %v, got %v", expectedLine, info.Currline.Moves)
	}
}

func assertUint(t *testing.T, name string, actual *uint, expected uint) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected %s to be non-nil, got nil", name)
	}
	if *actual != expected {
		t.Errorf("expected %s to be %d, got %d", name, expected, *actual)
	}
}

func assertInt(t *testing.T, name string, actual *int, expected int) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected %s to be non-nil, got nil", name)
	}
	if *actual != expected {
		t.Errorf("expected %s to be %d, got %d", name, expected, *actual)
	}
}

func TestInfoParsing_invalidBeforeInfo(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("nonsense info depth 8\n"))
	info := client.ReadInfo()

	if info.Depth == nil {
		t.Fatal("expected Depth to be non-nil, got nil")
	}
	if *info.Depth != 8 {
		t.Errorf("expected Depth to be 8, got %d", *info.Depth)
	}
}

func TestInfoParsing_invalidAfterValue(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info depth 5 unexpected\n"))
	info := client.ReadInfo()

	if info.Depth == nil {
		t.Fatal("expected Depth to be non-nil, got nil")
	}
	if *info.Depth != 5 {
		t.Errorf("expected Depth to be 5, got %d", *info.Depth)
	}
}

func TestInfoParsing_invalidAfterInfo(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info unexpected depth 5\n"))
	info := client.ReadInfo()

	if info.Depth == nil {
		t.Fatal("expected Depth to be non-nil, got nil")
	}
	if *info.Depth != 5 {
		t.Errorf("expected Depth to be 5, got %d", *info.Depth)
	}
}

func TestInfoParsing_WhitespaceBetweenTokens(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Extra tabs/spaces between keywords and values
	dummy.stdoutWriter.Write([]byte("info   depth\t\t42   seldepth   \t  18  nodes\t999999 \n"))
	info := client.ReadInfo()

	if info.Depth == nil {
		t.Fatal("expected Depth to be non-nil, got nil")
	}
	if *info.Depth != 42 {
		t.Errorf("expected Depth to be 42, got %d", *info.Depth)
	}

	if info.Seldepth == nil {
		t.Fatal("expected Seldepth to be non-nil, got nil")
	}
	if *info.Seldepth != 18 {
		t.Errorf("expected Seldepth to be 18, got %d", *info.Seldepth)
	}

	if info.Nodes == nil {
		t.Fatal("expected Nodes to be non-nil, got nil")
	}
	if *info.Nodes != 999999 {
		t.Errorf("expected Nodes to be 999999, got %d", *info.Nodes)
	}
}

func TestInfoParsing_WhitespaceAroundScore(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info   score\t cp\t 77   mate   -3\tlowerbound\t   upperbound\t\n"))
	info := client.ReadInfo()

	if info.Score == nil {
		t.Fatal("expected Score to be non-nil, got nil")
	}
	if info.Score.Cp == nil || *info.Score.Cp != 77 {
		t.Errorf("expected Score.Cp to be 77, got %v", info.Score.Cp)
	}
	if info.Score.Mate == nil || *info.Score.Mate != -3 {
		t.Errorf("expected Score.Mate to be -3, got %v", info.Score.Mate)
	}
	if info.Score.Lowerbound != true {
		t.Errorf("expected Score.Lowerbound to be true, got %v", info.Score.Lowerbound)
	}
	if info.Score.Upperbound != true {
		t.Errorf("expected Score.Upperbound to be true, got %v", info.Score.Upperbound)
	}
}

func TestInfoParsing_WhitespaceInPv(t *testing.T) {
	dummy := newDummyClientProgram()
	defer dummy.Kill()

	client, err := newClientFromClientProgram(dummy, ClientSettings{})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dummy.stdoutWriter.Write([]byte("info  pv   e2e4\t  e7e5   g1f3  \n"))
	info := client.ReadInfo()

	expected := []chess.Move{
		{FromSquare: chess.E2, ToSquare: chess.E4, Promotion: chess.NoPieceType},
		{FromSquare: chess.E7, ToSquare: chess.E5, Promotion: chess.NoPieceType},
		{FromSquare: chess.G1, ToSquare: chess.F3, Promotion: chess.NoPieceType},
	}
	if info.Pv == nil {
		t.Fatal("expected Pv to be non-nil, got nil")
	}
	if !slices.Equal(info.Pv, expected) {
		t.Errorf("expected Pv %v, got %v", expected, info.Pv)
	}
}
