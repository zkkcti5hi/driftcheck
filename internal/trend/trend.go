// Package trend analyses drift results over time to surface
// recurring or worsening drift patterns across scan history.
package trend

import (
	"sort"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Entry represents a single historical scan result set with a timestamp.
type Entry struct {
	Timestamp time.Time
	Results   []drift.Result
}

// ServiceTrend summarises drift behaviour for a single service over time.
type ServiceTrend struct {
	Service      string
	TotalScans   int
	DriftedScans int
	// DriftRate is the fraction of scans where drift was detected (0.0–1.0).
	DriftRate float64
	// Worsening is true when the most recent scan is drifted and the
	// previous scan was clean.
	Worsening bool
}

// Analyse computes per-service trend data from a slice of historical entries.
// Entries should be provided in chronological order (oldest first).
func Analyse(entries []Entry) []ServiceTrend {
	type record struct {
		total   int
		drifted int
		states  []bool // true = drifted, ordered chronologically
	}

	services := map[string]*record{}

	for _, e := range entries {
		for _, r := range e.Results {
			rec, ok := services[r.Service]
			if !ok {
				rec = &record{}
				services[r.Service] = rec
			}
			rec.total++
			if r.Drifted {
				rec.drifted++
			}
			rec.states = append(rec.states, r.Drifted)
		}
	}

	trends := make([]ServiceTrend, 0, len(services))
	for svc, rec := range services {
		rate := 0.0
		if rec.total > 0 {
			rate = float64(rec.drifted) / float64(rec.total)
		}

		worsening := false
		if len(rec.states) >= 2 {
			last := rec.states[len(rec.states)-1]
			prev := rec.states[len(rec.states)-2]
			worsening = last && !prev
		}

		trends = append(trends, ServiceTrend{
			Service:      svc,
			TotalScans:   rec.total,
			DriftedScans: rec.drifted,
			DriftRate:    rate,
			Worsening:    worsening,
		})
	}

	sort.Slice(trends, func(i, j int) bool {
		if trends[i].DriftRate != trends[j].DriftRate {
			return trends[i].DriftRate > trends[j].DriftRate
		}
		return trends[i].Service < trends[j].Service
	})

	return trends
}
