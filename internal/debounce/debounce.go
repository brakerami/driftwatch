// Package debounce provides a mechanism to suppress repeated drift findings
// within a configurable time window, reducing alert noise for flapping containers.
package debounce

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// key uniquely identifies a finding by container name and drift type.
type key struct {
	container string
	driftType drift.Type
}

// entry records the last time a finding was forwarded.
type entry struct {
	lastSent time.Time
}

// Debouncer suppresses duplicate findings that occur within a quiet window.
type Debouncer struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[key]entry
	now    func() time.Time
}

// New creates a Debouncer with the given quiet window.
// Returns an error if window is zero or negative.
func New(window time.Duration) (*Debouncer, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Debouncer{
		window: window,
		seen:   make(map[key]entry),
		now:    time.Now,
	}, nil
}

// Filter returns only those findings that have not been seen within the quiet
// window. Findings that pass through update the internal timestamp.
func (d *Debouncer) Filter(findings []drift.Finding) []drift.Finding {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	var out []drift.Finding

	for _, f := range findings {
		k := key{container: f.Container, driftType: f.Type}
		e, seen := d.seen[k]
		if !seen || now.Sub(e.lastSent) >= d.window {
			d.seen[k] = entry{lastSent: now}
			out = append(out, f)
		}
	}
	return out
}

// Flush removes all tracked entries, resetting the debounce state.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[key]entry)
}
