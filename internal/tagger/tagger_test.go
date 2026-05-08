package tagger_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tagger"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Container: "api", Type: drift.DriftTypeEnv, Field: "LOG_LEVEL", Want: "info", Got: "debug"},
		{Container: "worker", Type: drift.DriftTypeImage, Field: "image", Want: "v1", Got: "v2"},
		{Container: "api", Type: drift.DriftTypeLabel, Field: "team", Want: "platform", Got: ""},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := tagger.New(nil)
	if err == nil {
		t.Fatal("expected error for nil rules, got nil")
	}
	_, err = tagger.New([]tagger.Rule{})
	if err == nil {
		t.Fatal("expected error for empty rules, got nil")
	}
}

func TestNew_ValidRules_NoError(t *testing.T) {
	_, err := tagger.New([]tagger.Rule{{DriftType: "env", Tags: []string{"config"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTag_MatchesByType(t *testing.T) {
	tgr, _ := tagger.New([]tagger.Rule{
		{DriftType: "env", Tags: []string{"config", "sensitive"}},
		{DriftType: "image", Tags: []string{"deployment"}},
	})

	tags := tgr.Tag(sampleFindings()[0]) // env finding
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(tags), tags)
	}
}

func TestTag_ContainerFilter_Matches(t *testing.T) {
	tgr, _ := tagger.New([]tagger.Rule{
		{DriftType: "env", Container: "api", Tags: []string{"api-env"}},
		{DriftType: "env", Container: "worker", Tags: []string{"worker-env"}},
	})

	tags := tgr.Tag(sampleFindings()[0]) // api env
	if len(tags) != 1 || tags[0] != "api-env" {
		t.Fatalf("expected [api-env], got %v", tags)
	}
}

func TestTag_WildcardType_MatchesAll(t *testing.T) {
	tgr, _ := tagger.New([]tagger.Rule{
		{DriftType: "*", Tags: []string{"drift"}},
	})

	for _, f := range sampleFindings() {
		tags := tgr.Tag(f)
		if len(tags) == 0 {
			t.Errorf("expected tag for finding %+v, got none", f)
		}
	}
}

func TestTag_DeduplicatesTags(t *testing.T) {
	tgr, _ := tagger.New([]tagger.Rule{
		{DriftType: "env", Tags: []string{"config"}},
		{DriftType: "env", Tags: []string{"config", "extra"}},
	})

	tags := tgr.Tag(sampleFindings()[0])
	if len(tags) != 2 {
		t.Fatalf("expected 2 deduplicated tags, got %d: %v", len(tags), tags)
	}
}

func TestTagAll_ReturnsIndexedMap(t *testing.T) {
	tgr, _ := tagger.New([]tagger.Rule{
		{DriftType: "*", Tags: []string{"all"}},
	})

	result := tgr.TagAll(sampleFindings())
	if len(result) != 3 {
		t.Fatalf("expected 3 entries in map, got %d", len(result))
	}
	for i := 0; i < 3; i++ {
		if _, ok := result[i]; !ok {
			t.Errorf("missing index %d in result", i)
		}
	}
}
