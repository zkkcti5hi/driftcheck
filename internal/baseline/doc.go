// Package baseline provides functionality for capturing a known-good drift
// state and comparing future runs against it.
//
// Operators can save a baseline after reviewing and accepting current drift,
// then use FilterNew to surface only regressions introduced since the snapshot
// was taken. Baselines are persisted as JSON files on disk and can be committed
// to version control alongside manifests.
//
// Typical usage:
//
//	store := baseline.NewStore(".driftcheck/baseline.json")
//
//	// capture current state as the accepted baseline
//	_ = store.Save(results)
//
//	// on subsequent runs, suppress known drift
//	snap, _ := store.Load()
//	novel := baseline.FilterNew(results, snap)
package baseline
