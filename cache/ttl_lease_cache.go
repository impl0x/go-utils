package cache

import "time"

// same as [LeaseCache], but with a time to live,
// that is, items expire after the expiry runs out.
type TTLLeaseCache[K k, V v] struct {
	lc LeaseCache[K, ttlItem[V]]
}

func NewTTLLeaseCache[K k, V v](cleanupInterval time.Duration) *TTLLeaseCache[K, V] {
	t := TTLLeaseCache[K, V]{
		LeaseCache[K, ttlItem[V]]{m: make(map[K]ttlItem[V])},
	}
	go t.ttlLeaseCacheCleaner(cleanupInterval)
	return &t
}

func (t *TTLLeaseCache[K, V]) ttlLeaseCacheCleaner(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	for range ticker.C {
		t.lc.mu.Lock()
		for k, v := range t.lc.m {
			if v.expiresAt.Before(time.Now()) {
				delete(t.lc.m, k)
			}
		}
		t.lc.mu.Unlock()
	}
}

// Represents a time to live Leased Item from the pool cache
type TTLLeasedItem[K k, V v] struct {
	parent   *TTLLeaseCache[K, V]
	key      K
	ttlValue ttlItem[V]
}

// returns the key of the item
func (t TTLLeasedItem[K, V]) Key() K {
	return t.key
}

// returns the value of the item
func (t TTLLeasedItem[K, V]) Value() V {
	return t.ttlValue.value
}

// returns when the item will expire
func (t TTLLeasedItem[K, V]) ExpiresAt() time.Time {
	return t.ttlValue.expiresAt
}

// returns if a leased item's underlying value is nil
func (l TTLLeasedItem[K, V]) Nil() bool {
	return l.parent == nil
}

// releases the item back into the pool,
// attempting to acquire this item after release has been
// called will result it in not failing and returning the item
//
// this method calls its parent [TTLLeaseCache.Release] directly
func (t TTLLeasedItem[K, V]) Release() {
	if t.Nil() {
		return
	}
	t.parent.Release(t.key, t.ttlValue.value, t.ttlValue.expiresAt)
}

// Releases the item associated with this key and value.
//
// Explained briefly in [TTLLeasedItem.Release]
func (t *TTLLeaseCache[K, V]) Release(key K, value V, expiresAt time.Time) {
	// we cannot add a expired item back into the cache
	if expiresAt.Before(time.Now()) {
		return
	}
	t.lc.Release(key, ttlItem[V]{value, expiresAt})
}

// acquires the item with the associated key from the pool,
// if the item has currently been acquired then this will fail.
// Hence it will return a [TTLLeasedItem] whose underlying value is nil
// and false for the second return value
func (t *TTLLeaseCache[K, V]) Acquire(key K) (TTLLeasedItem[K, V], bool) {
	li, ok := t.lc.Acquire(key)
	if !ok {
		t = nil
	}
	return TTLLeasedItem[K, V]{t, li.key, li.value}, ok
}

// Adds a item to the pool cache, this is the absolute same as calling Release.
func (t *TTLLeaseCache[K, V]) Add(key K, value V, expiresAt time.Time) {
	t.Release(key, value, expiresAt)
}
