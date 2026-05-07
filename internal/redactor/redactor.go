// Package redactor masks sensitive values (passwords, tokens, secrets)
// before they appear in drift reports or log output.
package redactor

import (
	"regexp"
	"strings"
)

// DefaultPatterns is the set of env-var key substrings treated as sensitive.
var DefaultPatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private_key",
	"auth",
	"credential",
}

const redacted = "[REDACTED]"

// Redactor masks environment variable values whose keys match sensitive
// patterns.
type Redactor struct {
	patterns []*regexp.Regexp
}

// New returns a Redactor that masks keys matching any of the supplied
// case-insensitive substring patterns. If patterns is empty,
// DefaultPatterns is used.
func New(patterns []string) (*Redactor, error) {
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}
	regs := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		r, err := regexp.Compile("(?i)" + regexp.QuoteMeta(p))
		if err != nil {
			return nil, err
		}
		regs = append(regs, r)
	}
	return &Redactor{patterns: regs}, nil
}

// IsSensitive reports whether the given env-var key should be masked.
func (r *Redactor) IsSensitive(key string) bool {
	for _, re := range r.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskValue returns the original value unchanged, or [REDACTED] if the
// key is considered sensitive.
func (r *Redactor) MaskValue(key, value string) string {
	if r.IsSensitive(key) {
		return redacted
	}
	return value
}

// MaskEnv accepts a slice of "KEY=VALUE" strings and returns a new slice
// with sensitive values replaced.
func (r *Redactor) MaskEnv(env []string) []string {
	out := make([]string, len(env))
	for i, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			out[i] = entry
			continue
		}
		out[i] = parts[0] + "=" + r.MaskValue(parts[0], parts[1])
	}
	return out
}
