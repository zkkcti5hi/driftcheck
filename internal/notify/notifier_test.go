package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/notify"
)

func driftedResult(service string) drift.Result {
	return drift.Result{
		Service:     service,
		Drifted:     true,
		DriftReason: "image mismatch",
	}
}

func cleanResult(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestNotify_SendsDriftedResults(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wh := notify.NewWebhook(ts.URL, 5*time.Second)
	err := wh.Notify([]drift.Result{driftedResult("api"), cleanResult("db")})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received["drift_count"].(float64) != 1 {
		t.Errorf("expected drift_count=1, got %v", received["drift_count"])
	}
}

func TestNotify_SkipsWhenNoDrift(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wh := notify.NewWebhook(ts.URL, 5*time.Second)
	err := wh.Notify([]drift.Result{cleanResult("api"), cleanResult("db")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected webhook not to be called when there is no drift")
	}
}

func TestNotify_ReturnsErrorOnNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wh := notify.NewWebhook(ts.URL, 5*time.Second)
	err := wh.Notify([]drift.Result{driftedResult("svc")})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestNewWebhook_DefaultTimeout(t *testing.T) {
	wh := notify.NewWebhook("http://example.com", 0)
	if wh.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", wh.Timeout)
	}
}
