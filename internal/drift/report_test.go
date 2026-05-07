package drift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftcheck/internal/drift"
)

func TestSummarise(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "web", Drifted: false},
		{ServiceName: "db", Drifted: true, Issues: []string{"image mismatch"}},
		{ServiceName: "cache", Drifted: true, Issues: []string{"service not running"}},
	}
	s := drift.Summarise(results)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Clean != 1 {
		t.Errorf("expected Clean=1, got %d", s.Clean)
	}
	if s.Drifted != 2 {
		t.Errorf("expected Drifted=2, got %d", s.Drifted)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := drift.Summarise([]drift.DriftResult{})
	if s.Total != 0 || s.Clean != 0 || s.Drifted != 0 {
		t.Errorf("expected all zeros for empty input, got Total=%d Clean=%d Drifted=%d",
			s.Total, s.Clean, s.Drifted)
	}
}

func TestWriteReport_ContainsDriftedLabel(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "web", ContainerID: "abc123def456", Drifted: true,
			Issues: []string{`image mismatch: want "nginx:1.25", got "nginx:1.24"`}},
	}
	var buf bytes.Buffer
	drift.WriteReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "DRIFTED") {
		t.Error("expected DRIFTED in output")
	}
	if !strings.Contains(out, "image mismatch") {
		t.Error("expected issue text in output")
	}
	if !strings.Contains(out, "abc123def456") {
		t.Error("expected short container ID in output")
	}
}

func TestWriteReport_CleanService(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "api", ContainerID: "deadbeef0000", Drifted: false},
	}
	var buf bytes.Buffer
	drift.WriteReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "OK") {
		t.Error("expected OK in output for clean service")
	}
	if strings.Contains(out, "DRIFTED") {
		t.Error("did not expect DRIFTED in output for clean service")
	}
}

func TestWriteReport_NoContainerID(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "worker", Drifted: true, Issues: []string{"service not running"}},
	}
	var buf bytes.Buffer
	drift.WriteReport(&buf, results)
	if !strings.Contains(buf.String(), "<none>") {
		t.Error("expected <none> placeholder for missing container ID")
	}
}
