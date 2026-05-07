package output_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/output"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		wantFmt  output.Format
	}{
		{"text", output.FormatText},
		{"TEXT", output.FormatText},
		{"", output.FormatText},
		{"json", output.FormatJSON},
		{"JSON", output.FormatJSON},
		{"table", output.FormatTable},
		{"TABLE", output.FormatTable},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := output.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.wantFmt {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.wantFmt)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestFormats_ReturnsAll(t *testing.T) {
	fmts := output.Formats()
	if len(fmts) != 3 {
		t.Errorf("expected 3 formats, got %d", len(fmts))
	}
}
