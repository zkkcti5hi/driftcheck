// Package remediate provides suggestions for fixing detected configuration drift.
package remediate

import (
	"fmt"
	"strings"

	"github.com/yourusername/driftcheck/internal/drift"
)

// Suggestion holds a human-readable remediation hint for a single drift field.
type Suggestion struct {
	Service string
	Field   string
	Current string
	Expected string
	Hint    string
}

// Report is the full set of suggestions produced for a scan run.
type Report struct {
	Suggestions []Suggestion
}

// HasSuggestions returns true when at least one suggestion exists.
func (r Report) HasSuggestions() bool {
	return len(r.Suggestions) > 0
}

// Generate produces remediation suggestions for every drifted result.
func Generate(results []drift.Result) Report {
	var suggestions []Suggestion
	for _, res := range results {
		if !res.Drifted {
			continue
		}
		for _, f := range res.Fields {
			if f.Current == f.Expected {
				continue
			}
			suggestions = append(suggestions, Suggestion{
				Service:  res.Service,
				Field:    f.Name,
				Current:  f.Current,
				Expected: f.Expected,
				Hint:     buildHint(res.Service, f.Name, f.Current, f.Expected),
			})
		}
	}
	return Report{Suggestions: suggestions}
}

func buildHint(service, field, current, expected string) string {
	switch strings.ToLower(field) {
	case "image":
		return fmt.Sprintf(
			"Re-deploy service %q with image %q (running: %q). Run: docker compose up -d --force-recreate %s",
			service, expected, current, service,
		)
	case "replicas":
		return fmt.Sprintf(
			"Scale service %q to %s replicas (currently %s). Run: kubectl scale deployment/%s --replicas=%s",
			service, expected, current, service, expected,
		)
	case "env":
		return fmt.Sprintf(
			"Environment variable mismatch for service %q. Update your manifest or re-create the container with the correct env.",
			service,
		)
	default:
		return fmt.Sprintf(
			"Field %q for service %q differs (want %q, got %q). Reconcile the manifest and re-deploy.",
			field, service, expected, current,
		)
	}
}
