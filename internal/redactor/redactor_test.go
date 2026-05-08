package redactor_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/redactor"
)

func TestNew_DefaultPatterns(t *testing.T) {
	r, err := redactor.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_CustomPattern(t *testing.T) {
	r, err := redactor.New([]string{"mysecret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsSensitive("MY_MYSECRET_KEY") {
		t.Error("expected MY_MYSECRET_KEY to be sensitive")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	r, _ := redactor.New(nil)
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"db_password", true},
		{"API_TOKEN", true},
		{"AUTH_HEADER", true},
		{"PORT", false},
		{"LOG_LEVEL", false},
	}
	for _, tc := range cases {
		got := r.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestMaskValue_SensitiveReturnsRedacted(t *testing.T) {
	r, _ := redactor.New(nil)
	got := r.MaskValue("DB_PASSWORD", "super-secret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMaskValue_SafeKeyPassesThrough(t *testing.T) {
	r, _ := redactor.New(nil)
	got := r.MaskValue("PORT", "8080")
	if got != "8080" {
		t.Errorf("expected 8080, got %q", got)
	}
}

func TestMaskValue_EmptyValueSensitiveKey(t *testing.T) {
	r, _ := redactor.New(nil)
	// An empty value for a sensitive key should still be redacted to avoid
	// leaking the fact that the variable is unset vs. intentionally blank.
	got := r.MaskValue("DB_PASSWORD", "")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED] for empty sensitive value, got %q", got)
	}
}

func TestMaskEnv_MixedEntries(t *testing.T) {
	r, _ := redactor.New(nil)
	input := []string{
		"PORT=8080",
		"DB_PASSWORD=hunter2",
		"API_TOKEN=abc123",
		"LOG_LEVEL=debug",
	}
	got := r.MaskEnv(input)
	expected := []string{
		"PORT=8080",
		"DB_PASSWORD=[REDACTED]",
		"API_TOKEN=[REDACTED]",
		"LOG_LEVEL=debug",
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("entry %d: got %q, want %q", i, got[i], e)
		}
	}
}

func TestMaskEnv_MalformedEntry_PassedThrough(t *testing.T) {
	r, _ := redactor.New(nil)
	input := []string{"NOEQUALSIGN"}
	got := r.MaskEnv(input)
	if got[0] != "NOEQUALSIGN" {
		t.Errorf("expected malformed entry unchanged, got %q", got[0])
	}
}
