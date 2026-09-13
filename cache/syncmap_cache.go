package cache

import (
	"sync"
)

// Cache map implementation which uses sync.Map for its underlying map
type SyncMapCache[K k, V v] struct {
	m sync.Map
}

func NewSyncMapCache[K k, V v]()*SyncMapCache[K,V]{
	return &SyncMapCache[K, V]{}
}

func (sm *SyncMapCache[K, V]) Add(key K, value V) {
	sm.m.Store(key, value)
}

func (sm *SyncMapCache[K, V]) Get(key K) (V, bool) {
	v, ok := sm.m.Load(key)
	return v.(V), ok
}

func (sm *SyncMapCache[K, V]) Update(key K, value V) {
	sm.m.Swap(key, value)
}

func (sm *SyncMapCache[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

func (sm *SyncMapCache[K, V]) LoopFunc(fn func(key K, value V) bool) {
	sm.m.Range(func(key, value any) bool { return fn(key.(K), value.(V)) })
}
