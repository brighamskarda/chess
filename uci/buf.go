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

import "sync"

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
	mu     sync.Mutex
	cond   *sync.Cond
	buffer []T
}

func newConcBuf[T any]() *concurrentBuf[T] {
	cb := &concurrentBuf[T]{}
	cb.cond = sync.NewCond(&cb.mu)
	return cb
}

func (cb *concurrentBuf[T]) Push(t T) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.buffer = append(cb.buffer, t)
	cb.cond.Signal()
}

func (cb *concurrentBuf[T]) Next() T {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	for len(cb.buffer) == 0 {
		cb.cond.Wait()
	}

	value := cb.buffer[0]
	cb.buffer = cb.buffer[1:]
	return value
}
