// Package pipeline chains drift-processing stages into a single
// ordered execution path. Each stage receives the findings produced
// by the previous one, allowing arbitrary transformation, filtering,
// or enrichment without coupling individual components.
package pipeline

import (
	"errors"
	"fmt"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Stage is a single step in the pipeline. It receives a slice of
// findings and returns a (possibly modified) slice or an error.
type Stage interface {
	Process(findings []drift.Finding) ([]drift.Finding, error)
}

// StageFunc is a convenience adapter that lets a plain function
// satisfy the Stage interface.
type StageFunc func([]drift.Finding) ([]drift.Finding, error)

func (f StageFunc) Process(findings []drift.Finding) ([]drift.Finding, error) {
	return f(findings)
}

// Pipeline executes a fixed list of stages in order.
type Pipeline struct {
	stages []Stage
}

// New constructs a Pipeline from the provided stages. At least one
// stage must be supplied.
func New(stages ...Stage) (*Pipeline, error) {
	if len(stages) == 0 {
		return nil, errors.New("pipeline: at least one stage is required")
	}
	for i, s := range stages {
		if s == nil {
			return nil, fmt.Errorf("pipeline: stage at index %d is nil", i)
		}
	}
	return &Pipeline{stages: stages}, nil
}

// Run passes findings through every stage sequentially. If any stage
// returns an error the pipeline halts and propagates the error
// together with the index of the failing stage.
func (p *Pipeline) Run(findings []drift.Finding) ([]drift.Finding, error) {
	current := findings
	for i, s := range p.stages {
		var err error
		current, err = s.Process(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline: stage %d failed: %w", i, err)
		}
	}
	return current, nil
}

// Len returns the number of stages registered in the pipeline.
func (p *Pipeline) Len() int { return len(p.stages) }
