package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func samplePolicy() *Policy {
	return &Policy{
		Rules: []Rule{
			{Name: "env-critical", DriftType: "env", Severity: SeverityCritical, Enabled: true},
			{Name: "image-warning", DriftType: "image", Severity: SeverityWarning, Enabled: true},
			{Name: "label-info", DriftType: "label", Severity: SeverityInfo, Enabled: false},
			{
				Name:       "env-specific",
				DriftType:  "env",
				Containers: []string{"api"},
				Severity:   SeverityCritical,
				Enabled:    true,
			},
		},
	}
}

func TestNew_OnlyEnabledRulesLoaded(t *testing.T) {
	s := New(samplePolicy())
	if len(s.rules) != 3 {
		t.Fatalf("expected 3 enabled rules, got %d", len(s.rules))
	}
}

func TestEvaluate_MatchesDriftType(t *testing.T) {
	s := New(samplePolicy())
	got := s.Evaluate("image", "any-container")
	if got != SeverityWarning {
		t.Fatalf("expected warning, got %q", got)
	}
}

func TestEvaluate_NoMatch_ReturnsEmpty(t *testing.T) {
	s := New(samplePolicy())
	got := s.Evaluate("label", "any-container")
	if got != "" {
		t.Fatalf("expected empty severity, got %q", got)
	}
}

func TestEvaluate_ContainerFilter_Matches(t *testing.T) {
	s := New(samplePolicy())
	got := s.Evaluate("env", "API") // case-insensitive
	if got != SeverityCritical {
		t.Fatalf("expected critical, got %q", got)
	}
}

func TestEvaluate_ContainerFilter_NoMatch(t *testing.T) {
	// Only the global env rule should fire for an unlisted container.
	s := New(&Policy{
		Rules: []Rule{
			{Name: "env-specific", DriftType: "env", Containers: []string{"api"}, Severity: SeverityCritical, Enabled: true},
		},
	})
	got := s.Evaluate("env", "worker")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestLoadFile_Valid(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "policy.json")

	p := samplePolicy()
	data, _ := json.Marshal(p)
	_ = os.WriteFile(path, data, 0o644)

	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Rules) != len(p.Rules) {
		t.Fatalf("rule count mismatch: got %d", len(loaded.Rules))
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := LoadFile("/no/such/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
