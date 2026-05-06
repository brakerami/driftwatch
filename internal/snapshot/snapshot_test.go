package snapshot_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/snapshot"
)

func sampleSnap() snapshot.ContainerSnapshot {
	return snapshot.ContainerSnapshot{
		ContainerID:   "abc123",
		ContainerName: "web",
		Image:         "nginx:1.25",
		Env:           map[string]string{"PORT": "8080", "DEBUG": "false"},
		Labels:        map[string]string{"app": "web"},
		CapturedAt:    time.Now().UTC().Truncate(time.Second),
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	orig := sampleSnap()
	if err := store.Save(orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load(orig.ContainerName)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.ContainerID != orig.ContainerID {
		t.Errorf("ContainerID: got %q, want %q", loaded.ContainerID, orig.ContainerID)
	}
	if loaded.Image != orig.Image {
		t.Errorf("Image: got %q, want %q", loaded.Image, orig.Image)
	}
	if loaded.Env["PORT"] != orig.Env["PORT"] {
		t.Errorf("Env PORT: got %q, want %q", loaded.Env["PORT"], orig.Env["PORT"])
	}
}

func TestLoad_NotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got: %v", err)
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	snap := sampleSnap()
	_ = store.Save(snap)

	if err := store.Delete(snap.ContainerName); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := store.Load(snap.ContainerName)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected snapshot to be deleted, got: %v", err)
	}
}

func TestDelete_Idempotent(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	if err := store.Delete("ghost"); err != nil {
		t.Errorf("Delete of non-existent should be no-op, got: %v", err)
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/snapshots"
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore should create nested dirs: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
