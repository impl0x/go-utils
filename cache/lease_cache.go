package cache

import (
	"sync"
)

// A type of cache where you can acquire objects
// instead of just reading them, useful when you need to
// make sure that only 1 goroutine can work on a item at once.
//
// Calling [LeaseCache.Acquire] will retrieve a value with
// its associated key and remove it from the pool so that no other
// goroutine can acquire it at the same time.
//
// Only when [LeaseCache.Release] or [LeasedItem.Release] is called
// is the item released back into the pool where it can be acquired again.
//
// If an item is not released by the acquirer then it can be considered as deleted.
type LeaseCache[K k, V v] struct {
	m  map[K]V
	mu sync.Mutex
}

func NewLeaseCache[K k, V v]() *LeaseCache[K, V] {
	return &LeaseCache[K, V]{m: make(map[K]V)}
}

// Represents a Leased Item from the pool cache
type LeasedItem[K k, V v] struct { // 10 + size(K) + size(V) bytes, need to use pointers to struct because mem size is unknown and could be larger than 100 bytes which is inefficient to copy for each method call
	parent       *LeaseCache[K, V]
	key          K
	value        V
	release      bool
	releaseCount uint8
}

// returns the key of the item
func (l *LeasedItem[K, V]) Key() K {
	return l.key
}

// returns the value of the item
func (l *LeasedItem[K, V]) Value() V {
	return l.value
}

// Releases the item back into the pool only if [LeasedItem.Delete] has
// not been called previously.
//
// Multiple calls to Release or calling Release after calling Delete does
// not affect it, items are only released once and are not released if
// Delete was called previously. Therefore this function is safe to defer
// in your own functions where you acquire a item and work with it.
// You can then later call [LeasedItem.Delete] if wanting to prevent release
// which makes the item be deleted.
//
// this method calls its parent [LeaseCache.Release] directly
func (l *LeasedItem[K, V]) Release() {
	if l.releaseCount == 0 && l.release {
		l.releaseCount++
		l.parent.Release(l.key, l.value)
	}
}

// Sets the internal flag of being released to false so that
// any calls to [LeasedItem.Release] after this has been called
// prevents the item from being released.
func (l *LeasedItem[K, V]) Delete() {
	l.release = false
}

// acquires the item with the associated key from the pool,
// if the item has currently been acquired then this will fail.
// Hence it will return a [LeasedItem] whose underlying value is nil
// and false for the second return value
func (m *LeaseCache[K, V]) Acquire(key K) (*LeasedItem[K, V], bool) {
	m.mu.Lock()
	v, ok := m.m[key]
	delete(m.m, key)
	m.mu.Unlock()
	if !ok {
		m = nil // looks pretty illegal to do but I think it'll end up fine
	}
	return &LeasedItem[K, V]{m, key, v, true, 0}, ok
}

// Releases the item associated with this key and value.
//
// Explained briefly in [LeasedItem.Release].
// Note: do not call this multiple times or defer this,
// the Release for LeasedItem and LeaseCache are different.
// This method is the internal implementation of releasing an item back to the pool
func (m *LeaseCache[K, V]) Release(key K, value V) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

// Adds a item to the pool cache, this is the absolute same as calling Release.
func (m *LeaseCache[K, V]) Add(key K, value V) {
	m.Release(key, value)
}

// Loops over all elements in the cache and passes it the function provided.
// This does not include currently acquired items, those will not be present.
//
// Note: the loop runs with a mutex, so no methods should be called
// inside the function fn, or it will lead to a deadlock because the
// methods want a lock acquire too, therefore blocking anyone from using
// the cache, so make sure the function provided is not taking too long
// to finish each iteration, and try to break early when done.
//
// returning true will break early, while false will continue the loop
func (m *LeaseCache[K, V]) LoopFunc(fn func(k K, v V) bool) {
	m.mu.Lock()
	for k, v := range m.m {
		if r := fn(k, v); r == true {
			break
		}
	}
	m.mu.Unlock()
}
