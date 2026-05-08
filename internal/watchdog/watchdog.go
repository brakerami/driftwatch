// Package watchdog monitors container health and triggers drift checks
// when containers are restarted or recreated unexpectedly.
package watchdog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Event represents a container lifecycle event detected by the watchdog.
type Event struct {
	ContainerID   string
	ContainerName string
	Kind          EventKind
	OccurredAt    time.Time
}

// EventKind classifies the type of lifecycle event.
type EventKind string

const (
	EventRestart  EventKind = "restart"
	EventRecreate EventKind = "recreate"
	EventStop     EventKind = "stop"
)

// Handler is a function invoked when a container event is detected.
type Handler func(ctx context.Context, e Event) error

// Watchdog polls a source for container events and dispatches them to a handler.
type Watchdog struct {
	source   EventSource
	handler  Handler
	interval time.Duration
	mu       sync.Mutex
	seen     map[string]time.Time
}

// EventSource is implemented by anything that can return recent container events.
type EventSource interface {
	Events(ctx context.Context, since time.Time) ([]Event, error)
}

// New creates a Watchdog. interval controls how often the source is polled.
func New(source EventSource, handler Handler, interval time.Duration) (*Watchdog, error) {
	if source == nil {
		return nil, fmt.Errorf("watchdog: source must not be nil")
	}
	if handler == nil {
		return nil, fmt.Errorf("watchdog: handler must not be nil")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("watchdog: interval must be positive")
	}
	return &Watchdog{
		source:   source,
		handler:  handler,
		interval: interval,
		seen:     make(map[string]time.Time),
	}, nil
}

// Run starts the watchdog loop and blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	lastPoll := time.Now().Add(-w.interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.C:
			events, err := w.source.Events(ctx, lastPoll)
			if err != nil {
				lastPoll = t
				continue
			}
			for _, e := range events {
				_ = w.handler(ctx, e)
			}
			lastPoll = t
		}
	}
}
