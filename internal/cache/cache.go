// Package cache provides a simple file-backed cache for storing drift
// results between runs, allowing driftcheck to report on changes over time.
package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/user/driftcheck/internal/drift"
)

// Entry is a single cached drift result snapshot.
type Entry struct {
	Timestamp time.Time          `json:"timestamp"`
	Results   []drift.Result     `json:"results"`
}

// Cache persists drift results to a JSON file on disk.
type Cache struct {
	path string
}

// New returns a Cache that stores data at the given file path.
func New(path string) *Cache {
	return &Cache{path: path}
}

// Save writes the provided results to the cache file, overwriting any
// previous snapshot.
func (c *Cache) Save(results []drift.Result) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}

// Load reads the most recent cached entry from disk.
// Returns (nil, nil) when no cache file exists yet.
func (c *Cache) Load() (*Entry, error) {
	data, err := os.ReadFile(c.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Clear removes the cache file if it exists.
func (c *Cache) Clear() error {
	err := os.Remove(c.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
