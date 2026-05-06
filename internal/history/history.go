// Package history records drift scan results over time,
// allowing callers to query past findings and detect trends.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Entry represents a single recorded scan result.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Container string         `json:"container"`
	Findings  []drift.Finding `json:"findings"`
}

// Store persists scan history entries to a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that writes entries under dir.
// The directory is created if it does not exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Record saves an entry for the given container and findings.
func (s *Store) Record(container string, findings []drift.Finding) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Container: container,
		Findings:  findings,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("history: marshal entry: %w", err)
	}
	filename := fmt.Sprintf("%s_%d.json", sanitize(container), entry.Timestamp.UnixNano())
	path := filepath.Join(s.dir, filename)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %q: %w", path, err)
	}
	return nil
}

// List returns all recorded entries for the given container,
// sorted ascending by timestamp.
func (s *Store) List(container string) ([]Entry, error) {
	pattern := filepath.Join(s.dir, sanitize(container)+"_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %q: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %q: %w", m, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

// sanitize replaces characters unsafe for filenames with underscores.
func sanitize(name string) string {
	out := make([]byte, len(name))
	for i := range name {
		if name[i] == '/' || name[i] == ':' || name[i] == '\\' {
			out[i] = '_'
		} else {
			out[i] = name[i]
		}
	}
	return string(out)
}
