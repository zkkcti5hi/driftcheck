package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftcheck/internal/audit"
	"github.com/yourusername/driftcheck/internal/drift"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{Service: "web", Drifted: true},
		{Service: "db", Drifted: false},
		{Service: "cache", Error: os.ErrNotExist},
	}
}

func newTempLogger(t *testing.T) (*audit.Logger, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	return l, path
}

func TestRecord_WritesEvent(t *testing.T) {
	l, _ := newTempLogger(t)

	if err := l.Record("manual", "docker-compose.yml", sampleResults()); err != nil {
		t.Fatalf("Record: %v", err)
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	ev := events[0]
	if ev.Trigger != "manual" {
		t.Errorf("trigger: want manual, got %s", ev.Trigger)
	}
	if ev.Total != 3 {
		t.Errorf("total: want 3, got %d", ev.Total)
	}
	if ev.Drifted != 1 {
		t.Errorf("drifted: want 1, got %d", ev.Drifted)
	}
	if ev.Clean != 1 {
		t.Errorf("clean: want 1, got %d", ev.Clean)
	}
	if ev.Errored != 1 {
		t.Errorf("errored: want 1, got %d", ev.Errored)
	}
	if len(ev.DriftedSvcs) != 1 || ev.DriftedSvcs[0] != "web" {
		t.Errorf("drifted_services: want [web], got %v", ev.DriftedSvcs)
	}
}

func TestRecord_AppendsMultipleEvents(t *testing.T) {
	l, _ := newTempLogger(t)

	for i := 0; i < 3; i++ {
		if err := l.Record("scheduled", "compose.yml", nil); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("want 3 events, got %d", len(events))
	}
}

func TestReadAll_NoFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	l, err := audit.NewLogger(filepath.Join(dir, "audit.log"))
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("want 0 events, got %d", len(events))
	}
}

func TestRecord_TimestampIsRecent(t *testing.T) {
	l, _ := newTempLogger(t)
	before := time.Now().UTC()

	_ = l.Record("ci", "k8s.yaml", nil)

	events, _ := l.ReadAll()
	if len(events) == 0 {
		t.Fatal("no events recorded")
	}
	if events[0].Timestamp.Before(before) {
		t.Errorf("timestamp %v is before scan start %v", events[0].Timestamp, before)
	}
}
