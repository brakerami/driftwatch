// Package limiter provides a concurrency limiter that caps the number of
// simultaneous drift-check goroutines to prevent resource exhaustion when
// many containers are evaluated in parallel.
package limiter

import (
	"errors"
	"fmt"
)

// ErrInvalidCapacity is returned when a non-positive capacity is provided.
var ErrInvalidCapacity = errors.New("limiter: capacity must be greater than zero")

// Limiter controls how many concurrent operations may proceed at once.
type Limiter struct {
	sem chan struct{}
	cap int
}

// New creates a Limiter that allows at most cap concurrent acquisitions.
// It returns ErrInvalidCapacity if cap is less than 1.
func New(cap int) (*Limiter, error) {
	if cap < 1 {
		return nil, fmt.Errorf("%w: got %d", ErrInvalidCapacity, cap)
	}
	return &Limiter{
		sem: make(chan struct{}, cap),
		cap: cap,
	}, nil
}

// Acquire blocks until a slot is available, then reserves it.
// Callers must pair every Acquire with a Release.
func (l *Limiter) Acquire() {
	l.sem <- struct{}{}
}

// TryAcquire attempts to reserve a slot without blocking.
// It returns true if the slot was acquired, false if all slots are busy.
func (l *Limiter) TryAcquire() bool {
	select {
	case l.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	<-l.sem
}

// Cap returns the maximum concurrency configured for this Limiter.
func (l *Limiter) Cap() int {
	return l.cap
}

// InFlight returns the number of slots currently in use.
func (l *Limiter) InFlight() int {
	return len(l.sem)
}
