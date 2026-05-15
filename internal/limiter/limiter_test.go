package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0, got nil")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative capacity, got nil")
	}
}

func TestNew_ValidCapacity(t *testing.T) {
	l, err := New(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Cap() != 4 {
		t.Fatalf("expected cap 4, got %d", l.Cap())
	}
}

func TestAcquireRelease_InFlightTracking(t *testing.T) {
	l, _ := New(3)

	l.Acquire()
	l.Acquire()
	if l.InFlight() != 2 {
		t.Fatalf("expected 2 in-flight, got %d", l.InFlight())
	}
	l.Release()
	if l.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight after release, got %d", l.InFlight())
	}
	l.Release()
	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all released, got %d", l.InFlight())
	}
}

func TestTryAcquire_SucceedsWhenSlotAvailable(t *testing.T) {
	l, _ := New(1)
	if !l.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed on empty limiter")
	}
	if l.TryAcquire() {
		t.Fatal("expected TryAcquire to fail when limiter is full")
	}
	l.Release()
	if !l.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed after release")
	}
	l.Release()
}

func TestAcquire_BlocksUntilRelease(t *testing.T) {
	l, _ := New(1)
	l.Acquire()

	done := make(chan struct{})
	go func() {
		l.Acquire()
		close(done)
		l.Release()
	}()

	// goroutine should be blocked; give it a moment
	select {
	case <-done:
		t.Fatal("goroutine should have been blocked")
	case <-time.After(30 * time.Millisecond):
	}

	l.Release() // unblock the goroutine
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("goroutine did not unblock after release")
	}
}

func TestConcurrentAcquire_RespectsCapacity(t *testing.T) {
	const cap = 4
	l, _ := New(cap)

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Acquire()
			defer l.Release()
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()
	if peak > cap {
		t.Fatalf("peak concurrency %d exceeded cap %d", peak, cap)
	}
}
