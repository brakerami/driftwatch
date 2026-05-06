package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New(0, func(_ context.Context) error { return nil }, nil)
	if !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestNew_NilJob(t *testing.T) {
	_, err := New(time.Second, nil, nil)
	if !errors.Is(err, ErrNilJob) {
		t.Fatalf("expected ErrNilJob, got %v", err)
	}
}

func TestNew_ValidParams(t *testing.T) {
	s, err := New(time.Second, func(_ context.Context) error { return nil }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

// TestRun_JobCalledImmediately verifies the job fires on the first tick
// without waiting a full interval.
func TestRun_JobCalledImmediately(t *testing.T) {
	var count atomic.Int32

	s, err := New(10*time.Second, func(_ context.Context) error {
		count.Add(1)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	// Give the goroutine a moment to execute the immediate run.
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	if count.Load() < 1 {
		t.Error("expected job to be called at least once immediately")
	}
}

// TestRun_JobCalledOnTick verifies the job fires again after the interval.
func TestRun_JobCalledOnTick(t *testing.T) {
	var count atomic.Int32

	s, err := New(30*time.Millisecond, func(_ context.Context) error {
		count.Add(1)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	// Wait long enough for at least two executions (immediate + 1 tick).
	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	if count.Load() < 2 {
		t.Errorf("expected at least 2 job calls, got %d", count.Load())
	}
}

// TestRun_JobErrorDoesNotStop verifies that a job error is logged but the
// scheduler keeps running.
func TestRun_JobErrorDoesNotStop(t *testing.T) {
	var count atomic.Int32
	boom := errors.New("boom")

	s, err := New(30*time.Millisecond, func(_ context.Context) error {
		count.Add(1)
		return boom
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	if count.Load() < 2 {
		t.Errorf("expected scheduler to continue despite errors, got %d calls", count.Load())
	}
}
