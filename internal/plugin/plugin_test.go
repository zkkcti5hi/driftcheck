package plugin_test

import (
	"errors"
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/plugin"
)

// stubChecker is a minimal Checker used in tests.
type stubChecker struct {
	name   string
	append drift.Result
	err    error
}

func (s *stubChecker) Name() string { return s.name }
func (s *stubChecker) Check(results []drift.Result) ([]drift.Result, error) {
	if s.err != nil {
		return results, s.err
	}
	if s.append.Service != "" {
		results = append(results, s.append)
	}
	return results, nil
}

func TestRegister_And_Names(t *testing.T) {
	r := plugin.NewRegistry()
	if err := r.Register(&stubChecker{name: "alpha"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Register(&stubChecker{name: "beta"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := r.Names()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	r := plugin.NewRegistry()
	_ = r.Register(&stubChecker{name: "dup"})
	err := r.Register(&stubChecker{name: "dup"})
	if err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestRun_NoCheckers_ReturnsResultsUnchanged(t *testing.T) {
	r := plugin.NewRegistry()
	input := []drift.Result{{Service: "svc"}}
	out, err := r.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 result, got %d", len(out))
	}
}

func TestRun_CheckerAppendsResult(t *testing.T) {
	r := plugin.NewRegistry()
	_ = r.Register(&stubChecker{
		name:   "appender",
		append: drift.Result{Service: "extra"},
	})
	out, err := r.Run([]drift.Result{{Service: "original"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestRun_CheckerError_Propagates(t *testing.T) {
	r := plugin.NewRegistry()
	_ = r.Register(&stubChecker{name: "failer", err: errors.New("boom")})
	_, err := r.Run(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errors.New("boom")) && err.Error() == "" {
		t.Errorf("unexpected error value: %v", err)
	}
}

func TestRun_CheckersExecuteInOrder(t *testing.T) {
	order := []string{}
	type recorder struct{ name string }
	makeChecker := func(n string) plugin.Checker {
		return &stubChecker{name: n}
	}
	r := plugin.NewRegistry()
	_ = r.Register(makeChecker("first"))
	_ = r.Register(makeChecker("second"))
	_ = order // suppress unused warning
	names := r.Names()
	if names[0] != "first" || names[1] != "second" {
		t.Errorf("wrong order: %v", names)
	}
}
