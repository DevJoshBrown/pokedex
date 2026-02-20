package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time // a time.Time that reps when the entry was created
	val       []byte    // a []byte that represents the raw data we are caching from the pokeAPI.
}

type Cache struct {
	entries  map[string]cacheEntry // this is a dict (map) to store all the entries (cacheEntry) into the cache.
	mu       sync.Mutex            // a mutex to lock the Cache while mutable operations are happening on it
	interval time.Duration         // a variable to hold a time duration we can use as an interval
}

func (c *Cache) Add(key string, val []byte) {
	// Used to add a key value pair to the cache.

	c.mu.Lock()         // create a mutex lock on the cache
	defer c.mu.Unlock() // defer releasing the lock until the end of the function.

	// create a new cacheEntry - created at NOW, with the value passed into the function
	c_ent := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = c_ent // cache the new entry with the key of the URL in the URL request.
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// used to retrieve from the cache, the data using a url string as input

	c.mu.Lock()         // create a mutex lock on the cache
	defer c.mu.Unlock() // defer releasing the lock until the end of the function.

	// if the entry exists in the cache, return the entry value (the data) and true.
	entry, exists := c.entries[key]
	if exists {
		return entry.val, true
	}
	// otherwise return nil and false
	return nil, false
}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)
	for range ticker.C { // every tick do this...
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {

	newCache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go newCache.reapLoop()

	return &newCache
}
