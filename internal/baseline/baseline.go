// Package baseline provides functionality for capturing and comparing
// a known-good state of container configurations to detect future drift.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Entry represents a saved baseline for a single container.
type Entry struct {
	ContainerID   string            `json:"container_id"`
	ContainerName string            `json:"container_name"`
	Env           map[string]string `json:"env"`
	Image         string            `json:"image"`
	CapturedAt    time.Time         `json:"captured_at"`
}

// Store persists and retrieves baseline entries on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating it if necessary.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(containerName string) string {
	return filepath.Join(s.dir, sanitize(containerName)+".json")
}

// Save writes an Entry to disk, overwriting any previous baseline.
func (s *Store) Save(e Entry) error {
	e.CapturedAt = time.Now().UTC()
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(s.path(e.ContainerName), data, 0o644); err != nil {
		return fmt.Errorf("baseline: write: %w", err)
	}
	return nil
}

// Load retrieves the baseline Entry for the named container.
// Returns os.ErrNotExist if no baseline has been captured.
func (s *Store) Load(containerName string) (Entry, error) {
	data, err := os.ReadFile(s.path(containerName))
	if err != nil {
		return Entry{}, fmt.Errorf("baseline: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return e, nil
}

// Delete removes the baseline for the named container. It is a no-op if none exists.
func (s *Store) Delete(containerName string) error {
	err := os.Remove(s.path(containerName))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("baseline: delete: %w", err)
	}
	return nil
}

// DiffFromBaseline returns drift findings by comparing current env/image
// against the stored baseline rather than a manifest.
func DiffFromBaseline(base Entry, currentEnv map[string]string, currentImage string) []drift.Finding {
	var findings []drift.Finding
	if base.Image != currentImage {
		findings = append(findings, drift.Finding{
			Container: base.ContainerName,
			Type:      drift.TypeImageMismatch,
			Expected:  base.Image,
			Actual:    currentImage,
		})
	}
	for k, want := range base.Env {
		got, ok := currentEnv[k]
		if !ok {
			findings = append(findings, drift.Finding{
				Container: base.ContainerName,
				Type:      drift.TypeEnvMissing,
				Expected:  fmt.Sprintf("%s=%s", k, want),
				Actual:    "",
			})
		} else if got != want {
			findings = append(findings, drift.Finding{
				Container: base.ContainerName,
				Type:      drift.TypeEnvMismatch,
				Expected:  fmt.Sprintf("%s=%s", k, want),
				Actual:    fmt.Sprintf("%s=%s", k, got),
			})
		}
	}
	return findings
}

func sanitize(name string) string {
	out := make([]byte, len(name))
	for i := range name {
		if name[i] == '/' || name[i] == ':' || name[i] == ' ' {
			out[i] = '_'
		} else {
			out[i] = name[i]
		}
	}
	return string(out)
}
