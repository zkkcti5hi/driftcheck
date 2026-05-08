package builtin_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/plugin/builtin"
)

func resultWithLabels(id string, labels map[string]string) drift.Result {
	return drift.Result{
		Service:     "svc",
		ContainerID: id,
		Labels:      labels,
	}
}

func TestLabelChecker_Name(t *testing.T) {
	c := builtin.NewLabelChecker("app")
	if c.Name() != "builtin/label-checker" {
		t.Errorf("unexpected name: %s", c.Name())
	}
}

func TestLabelChecker_EmptyLabel_NoOp(t *testing.T) {
	c := builtin.NewLabelChecker("")
	input := []drift.Result{resultWithLabels("abc", nil)}
	out, err := c.Check(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Drifted {
		t.Error("expected no drift when RequiredLabel is empty")
	}
}

func TestLabelChecker_MissingLabel_MarksDrifted(t *testing.T) {
	c := builtin.NewLabelChecker("com.example/owner")
	input := []drift.Result{resultWithLabels("abc123", map[string]string{"other": "val"})}
	out, err := c.Check(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out[0].Drifted {
		t.Error("expected result to be marked drifted")
	}
	if len(out[0].Reasons) == 0 {
		t.Error("expected at least one reason")
	}
}

func TestLabelChecker_LabelPresent_NoDrift(t *testing.T) {
	c := builtin.NewLabelChecker("com.example/owner")
	input := []drift.Result{
		resultWithLabels("abc123", map[string]string{"com.example/owner": "team-a"}),
	}
	out, err := c.Check(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Drifted {
		t.Error("expected no drift when label is present")
	}
}

func TestLabelChecker_NoContainerID_Skipped(t *testing.T) {
	c := builtin.NewLabelChecker("required")
	input := []drift.Result{{Service: "orphan", ContainerID: ""}}
	out, err := c.Check(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Drifted {
		t.Error("container without ID should be skipped")
	}
}

func TestLabelChecker_MultipleResults(t *testing.T) {
	c := builtin.NewLabelChecker("env")
	input := []drift.Result{
		resultWithLabels("id1", map[string]string{"env": "prod"}),
		resultWithLabels("id2", map[string]string{}),
	}
	out, err := c.Check(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Drifted {
		t.Error("first result should not be drifted")
	}
	if !out[1].Drifted {
		t.Error("second result should be drifted")
	}
}
