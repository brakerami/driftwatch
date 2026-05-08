// Package cache implements a lightweight, generic, TTL-based in-memory cache
// used by driftwatch to avoid re-inspecting containers that were already
// evaluated within the current scan window.
//
// Entries expire automatically based on the TTL supplied at construction time.
// Callers may also explicitly delete entries or purge all expired keys in bulk.
//
// The cache is safe for concurrent use by multiple goroutines.
package cache
