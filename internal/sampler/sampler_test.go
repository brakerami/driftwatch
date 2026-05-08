package sampler

import (
	"testing"

	"github.com/driftwatch/internal/drift"
)

func sampleFindings(n int) []drift.Finding {
	out := make([]drift.Finding, n)
	for i := range out {
		out[i] = drift.Finding{
			Container: "web",
			Type:      drift.DriftTypeEnv,
			Field:     "PORT",
			Expected:  "8080",
			Actual:    "9090",
		}
	}
	return out
}

func TestNew_InvalidRate_TooHigh(t *testing.T) {
	_, err := New(1.5)
	if err == nil {
		t.Fatal("expected error for rate > 1.0")
	}
}

func TestNew_InvalidRate_Negative(t *testing.T) {
	_, err := New(-0.1)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Errorf("expected rate 0.5, got %v", s.Rate())
	}
}

func TestSample_RateZero_DropsAll(t *testing.T) {
	s, _ := New(0.0)
	result := s.Sample(sampleFindings(20))
	if len(result) != 0 {
		t.Errorf("expected 0 findings, got %d", len(result))
	}
}

func TestSample_RateOne_PassesAll(t *testing.T) {
	s, _ := New(1.0)
	input := sampleFindings(10)
	result := s.Sample(input)
	if len(result) != len(input) {
		t.Errorf("expected %d findings, got %d", len(input), len(result))
	}
}

func TestSample_EmptyInput_ReturnsNil(t *testing.T) {
	s, _ := New(0.5)
	result := s.Sample(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSample_PartialRate_ReducesCount(t *testing.T) {
	s, _ := New(0.5)
	input := sampleFindings(1000)
	result := s.Sample(input)
	// With 1000 samples at 50%, expect roughly 300–700 to pass.
	if len(result) < 300 || len(result) > 700 {
		t.Errorf("unexpected sample count %d (expected ~500)", len(result))
	}
}
