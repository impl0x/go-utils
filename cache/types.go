package cache

type k comparable // type for the key, it should be comparable as map keys must be comparable
type v any        // type for the value, it can be anything

// interface for caching
type Cacher[K k, V v] interface {
	// Adds a key value pair to the cache
	Add(key K, value V)
	// Gets a value with the key provided, returning bool to convey wether the key exists
	Get(key K) (V, bool)
	// Updates a value with the key provided. Does nothing if key doesn't exist
	Update(key K, value V)
	// Deletes a value with the key provided. Does nothing if key doesn't exist
	Delete(key K)
	// Loops over all the items in the list and passes the key and value to the function provided
	//
	// the function must return true to stop the iteration
	//   - return false to continue the loop
	//   - return true to exit the loop
	//
	// the exact locking and mutex implementation depends on the type implementing this interface
	LoopFunc(fn func(key K, value V) bool)
}

type cache int

const (
	BasicCacheType cache = iota
	SyncMapCacheType
)
