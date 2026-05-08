// Package cache provides a short-lived in-memory result cache for drift
// findings, reducing redundant container inspections within a single scan cycle.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value along with its expiry time.
type Entry[V any] struct {
	Value     V
	ExpiresAt time.Time
}

// Expired reports whether the entry has passed its TTL.
func (e Entry[V]) Expired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// Cache is a generic TTL-based in-memory cache keyed by string.
type Cache[V any] struct {
	mu      sync.RWMutex
	items   map[string]Entry[V]
	ttl     time.Duration
	nowFunc func() time.Time
}

// New creates a Cache with the given TTL. Returns an error if TTL is <= 0.
func New[V any](ttl time.Duration) (*Cache[V], error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Cache[V]{
		items:   make(map[string]Entry[V]),
		ttl:     ttl,
		nowFunc: time.Now,
	}, nil
}

// Set stores a value under key, overwriting any existing entry.
func (c *Cache[V]) Set(key string, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry[V]{
		Value:     value,
		ExpiresAt: c.nowFunc().Add(c.ttl),
	}
}

// Get retrieves a value by key. Returns the zero value and false if the key
// is absent or the entry has expired.
func (c *Cache[V]) Get(key string) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || entry.Expired(c.nowFunc()) {
		var zero V
		return zero, false
	}
	return entry.Value, true
}

// Delete removes a single key from the cache.
func (c *Cache[V]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Purge removes all expired entries and returns the number removed.
func (c *Cache[V]) Purge() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.nowFunc()
	removed := 0
	for k, e := range c.items {
		if e.Expired(now) {
			delete(c.items, k)
			removed++
		}
	}
	return removed
}

// Len returns the total number of entries, including expired ones.
func (c *Cache[V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
