package digest_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/digest"
)

func baseState() digest.ContainerState {
	return digest.ContainerState{
		Image: "nginx:1.25",
		Env:   map[string]string{"PORT": "8080", "ENV": "prod"},
		Labels: map[string]string{"app": "web"},
		Ports:  []string{"443/tcp", "80/tcp"},
	}
}

func TestCompute_ReturnsSHA256Hex(t *testing.T) {
	s := baseState()
	d, err := digest.Compute(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d) != 64 {
		t.Fatalf("expected 64-char hex digest, got %d chars: %s", len(d), d)
	}
}

func TestCompute_Deterministic(t *testing.T) {
	s := baseState()
	d1, _ := digest.Compute(s)
	d2, _ := digest.Compute(s)
	if d1 != d2 {
		t.Fatalf("digest not deterministic: %s vs %s", d1, d2)
	}
}

func TestCompute_MapOrderIndependent(t *testing.T) {
	a := digest.ContainerState{
		Image:  "alpine:3.19",
		Env:    map[string]string{"A": "1", "B": "2"},
		Labels: map[string]string{},
		Ports:  nil,
	}
	// Go map iteration is random; computing twice should still match.
	d1, _ := digest.Compute(a)
	d2, _ := digest.Compute(a)
	if d1 != d2 {
		t.Fatalf("map ordering caused digest mismatch")
	}
}

func TestEqual_SameState(t *testing.T) {
	s := baseState()
	ok, err := digest.Equal(s, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected equal states to match")
	}
}

func TestEqual_DifferentImage(t *testing.T) {
	a := baseState()
	b := baseState()
	b.Image = "nginx:1.26"
	ok, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected states with different images to differ")
	}
}

func TestEqual_DifferentEnv(t *testing.T) {
	a := baseState()
	b := baseState()
	b.Env["PORT"] = "9090"
	ok, _ := digest.Equal(a, b)
	if ok {
		t.Fatal("expected env change to produce different digest")
	}
}

func TestEqual_PortOrderIrrelevant(t *testing.T) {
	a := baseState()
	b := baseState()
	b.Ports = []string{"80/tcp", "443/tcp"} // reversed
	ok, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected port order to be irrelevant for digest equality")
	}
}
