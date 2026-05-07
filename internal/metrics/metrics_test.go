package metrics_test

import (
	"sync"
	"testing"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestCounter_IncAndValue(t *testing.T) {
	r := metrics.New()
	c := r.Counter("scans_total")
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounter_Add(t *testing.T) {
	r := metrics.New()
	c := r.Counter("findings_total")
	c.Add(5)
	if c.Value() != 5 {
		t.Fatalf("expected 5, got %d", c.Value())
	}
}

func TestRegistry_SameNameReturnsSameCounter(t *testing.T) {
	r := metrics.New()
	a := r.Counter("x")
	b := r.Counter("x")
	a.Inc()
	if b.Value() != 1 {
		t.Fatal("expected same counter instance")
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := metrics.New()
	r.Counter("a").Add(3)
	r.Counter("b").Inc()
	snap := r.Snapshot()
	if snap["a"] != 3 || snap["b"] != 1 {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestRegistry_Reset(t *testing.T) {
	r := metrics.New()
	r.Counter("a").Add(10)
	r.Reset()
	if r.Counter("a").Value() != 0 {
		t.Fatal("expected counter to be reset to 0")
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	r := metrics.New()
	c := r.Counter("concurrent")
	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, c.Value())
	}
}
