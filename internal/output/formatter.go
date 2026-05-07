// Package output provides formatting utilities for drift detection results.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatTable Format = "table"
)

// Write writes the drift results to w in the specified format.
func Write(w io.Writer, results []drift.Result, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, results)
	case FormatTable:
		return writeTable(w, results)
	default:
		return writeText(w, results)
	}
}

func writeJSON(w io.Writer, results []drift.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func writeText(w io.Writer, results []drift.Result) error {
	for _, r := range results {
		status := "OK"
		if r.Drifted {
			status = "DRIFTED"
		}
		fmt.Fprintf(w, "[%s] %s\n", status, r.ServiceName)
		for _, d := range r.Differences {
			fmt.Fprintf(w, "  - %s\n", d)
		}
	}
	return nil
}

func writeTable(w io.Writer, results []drift.Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTATUS\tDIFFERENCES")
	fmt.Fprintln(tw, "-------\t------\t-----------")
	for _, r := range results {
		status := "OK"
		if r.Drifted {
			status = "DRIFTED"
		}
		diffs := "-"
		if len(r.Differences) > 0 {
			diffs = r.Differences[0]
			if len(r.Differences) > 1 {
				diffs += fmt.Sprintf(" (+%d more)", len(r.Differences)-1)
			}
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.ServiceName, status, diffs)
	}
	return tw.Flush()
}
