package enricher

import "errors"

// ErrNoConfig is returned when New is called with a zero-value Config that
// would produce a no-op enricher and the caller has explicitly opted in to
// strict mode (reserved for future use).
var ErrNoConfig = errors.New("enricher: no configuration provided")
