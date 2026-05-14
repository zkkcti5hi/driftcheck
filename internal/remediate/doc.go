// Package remediate analyses drift results and produces actionable remediation
// suggestions to help operators reconcile running containers with their
// declared manifests.
//
// Usage:
//
//	results := detector.Detect(...)
//	report  := remediate.Generate(results)
//	remediate.Write(os.Stdout, report, remediate.FormatText)
//
// Three output formats are supported:
//
//	- text  – plain human-readable hints (default)
//	- json  – machine-readable JSON array
//	- table – aligned tabular summary
package remediate
