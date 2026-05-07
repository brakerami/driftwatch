// Package metrics exposes lightweight runtime counters for driftwatch.
// Counters are safe for concurrent use and can be serialised to a simple
// key/value map for reporting or Prometheus-style scraping.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing uint64 counter.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.value, 1) }

// Add increments the counter by n.
func (c *Counter) Add(n uint64) { atomic.AddUint64(&c.value, n) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.value) }

// Registry holds a named set of counters.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{counters: make(map[string]*Counter)}
}

// Counter returns the named counter, creating it if it does not exist.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a point-in-time copy of all counter values.
func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]uint64, len(r.counters))
	for name, c := range r.counters {
		out[name] = c.Value()
	}
	return out
}

// Reset sets every counter back to zero.
func (r *Registry) Reset() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.counters {
		atomic.StoreUint64(&c.value, 0)
	}
}
