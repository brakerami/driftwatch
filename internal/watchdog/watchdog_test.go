package watchdog_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/org/driftwatch/internal/watchdog"
)

type stubSource struct {
	mu     sync.Mutex
	events []watchdog.Event
	err    error
}

func (s *stubSource) Events(_ context.Context, _ time.Time) ([]watchdog.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.events, s.err
}

func TestNew_NilSource_ReturnsError(t *testing.T) {
	_, err := watchdog.New(nil, func(_ context.Context, _ watchdog.Event) error { return nil }, time.Second)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestNew_NilHandler_ReturnsError(t *testing.T) {
	_, err := watchdog.New(&stubSource{}, nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestNew_ZeroInterval_ReturnsError(t *testing.T) {
	_, err := watchdog.New(&stubSource{}, func(_ context.Context, _ watchdog.Event) error { return nil }, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_ValidParams_NoError(t *testing.T) {
	_, err := watchdog.New(&stubSource{}, func(_ context.Context, _ watchdog.Event) error { return nil }, time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_DispatchesEvents(t *testing.T) {
	src := &stubSource{
		events: []watchdog.Event{
			{ContainerID: "abc", ContainerName: "web", Kind: watchdog.EventRestart, OccurredAt: time.Now()},
		},
	}
	var mu sync.Mutex
	var received []watchdog.Event
	handler := func(_ context.Context, e watchdog.Event) error {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e)
		return nil
	}
	w, _ := watchdog.New(src, handler, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)
	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected at least one event dispatched")
	}
}

func TestRun_SourceError_ContinuesPolling(t *testing.T) {
	src := &stubSource{err: errors.New("docker unavailable")}
	called := 0
	handler := func(_ context.Context, _ watchdog.Event) error {
		called++
		return nil
	}
	w, _ := watchdog.New(src, handler, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)
	if called != 0 {
		t.Fatalf("handler should not be called when source errors, got %d calls", called)
	}
}
