package explain_test

import (
	"testing"

	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/explain"
)

func driftedResult(service string, fields ...drift.FieldDiff) drift.Result {
	return drift.Result{
		Service: service,
		Drifted: true,
		Fields:  fields,
	}
}

func cleanResult(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func field(name, want, got string) drift.FieldDiff {
	return drift.FieldDiff{Name: name, Want: want, Got: got}
}

func TestExplain_SkipsCleanResults(t *testing.T) {
	results := []drift.Result{cleanResult("web"), cleanResult("db")}
	got := explain.Explain(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 explanations for clean results, got %d", len(got))
	}
}

func TestExplain_ImageField(t *testing.T) {
	results := []drift.Result{
		driftedResult("api", field("image", "nginx:1.24", "nginx:1.21")),
	}
	got := explain.Explain(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(got))
	}
	e := got[0]
	if e.Service != "api" {
		t.Errorf("service: want %q got %q", "api", e.Service)
	}
	if e.Field != "image" {
		t.Errorf("field: want %q got %q", "image", e.Field)
	}
	if e.Reason == "" {
		t.Error("expected non-empty Reason")
	}
	if e.Remediation == "" {
		t.Error("expected non-empty Remediation")
	}
}

func TestExplain_EnvField(t *testing.T) {
	results := []drift.Result{
		driftedResult("worker", field("env", "LOG_LEVEL=info", "LOG_LEVEL=debug")),
	}
	got := explain.Explain(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(got))
	}
	if got[0].Field != "env" {
		t.Errorf("expected field env, got %q", got[0].Field)
	}
}

func TestExplain_ReplicasField(t *testing.T) {
	results := []drift.Result{
		driftedResult("svc", field("replicas", "3", "1")),
	}
	got := explain.Explain(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(got))
	}
	if got[0].Remediation == "" {
		t.Error("expected non-empty Remediation for replicas field")
	}
}

func TestExplain_UnknownField_FallsBackToGeneric(t *testing.T) {
	results := []drift.Result{
		driftedResult("cache", field("memory_limit", "512m", "256m")),
	}
	got := explain.Explain(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(got))
	}
	if got[0].Reason == "" {
		t.Error("expected non-empty Reason for unknown field")
	}
}

func TestExplain_MultipleFields_MultipleExplanations(t *testing.T) {
	results := []drift.Result{
		driftedResult("web",
			field("image", "nginx:1.24", "nginx:1.21"),
			field("env", "PORT=80", "PORT=8080"),
		),
	}
	got := explain.Explain(results)
	if len(got) != 2 {
		t.Fatalf("expected 2 explanations, got %d", len(got))
	}
}
