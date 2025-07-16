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
	"context"
	"errors"
)

// concurrentCircBuf is a thread-safe circular buffer that overwrites old values,
// and blocks on Next() if nothing is available.
type concurrentCircBuf[T any] struct {
	contents chan T
}

func newCircBuf[T any](size int) *concurrentCircBuf[T] {
	return &concurrentCircBuf[T]{
		contents: make(chan T, size),
	}
}

func (cb *concurrentCircBuf[T]) Next() T {
	return <-cb.contents // blocks until something is available
}

func (cb *concurrentCircBuf[T]) Push(t T) {
	select {
	case cb.contents <- t:
		// success
	default:
		// channel is full, discard oldest
		<-cb.contents
		cb.contents <- t
	}
}

// concurrentBuf is a thread-safe buffer. Unlike circular buf it doesn't drop old values. It will just keep growing to accommodate unread values. Blocks next if nothing is available.
type concurrentBuf[T any] struct {
	inCh  chan T
	outCh chan T
}

func newConcBuf[T any](ctx context.Context) *concurrentBuf[T] {
	cb := &concurrentBuf[T]{
		inCh:  make(chan T),
		outCh: make(chan T),
	}
	go cb.run(ctx)
	return cb
}

func (cb *concurrentBuf[T]) run(ctx context.Context) {
	var buffer []T
	var outCh chan T
	var next T
	done := ctx.Done()

	for {
		// If we have something to send, prepare the output
		if len(buffer) > 0 {
			outCh = cb.outCh
			next = buffer[0]
		} else {
			outCh = nil // No value to send
		}

		select {
		case item := <-cb.inCh:
			buffer = append(buffer, item)
		case outCh <- next: // blocks forever if nil
			buffer = buffer[1:]
		case <-done:
			return
		}
	}
}

func (cb *concurrentBuf[T]) Push(t T) {
	cb.inCh <- t
}

func (cb *concurrentBuf[T]) Next() T {
	return <-cb.outCh
}

// NextWithContext blocks until the context expires, at which point it returns an error if no value was available.
func (cb *concurrentBuf[T]) NextWithContext(ctx context.Context) (T, error) {
	select {
	case val := <-cb.outCh:
		return val, nil
	case <-ctx.Done():
		var zero T
		return zero, errors.New("could not get next value, context expired")
	}
}
