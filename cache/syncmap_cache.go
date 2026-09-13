package cache

import (
	"sync"
)

// Cache map implementation which uses sync.Map for its underlying map
type SyncMapCache[K k, V v] struct {
	m sync.Map
}

// Instantiates a new sync map cache using the generic types provided to this constructor function
func NewSyncMapCache[K k, V v]()*SyncMapCache[K,V]{
	return &SyncMapCache[K, V]{}
}

// Adds a key value pair to the map
func (sm *SyncMapCache[K, V]) Add(key K, value V) {
	sm.m.Store(key, value)
}

// Gets a value with the key provided, returning bool to convey wether the key exists
func (sm *SyncMapCache[K, V]) Get(key K) (V, bool) {
	v, ok := sm.m.Load(key)
	return v.(V), ok
}

// Updates a value with the key provided. Does nothing if key doesn't exist
func (sm *SyncMapCache[K, V]) Update(key K, value V) {
	sm.m.Swap(key, value)
}

// Deletes a value with the key provided. Does nothing if key doesn't exist
func (sm *SyncMapCache[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

// Loops over all the items in the list and passes the key and value to the function provided
//
// the function must return true to stop the iteration
//   - return false to continue the loop
//   - return true to exit the loop
//
// locking mechanism specified in [sync.Map.Range]
func (sm *SyncMapCache[K, V]) LoopFunc(fn func(key K, value V) bool) {
	sm.m.Range(func(key, value any) bool { return fn(key.(K), value.(V)) })
}
