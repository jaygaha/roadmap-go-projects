package cache

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// CacheEntry holds the data for a cached HTTP response
type CacheEntry struct {
	Headers      http.Header
	ResponseData []byte
	StatusCode   int
	CreatedAt    time.Time
}

var (
	// store is the in-memory cache store
	store            = make(map[string]CacheEntry) // the key is the request URL, the value is the CacheEntry struct
	mutex            = &sync.RWMutex{}             // it protects concurrent access to the store map
	internalCacheTTL time.Duration                 // it is set by Init and used for expiration checks
)

// Init initializes the cache with necessary configurations
func Init(ttl time.Duration) {
	internalCacheTTL = ttl
	// log.Println("Cache initialized with TTL:", internalCacheTTL)
}

// Set adds or updates a cache entry
func Set(key string, entry CacheEntry) {
	mutex.Lock() // Acquire the write lock
	defer mutex.Unlock()

	store[key] = entry
	// log.Printf("Cached response for %s (Status: %d)", key, entry.StatusCode)
}

// Get attemps to retrieve an entry from the cache
func Get(key string) (CacheEntry, bool) {
	mutex.RLock() // Acquire a read lock for safe concurrent access
	defer mutex.RUnlock()

	entry, found := store[key]
	if found && time.Since(entry.CreatedAt) > internalCacheTTL {
		log.Printf("Cache STALE for key: %s (older than %v)", key, internalCacheTTL)
		return CacheEntry{}, false // Entry found but expired
	}

	return entry, found
}

// Clear removes all entries from the cache
func Clear() {
	mutex.Lock()
	defer mutex.Unlock()

	store = make(map[string]CacheEntry) //Replace with a new empty map to clear the cache
	log.Println("Cache cleared.")
}
