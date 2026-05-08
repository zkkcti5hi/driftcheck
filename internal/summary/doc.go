// Package summary aggregates drift scan results into a high-level report.
//
// Use Compute to derive statistics such as total services checked, how many
// are drifted or errored, the overall drift-rate percentage, and a ranked list
// of the most-frequently drifted services.
//
// Example:
//
//	report := summary.Compute(results)
//	fmt.Printf("Drift rate: %.1f%%\n", report.Stats.DriftRate)
//	for _, td := range report.TopDrifted {
//		fmt.Printf("  %s: %d occurrences\n", td.Service, td.Count)
//	}
package summary
