package throttle

import (
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_InvalidMaxTokens(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero maxTokens")
	}
}

func TestNew_ValidParams(t *testing.T) {
	th, err := New(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}

func TestAllow_WithinBudget(t *testing.T) {
	th, _ := New(time.Second, 3)
	for i := 0; i < 3; i++ {
		if err := th.Allow(); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsBudget(t *testing.T) {
	th, _ := New(time.Second, 2)
	_ = th.Allow()
	_ = th.Allow()
	if err := th.Allow(); err != ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	th, _ := New(50*time.Millisecond, 1)
	th.now = func() time.Time { return now }

	_ = th.Allow()
	if err := th.Allow(); err != ErrThrottled {
		t.Fatal("expected throttle before window expires")
	}

	// Advance past the window.
	th.now = func() time.Time { return now.Add(60 * time.Millisecond) }
	if err := th.Allow(); err != nil {
		t.Fatalf("expected allow after window reset, got %v", err)
	}
}

func TestRemaining_DecrementsOnAllow(t *testing.T) {
	th, _ := New(time.Second, 4)
	if th.Remaining() != 4 {
		t.Fatalf("expected 4 remaining, got %d", th.Remaining())
	}
	_ = th.Allow()
	if th.Remaining() != 3 {
		t.Fatalf("expected 3 remaining, got %d", th.Remaining())
	}
}

func TestResetAt_IsInFuture(t *testing.T) {
	th, _ := New(time.Minute, 5)
	if !th.ResetAt().After(time.Now()) {
		t.Fatal("expected ResetAt to be in the future")
	}
}
