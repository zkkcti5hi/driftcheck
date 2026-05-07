package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/output"
)

var sampleResults = []drift.Result{
	{
		ServiceName: "web",
		Drifted:     true,
		Differences: []string{"image mismatch: want nginx:1.25, got nginx:1.24"},
	},
	{
		ServiceName: "db",
		Drifted:     false,
		Differences: nil,
	},
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, output.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFTED] web") {
		t.Errorf("expected DRIFTED label for web, got:\n%s", out)
	}
	if !strings.Contains(out, "[OK] db") {
		t.Errorf("expected OK label for db, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, output.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"ServiceName": "web"`) {
		t.Errorf("expected service name in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"Drifted": true`) {
		t.Errorf("expected drifted flag in JSON output, got:\n%s", out)
	}
}

func TestWrite_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, output.FormatTable); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SERVICE") {
		t.Errorf("expected table header, got:\n%s", out)
	}
	if !strings.Contains(out, "DRIFTED") {
		t.Errorf("expected DRIFTED in table, got:\n%s", out)
	}
}

func TestWrite_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	// Empty string should fall back to text
	if err := output.Write(&buf, sampleResults, output.Format("")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for default format")
	}
}
