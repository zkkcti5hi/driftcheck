package filter_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/filter"
)

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	candidates := []string{"web", "db", "cache"}
	got := filter.Filter(candidates, filter.Options{})
	if len(got) != len(candidates) {
		t.Fatalf("expected %d services, got %d", len(candidates), len(got))
	}
}

func TestFilter_ServiceAllowList(t *testing.T) {
	candidates := []string{"web", "db", "cache"}
	opts := filter.Options{Services: []string{"web", "cache"}}
	got := filter.Filter(candidates, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(got), got)
	}
	for _, s := range got {
		if s == "db" {
			t.Error("db should have been filtered out")
		}
	}
}

func TestFilter_ServiceNotInCandidates(t *testing.T) {
	candidates := []string{"web", "db"}
	opts := filter.Options{Services: []string{"worker"}}
	got := filter.Filter(candidates, opts)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestFilter_EmptyCandidates(t *testing.T) {
	got := filter.Filter(nil, filter.Options{Services: []string{"web"}})
	if got != nil && len(got) != 0 {
		t.Fatalf("expected nil/empty, got %v", got)
	}
}

// TestFilter_PreservesOrder ensures that Filter returns results in the same
// order as they appear in the candidates slice, not the allow-list order.
func TestFilter_PreservesOrder(t *testing.T) {
	candidates := []string{"alpha", "beta", "gamma", "delta"}
	opts := filter.Options{Services: []string{"delta", "alpha"}}
	got := filter.Filter(candidates, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d: %v", len(got), got)
	}
	if got[0] != "alpha" || got[1] != "delta" {
		t.Errorf("expected [alpha delta] in candidate order, got %v", got)
	}
}

// TestFilter_DuplicatesInAllowList ensures that duplicate entries in the
// allow-list do not cause a candidate to appear more than once in the result.
func TestFilter_DuplicatesInAllowList(t *testing.T) {
	candidates := []string{"web", "db", "cache"}
	opts := filter.Options{Services: []string{"web", "web", "db"}}
	got := filter.Filter(candidates, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 results (no duplicates), got %d: %v", len(got), got)
	}
}

func TestMatchLabels_AllMatch(t *testing.T) {
	container := map[string]string{"env": "prod", "team": "platform"}
	required := map[string]string{"env": "prod"}
	if !filter.MatchLabels(container, required) {
		t.Error("expected labels to match")
	}
}

func TestMatchLabels_MissingKey(t *testing.T) {
	container := map[string]string{"env": "prod"}
	required := map[string]string{"env": "prod", "team": "platform"}
	if filter.MatchLabels(container, required) {
		t.Error("expected labels NOT to match")
	}
}

func TestMatchLabels_WrongValue(t *testing.T) {
	container := map[string]string{"env": "staging"}
	required := map[string]string{"env": "prod"}
	if filter.MatchLabels(container, required) {
		t.Error("expected labels NOT to match due to wrong value")
	}
}

func TestMatchLabels_EmptyRequired(t *testing.T) {
	container := map[string]string{"env": "prod"}
	if !filter.MatchLabels(container, map[string]string{}) {
		t.Error("empty required should always match")
	}
}
