package cache

import "time"

// Cache with expiration time of items, automatically clears items after it has expired. Written on top of [BasicCache]
type TTLCache[K k, V v] struct {
	cache         Cacher[K, ttlItem[V]]
	passiveDelete bool
}

// Config for [TTLCache], reminder that this is a initialization config and when passed the value is set permanently for the instance of [TTLCache]
type TTLCacheConfig struct {
	CacheType     cache // BasicCache or SyncMapCache
	CacheConfig   any   // the config for the cache, can be nil for the types which dont have a separate config
	PassiveDelete bool  // Deletes an item on a Get call if it is expired, before the cleaner does it.
}

var DefaultTTLCacheConfig = TTLCacheConfig{BasicCacheType, DefaultBasicCacheConfig, true} // uses BasicCacheType for CacheType, DefaultBasicCacheConfig for CacheConfig and true for passive delete

// Wrapper for value V with a expiresAt field
type ttlItem[V any] struct {
	value     V
	expiresAt time.Time
}

func newTTLCache[K k, V v](config TTLCacheConfig) *TTLCache[K, V] {
	var cacher Cacher[K, ttlItem[V]]
	switch config.CacheType {
	case BasicCacheType:
		cacher = newBasicCache[K, ttlItem[V]](config.CacheConfig.(BasicCacheConfig))
	case SyncMapCacheType:
		cacher = NewSyncMapCache[K, ttlItem[V]]()
	}
	return &TTLCache[K, V]{cacher, config.PassiveDelete}
}

// Instantiates a new ttl cache using the generic types provided to this constructor function and [DefaultTTLCacheConfig]
func NewTTLCache[K k, V v](cleanupInterval time.Duration) *TTLCache[K, V] {
	tc := newTTLCache[K, V](DefaultTTLCacheConfig)
	go cleaner(tc, cleanupInterval)
	return tc
}

// Instantiates a new ttl cache with the config provided
func NewTTLCacheWithConfig[K k, V v](cleanupInterval time.Duration, config TTLCacheConfig) *TTLCache[K, V] {
	tc := newTTLCache[K, V](config)
	if config.CacheType == BasicCacheType {
		go cleaner(tc, cleanupInterval)
	}
	return tc
}

// Adds a key value pair to the cache
func (tc *TTLCache[K, V]) Add(key K, value V, expiresAt time.Time) {
	tc.cache.Add(key, ttlItem[V]{value, expiresAt})
}

// Gets a value with the key provided, returning bool to convey wether the key exists, and the expiresAt time
//
// If passive delete was enabled at startup then expired items are deleted automatically here
func (tc *TTLCache[K, V]) Get(key K) (V, bool, time.Time) {
	r, ok := tc.cache.Get(key)
	if ok && tc.passiveDelete && r.expiresAt.Before(time.Now()) {
		tc.Delete(key)
	}
	return r.value, ok, r.expiresAt
}

// Updates a value with the key provided. Does nothing if key doesn't exist
//
// expiresAt parameter is optional, supply an empty time.Time{} if you want time to not be updated
func (tc *TTLCache[K, V]) Update(key K, value V, expiresAt time.Time) {
	item := ttlItem[V]{value: value}
	if !expiresAt.IsZero() {
		item.expiresAt = expiresAt
	}
	tc.cache.Update(key, item)
}

// Deletes an item from the cache
func (tc *TTLCache[K, V]) Delete(key K) {
	tc.cache.Delete(key)
}

// Loops over all the items in the list and passes the key, value and expiresAt to the function provided
//
// the function must return true to stop the iteration
//   - return false to continue the loop
//   - return true to exit the loop
//
// the read lock unlocks itself only after this function call has ended, do not run any other methods on this instance of [TTLCache] inside the passed function fn
func (tc *TTLCache[K, V]) LoopFunc(fn func(key K, value V, expiresAt time.Time) bool) {
	tc.cache.LoopFunc(
		func(key K, value ttlItem[V]) bool {
			return fn(key, value.value, value.expiresAt)
		},
	)
}

// cleaner function to clean up expired items on a set time interval
//
// ! warning: only use with [BasicCache], otherwise panics, currently support for other caches are not present
func cleaner[K k, V v](tc *TTLCache[K, V], interval time.Duration) {
	basicCache, ok := tc.cache.(*BasicCache[K, ttlItem[V]])
	if !ok {
		panic("go-utils/cache ttl_cache.cleaner(): unsupported cache passed in for cleaning, only BasicCache is supported for cleaning.")
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		basicCache.mu.Lock()
		for k, v := range basicCache.m {
			if v.expiresAt.Before(now) {
				delete(basicCache.m, k)
			}
		}
		basicCache.mu.Unlock()
	}
}
