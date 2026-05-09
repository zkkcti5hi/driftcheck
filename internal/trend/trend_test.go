package trend_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/trend"
)

func makeResult(service string, drifted bool) drift.Result {
	return drift.Result{Service: service, Drifted: drifted}
}

func ts(offset int) time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(offset) * time.Hour)
}

func TestAnalyse_Empty(t *testing.T) {
	trends := trend.Analyse(nil)
	if len(trends) != 0 {
		t.Fatalf("expected 0 trends, got %d", len(trends))
	}
}

func TestAnalyse_AllClean(t *testing.T) {
	entries := []trend.Entry{
		{Timestamp: ts(0), Results: []drift.Result{makeResult("web", false)}},
		{Timestamp: ts(1), Results: []drift.Result{makeResult("web", false)}},
	}
	trends := trend.Analyse(entries)
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(trends))
	}
	if trends[0].DriftRate != 0.0 {
		t.Errorf("expected DriftRate 0.0, got %f", trends[0].DriftRate)
	}
	if trends[0].Worsening {
		t.Error("expected Worsening=false")
	}
}

func TestAnalyse_DriftRate(t *testing.T) {
	entries := []trend.Entry{
		{Timestamp: ts(0), Results: []drift.Result{makeResult("api", true)}},
		{Timestamp: ts(1), Results: []drift.Result{makeResult("api", false)}},
		{Timestamp: ts(2), Results: []drift.Result{makeResult("api", true)}},
		{Timestamp: ts(3), Results: []drift.Result{makeResult("api", true)}},
	}
	trends := trend.Analyse(entries)
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend")
	}
	const want = 0.75
	if trends[0].DriftRate != want {
		t.Errorf("DriftRate: want %f, got %f", want, trends[0].DriftRate)
	}
}

func TestAnalyse_Worsening(t *testing.T) {
	entries := []trend.Entry{
		{Timestamp: ts(0), Results: []drift.Result{makeResult("db", false)}},
		{Timestamp: ts(1), Results: []drift.Result{makeResult("db", true)}},
	}
	trends := trend.Analyse(entries)
	if !trends[0].Worsening {
		t.Error("expected Worsening=true when last drifted and previous clean")
	}
}

func TestAnalyse_NotWorsening_WhenBothDrifted(t *testing.T) {
	entries := []trend.Entry{
		{Timestamp: ts(0), Results: []drift.Result{makeResult("cache", true)}},
		{Timestamp: ts(1), Results: []drift.Result{makeResult("cache", true)}},
	}
	trends := trend.Analyse(entries)
	if trends[0].Worsening {
		t.Error("expected Worsening=false when both scans drifted")
	}
}

func TestAnalyse_SortedByDriftRateDesc(t *testing.T) {
	entries := []trend.Entry{
		{Timestamp: ts(0), Results: []drift.Result{
			makeResult("low", false),
			makeResult("high", true),
		}},
	}
	trends := trend.Analyse(entries)
	if trends[0].Service != "high" {
		t.Errorf("expected high-drift service first, got %s", trends[0].Service)
	}
}
