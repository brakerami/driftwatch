package deduplicator

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Container: "web", Type: drift.DriftTypeEnv, Field: "PORT", Expected: "8080", Actual: "9090"},
		{Container: "db", Type: drift.DriftTypeImage, Field: "image", Expected: "postgres:14", Actual: "postgres:13"},
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ValidWindow(t *testing.T) {
	s, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestFilter_FirstCallPassesThrough(t *testing.T) {
	s, _ := New(time.Minute)
	out := s.Filter(sampleFindings())
	if len(out) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out))
	}
}

func TestFilter_DuplicateWithinWindowSuppressed(t *testing.T) {
	s, _ := New(time.Minute)
	s.Filter(sampleFindings())
	out := s.Filter(sampleFindings())
	if len(out) != 0 {
		t.Fatalf("expected 0 findings on second call, got %d", len(out))
	}
}

func TestFilter_AfterWindowExpiry_PassesThrough(t *testing.T) {
	s, _ := New(50 * time.Millisecond)
	fixed := time.Now()
	s.nowFunc = func() time.Time { return fixed }
	s.Filter(sampleFindings())

	// Advance past the window.
	s.nowFunc = func() time.Time { return fixed.Add(100 * time.Millisecond) }
	out := s.Filter(sampleFindings())
	if len(out) != 2 {
		t.Fatalf("expected 2 findings after window expiry, got %d", len(out))
	}
}

func TestFilter_PartialDuplicates(t *testing.T) {
	s, _ := New(time.Minute)
	s.Filter(sampleFindings()[:1]) // record only first
	out := s.Filter(sampleFindings()) // second call with both
	if len(out) != 1 {
		t.Fatalf("expected 1 new finding, got %d", len(out))
	}
	if out[0].Container != "db" {
		t.Fatalf("expected db finding to pass through, got %s", out[0].Container)
	}
}

func TestLen_TracksEntries(t *testing.T) {
	s, _ := New(time.Minute)
	if s.Len() != 0 {
		t.Fatal("expected empty store")
	}
	s.Filter(sampleFindings())
	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}
}
