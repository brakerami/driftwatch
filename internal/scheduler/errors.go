package scheduler

import "errors"

// Sentinel errors returned by New.
var (
	// ErrInvalidInterval is returned when a non-positive interval is supplied.
	ErrInvalidInterval = errors.New("scheduler: interval must be greater than zero")

	// ErrNilJob is returned when a nil Job function is supplied.
	ErrNilJob = errors.New("scheduler: job must not be nil")
)
