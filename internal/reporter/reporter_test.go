package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/reporter"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{ContainerID: "abc123", Type: drift.DriftTypeEnv, Detail: "ENV FOO: want bar got baz"},
		{ContainerID: "def456", Type: drift.DriftTypeImage, Detail: "image tag mismatch"},
	}
}

func TestWrite_TextFormat_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(sampleFindings()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"CONTAINER", "TYPE", "DETAIL", "abc123", "def456"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Write(sampleFindings()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rep reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &rep); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if rep.TotalDrift != 2 {
		t.Errorf("expected TotalDrift=2, got %d", rep.TotalDrift)
	}
}

func TestWrite_NoFindings_TextSaysNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected 'No drift detected' message, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Passing nil should not panic; we just verify construction succeeds.
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
