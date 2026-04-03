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

// Optional is a generic value based optional type specialized for uci chess engines.
//
// By being value based this type introduces zero garbage collection overhead which is ideal for high perfomance scenarios. Consequently, performance may suffer if the type being passed in is very large.
//
// In many cases a zero value or nil should be used instead. But for situations requiring lots of optional integers this struct may be appropriate.
//
// The zero value is usable as an empty optional.
type Optional[T any] struct {
	value   T
	present bool
}

// HasValue returns true if the optional is not empty.
func (opt Optional[T]) HasValue() bool {
	return opt.present
}

// Value returns the value of the optional. Panics if the optional is empty.
func (opt Optional[T]) Value() T {
	if !opt.present {
		panic("invalid access of Value from empty optional")
	}
	return opt.value
}

// OptionalOf creates an optional with the given value.
func OptionalOf[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}

// OptionalEmpty creates an optional without a value.
func OptionalEmpty[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}
