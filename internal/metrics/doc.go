// Package metrics exposes lightweight, goroutine-safe counters that accumulate
// statistics across driftcheck scan runs.
//
// # Usage
//
// Record the outcome of each scan via the package-level Global singleton:
//
//	metrics.FromResults(metrics.Global(), results)
//
// Print a human-readable summary at any time:
//
//	metrics.Global().Write(os.Stdout)
//
// In tests, call metrics.Reset() at the start of each test function to
// obtain a clean slate without sharing state between test cases.
package metrics
