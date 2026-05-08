package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftcheck/internal/baseline"
	"github.com/yourorg/driftcheck/internal/diff"
	"github.com/yourorg/driftcheck/internal/drift"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{
			Service: "web",
			Drifts: []diff.Field{
				{Name: "image", Manifest: "nginx:1.24", Running: "nginx:1.25"},
			},
		},
		{
			Service: "db",
			Drifts:  nil,
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := baseline.NewStore(filepath.Join(dir, "baseline.json"))

	results := sampleResults()
	if err := store.Save(results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_NoBaseline(t *testing.T) {
	store := baseline.NewStore(filepath.Join(t.TempDir(), "missing.json"))
	_, err := store.Load()
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "baseline.json")
	_ = os.WriteFile(p, []byte("not json{"), 0o644)
	store := baseline.NewStore(p)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestFilterNew_AllNew(t *testing.T) {
	snap := &baseline.Snapshot{Results: []drift.Result{}}
	novel := baseline.FilterNew(sampleResults(), snap)
	if len(novel) != 2 {
		t.Errorf("expected 2 novel results, got %d", len(novel))
	}
}

func TestFilterNew_NoneNew(t *testing.T) {
	results := sampleResults()
	snap := &baseline.Snapshot{Results: results}
	novel := baseline.FilterNew(results, snap)
	if len(novel) != 0 {
		t.Errorf("expected 0 novel results, got %d", len(novel))
	}
}

func TestFilterNew_PartialOverlap(t *testing.T) {
	results := sampleResults()
	// baseline only contains the first result
	snap := &baseline.Snapshot{Results: results[:1]}
	novel := baseline.FilterNew(results, snap)
	if len(novel) != 1 {
		t.Errorf("expected 1 novel result, got %d", len(novel))
	}
	if novel[0].Service != "db" {
		t.Errorf("expected service 'db', got %q", novel[0].Service)
	}
}
