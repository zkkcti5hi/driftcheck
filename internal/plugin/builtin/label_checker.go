// Package builtin ships the checkers that are bundled with driftcheck.
package builtin

import (
	"fmt"

	"github.com/yourorg/driftcheck/internal/drift"
)

// LabelChecker flags running containers that are missing a required label.
type LabelChecker struct {
	// RequiredLabel is the Docker label key every container must carry.
	RequiredLabel string
}

// NewLabelChecker returns a LabelChecker that enforces the given label key.
func NewLabelChecker(label string) *LabelChecker {
	return &LabelChecker{RequiredLabel: label}
}

// Name implements plugin.Checker.
func (l *LabelChecker) Name() string { return "builtin/label-checker" }

// Check implements plugin.Checker. It marks any result whose running
// container is missing RequiredLabel as drifted.
func (l *LabelChecker) Check(results []drift.Result) ([]drift.Result, error) {
	if l.RequiredLabel == "" {
		return results, nil
	}
	for i, r := range results {
		if r.ContainerID == "" {
			// No running container – nothing to inspect.
			continue
		}
		if _, ok := r.Labels[l.RequiredLabel]; !ok {
			results[i].Drifted = true
			results[i].Reasons = append(
				results[i].Reasons,
				fmt.Sprintf("missing required label %q", l.RequiredLabel),
			)
		}
	}
	return results, nil
}
