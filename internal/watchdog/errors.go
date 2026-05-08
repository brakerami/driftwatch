package watchdog

import "errors"

// Sentinel errors returned by the watchdog package.
var (
	// ErrNilSource is returned when New receives a nil EventSource.
	ErrNilSource = errors.New("watchdog: source must not be nil")

	// ErrNilHandler is returned when New receives a nil Handler.
	ErrNilHandler = errors.New("watchdog: handler must not be nil")

	// ErrInvalidInterval is returned when the polling interval is non-positive.
	ErrInvalidInterval = errors.New("watchdog: interval must be positive")
)
