package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero cooldown, got nil")
	}
	_, err = New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative cooldown, got nil")
	}
}

func TestNew_ValidCooldown(t *testing.T) {
	l, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestAllow_FirstCallAllowed(t *testing.T) {
	l, _ := New(10 * time.Second)
	if !l.Allow("container-abc") {
		t.Error("first call should be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	l, _ := New(10 * time.Second)
	l.Allow("container-abc")
	if l.Allow("container-abc") {
		t.Error("second call within cooldown should be blocked")
	}
}

func TestAllow_AfterCooldown_Allowed(t *testing.T) {
	l, _ := New(50 * time.Millisecond)
	l.Allow("container-abc")
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("container-abc") {
		t.Error("call after cooldown should be allowed")
	}
}

func TestAllow_DifferentContainersIndependent(t *testing.T) {
	l, _ := New(10 * time.Second)
	l.Allow("container-1")
	if !l.Allow("container-2") {
		t.Error("different container should be allowed independently")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l, _ := New(10 * time.Second)
	l.Allow("container-abc")
	l.Reset("container-abc")
	if !l.Allow("container-abc") {
		t.Error("after reset, container should be allowed again")
	}
}

func TestPurge_RemovesStaleEntries(t *testing.T) {
	l, _ := New(50 * time.Millisecond)
	l.Allow("container-stale")
	time.Sleep(60 * time.Millisecond)
	l.Purge()

	l.mu.Lock()
	_, exists := l.lastSent["container-stale"]
	l.mu.Unlock()

	if exists {
		t.Error("stale entry should have been purged")
	}
}
