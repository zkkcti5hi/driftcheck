// Package notify implements outbound alerting for driftcheck.
//
// When drift is detected between running containers and their source
// manifests, the notifier can push a structured JSON payload to any
// HTTP webhook endpoint (e.g. Slack incoming webhooks, PagerDuty
// event API, or a custom receiver).
//
// Only results where Drifted == true are included in the payload,
// so callers may safely pass the full result slice without
// pre-filtering.
//
// Usage:
//
//	wh := notify.NewWebhook("https://hooks.example.com/drift", 0)
//	if err := wh.Notify(results); err != nil {
//		log.Printf("webhook notification failed: %v", err)
//	}
package notify
