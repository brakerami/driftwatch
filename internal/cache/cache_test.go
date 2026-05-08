package cache

import (
	"testing"
	"time"
)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := New[string](0)
	if err != ErrInvalidTTL {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestNew_ValidTTL(t *testing.T) {
	c, err := New[string](time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestSetAndGet_HitBeforeExpiry(t *testing.T) {
	c, _ := New[int](time.Minute)
	c.Set("k", 42)
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestGet_MissOnUnknownKey(t *testing.T) {
	c, _ := New[int](time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestGet_MissAfterExpiry(t *testing.T) {
	c, _ := New[string](time.Millisecond)
	now := time.Now()
	c.nowFunc = func() time.Time { return now }
	c.Set("x", "hello")
	// advance clock past TTL
	c.nowFunc = func() time.Time { return now.Add(2 * time.Millisecond) }
	_, ok := c.Get("x")
	if ok {
		t.Fatal("expected cache miss after expiry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c, _ := New[string](time.Minute)
	c.Set("key", "val")
	c.Delete("key")
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestPurge_RemovesExpiredOnly(t *testing.T) {
	c, _ := New[string](time.Millisecond)
	now := time.Now()
	c.nowFunc = func() time.Time { return now }
	c.Set("old", "a")
	c.Set("fresh", "b")
	// expire only "old"
	c.items["old"] = Entry[string]{Value: "a", ExpiresAt: now.Add(-time.Second)}
	removed := c.Purge()
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", c.Len())
	}
}

func TestLen_ReflectsAllEntries(t *testing.T) {
	c, _ := New[bool](time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("a", true)
	c.Set("b", false)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
