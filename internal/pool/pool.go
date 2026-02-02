// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package pool

import "sync"

// Pool is a type-safe wrapper around sync.Pool using generics.
// It provides compile-time type safety for pooled objects.
type Pool[T any] struct {
	pool sync.Pool
}

// New creates a new typed Pool with the given factory function.
// The factory function is called when the pool needs to create a new object.
func New[T any](factory func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return factory()
			},
		},
	}
}

// Get retrieves an object from the pool.
// If the pool is empty, it uses the factory function to create a new object.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put adds an object back to the pool.
// The object may be reused by future Get calls.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
