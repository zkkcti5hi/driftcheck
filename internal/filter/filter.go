// Package filter provides utilities for narrowing drift detection
// to a specific subset of services or namespaces.
package filter

import "strings"

// Options holds the criteria used to filter services before drift
// detection is performed.
type Options struct {
	// Services is an explicit allow-list of service names. When empty
	// every service is included.
	Services []string

	// Labels is a map of key/value pairs; a service must carry ALL of
	// them to be included. When empty the check is skipped.
	Labels map[string]string
}

// Filter returns only the service names from candidates that satisfy
// the Options criteria.
func Filter(candidates []string, opts Options) []string {
	if len(opts.Services) == 0 && len(opts.Labels) == 0 {
		return candidates
	}

	allowSet := make(map[string]struct{}, len(opts.Services))
	for _, s := range opts.Services {
		allowSet[strings.TrimSpace(s)] = struct{}{}
	}

	var out []string
	for _, c := range candidates {
		if len(allowSet) > 0 {
			if _, ok := allowSet[c]; !ok {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}

// MatchLabels reports whether containerLabels satisfies every
// key/value pair in required.
func MatchLabels(containerLabels, required map[string]string) bool {
	for k, v := range required {
		if containerLabels[k] != v {
			return false
		}
	}
	return true
}
