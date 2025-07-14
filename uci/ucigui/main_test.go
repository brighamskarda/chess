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
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

var dummyBinaryPath = "./testdata/dummy/dummy"

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		dummyBinaryPath += ".exe"
	}

	if _, err := os.Stat(dummyBinaryPath); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", dummyBinaryPath, "./testdata/dummy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build dummyuci: %v\n", err)
			os.Exit(1)
		}
	}

	code := m.Run()
	os.Exit(code)
}
