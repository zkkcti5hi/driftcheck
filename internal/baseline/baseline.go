// Package baseline captures and compares drift results against a known-good
// baseline snapshot, allowing operators to suppress expected drift.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Snapshot holds a saved baseline of drift results.
type Snapshot struct {
	CapturedAt time.Time          `json:"captured_at"`
	Results    []drift.Result     `json:"results"`
}

// Store manages reading and writing baseline snapshots to disk.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the supplied results as the current baseline.
func (s *Store) Save(results []drift.Result) error {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the baseline snapshot from disk.
// Returns ErrNoBaseline if no baseline has been saved yet.
func (s *Store) Load() (*Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoBaseline
		}
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// ErrNoBaseline is returned when no baseline file exists.
var ErrNoBaseline = errors.New("no baseline snapshot found")

// FilterNew returns only those results that represent NEW drift not present
// in the baseline snapshot. A result is considered known if the service name
// and every drift field match.
func FilterNew(current []drift.Result, snap *Snapshot) []drift.Result {
	known := make(map[string]struct{}, len(snap.Results))
	for _, r := range snap.Results {
		known[baselineKey(r)] = struct{}{}
	}
	var novel []drift.Result
	for _, r := range current {
		if _, ok := known[baselineKey(r)]; !ok {
			novel = append(novel, r)
		}
	}
	return novel
}

func baselineKey(r drift.Result) string {
	b, _ := json.Marshal(struct {
		Service string
		Drifts  interface{}
	}{r.Service, r.Drifts})
	return string(b)
}
