package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/driftcheck/internal/drift"
	"github.com/user/driftcheck/internal/history"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{Service: "web", Drifted: true, Reasons: []string{"image mismatch"}},
		{Service: "db", Drifted: false},
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/history"
	_, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_And_List(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	if err := store.Save(sampleResults()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.DriftedCount != 1 {
		t.Errorf("DriftedCount: want 1, got %d", e.DriftedCount)
	}
	if e.CleanCount != 1 {
		t.Errorf("CleanCount: want 1, got %d", e.CleanCount)
	}
	if len(e.Results) != 2 {
		t.Errorf("Results: want 2, got %d", len(e.Results))
	}
}

func TestList_Empty(t *testing.T) {
	store, _ := history.NewStore(t.TempDir())
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestList_SortedByTimestamp(t *testing.T) {
	store, _ := history.NewStore(t.TempDir())

	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Millisecond)
		if err := store.Save(sampleResults()); err != nil {
			t.Fatalf("Save %d: %v", i, err)
		}
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if !entries[i].Timestamp.After(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted: index %d timestamp not after index %d", i, i-1)
		}
	}
}
