package pool

import "sync"

// Generic wrapper over [sync.Pool],
// Methods use generic types therefore giving type safety at compile time.
//
// Can be used by directly using this struct and populating the Pool field
// with a sync.Pool instance with the new field having a function wrapper
//
// or use [NewPool] to make it easier
type Pool[T any] struct {
	Pool sync.Pool
}

func NewPool[T any](new func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any { return new() },
		},
	}
}

// Gets an item from the pool, returns variable of type T
func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

// Puts an item back into the pool, of type T
func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}
