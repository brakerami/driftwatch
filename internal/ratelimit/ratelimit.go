// Package ratelimit provides a simple per-container rate limiter to prevent
// alert flooding when the same drift findings are detected repeatedly.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per container and suppresses repeated
// alerts that occur within the configured cooldown window.
type Limiter struct {
	cooldown time.Duration
	mu       sync.Mutex
	lastSent map[string]time.Time
}

// New creates a Limiter with the given cooldown duration. Returns an error if
// cooldown is zero or negative.
func New(cooldown time.Duration) (*Limiter, error) {
	if cooldown <= 0 {
		return nil, ErrInvalidCooldown
	}
	return &Limiter{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}, nil
}

// Allow reports whether an alert for the given container should be sent.
// It returns true the first time a container is seen, and again only after
// the cooldown window has elapsed since the last allowed alert.
func (l *Limiter) Allow(containerID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	last, seen := l.lastSent[containerID]
	if !seen || time.Since(last) >= l.cooldown {
		l.lastSent[containerID] = time.Now()
		return true
	}
	return false
}

// Reset clears the rate-limit state for a specific container. Useful when a
// container is restarted or removed.
func (l *Limiter) Reset(containerID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.lastSent, containerID)
}

// Purge removes all entries older than the cooldown window to keep memory
// usage bounded during long-running daemon operation.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for id, t := range l.lastSent {
		if now.Sub(t) >= l.cooldown {
			delete(l.lastSent, id)
		}
	}
}
