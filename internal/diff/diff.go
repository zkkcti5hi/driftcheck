// Package diff provides helpers for producing human-readable field-level
// differences between a manifest service definition and a running container.
package diff

import (
	"fmt"
	"strings"
)

// Field represents a single differing attribute between the manifest and the
// running container.
type Field struct {
	Name     string
	Expected string
	Actual   string
}

// String returns a compact one-line representation of the field diff.
func (f Field) String() string {
	return fmt.Sprintf("%s: expected %q, got %q", f.Name, f.Expected, f.Actual)
}

// Result holds all field-level differences for a single service.
type Result struct {
	Service string
	Fields  []Field
}

// HasDrift reports whether any fields differ.
func (r Result) HasDrift() bool {
	return len(r.Fields) > 0
}

// Summary returns a multi-line human-readable summary of all diffs.
func (r Result) Summary() string {
	if !r.HasDrift() {
		return fmt.Sprintf("%s: no drift detected", r.Service)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s: %d field(s) differ\n", r.Service, len(r.Fields))
	for _, f := range r.Fields {
		fmt.Fprintf(&sb, "  - %s\n", f.String())
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Compare builds a Result by comparing expected vs actual string maps keyed by
// field name. Only entries present in expected are evaluated.
func Compare(service string, expected, actual map[string]string) Result {
	r := Result{Service: service}
	for name, want := range expected {
		got, ok := actual[name]
		if !ok {
			got = "<missing>"
		}
		if want != got {
			r.Fields = append(r.Fields, Field{Name: name, Expected: want, Actual: got})
		}
	}
	return r
}
