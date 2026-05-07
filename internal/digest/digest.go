// Package digest computes and compares deterministic fingerprints for
// container runtime state, allowing quick detection of config drift without
// a full field-by-field comparison on every poll cycle.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// ContainerState holds the subset of runtime state used for fingerprinting.
type ContainerState struct {
	Image  string            `json:"image"`
	Env    map[string]string `json:"env"`
	Labels map[string]string `json:"labels"`
	Ports  []string          `json:"ports"`
}

// Compute returns a stable SHA-256 hex digest of the given ContainerState.
// Maps are sorted by key before hashing to ensure determinism.
func Compute(s ContainerState) (string, error) {
	norm := struct {
		Image  string     `json:"image"`
		Env    [][]string `json:"env"`
		Labels [][]string `json:"labels"`
		Ports  []string   `json:"ports"`
	}{
		Image:  s.Image,
		Env:    sortedPairs(s.Env),
		Labels: sortedPairs(s.Labels),
		Ports:  sortedStrings(s.Ports),
	}

	data, err := json.Marshal(norm)
	if err != nil {
		return "", fmt.Errorf("digest: marshal: %w", err)
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Equal returns true when both states produce the same digest.
// The error from either Compute call is surfaced if non-nil.
func Equal(a, b ContainerState) (bool, error) {
	da, err := Compute(a)
	if err != nil {
		return false, err
	}
	db, err := Compute(b)
	if err != nil {
		return false, err
	}
	return da == db, nil
}

func sortedPairs(m map[string]string) [][]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([][]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, []string{k, m[k]})
	}
	return pairs
}

func sortedStrings(ss []string) []string {
	copy_ := append([]string(nil), ss...)
	sort.Strings(copy_)
	return copy_
}
