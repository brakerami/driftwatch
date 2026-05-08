package policy

import (
	"fmt"
	"os"
)

// Loader resolves a Policy from a file path, returning a ready-to-use Store.
// When no path is provided it falls back to a permissive default policy that
// assigns SeverityInfo to every drift type.
type Loader struct {
	path string
}

// NewLoader creates a Loader targeting the given file path.
// Pass an empty string to use the built-in default policy.
func NewLoader(path string) *Loader {
	return &Loader{path: path}
}

// Load reads the policy file (if configured) and returns a Store.
// If the path is empty the default policy is returned.
func (l *Loader) Load() (*Store, error) {
	if l.path == "" {
		return New(defaultPolicy()), nil
	}
	if _, err := os.Stat(l.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("policy: file not found: %q", l.path)
	}
	p, err := LoadFile(l.path)
	if err != nil {
		return nil, err
	}
	return New(p), nil
}

// defaultPolicy returns a minimal policy that surfaces all drift types at info
// level so that operators always receive some signal even without configuration.
func defaultPolicy() *Policy {
	return &Policy{
		Rules: []Rule{
			{Name: "default-env", DriftType: "env", Severity: SeverityInfo, Enabled: true},
			{Name: "default-image", DriftType: "image", Severity: SeverityInfo, Enabled: true},
			{Name: "default-label", DriftType: "label", Severity: SeverityInfo, Enabled: true},
		},
	}
}
