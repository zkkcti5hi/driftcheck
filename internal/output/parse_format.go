package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string flag value into a Format constant.
// It returns an error if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "table":
		return FormatTable, nil
	default:
		return FormatText, fmt.Errorf(
			"unknown output format %q: must be one of text, json, table", s,
		)
	}
}

// Formats returns all supported format names as a slice.
func Formats() []string {
	return []string{"text", "json", "table"}
}
