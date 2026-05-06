package alerter_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/driftwatch/internal/alerter"
	"github.com/user/driftwatch/internal/drift"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Type: drift.DriftTypeEnv, Field: "LOG_LEVEL", Expected: "info", Actual: "debug"},
	}
}

func TestNew_EmptyURL_ReturnsError(t *testing.T) {
	_, err := alerter.New("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNew_ValidURL_NoError(t *testing.T) {
	_, err := alerter.New("http://example.com/hook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSend_NoFindings_SkipsRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	a, _ := alerter.New(ts.URL)
	if err := a.Send(context.Background(), "myapp", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty findings")
	}
}

func TestSend_PostsCorrectPayload(t *testing.T) {
	var received alerter.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a, _ := alerter.New(ts.URL)
	findings := sampleFindings()
	if err := a.Send(context.Background(), "web", findings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Container != "web" {
		t.Errorf("container: got %q, want %q", received.Container, "web")
	}
	if len(received.Findings) != 1 {
		t.Errorf("findings count: got %d, want 1", len(received.Findings))
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestSend_Non2xx_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a, _ := alerter.New(ts.URL)
	err := a.Send(context.Background(), "svc", sampleFindings())
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}
