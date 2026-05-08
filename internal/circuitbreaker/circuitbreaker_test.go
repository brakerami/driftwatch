package circuitbreaker

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Second, 1)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(3, 0, 1)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}

func TestNew_InvalidProbeSuccess(t *testing.T) {
	_, err := New(3, time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero probeSuccess")
	}
}

func TestNew_ValidParams(t *testing.T) {
	b, err := New(3, time.Second, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.CurrentState() != StateClosed {
		t.Errorf("expected StateClosed, got %v", b.CurrentState())
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b, _ := New(3, time.Second, 1)
	if !b.Allow() {
		t.Error("expected Allow() == true in closed state")
	}
}

func TestRecordFailure_TripsAfterThreshold(t *testing.T) {
	b, _ := New(3, time.Second, 1)
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != StateClosed {
		t.Error("should still be closed after 2 failures with threshold 3")
	}
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Errorf("expected StateOpen after threshold failures, got %v", b.CurrentState())
	}
}

func TestAllow_OpenStateRejectsCalls(t *testing.T) {
	b, _ := New(1, 10*time.Second, 1)
	b.RecordFailure()
	if b.Allow() {
		t.Error("expected Allow() == false in open state")
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b, _ := New(1, 10*time.Millisecond, 1)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if !b.Allow() {
		t.Error("expected Allow() == true after cooldown (half-open)")
	}
	if b.CurrentState() != StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", b.CurrentState())
	}
}

func TestRecordSuccess_ClosesFromHalfOpen(t *testing.T) {
	b, _ := New(1, 10*time.Millisecond, 2)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	b.Allow() // transition to half-open
	b.RecordSuccess()
	if b.CurrentState() != StateHalfOpen {
		t.Error("should still be half-open after 1 success with probeSuccess=2")
	}
	b.RecordSuccess()
	if b.CurrentState() != StateClosed {
		t.Errorf("expected StateClosed after probe successes, got %v", b.CurrentState())
	}
}

func TestRecordFailure_InHalfOpen_ReOpens(t *testing.T) {
	b, _ := New(1, 10*time.Millisecond, 2)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	b.Allow() // half-open
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Errorf("expected StateOpen after failure in half-open, got %v", b.CurrentState())
	}
}
