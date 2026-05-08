package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftcheck/internal/cache"
	"github.com/user/driftcheck/internal/drift"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{Service: "web", Drifted: true, Reasons: []string{"image mismatch"}},
		{Service: "db", Drifted: false},
	}
}

func TestCache_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(filepath.Join(dir, "drift.cache.json"))

	before := time.Now().UTC().Truncate(time.Second)
	if err := c.Save(sampleResults()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := c.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if !entry.Timestamp.After(before) && !entry.Timestamp.Equal(before) {
		t.Errorf("timestamp %v is before save time %v", entry.Timestamp, before)
	}
	if len(entry.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(entry.Results))
	}
	if entry.Results[0].Service != "web" || !entry.Results[0].Drifted {
		t.Errorf("unexpected first result: %+v", entry.Results[0])
	}
}

func TestCache_Load_NoFile(t *testing.T) {
	c := cache.New(filepath.Join(t.TempDir(), "missing.json"))
	entry, err := c.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entry != nil {
		t.Fatalf("expected nil entry for missing file, got %+v", entry)
	}
}

func TestCache_Clear(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "drift.cache.json")
	c := cache.New(path)

	if err := c.Save(sampleResults()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("cache file should exist after Save: %v", err)
	}
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("cache file should be removed after Clear")
	}
}

func TestCache_Clear_NoFile(t *testing.T) {
	c := cache.New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear on missing file should not error, got %v", err)
	}
}

func TestCache_Save_CreatesIntermediateDirs(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(filepath.Join(dir, "sub", "dir", "drift.json"))
	if err := c.Save(sampleResults()); err != nil {
		t.Fatalf("Save should create intermediate dirs: %v", err)
	}
}
