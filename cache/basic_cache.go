package cache

import "sync"

// Config for [BasicCache], reminder that this is a initialization config and when passed the value is set permanently for the instance of [BasicCache]
type BasicCacheConfig struct {
	InitialMapCapacity int // the initial map capacity used to make the map using make(map[K]V, InitialMapCapacity)
}

var DefaultBasicCacheConfig = BasicCacheConfig{} // Default config where the InitialMapCapacity is set to 0 same as the default map initiation in go

// Basic Cache implementation with generic map and read write mutex
type BasicCache[K k, V v] struct {
	m  map[K]V
	mu sync.RWMutex
}

func newBasicCache[K k, V v](config BasicCacheConfig) *BasicCache[K, V] {
	return &BasicCache[K, V]{
		m: make(map[K]V, config.InitialMapCapacity),
	}
}

// Instantiates a new cache using the generic types provided to this constructor function and default [DefaultCacheConfig]
func NewBasicCache[K comparable, V any]() *BasicCache[K, V] {
	return newBasicCache[K, V](DefaultBasicCacheConfig)
}

// Instantiates a new cache using the config provided and generics
func NewCacheWithConfig[K comparable, V any](config BasicCacheConfig) *BasicCache[K, V] {
	return newBasicCache[K, V](config)
}

// Adds a key value pair to the cache
func (c *BasicCache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Gets a value with the key provided, returning bool to convey wether the key exists
func (c *BasicCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

// Updates a value with the key provided. Does nothing if key doesn't exist
func (c *BasicCache[K, V]) Update(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Deletes a value with the key provided. Does nothing if key doesn't exist
func (c *BasicCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

// Loops over all the items in the list and passes the key and value to the function provided
//
// the function must return true to stop the iteration
//   - return false to continue the loop
//   - return true to exit the loop
//
// the read lock unlocks itself only after this function call has ended, do not run any other methods on this instance of [BasicCache] inside the provided function fn
func (c *BasicCache[K, V]) LoopFunc(fn func(key K, value V) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.m {
		if fn(k, v) {
			return
		}
	}
}
