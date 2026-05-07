package baseline

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

func sampleEntry(name string) Entry {
	return Entry{
		ContainerID:   "abc123",
		ContainerName: name,
		Image:         "nginx:1.25",
		Env:           map[string]string{"PORT": "8080", "ENV": "prod"},
		CapturedAt:    time.Now().UTC(),
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	e := sampleEntry("web")
	if err := store.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load("web")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Image != e.Image {
		t.Errorf("image: want %q got %q", e.Image, got.Image)
	}
	if got.Env["PORT"] != "8080" {
		t.Errorf("env PORT: want 8080 got %q", got.Env["PORT"])
	}
}

func TestLoad_NotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
	if !os.IsNotExist(err) {
		// wrapped error — just ensure it's non-nil (already checked)
	}
}

func TestDelete_RemovesBaseline(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)
	e := sampleEntry("api")
	_ = store.Save(e)
	if err := store.Delete("api"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("api")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDelete_Idempotent(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)
	if err := store.Delete("ghost"); err != nil {
		t.Fatalf("Delete on missing should not error: %v", err)
	}
}

func TestDiffFromBaseline_ImageMismatch(t *testing.T) {
	base := sampleEntry("web")
	findings := DiffFromBaseline(base, base.Env, "nginx:1.26")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Type != drift.TypeImageMismatch {
		t.Errorf("expected TypeImageMismatch, got %v", findings[0].Type)
	}
}

func TestDiffFromBaseline_EnvMismatch(t *testing.T) {
	base := sampleEntry("web")
	env := map[string]string{"PORT": "9090", "ENV": "prod"}
	findings := DiffFromBaseline(base, env, base.Image)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Type != drift.TypeEnvMismatch {
		t.Errorf("expected TypeEnvMismatch, got %v", findings[0].Type)
	}
}

func TestDiffFromBaseline_NoDrift(t *testing.T) {
	base := sampleEntry("web")
	findings := DiffFromBaseline(base, base.Env, base.Image)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}
