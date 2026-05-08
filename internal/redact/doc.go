// Package redact provides utilities for scrubbing sensitive values
// (such as passwords, tokens, and API keys) from environment variable
// slices before they are written to reports, logs, or webhook payloads.
//
// Usage:
//
//	r := redact.New(nil)           // uses default sensitive key substrings
//	clean := r.Env(container.Env)  // returns redacted copy
package redact
