package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Severity represents the severity level of a policy rule.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Rule defines a single drift policy rule.
type Rule struct {
	Name      string   `json:"name"`
	DriftType string   `json:"drift_type"`
	Containers []string `json:"containers,omitempty"`
	Severity  Severity `json:"severity"`
	Enabled   bool     `json:"enabled"`
}

// Policy holds a collection of rules loaded from a policy file.
type Policy struct {
	Rules []Rule `json:"rules"`
}

// LoadFile reads a JSON policy file from the given path.
func LoadFile(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("policy: open %q: %w", path, err)
	}
	defer f.Close()

	var p Policy
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("policy: decode %q: %w", path, err)
	}
	return &p, nil
}

// Store evaluates findings against policy rules.
type Store struct {
	rules []Rule
}

// New creates a Store from a Policy, keeping only enabled rules.
func New(p *Policy) *Store {
	var enabled []Rule
	for _, r := range p.Rules {
		if r.Enabled {
			enabled = append(enabled, r)
		}
	}
	return &Store{rules: enabled}
}

// Evaluate returns the highest Severity that matches the given drift type and
// container name, or an empty string when no rule matches.
func (s *Store) Evaluate(driftType, container string) Severity {
	var best Severity
	for _, r := range s.rules {
		if !strings.EqualFold(r.DriftType, driftType) {
			continue
		}
		if len(r.Containers) > 0 && !containsCI(r.Containers, container) {
			continue
		}
		if rank(r.Severity) > rank(best) {
			best = r.Severity
		}
	}
	return best
}

func rank(s Severity) int {
	switch s {
	case SeverityCritical:
		return 3
	case SeverityWarning:
		return 2
	case SeverityInfo:
		return 1
	}
	return 0
}

func containsCI(list []string, target string) bool {
	for _, v := range list {
		if strings.EqualFold(v, target) {
			return true
		}
	}
	return false
}
