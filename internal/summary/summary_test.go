package summary_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/summary"
)

func makeResult(service string, drifted bool, errMsg string) drift.Result {
	return drift.Result{
		Service: service,
		Drifted: drifted,
		Error:   errMsg,
	}
}

func TestCompute_Empty(t *testing.T) {
	r := summary.Compute(nil)
	if r.Stats.Total != 0 {
		t.Fatalf("expected 0 total, got %d", r.Stats.Total)
	}
	if r.Stats.DriftRate != 0 {
		t.Fatalf("expected 0 drift rate, got %f", r.Stats.DriftRate)
	}
	if len(r.TopDrifted) != 0 {
		t.Fatalf("expected empty top-drifted list")
	}
}

func TestCompute_AllClean(t *testing.T) {
	results := []drift.Result{
		makeResult("web", false, ""),
		makeResult("db", false, ""),
	}
	r := summary.Compute(results)
	if r.Stats.Total != 2 || r.Stats.Clean != 2 || r.Stats.Drifted != 0 {
		t.Fatalf("unexpected stats: %+v", r.Stats)
	}
	if r.Stats.DriftRate != 0 {
		t.Fatalf("expected 0 drift rate")
	}
}

func TestCompute_MixedResults(t *testing.T) {
	results := []drift.Result{
		makeResult("web", true, ""),
		makeResult("web", true, ""),
		makeResult("db", false, ""),
		makeResult("cache", false, ""),
		makeResult("worker", false, "connection refused"),
	}
	r := summary.Compute(results)
	if r.Stats.Total != 5 {
		t.Fatalf("expected 5 total, got %d", r.Stats.Total)
	}
	if r.Stats.Drifted != 2 {
		t.Fatalf("expected 2 drifted, got %d", r.Stats.Drifted)
	}
	if r.Stats.Clean != 2 {
		t.Fatalf("expected 2 clean, got %d", r.Stats.Clean)
	}
	if r.Stats.Errored != 1 {
		t.Fatalf("expected 1 errored, got %d", r.Stats.Errored)
	}
	if r.Stats.DriftRate != 40 {
		t.Fatalf("expected 40%% drift rate, got %f", r.Stats.DriftRate)
	}
}

func TestCompute_TopDrifted_Ordered(t *testing.T) {
	results := []drift.Result{
		makeResult("api", true, ""),
		makeResult("web", true, ""),
		makeResult("web", true, ""),
		makeResult("web", true, ""),
		makeResult("api", true, ""),
	}
	r := summary.Compute(results)
	if len(r.TopDrifted) == 0 {
		t.Fatal("expected non-empty top-drifted")
	}
	if r.TopDrifted[0].Service != "web" || r.TopDrifted[0].Count != 3 {
		t.Fatalf("expected web first with 3, got %+v", r.TopDrifted[0])
	}
	if r.TopDrifted[1].Service != "api" || r.TopDrifted[1].Count != 2 {
		t.Fatalf("expected api second with 2, got %+v", r.TopDrifted[1])
	}
}

func TestCompute_ComputedAt_NonZero(t *testing.T) {
	r := summary.Compute([]drift.Result{makeResult("svc", false, "")})
	if r.Stats.ComputedAt.IsZero() {
		t.Fatal("expected non-zero ComputedAt timestamp")
	}
}
