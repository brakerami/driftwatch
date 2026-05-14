package enricher

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{
			Container:       "api",
			Type:            drift.TypeEnv,
			Field:           "LOG_LEVEL",
			Expected:        "info",
			Actual:          "debug",
			ContainerLabels: map[string]string{"team": "platform", "env": "prod"},
			Metadata:        map[string]string{"existing": "value"},
		},
		{
			Container:       "worker",
			Type:            drift.TypeImage,
			Field:           "image",
			Expected:        "worker:1.0",
			Actual:          "worker:latest",
			ContainerLabels: map[string]string{},
			Metadata:        nil,
		},
	}
}

func TestNew_DefaultsHostname(t *testing.T) {
	e, err := New(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.hostname == "" {
		t.Error("hostname should not be empty")
	}
}

func TestEnrich_AddsHost(t *testing.T) {
	e, _ := New(Config{})
	e.hostname = "testnode"

	out := e.Enrich(sampleFindings())
	for _, f := range out {
		if got := f.Metadata["host"]; got != "testnode" {
			t.Errorf("expected host=testnode, got %q", got)
		}
	}
}

func TestEnrich_StaticTagsPresent(t *testing.T) {
	e, _ := New(Config{
		StaticTags: map[string]string{"env": "staging", "region": "us-east-1"},
	})
	out := e.Enrich(sampleFindings())
	for _, f := range out {
		if f.Metadata["env"] != "staging" {
			t.Errorf("expected env=staging in metadata, got %q", f.Metadata["env"])
		}
		if f.Metadata["region"] != "us-east-1" {
			t.Errorf("expected region=us-east-1, got %q", f.Metadata["region"])
		}
	}
}

func TestEnrich_LabelKeysPromoted(t *testing.T) {
	e, _ := New(Config{LabelKeys: []string{"team"}})
	out := e.Enrich(sampleFindings())

	// first finding has label "team"="platform"
	if got := out[0].Metadata["label.team"]; got != "platform" {
		t.Errorf("expected label.team=platform, got %q", got)
	}
	// second finding has no "team" label — key should be absent
	if _, ok := out[1].Metadata["label.team"]; ok {
		t.Error("label.team should not be present when container has no such label")
	}
}

func TestEnrich_PreservesExistingMetadata(t *testing.T) {
	e, _ := New(Config{})
	out := e.Enrich(sampleFindings())
	if out[0].Metadata["existing"] != "value" {
		t.Error("pre-existing metadata key should be preserved")
	}
}

func TestEnrich_DoesNotMutateOriginal(t *testing.T) {
	original := sampleFindings()
	e, _ := New(Config{StaticTags: map[string]string{"injected": "yes"}})
	e.Enrich(original)
	if _, ok := original[0].Metadata["injected"]; ok {
		t.Error("Enrich must not mutate the original findings slice")
	}
}

func TestEnrich_EmptyFindings_ReturnsEmpty(t *testing.T) {
	e, _ := New(Config{})
	out := e.Enrich([]drift.Finding{})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d findings", len(out))
	}
}
