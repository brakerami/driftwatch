package tagger

import "errors"

// ErrNoRules is returned when a Tagger is constructed with no rules.
var ErrNoRules = errors.New("tagger: at least one rule is required")
