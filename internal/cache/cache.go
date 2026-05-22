package cache

import (
	"sync"
	"time"
)

// Entry holds a cached index along with metadata about when it was built.
type Entry struct {
	FilePath    string
	FileSize    int64
	FileModTime time.Time
	BuiltAt     time.Time
	Payload     interface{}
}

// Cache is a simple in-memory store for index entries keyed by file path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a new Cache with the given TTL for entries.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Put stores an entry in the cache, replacing any existing entry for the key.
func (c *Cache) Put(key string, entry *Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry.BuiltAt = time.Now()
	c.entries[key] = entry
}

// Get retrieves a cache entry. Returns nil if not found or if the entry has
// expired, the file size changed, or the modification time changed.
func (c *Cache) Get(key string, currentSize int64, currentModTime time.Time) *Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.entries[key]
	if !ok {
		return nil
	}
	if time.Since(e.BuiltAt) > c.ttl {
		return nil
	}
	if e.FileSize != currentSize || !e.FileModTime.Equal(currentModTime) {
		return nil
	}
	return e
}

// Invalidate removes the entry for the given key, if present.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of entries currently held in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
