package baseline_test

import (
	"path/filepath"
	"testing"

	"github.com/yourorg/driftcheck/internal/baseline"
	"github.com/yourorg/driftcheck/internal/diff"
	"github.com/yourorg/driftcheck/internal/drift"
)

// TestRoundTrip_PreservesFieldValues ensures that field-level drift details
// survive a full save→load→filter cycle without mutation.
func TestRoundTrip_PreservesFieldValues(t *testing.T) {
	dir := t.TempDir()
	store := baseline.NewStore(filepath.Join(dir, "bl.json"))

	original := []drift.Result{
		{
			Service:     "api",
			ContainerID: "abc123",
			Drifts: []diff.Field{
				{Name: "image", Manifest: "myapp:v1", Running: "myapp:v2"},
				{Name: "env.LOG_LEVEL", Manifest: "info", Running: "debug"},
			},
		},
	}

	if err := store.Save(original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Same results should yield zero novel entries
	novel := baseline.FilterNew(original, snap)
	if len(novel) != 0 {
		t.Errorf("expected no novel results after round-trip, got %d", len(novel))
	}

	// A new result with different drift should surface as novel
	newResult := drift.Result{
		Service:     "api",
		ContainerID: "abc123",
		Drifts: []diff.Field{
			{Name: "image", Manifest: "myapp:v1", Running: "myapp:v3"},
		},
	}
	novel = baseline.FilterNew([]drift.Result{newResult}, snap)
	if len(novel) != 1 {
		t.Errorf("expected 1 novel result for changed drift, got %d", len(novel))
	}
}

// TestFilterNew_EmptyCurrentResults returns empty slice gracefully.
func TestFilterNew_EmptyCurrentResults(t *testing.T) {
	snap := &baseline.Snapshot{Results: sampleResults()}
	novel := baseline.FilterNew([]drift.Result{}, snap)
	if novel != nil && len(novel) != 0 {
		t.Errorf("expected empty slice, got %v", novel)
	}
}
