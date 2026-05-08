package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func setup(t *testing.T) *Counters {
	t.Helper()
	Reset()
	return Global()
}

func TestRecordScan_IncrementsCounters(t *testing.T) {
	c := setup(t)
	c.RecordScan(3, 5, 1)

	if got := c.ScansTotal.Load(); got != 1 {
		t.Errorf("ScansTotal = %d, want 1", got)
	}
	if got := c.DriftedTotal.Load(); got != 3 {
		t.Errorf("DriftedTotal = %d, want 3", got)
	}
	if got := c.CleanTotal.Load(); got != 5 {
		t.Errorf("CleanTotal = %d, want 5", got)
	}
	if got := c.ErrorsTotal.Load(); got != 1 {
		t.Errorf("ErrorsTotal = %d, want 1", got)
	}
}

func TestRecordScan_Accumulates(t *testing.T) {
	c := setup(t)
	c.RecordScan(1, 2, 0)
	c.RecordScan(2, 3, 1)

	if got := c.ScansTotal.Load(); got != 2 {
		t.Errorf("ScansTotal = %d, want 2", got)
	}
	if got := c.DriftedTotal.Load(); got != 3 {
		t.Errorf("DriftedTotal = %d, want 3", got)
	}
}

func TestLastScan_ZeroBeforeFirstScan(t *testing.T) {
	c := setup(t)
	if !c.LastScan().IsZero() {
		t.Error("expected zero time before any scan")
	}
}

func TestLastScan_SetAfterScan(t *testing.T) {
	c := setup(t)
	before := time.Now().UTC().Add(-time.Second)
	c.RecordScan(0, 1, 0)
	if c.LastScan().Before(before) {
		t.Error("LastScan should be recent")
	}
}

func TestLastDrifted_ZeroWhenNoDrift(t *testing.T) {
	c := setup(t)
	c.RecordScan(0, 5, 0)
	if !c.LastDrifted().IsZero() {
		t.Error("LastDrifted should be zero when no drift recorded")
	}
}

func TestLastDrifted_SetWhenDrift(t *testing.T) {
	c := setup(t)
	c.RecordScan(1, 0, 0)
	if c.LastDrifted().IsZero() {
		t.Error("LastDrifted should be set after drift recorded")
	}
}

func TestWrite_ContainsExpectedFields(t *testing.T) {
	c := setup(t)
	c.RecordScan(2, 4, 1)

	var buf bytes.Buffer
	c.Write(&buf)
	out := buf.String()

	for _, want := range []string{"scans_total", "drifted_total", "clean_total", "errors_total", "last_scan_at"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing field %q\ngot:\n%s", want, out)
		}
	}
}

func TestGlobal_ReturnsSameInstance(t *testing.T) {
	Reset()
	if Global() != Global() {
		t.Error("Global() should return the same pointer")
	}
}
