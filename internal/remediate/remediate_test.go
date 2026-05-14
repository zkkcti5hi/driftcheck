package remediate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/driftcheck/internal/drift"
	"github.com/yourusername/driftcheck/internal/remediate"
)

func driftedResult(service, field, current, expected string) drift.Result {
	return drift.Result{
		Service: service,
		Drifted: true,
		Fields: []drift.Field{
			{Name: field, Current: current, Expected: expected},
		},
	}
}

func cleanResult(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestGenerate_SkipsCleanResults(t *testing.T) {
	report := remediate.Generate([]drift.Result{cleanResult("web")})
	if report.HasSuggestions() {
		t.Fatalf("expected no suggestions for clean result, got %d", len(report.Suggestions))
	}
}

func TestGenerate_ImageField(t *testing.T) {
	report := remediate.Generate([]drift.Result{
		driftedResult("api", "image", "nginx:1.24", "nginx:1.25"),
	})
	if !report.HasSuggestions() {
		t.Fatal("expected suggestions")
	}
	s := report.Suggestions[0]
	if s.Service != "api" {
		t.Errorf("service: want api, got %s", s.Service)
	}
	if !strings.Contains(s.Hint, "force-recreate") {
		t.Errorf("image hint should mention force-recreate, got: %s", s.Hint)
	}
}

func TestGenerate_ReplicasField(t *testing.T) {
	report := remediate.Generate([]drift.Result{
		driftedResult("worker", "replicas", "1", "3"),
	})
	s := report.Suggestions[0]
	if !strings.Contains(s.Hint, "kubectl scale") {
		t.Errorf("replicas hint should mention kubectl scale, got: %s", s.Hint)
	}
}

func TestGenerate_UnknownField(t *testing.T) {
	report := remediate.Generate([]drift.Result{
		driftedResult("db", "port", "5432", "5433"),
	})
	s := report.Suggestions[0]
	if !strings.Contains(s.Hint, "port") {
		t.Errorf("generic hint should mention field name, got: %s", s.Hint)
	}
}

func TestGenerate_MultipleFields(t *testing.T) {
	res := drift.Result{
		Service: "svc",
		Drifted: true,
		Fields: []drift.Field{
			{Name: "image", Current: "a:1", Expected: "a:2"},
			{Name: "env", Current: "X=1", Expected: "X=2"},
		},
	}
	report := remediate.Generate([]drift.Result{res})
	if len(report.Suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(report.Suggestions))
	}
}
