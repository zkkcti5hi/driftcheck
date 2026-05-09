// Package snapshot captures a point-in-time view of drift results and
// allows comparing two snapshots to surface changes between scans.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your-org/driftcheck/internal/drift"
)

// Snapshot holds drift results captured at a specific moment.
type Snapshot struct {
	CapturedAt time.Time          `json:"captured_at"`
	Results    []drift.Result     `json:"results"`
}

// Delta describes the change in drift status for a single service
// between two snapshots.
type Delta struct {
	Service  string `json:"service"`
	Previous string `json:"previous"` // "drifted", "clean", "missing"
	Current  string `json:"current"`
}

// Take creates a new Snapshot from the provided results.
func Take(results []drift.Result) Snapshot {
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}
}

// Save writes the snapshot to the given file path as JSON.
func Save(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return s, nil
}

// Compare returns the list of services whose drift status changed between
// the previous and current snapshot.
func Compare(prev, curr Snapshot) []Delta {
	prevStatus := statusMap(prev.Results)
	currStatus := statusMap(curr.Results)

	seen := make(map[string]struct{})
	var deltas []Delta

	for svc, cur := range currStatus {
		seen[svc] = struct{}{}
		pre, ok := prevStatus[svc]
		if !ok {
			pre = "missing"
		}
		if pre != cur {
			deltas = append(deltas, Delta{Service: svc, Previous: pre, Current: cur})
		}
	}

	for svc, pre := range prevStatus {
		if _, ok := seen[svc]; !ok {
			deltas = append(deltas, Delta{Service: svc, Previous: pre, Current: "missing"})
		}
	}

	return deltas
}

func statusMap(results []drift.Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		if r.Drifted {
			m[r.Service] = "drifted"
		} else {
			m[r.Service] = "clean"
		}
	}
	return m
}
