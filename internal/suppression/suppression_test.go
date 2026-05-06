package suppression_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/driftwatch/internal/suppression"
)

func TestIsSuppressed_EmptyStore(t *testing.T) {
	s := suppression.New(nil)
	if s.IsSuppressed("web", "env") {
		t.Error("expected no suppression with empty store")
	}
}

func TestIsSuppressed_ExactMatch(t *testing.T) {
	s := suppression.New([]suppression.Rule{
		{ContainerName: "web", DriftType: "env", Reason: "known"},
	})
	if !s.IsSuppressed("web", "env") {
		t.Error("expected suppression for web/env")
	}
}

func TestIsSuppressed_CaseInsensitive(t *testing.T) {
	s := suppression.New([]suppression.Rule{
		{ContainerName: "Web", DriftType: "ENV"},
	})
	if !s.IsSuppressed("web", "env") {
		t.Error("expected case-insensitive match")
	}
}

func TestIsSuppressed_WildcardContainer(t *testing.T) {
	s := suppression.New([]suppression.Rule{
		{ContainerName: "", DriftType: "image"},
	})
	if !s.IsSuppressed("any-container", "image") {
		t.Error("expected wildcard container to match any name")
	}
}

func TestIsSuppressed_NoMatchDifferentType(t *testing.T) {
	s := suppression.New([]suppression.Rule{
		{ContainerName: "web", DriftType: "env"},
	})
	if s.IsSuppressed("web", "image") {
		t.Error("expected no match for different drift type")
	}
}

func TestLoadFile_InvalidPath(t *testing.T) {
	_, err := suppression.LoadFile("/nonexistent/rules.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_ValidJSON(t *testing.T) {
	rules := []suppression.Rule{
		{ContainerName: "db", DriftType: "env", Reason: "intentional"},
	}
	f, err := os.CreateTemp(t.TempDir(), "rules*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(rules); err != nil {
		t.Fatal(err)
	}
	f.Close()

	s, err := suppression.LoadFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules()) != 1 {
		t.Errorf("expected 1 rule, got %d", len(s.Rules()))
	}
}
