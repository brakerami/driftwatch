package rollup_test

import (
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/rollup"
)

func makeFindings(types ...drift.DriftType) []drift.Finding {
	out := make([]drift.Finding, len(types))
	for i, t := range types {
		out[i] = drift.Finding{Type: t, Field: "x", Expected: "a", Actual: "b"}
	}
	return out
}

func TestAggregate_EmptyResults(t *testing.T) {
	s := rollup.Aggregate(map[string][]drift.Finding{})
	if s.TotalContainers != 0 || s.DriftedContainers != 0 {
		t.Fatalf("expected zero counts, got %+v", s)
	}
}

func TestAggregate_NoDrift(t *testing.T) {
	results := map[string][]drift.Finding{
		"web": {},
		"db":  {},
	}
	s := rollup.Aggregate(results)
	if s.TotalContainers != 2 {
		t.Fatalf("expected 2 total, got %d", s.TotalContainers)
	}
	if s.DriftedContainers != 0 {
		t.Fatalf("expected 0 drifted, got %d", s.DriftedContainers)
	}
}

func TestAggregate_CountsByType(t *testing.T) {
	results := map[string][]drift.Finding{
		"alpha": makeFindings(drift.DriftTypeEnv, drift.DriftTypeEnv, drift.DriftTypeImage),
		"beta":  makeFindings(drift.DriftTypeEnv),
	}
	s := rollup.Aggregate(results)
	if s.DriftedContainers != 2 {
		t.Fatalf("expected 2 drifted, got %d", s.DriftedContainers)
	}
	if s.ByType[string(drift.DriftTypeEnv)] != 3 {
		t.Fatalf("expected 3 env findings, got %d", s.ByType[string(drift.DriftTypeEnv)])
	}
	if s.ByType[string(drift.DriftTypeImage)] != 1 {
		t.Fatalf("expected 1 image finding, got %d", s.ByType[string(drift.DriftTypeImage)])
	}
}

func TestAggregate_TopOffendersOrdered(t *testing.T) {
	results := map[string][]drift.Finding{
		"low":    makeFindings(drift.DriftTypeEnv),
		"high":   makeFindings(drift.DriftTypeEnv, drift.DriftTypeEnv, drift.DriftTypeImage),
		"medium": makeFindings(drift.DriftTypeEnv, drift.DriftTypeImage),
	}
	s := rollup.Aggregate(results)
	if len(s.TopOffenders) != 3 {
		t.Fatalf("expected 3 offenders, got %d", len(s.TopOffenders))
	}
	if s.TopOffenders[0].Name != "high" {
		t.Fatalf("expected 'high' first, got %s", s.TopOffenders[0].Name)
	}
}

func TestSummary_String_ContainsHeaders(t *testing.T) {
	results := map[string][]drift.Finding{
		"app": makeFindings(drift.DriftTypeEnv),
	}
	s := rollup.Aggregate(results)
	out := s.String()
	for _, want := range []string{"Containers checked", "Findings by type", "Top offenders"} {
		if !strings.Contains(out, want) {
			t.Errorf("summary missing %q\n%s", want, out)
		}
	}
}
