package metrics

import (
	"github.com/you/driftcheck/internal/drift"
)

// FromResults derives scan counters from a slice of drift results and records
// them against c. It returns the number of drifted, clean, and errored
// services so callers can act on them without re-iterating the slice.
func FromResults(c *Counters, results []drift.Result) (drifted, clean, errors int) {
	for _, r := range results {
		switch {
		case r.Error != nil:
			errors++
		case r.Drifted:
			drifted++
		default:
			clean++
		}
	}
	c.RecordScan(drifted, clean, errors)
	return
}
