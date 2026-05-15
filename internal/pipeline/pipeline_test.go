package pipeline_test

import (
	"errors"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/pipeline"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Container: "web", Type: drift.DriftTypeEnv, Field: "PORT", Expected: "8080", Actual: "9090"},
		{Container: "db", Type: drift.DriftTypeImage, Field: "image", Expected: "postgres:15", Actual: "postgres:14"},
	}
}

func TestNew_NoStages_ReturnsError(t *testing.T) {
	_, err := pipeline.New()
	if err == nil {
		t.Fatal("expected error for empty stage list")
	}
}

func TestNew_NilStage_ReturnsError(t *testing.T) {
	_, err := pipeline.New(nil)
	if err == nil {
		t.Fatal("expected error for nil stage")
	}
}

func TestNew_ValidStages_NoError(t *testing.T) {
	stage := pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) { return f, nil })
	p, err := pipeline.New(stage)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestRun_PassesThroughAllStages(t *testing.T) {
	var order []int
	makeStage := func(id int) pipeline.Stage {
		return pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) {
			order = append(order, id)
			return f, nil
		})
	}
	p, _ := pipeline.New(makeStage(1), makeStage(2), makeStage(3))
	_, err := p.Run(sampleFindings())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected execution order: %v", order)
	}
}

func TestRun_StageCanFilterFindings(t *testing.T) {
	filter := pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) {
		var out []drift.Finding
		for _, x := range f {
			if x.Type == drift.DriftTypeEnv {
				out = append(out, x)
			}
		}
		return out, nil
	})
	p, _ := pipeline.New(filter)
	results, err := p.Run(sampleFindings())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Type != drift.DriftTypeEnv {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestRun_StageError_HaltsPipeline(t *testing.T) {
	boom := pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) {
		return nil, errors.New("stage failure")
	})
	neverCalled := pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) {
		t.Error("subsequent stage must not be called after error")
		return f, nil
	})
	p, _ := pipeline.New(boom, neverCalled)
	_, err := p.Run(sampleFindings())
	if err == nil {
		t.Fatal("expected error from failing stage")
	}
}

func TestRun_EmptyInput_NoError(t *testing.T) {
	passthrough := pipeline.StageFunc(func(f []drift.Finding) ([]drift.Finding, error) { return f, nil })
	p, _ := pipeline.New(passthrough)
	results, err := p.Run(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}
