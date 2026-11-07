package cache

import "sync"

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

// Represents a Leased Item from the pool cache
type LeasedItem[K k, V v] struct {
	parent *LeaseCache[K, V]
	key    K
	value  V
}

// returns the key of the item
func (l LeasedItem[K, V]) Key() K {
	return l.key
}

// returns the value of the item
func (l LeasedItem[K, V]) Value() V {
	return l.value
}

// returns if a leased item's underlying value is nil
func (l LeasedItem[K, V]) Nil() bool {
	return l.parent == nil
}

// releases the item back into the pool,
// attempting to acquire this item after release has been
// called will result it in not failing and returning the item
//
// this method calls its parent [LeaseCache.Release] directly
func (l LeasedItem[K, V]) Release() {
	l.parent.Release(l.key, l.value)
}

// acquires the item with the associated key from the pool,
// if the item has currently been acquired then this will fail.
// Hence it will return a [LeasedItem] whose underlying value is nil
// and false for the second return value
func (m *LeaseCache[K, V]) Acquire(key K) (LeasedItem[K, V], bool) {
	m.mu.Lock()
	v, ok := m.m[key]
	delete(m.m, key)
	m.mu.Unlock()
	if !ok {
		m = nil // looks pretty illegal to do but I think it'll end up fine
	}
	return LeasedItem[K, V]{m, key, v}, ok
}

// Releases the item associated with this key and value.
//
// Explained briefly in [LeasedItem.Release]
func (m *LeaseCache[K, V]) Release(key K, value V) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

// Adds a item to the pool cache, this is the absolute same as calling Release.
func (m *LeaseCache[K, V]) Add(key K, value V) {
	m.Release(key, value)
}
