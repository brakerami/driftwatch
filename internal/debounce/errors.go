package debounce

import "errors"

// ErrInvalidWindow is returned when a non-positive window duration is provided.
var ErrInvalidWindow = errors.New("debounce: window must be greater than zero")
