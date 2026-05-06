// Package scheduler provides a periodic tick-based runner that triggers
// drift detection at a configured interval.
package scheduler

import (
	"context"
	"log"
	"time"
)

// Job is a function that performs a single drift-detection cycle.
// It receives the context so it can respect cancellation.
type Job func(ctx context.Context) error

// Scheduler runs a Job on a fixed interval until the context is cancelled.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *log.Logger
}

// New creates a Scheduler that will call job every interval.
// interval must be positive; if logger is nil the default logger is used.
func New(interval time.Duration, job Job, logger *log.Logger) (*Scheduler, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if job == nil {
		return nil, ErrNilJob
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}, nil
}

// Run starts the scheduler loop. It executes the job immediately on the first
// tick and then repeats every s.interval. Run blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Printf("scheduler: starting, interval=%s", s.interval)

	// Run once immediately before waiting for the first tick.
	s.runOnce(ctx)

	for {
		select {
		case <-ticker.C:
			s.runOnce(ctx)
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	if err := s.job(ctx); err != nil {
		s.logger.Printf("scheduler: job error: %v", err)
	}
}
