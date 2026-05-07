// Package throttle provides a token-bucket style throttle that limits
// how many drift-check cycles can be dispatched within a rolling window.
package throttle

import (
	"errors"
	"sync"
	"time"
)

// ErrThrottled is returned when the caller exceeds the allowed burst.
var ErrThrottled = errors.New("throttle: request rate exceeded")

// Throttle enforces a maximum number of allowed calls within a window.
type Throttle struct {
	mu        sync.Mutex
	window    time.Duration
	maxTokens int
	tokens    int
	resetAt   time.Time
	now       func() time.Time
}

// New creates a Throttle that permits at most maxTokens calls per window.
func New(window time.Duration, maxTokens int) (*Throttle, error) {
	if window <= 0 {
		return nil, errors.New("throttle: window must be positive")
	}
	if maxTokens <= 0 {
		return nil, errors.New("throttle: maxTokens must be positive")
	}
	now := time.Now()
	return &Throttle{
		window:    window,
		maxTokens: maxTokens,
		tokens:    maxTokens,
		resetAt:   now.Add(window),
		now:       time.Now,
	}, nil
}

// Allow returns nil if the call is permitted, or ErrThrottled if the
// budget for the current window is exhausted.
func (t *Throttle) Allow() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if now.After(t.resetAt) {
		t.tokens = t.maxTokens
		t.resetAt = now.Add(t.window)
	}

	if t.tokens <= 0 {
		return ErrThrottled
	}
	t.tokens--
	return nil
}

// Remaining returns the number of tokens left in the current window.
func (t *Throttle) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if now.After(t.resetAt) {
		return t.maxTokens
	}
	return t.tokens
}

// ResetAt returns the time at which the current window expires.
func (t *Throttle) ResetAt() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.resetAt
}
