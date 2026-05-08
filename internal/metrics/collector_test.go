package metrics

import (
	"errors"
	"testing"

	"github.com/you/driftcheck/internal/drift"
)

func makeResult(drifted bool, err error) drift.Result {
	return drift.Result{
		Service: "svc",
		Drifted: drifted,
		Error:   err,
	}
}

func TestFromResults_AllClean(t *testing.T) {
	Reset()
	c := Global()
	results := []drift.Result{
		makeResult(false, nil),
		makeResult(false, nil),
	}
	d, cl, e := FromResults(c, results)
	if d != 0 || cl != 2 || e != 0 {
		t.Errorf("got drifted=%d clean=%d errors=%d, want 0 2 0", d, cl, e)
	}
	if c.CleanTotal.Load() != 2 {
		t.Errorf("CleanTotal = %d, want 2", c.CleanTotal.Load())
	}
}

func TestFromResults_MixedResults(t *testing.T) {
	Reset()
	c := Global()
	results := []drift.Result{
		makeResult(true, nil),
		makeResult(false, nil),
		makeResult(false, errors.New("boom")),
	}
	d, cl, e := FromResults(c, results)
	if d != 1 || cl != 1 || e != 1 {
		t.Errorf("got drifted=%d clean=%d errors=%d, want 1 1 1", d, cl, e)
	}
}

func TestFromResults_EmptySlice(t *testing.T) {
	Reset()
	c := Global()
	d, cl, e := FromResults(c, nil)
	if d != 0 || cl != 0 || e != 0 {
		t.Error("expected all zeros for empty results")
	}
	if c.ScansTotal.Load() != 1 {
		t.Error("RecordScan should still be called once")
	}
}

func TestFromResults_ErrorTakesPrecedenceOverDrift(t *testing.T) {
	// A result with both Drifted=true and an Error should be counted as an error.
	Reset()
	c := Global()
	r := drift.Result{Service: "svc", Drifted: true, Error: errors.New("oops")}
	d, _, e := FromResults(c, []drift.Result{r})
	if d != 0 || e != 1 {
		t.Errorf("error result should not be counted as drift: drifted=%d errors=%d", d, e)
	}
}
