// Package cache provides lightweight persistence for drift check results.
//
// A Cache serialises []drift.Result snapshots to a JSON file on disk so that
// successive runs of driftcheck can compare the current state against the
// previously recorded state and surface newly-introduced drift.
//
// Typical usage:
//
//	c := cache.New("/var/lib/driftcheck/last-run.json")
//
//	// Load previous snapshot (nil if first run).
//	prev, err := c.Load()
//
//	// … run detection …
//
//	// Persist current results for next run.
//	if err := c.Save(results); err != nil {
//		log.Fatalf("cache write: %v", err)
//	}
package cache
