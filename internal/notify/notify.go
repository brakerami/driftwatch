// Package notify provides a pluggable notification dispatch layer
// that routes drift findings to one or more registered sinks.
package notify

import (
	"context"
	"fmt"
	"log"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Sink is any destination that can receive drift findings.
type Sink interface {
	// Name returns a human-readable identifier for the sink.
	Name() string
	// Send dispatches findings to the sink.
	Send(ctx context.Context, findings []drift.Finding) error
}

// Dispatcher fans out findings to all registered sinks.
type Dispatcher struct {
	sinks  []Sink
	logger *log.Logger
}

// New creates a Dispatcher with the provided sinks.
// At least one sink must be supplied.
func New(logger *log.Logger, sinks ...Sink) (*Dispatcher, error) {
	if len(sinks) == 0 {
		return nil, fmt.Errorf("notify: at least one sink is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("notify: logger must not be nil")
	}
	return &Dispatcher{sinks: sinks, logger: logger}, nil
}

// Dispatch sends findings to every registered sink.
// Errors from individual sinks are logged but do not abort delivery to
// remaining sinks. Returns an error only when ALL sinks fail.
func (d *Dispatcher) Dispatch(ctx context.Context, findings []drift.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	failed := 0
	for _, s := range d.sinks {
		if err := s.Send(ctx, findings); err != nil {
			d.logger.Printf("notify: sink %q error: %v", s.Name(), err)
			failed++
		}
	}

	if failed == len(d.sinks) {
		return fmt.Errorf("notify: all %d sink(s) failed", failed)
	}
	return nil
}

// SinkCount returns the number of registered sinks.
func (d *Dispatcher) SinkCount() int { return len(d.sinks) }
