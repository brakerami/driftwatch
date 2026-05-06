// Package snapshot provides functionality for capturing and persisting
// container state snapshots to detect drift over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContainerSnapshot holds a point-in-time record of a container's
// observed state, used for comparing against future inspections.
type ContainerSnapshot struct {
	ContainerID   string            `json:"container_id"`
	ContainerName string            `json:"container_name"`
	Image         string            `json:"image"`
	Env           map[string]string `json:"env"`
	Labels        map[string]string `json:"labels"`
	CapturedAt    time.Time         `json:"captured_at"`
}

// Store manages reading and writing snapshots to a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that persists snapshots under dir.
// The directory is created if it does not already exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes snap to disk, keyed by container name.
func (s *Store) Save(snap ContainerSnapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := s.pathFor(snap.ContainerName)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads the most recent snapshot for the given container name.
// Returns os.ErrNotExist if no snapshot has been saved yet.
func (s *Store) Load(containerName string) (ContainerSnapshot, error) {
	var snap ContainerSnapshot
	data, err := os.ReadFile(s.pathFor(containerName))
	if err != nil {
		return snap, fmt.Errorf("snapshot: read %s: %w", containerName, err)
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, fmt.Errorf("snapshot: unmarshal %s: %w", containerName, err)
	}
	return snap, nil
}

// Delete removes the snapshot file for the given container name.
// It is a no-op if the file does not exist.
func (s *Store) Delete(containerName string) error {
	path := s.pathFor(containerName)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete %s: %w", path, err)
	}
	return nil
}

func (s *Store) pathFor(containerName string) string {
	return filepath.Join(s.dir, containerName+".json")
}
