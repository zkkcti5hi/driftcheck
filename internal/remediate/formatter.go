package remediate

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format controls the output style for remediation reports.
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

// Write renders the remediation Report to w in the requested format.
func Write(w io.Writer, r Report, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, r)
	case FormatTable:
		return writeTable(w, r)
	default:
		return writeText(w, r)
	}
}

func writeText(w io.Writer, r Report) error {
	if !r.HasSuggestions() {
		_, err := fmt.Fprintln(w, "No remediation needed — all services are in sync.")
		return err
	}
	for _, s := range r.Suggestions {
		if _, err := fmt.Fprintf(w, "[%s] %s\n  → %s\n", s.Service, s.Field, s.Hint); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.Suggestions)
}

func writeTable(w io.Writer, r Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tCURRENT\tEXPECTED")
	for _, s := range r.Suggestions {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Service, s.Field, s.Current, s.Expected)
	}
	return tw.Flush()
}
