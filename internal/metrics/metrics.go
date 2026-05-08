// Package metrics provides lightweight run-time counters that track
// drift-check scan activity. Counters are safe for concurrent use.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

// Counters holds cumulative statistics for a driftcheck process lifetime.
type Counters struct {
	ScansTotal    atomic.Int64
	DriftedTotal  atomic.Int64
	CleanTotal    atomic.Int64
	ErrorsTotal   atomic.Int64
	LastScanAt    atomic.Value // stores time.Time
	LastDriftedAt atomic.Value // stores time.Time
}

var (
	mu      sync.Mutex
	global  = &Counters{}
)

// Global returns the process-wide Counters singleton.
func Global() *Counters { return global }

// Reset replaces the global counters with a fresh instance (useful in tests).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	global = &Counters{}
}

// RecordScan updates counters after a single scan completes.
func (c *Counters) RecordScan(drifted, clean, errors int) {
	c.ScansTotal.Add(1)
	c.DriftedTotal.Add(int64(drifted))
	c.CleanTotal.Add(int64(clean))
	c.ErrorsTotal.Add(int64(errors))
	now := time.Now().UTC()
	c.LastScanAt.Store(now)
	if drifted > 0 {
		c.LastDriftedAt.Store(now)
	}
}

// LastScan returns the time of the most recent scan, or zero if none yet.
func (c *Counters) LastScan() time.Time {
	if v := c.LastScanAt.Load(); v != nil {
		return v.(time.Time)
	}
	return time.Time{}
}

// LastDrifted returns the time drift was last detected, or zero if never.
func (c *Counters) LastDrifted() time.Time {
	if v := c.LastDriftedAt.Load(); v != nil {
		return v.(time.Time)
	}
	return time.Time{}
}

// Write prints a human-readable summary of the counters to w.
func (c *Counters) Write(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Metric\tValue")
	fmt.Fprintf(tw, "scans_total\t%d\n", c.ScansTotal.Load())
	fmt.Fprintf(tw, "drifted_total\t%d\n", c.DriftedTotal.Load())
	fmt.Fprintf(tw, "clean_total\t%d\n", c.CleanTotal.Load())
	fmt.Fprintf(tw, "errors_total\t%d\n", c.ErrorsTotal.Load())
	if ls := c.LastScan(); !ls.IsZero() {
		fmt.Fprintf(tw, "last_scan_at\t%s\n", ls.Format(time.RFC3339))
	}
	if ld := c.LastDrifted(); !ld.IsZero() {
		fmt.Fprintf(tw, "last_drifted_at\t%s\n", ld.Format(time.RFC3339))
	}
	_ = tw.Flush()
}
