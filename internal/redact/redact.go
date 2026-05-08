// Package redact provides utilities for scrubbing sensitive values
// (e.g. environment variable secrets) from drift results before they
// are written to reports, webhooks, or history.
package redact

import "strings"

const placeholder = "***REDACTED***"

// DefaultSensitiveKeys is the built-in list of environment-variable name
// substrings that are considered sensitive.
var DefaultSensitiveKeys = []string{
	"PASSWORD",
	"PASSWD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
}

// Redactor scrubs sensitive values from key=value environment strings.
type Redactor struct {
	keys []string
}

// New returns a Redactor that treats any env-var whose name contains one
// of the provided substrings (case-insensitive) as sensitive.
// If keys is empty, DefaultSensitiveKeys is used.
func New(keys []string) *Redactor {
	if len(keys) == 0 {
		keys = DefaultSensitiveKeys
	}
	upper := make([]string, len(keys))
	for i, k := range keys {
		upper[i] = strings.ToUpper(k)
	}
	return &Redactor{keys: upper}
}

// Env redacts a slice of "KEY=value" strings, replacing the value portion
// of any sensitive variable with the placeholder.
func (r *Redactor) Env(pairs []string) []string {
	out := make([]string, len(pairs))
	for i, pair := range pairs {
		out[i] = r.redactPair(pair)
	}
	return out
}

// IsSensitive reports whether the given environment variable name should be
// considered sensitive.
func (r *Redactor) IsSensitive(name string) bool {
	upper := strings.ToUpper(name)
	for _, k := range r.keys {
		if strings.Contains(upper, k) {
			return true
		}
	}
	return false
}

func (r *Redactor) redactPair(pair string) string {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return pair
	}
	name := pair[:idx]
	if r.IsSensitive(name) {
		return name + "=" + placeholder
	}
	return pair
}
