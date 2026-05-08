// Package history provides functionality for tracking drift check runs
// over time, enabling trend analysis and change detection.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/user/driftcheck/internal/drift"
)

// Entry represents a single drift check run stored in history.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Results   []drift.Result    `json:"results"`
	DriftedCount int            `json:"drifted_count"`
	CleanCount   int            `json:"clean_count"`
}

// Store manages persisted history entries on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that persists entries under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a new history entry derived from results to disk.
func (s *Store) Save(results []drift.Result) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	}
	for _, r := range results {
		if r.Drifted {
			entry.DriftedCount++
		} else {
			entry.CleanCount++
		}
	}

	filename := entry.Timestamp.Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, filename)

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// List returns all stored history entries sorted by timestamp ascending.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}

	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: parse %s: %w", m, err)
		}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}
