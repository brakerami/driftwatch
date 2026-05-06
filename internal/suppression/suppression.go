// Package suppression provides a mechanism to suppress known or accepted
// drift findings so they do not appear in reports or trigger alerts.
package suppression

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Rule describes a single suppression rule. A finding is suppressed when
// both ContainerName and DriftType match (case-insensitive). An empty
// ContainerName matches any container.
type Rule struct {
	ContainerName string `json:"container_name"`
	DriftType     string `json:"drift_type"`
	Reason        string `json:"reason"`
}

// Store holds a collection of suppression rules loaded from a file.
type Store struct {
	rules []Rule
}

// LoadFile reads suppression rules from a JSON file at the given path.
func LoadFile(path string) (*Store, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("suppression: open %s: %w", path, err)
	}
	defer f.Close()

	var rules []Rule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, fmt.Errorf("suppression: decode %s: %w", path, err)
	}
	return &Store{rules: rules}, nil
}

// New creates a Store from an in-memory slice of rules.
func New(rules []Rule) *Store {
	return &Store{rules: rules}
}

// IsSuppressed reports whether a finding identified by containerName and
// driftType should be suppressed.
func (s *Store) IsSuppressed(containerName, driftType string) bool {
	for _, r := range s.rules {
		nameMatch := r.ContainerName == "" ||
			strings.EqualFold(r.ContainerName, containerName)
		typeMatch := strings.EqualFold(r.DriftType, driftType)
		if nameMatch && typeMatch {
			return true
		}
	}
	return false
}

// Rules returns a copy of the current rule list.
func (s *Store) Rules() []Rule {
	out := make([]Rule, len(s.rules))
	copy(out, s.rules)
	return out
}
