package throttle

import "errors"

// Sentinel errors exposed so callers can use errors.Is without importing
// the full package internals.
var (
	// ErrInvalidWindow is returned when New receives a non-positive window.
	ErrInvalidWindow = errors.New("throttle: window must be positive")

	// ErrInvalidMaxTokens is returned when New receives a non-positive maxTokens.
	ErrInvalidMaxTokens = errors.New("throttle: maxTokens must be positive")
)
