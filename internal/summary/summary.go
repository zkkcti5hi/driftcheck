// Package summary provides aggregated statistics across multiple drift scan results.
package summary

import (
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Stats holds aggregated metrics for a collection of drift results.
type Stats struct {
	Total     int       `json:"total"`
	Drifted   int       `json:"drifted"`
	Clean     int       `json:"clean"`
	Errored   int       `json:"errored"`
	DriftRate float64   `json:"drift_rate_pct"`
	ComputedAt time.Time `json:"computed_at"`
}

// TopDrifted holds a service name and the number of drift occurrences.
type TopDrifted struct {
	Service string `json:"service"`
	Count   int    `json:"count"`
}

// Report is the full summary output.
type Report struct {
	Stats      Stats        `json:"stats"`
	TopDrifted []TopDrifted `json:"top_drifted"`
}

// Compute builds a Report from a slice of drift results.
func Compute(results []drift.Result) Report {
	counts := make(map[string]int)
	var drifted, errored int

	for _, r := range results {
		switch {
		case r.Error != "":
			errored++
		case r.Drifted:
			drifted++
			counts[r.Service]++
		}
	}

	total := len(results)
	clean := total - drifted - errored

	var rate float64
	if total > 0 {
		rate = float64(drifted) / float64(total) * 100
	}

	return Report{
		Stats: Stats{
			Total:      total,
			Drifted:    drifted,
			Clean:      clean,
			Errored:    errored,
			DriftRate:  rate,
			ComputedAt: time.Now().UTC(),
		},
		TopDrifted: rankDrifted(counts),
	}
}

// rankDrifted converts a frequency map into a sorted slice (descending).
func rankDrifted(counts map[string]int) []TopDrifted {
	result := make([]TopDrifted, 0, len(counts))
	for svc, n := range counts {
		result = append(result, TopDrifted{Service: svc, Count: n})
	}
	// simple insertion sort — result sets are small
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Count > result[j-1].Count; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result
}
