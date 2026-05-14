package remediate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/driftcheck/internal/remediate"
)

func sampleReport() remediate.Report {
	return remediate.Report{
		Suggestions: []remediate.Suggestion{
			{
				Service:  "api",
				Field:    "image",
				Current:  "nginx:1.24",
				Expected: "nginx:1.25",
				Hint:     "Re-deploy api with nginx:1.25",
			},
		},
	}
}

func TestWrite_TextFormat_NoSuggestions(t *testing.T) {
	var buf bytes.Buffer
	if err := remediate.Write(&buf, remediate.Report{}, remediate.FormatText); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "in sync") {
		t.Errorf("expected 'in sync' message, got: %s", buf.String())
	}
}

func TestWrite_TextFormat_WithSuggestions(t *testing.T) {
	var buf bytes.Buffer
	if err := remediate.Write(&buf, sampleReport(), remediate.FormatText); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[api]") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	if !strings.Contains(out, "Re-deploy") {
		t.Errorf("expected hint in output, got: %s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := remediate.Write(&buf, sampleReport(), remediate.FormatJSON); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"Service"`) && !strings.Contains(buf.String(), `"service"`) {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestWrite_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := remediate.Write(&buf, sampleReport(), remediate.FormatTable); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "SERVICE") {
		t.Errorf("expected table header, got: %s", out)
	}
	if !strings.Contains(out, "nginx:1.25") {
		t.Errorf("expected expected image in table, got: %s", out)
	}
}

func TestWrite_DefaultFormat_FallsBackToText(t *testing.T) {
	var buf bytes.Buffer
	if err := remediate.Write(&buf, remediate.Report{}, "unknown"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "in sync") {
		t.Errorf("default should fall back to text, got: %s", buf.String())
	}
}
