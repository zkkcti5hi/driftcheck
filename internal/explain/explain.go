// Package explain provides human-readable explanations for detected drift,
// describing why a particular field is considered drifted and suggesting
// remediation steps.
package explain

import (
	"fmt"
	"strings"

	"github.com/driftcheck/internal/drift"
)

// Explanation holds a human-readable description of a single drift finding
// along with a suggested remediation action.
type Explanation struct {
	Service     string
	Field       string
	Reason      string
	Remediation string
}

// Explain produces Explanation entries for every drifted field found in the
// provided drift results. Results that have no drift are silently skipped.
func Explain(results []drift.Result) []Explanation {
	var out []Explanation
	for _, r := range results {
		if !r.Drifted {
			continue
		}
		for _, f := range r.Fields {
			out = append(out, buildExplanation(r.Service, f))
		}
	}
	return out
}

// buildExplanation constructs an Explanation for a single drifted field.
func buildExplanation(service string, field drift.FieldDiff) Explanation {
	switch strings.ToLower(field.Name) {
	case "image":
		return Explanation{
			Service: service,
			Field:   field.Name,
			Reason: fmt.Sprintf(
				"running image %q does not match manifest image %q",
				field.Got, field.Want,
			),
			Remediation: "re-deploy the service to pull the image declared in the manifest",
		}
	case "env":
		return Explanation{
			Service: service,
			Field:   field.Name,
			Reason: fmt.Sprintf(
				"environment variables differ: got %q, want %q",
				field.Got, field.Want,
			),
			Remediation: "update the running container's environment to match the manifest, then restart",
		}
	case "replicas":
		return Explanation{
			Service: service,
			Field:   field.Name,
			Reason: fmt.Sprintf(
				"replica count is %q but manifest specifies %q",
				field.Got, field.Want,
			),
			Remediation: "scale the deployment to match the replica count in the manifest",
		}
	default:
		return Explanation{
			Service: service,
			Field:   field.Name,
			Reason: fmt.Sprintf(
				"field %q has value %q but manifest declares %q",
				field.Name, field.Got, field.Want,
			),
			Remediation: "reconcile the running configuration with the manifest and redeploy if necessary",
		}
	}
}
