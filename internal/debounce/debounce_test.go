package debounce

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Container: "web", Type: drift.TypeEnv, Field: "PORT", Expected: "8080", Actual: "9090"},
		{Container: "db", Type: drift.TypeImage, Field: "image", Expected: "postgres:14", Actual: "postgres:13"},
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_ValidWindow(t *testing.T) {
	d, err := New(5 * time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestFilter_FirstCallPassesThrough(t *testing.T) {
	d, _ := New(1 * time.Minute)
	out := d.Filter(sampleFindings())
	if len(out) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out))
	}
}

func TestFilter_SecondCallWithinWindowSuppressed(t *testing.T) {
	d, _ := New(1 * time.Minute)
	d.Filter(sampleFindings())
	out := d.Filter(sampleFindings())
	if len(out) != 0 {
		t.Fatalf("expected 0 findings within window, got %d", len(out))
	}
}

func TestFilter_AfterWindowExpiry_PassesThrough(t *testing.T) {
	d, _ := New(1 * time.Minute)

	fixed := time.Now()
	d.now = func() time.Time { return fixed }
	d.Filter(sampleFindings())

	// Advance clock past the window.
	d.now = func() time.Time { return fixed.Add(2 * time.Minute) }
	out := d.Filter(sampleFindings())
	if len(out) != 2 {
		t.Fatalf("expected 2 findings after window expiry, got %d", len(out))
	}
}

func TestFlush_ResetsState(t *testing.T) {
	d, _ := New(1 * time.Minute)
	d.Filter(sampleFindings())
	d.Flush()
	out := d.Filter(sampleFindings())
	if len(out) != 2 {
		t.Fatalf("expected 2 findings after flush, got %d", len(out))
	}
}

func TestFilter_EmptyInput_ReturnsNil(t *testing.T) {
	d, _ := New(1 * time.Minute)
	out := d.Filter(nil)
	if out != nil {
		t.Fatalf("expected nil for empty input, got %v", out)
	}
}
