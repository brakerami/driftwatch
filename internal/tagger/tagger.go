// Package tagger assigns semantic tags to drift findings based on
// configurable rules, enabling downstream grouping and routing.
package tagger

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Rule maps a drift type and optional container name pattern to a set of tags.
type Rule struct {
	DriftType string   `json:"drift_type"` // e.g. "env", "image", "label"
	Container string   `json:"container"`  // empty means match all
	Tags      []string `json:"tags"`
}

// Tagger applies tag rules to findings.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger from the provided rules.
// Returns an error if rules is nil or empty.
func New(rules []Rule) (*Tagger, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	return &Tagger{rules: rules}, nil
}

// Tag returns the deduplicated set of tags that apply to the given finding.
func (t *Tagger) Tag(f drift.Finding) []string {
	seen := make(map[string]struct{})
	var out []string

	for _, r := range t.rules {
		if !matchesType(r.DriftType, string(f.Type)) {
			continue
		}
		if !matchesContainer(r.Container, f.Container) {
			continue
		}
		for _, tag := range r.Tags {
			if _, ok := seen[tag]; !ok {
				seen[tag] = struct{}{}
				out = append(out, tag)
			}
		}
	}
	return out
}

// TagAll applies Tag to every finding and returns a map keyed by finding index.
func (t *Tagger) TagAll(findings []drift.Finding) map[int][]string {
	result := make(map[int][]string, len(findings))
	for i, f := range findings {
		if tags := t.Tag(f); len(tags) > 0 {
			result[i] = tags
		}
	}
	return result
}

func matchesType(ruleType, findingType string) bool {
	if ruleType == "" || ruleType == "*" {
		return true
	}
	return strings.EqualFold(ruleType, findingType)
}

func matchesContainer(pattern, name string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	return strings.EqualFold(pattern, name)
}
