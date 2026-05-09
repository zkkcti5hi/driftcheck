package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/driftcheck/internal/drift"
	"github.com/your-org/driftcheck/internal/snapshot"
)

func makeResult(service string, drifted bool) drift.Result {
	return drift.Result{Service: service, Drifted: drifted}
}

func TestTake_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take([]drift.Result{makeResult("web", false)})
	after := time.Now().UTC()

	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v out of expected range [%v, %v]", s.CapturedAt, before, after)
	}
	if len(s.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(s.Results))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.Take([]drift.Result{
		makeResult("api", true),
		makeResult("db", false),
	})

	if err := snapshot.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(loaded.Results))
	}
	if loaded.Results[0].Service != "api" || !loaded.Results[0].Drifted {
		t.Errorf("unexpected first result: %+v", loaded.Results[0])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
}

func TestCompare_DetectsChanges(t *testing.T) {
	prev := snapshot.Take([]drift.Result{
		makeResult("web", false),
		makeResult("api", true),
	})
	curr := snapshot.Take([]drift.Result{
		makeResult("web", true),  // newly drifted
		makeResult("api", false), // resolved
	})

	deltas := snapshot.Compare(prev, curr)
	if len(deltas) != 2 {
		t.Fatalf("expected 2 deltas, got %d", len(deltas))
	}
}

func TestCompare_NoDrift_NoDeltas(t *testing.T) {
	results := []drift.Result{makeResult("web", false), makeResult("api", false)}
	prev := snapshot.Take(results)
	curr := snapshot.Take(results)

	deltas := snapshot.Compare(prev, curr)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas, got %d", len(deltas))
	}
}

func TestCompare_ServiceDisappears(t *testing.T) {
	prev := snapshot.Take([]drift.Result{makeResult("web", false), makeResult("cache", false)})
	curr := snapshot.Take([]drift.Result{makeResult("web", false)})

	deltas := snapshot.Compare(prev, curr)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta for disappeared service, got %d", len(deltas))
	}
	if deltas[0].Service != "cache" || deltas[0].Current != "missing" {
		t.Errorf("unexpected delta: %+v", deltas[0])
	}
}
