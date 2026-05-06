package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/history"
)

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Type: drift.DriftTypeEnv, Field: "LOG_LEVEL", Expected: "info", Actual: "debug"},
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subDir := dir + "/nested/history"
	_, err := history.NewStore(subDir)
	if err != nil {
		t.Fatalf("NewStore: unexpected error: %v", err)
	}
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestRecord_And_List_RoundTrip(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	const container = "web"
	findings := sampleFindings()

	if err := store.Record(container, findings); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := store.List(container)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Container != container {
		t.Errorf("container: want %q, got %q", container, entries[0].Container)
	}
	if len(entries[0].Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(entries[0].Findings))
	}
}

func TestList_Empty_ReturnsNil(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	entries, err := store.List("nonexistent")
	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestList_SortedByTimestamp(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	const container = "worker"
	for i := 0; i < 3; i++ {
		if err := store.Record(container, sampleFindings()); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Millisecond)
	}

	entries, err := store.List(container)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.Before(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}
