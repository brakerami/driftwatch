// Package deduplicator suppresses duplicate drift findings within a
// configurable time window, preventing the same finding from being
// reported repeatedly across successive scan cycles.
package deduplicator

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Store tracks recently seen findings by a content-derived key.
type Store struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// New creates a Store that suppresses duplicate findings for the given
// window duration. Returns an error if window is non-positive.
func New(window time.Duration) (*Store, error) {
	if window <= 0 {
		return nil, fmt.Errorf("deduplicator: window must be positive, got %s", window)
	}
	return &Store{
		seen:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}, nil
}

// Filter returns only those findings that have not been seen within the
// configured window. Findings that pass through are recorded.
func (s *Store) Filter(findings []drift.Finding) []drift.Finding {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.nowFunc()
	s.evict(now)

	var out []drift.Finding
	for _, f := range findings {
		key := fingerprintFinding(f)
		if _, exists := s.seen[key]; !exists {
			s.seen[key] = now
			out = append(out, f)
		}
	}
	return out
}

// Len returns the number of entries currently tracked in the window.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.seen)
}

// evict removes entries whose window has expired. Caller must hold s.mu.
func (s *Store) evict(now time.Time) {
	for key, ts := range s.seen {
		if now.Sub(ts) >= s.window {
			delete(s.seen, key)
		}
	}
}

// fingerprintFinding produces a stable hash key for a Finding.
func fingerprintFinding(f drift.Finding) string {
	raw := fmt.Sprintf("%s|%s|%s|%s", f.Container, f.Type, f.Field, f.Expected)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}
