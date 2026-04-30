package drift

import (
	"fmt"
	"io"
	"strings"
)

// Summary holds aggregated statistics for a drift check run.
type Summary struct {
	Total   int
	Drifted int
	Clean   int
}

// Summarise computes a Summary from a slice of DriftResult.
func Summarise(results []DriftResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
	}
	return s
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, results []DriftResult) {
	for _, r := range results {
		status := "OK"
		if r.Drifted {
			status = "DRIFTED"
		}
		fmt.Fprintf(w, "[%s] service=%s container=%s\n", status, r.ServiceName, shortID(r.ContainerID))
		for _, issue := range r.Issues {
			fmt.Fprintf(w, "  - %s\n", issue)
		}
	}

	s := Summarise(results)
	fmt.Fprintf(w, "\nSummary: %d total, %d clean, %d drifted\n", s.Total, s.Clean, s.Drifted)
}

// shortID returns the first 12 characters of a container ID, or the full
// string if it is shorter.
func shortID(id string) string {
	if len(id) <= 12 {
		if id == "" {
			return "<none>"
		}
		return id
	}
	return strings.ToLower(id[:12])
}
