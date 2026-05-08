// Package diff compares two snapshots of a service's runtime state and
// produces a structured result describing which fields have changed.
//
// Each Field records the field name together with the observed (running)
// value and the expected (manifest) value. A Result aggregates all
// differing fields for a single service and exposes helpers for
// summarising drift in human-readable or machine-readable form.
//
// Example:
//
//	result := diff.Compare("web", expected, observed)
//	if result.HasDrift() {
//		fmt.Println(result.Summary())
//	}
package diff
