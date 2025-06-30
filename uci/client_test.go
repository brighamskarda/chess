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

package uci

import "testing"

func TestNewClient_ErrorOnInvalidBinary(t *testing.T) {
	_, err := NewClient("./dkfdks.exe", ClientSettings{})
	if err == nil {
		t.Error("did not get error on invalid binary")
	}
}

func TestNewClient_NoErrorOnValidBinary(t *testing.T) {
	_, err := NewClient(dummyBinaryPath, ClientSettings{})
	if err != nil {
		t.Errorf("%v", err)
	}
}
