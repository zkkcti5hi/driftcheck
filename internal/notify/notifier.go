// Package notify provides alerting integrations for drift detection results.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/driftcheck/internal/drift"
)

// Webhook represents a generic webhook notification target.
type Webhook struct {
	URL     string
	Timeout time.Duration
	client  *http.Client
}

// NewWebhook creates a Webhook notifier with the given URL.
// If timeout is zero, a 10-second default is used.
func NewWebhook(url string, timeout time.Duration) *Webhook {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Webhook{
		URL:     url,
		Timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// payload is the JSON body sent to the webhook endpoint.
type payload struct {
	Timestamp  string         `json:"timestamp"`
	DriftCount int            `json:"drift_count"`
	Results    []drift.Result `json:"results"`
}

// Notify sends drifted results to the configured webhook URL.
// It returns an error if the HTTP request fails or the server
// responds with a non-2xx status code.
func (w *Webhook) Notify(results []drift.Result) error {
	drifted := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if r.Drifted {
			drifted = append(drifted, r)
		}
	}
	if len(drifted) == 0 {
		return nil
	}

	p := payload{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		DriftCount: len(drifted),
		Results:    drifted,
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post to %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
