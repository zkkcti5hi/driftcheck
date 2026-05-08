// Package audit provides structured event logging for drift scan activity,
// recording who triggered a scan, when, and what the outcome was.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/driftcheck/internal/drift"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Trigger    string    `json:"trigger"`              // "manual", "scheduled", "ci"
	Manifest   string    `json:"manifest"`             // path to manifest file
	Total      int       `json:"total"`
	Drifted    int       `json:"drifted"`
	Clean      int       `json:"clean"`
	Errored    int       `json:"errored"`
	DriftedSvcs []string `json:"drifted_services,omitempty"`
}

// Logger writes audit events to a newline-delimited JSON file.
type Logger struct {
	path string
}

// NewLogger returns a Logger that appends events to the file at path.
// The directory is created if it does not exist.
func NewLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create directory: %w", err)
	}
	return &Logger{path: path}, nil
}

// Record builds an Event from the provided results and appends it to the log.
func (l *Logger) Record(trigger, manifest string, results []drift.Result) error {
	ev := buildEvent(trigger, manifest, results)
	return l.append(ev)
}

func (l *Logger) append(ev Event) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(ev); err != nil {
		return fmt.Errorf("audit: encode event: %w", err)
	}
	return nil
}

// ReadAll returns all events stored in the audit log.
func (l *Logger) ReadAll() ([]Event, error) {
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var events []Event
	dec := json.NewDecoder(
		// wrap raw bytes in a reader via strings package would add an import;
		// write to a temp buffer instead.
		newBytesReader(data),
	)
	for dec.More() {
		var ev Event
		if err := dec.Decode(&ev); err != nil {
			return nil, fmt.Errorf("audit: decode event: %w", err)
		}
		events = append(events, ev)
	}
	return events, nil
}

func buildEvent(trigger, manifest string, results []drift.Result) Event {
	ev := Event{
		Timestamp: time.Now().UTC(),
		Trigger:   trigger,
		Manifest:  manifest,
		Total:     len(results),
	}
	for _, r := range results {
		switch {
		case r.Error != nil:
			ev.Errored++
		case r.Drifted:
			ev.Drifted++
			ev.DriftedSvcs = append(ev.DriftedSvcs, r.Service)
		default:
			ev.Clean++
		}
	}
	return ev
}
