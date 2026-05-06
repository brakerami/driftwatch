package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/filter"
)

func TestMatchName_EmptyAllowsAll(t *testing.T) {
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.MatchName("anything") {
		t.Error("empty name list should match any name")
	}
}

func TestMatchName_ExplicitList(t *testing.T) {
	f, _ := filter.New(filter.Options{Names: []string{"web", "db"}})

	if !f.MatchName("web") {
		t.Error("expected 'web' to match")
	}
	if f.MatchName("cache") {
		t.Error("expected 'cache' not to match")
	}
}

func TestMatchLabels_EmptySelectorMatchesAll(t *testing.T) {
	f, _ := filter.New(filter.Options{})
	if !f.MatchLabels(map[string]string{"env": "prod"}) {
		t.Error("empty selector should match any labels")
	}
}

func TestMatchLabels_AllMustMatch(t *testing.T) {
	f, err := filter.New(filter.Options{LabelSelector: "env=prod,team=platform"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	good := map[string]string{"env": "prod", "team": "platform", "extra": "ok"}
	if !f.MatchLabels(good) {
		t.Error("expected full match")
	}

	bad := map[string]string{"env": "prod", "team": "infra"}
	if f.MatchLabels(bad) {
		t.Error("expected mismatch on team label")
	}
}

func TestMatchLabels_MissingKey(t *testing.T) {
	f, _ := filter.New(filter.Options{LabelSelector: "env=prod"})
	if f.MatchLabels(map[string]string{}) {
		t.Error("missing key should not match")
	}
}

func TestParseLabels_InvalidSegment(t *testing.T) {
	_, err := filter.New(filter.Options{LabelSelector: "noequals"})
	if err == nil {
		t.Fatal("expected parse error for segment without '='")
	}
	if _, ok := err.(*filter.ParseError); !ok {
		t.Errorf("expected *filter.ParseError, got %T", err)
	}
}

func TestParseLabels_EmptyValueAllowed(t *testing.T) {
	f, err := filter.New(filter.Options{LabelSelector: "env="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.MatchLabels(map[string]string{"env": ""}) {
		t.Error("empty value should match empty value")
	}
}
