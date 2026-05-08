// Package circuitbreaker provides a simple circuit breaker that stops
// forwarding calls to a downstream dependency after a threshold of
// consecutive failures and automatically resets after a cooldown period.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are being rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota // normal operation
	StateOpen                // failing fast
	StateHalfOpen            // testing recovery
)

// Breaker is a circuit breaker instance.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	openedAt     time.Time
	successes    int
	probeSuccess int
}

// New creates a Breaker that trips after threshold consecutive failures and
// attempts recovery after cooldown. probeSuccess is the number of consecutive
// successes in half-open state required to close the circuit again.
func New(threshold int, cooldown time.Duration, probeSuccess int) (*Breaker, error) {
	if threshold <= 0 {
		return nil, errors.New("threshold must be greater than zero")
	}
	if cooldown <= 0 {
		return nil, errors.New("cooldown must be greater than zero")
	}
	if probeSuccess <= 0 {
		return nil, errors.New("probeSuccess must be greater than zero")
	}
	return &Breaker{
		threshold:    threshold,
		cooldown:     cooldown,
		probeSuccess: probeSuccess,
	}, nil
}

// Allow reports whether the caller is permitted to proceed.
// It transitions the breaker between states as appropriate.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			b.successes = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful call. In half-open state, enough
// successes will close the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.successes++
		if b.successes >= b.probeSuccess {
			b.state = StateClosed
			b.failures = 0
		}
	} else {
		b.failures = 0
	}
}

// RecordFailure records a failed call. Enough failures will open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
		b.failures = 0
	}
}

// CurrentState returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
