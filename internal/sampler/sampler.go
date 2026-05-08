// Package sampler provides probabilistic sampling for drift findings,
// allowing high-volume environments to reduce noise by forwarding only
// a representative fraction of detected drift events.
package sampler

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/driftwatch/internal/drift"
)

// ErrInvalidRate is returned when a sampling rate outside [0.0, 1.0] is provided.
var ErrInvalidRate = errors.New("sampler: rate must be between 0.0 and 1.0")

// Sampler probabilistically filters drift findings.
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New creates a Sampler with the given sampling rate.
// A rate of 1.0 passes all findings; 0.0 drops all findings.
func New(rate float64) (*Sampler, error) {
	if rate < 0.0 || rate > 1.0 {
		return nil, ErrInvalidRate
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(rand.Int63())), //nolint:gosec
	}, nil
}

// Sample returns a subset of findings based on the configured rate.
// Each finding is independently included with probability equal to the rate.
func (s *Sampler) Sample(findings []drift.Finding) []drift.Finding {
	if len(findings) == 0 {
		return nil
	}
	if s.rate == 1.0 {
		return findings
	}
	if s.rate == 0.0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]drift.Finding, 0, len(findings))
	for _, f := range findings {
		if s.rng.Float64() < s.rate {
			out = append(out, f)
		}
	}
	return out
}

// Rate returns the current sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
