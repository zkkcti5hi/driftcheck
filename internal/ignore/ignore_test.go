package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftcheck/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "drift-ignore.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

const sampleIgnore = `
ignore:
  - service: web
    fields:
      - image
      - env
  - service: worker
    fields:
      - ports
`

func TestLoad_EmptyPath(t *testing.T) {
	s, err := ignore.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 0 {
		t.Errorf("expected no rules, got %d", len(s.Rules))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	s, err := ignore.Load("/nonexistent/path/drift-ignore.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.Rules) != 0 {
		t.Errorf("expected no rules, got %d", len(s.Rules))
	}
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeTempIgnore(t, sampleIgnore)
	s, err := ignore.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(s.Rules))
	}
	if s.Rules[0].Service != "web" {
		t.Errorf("expected service 'web', got %q", s.Rules[0].Service)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTempIgnore(t, ": bad: yaml: [")
	_, err := ignore.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSuppressed_MatchesFieldAndService(t *testing.T) {
	p := writeTempIgnore(t, sampleIgnore)
	s, _ := ignore.Load(p)

	if !s.Suppressed("web", "image") {
		t.Error("expected web/image to be suppressed")
	}
	if !s.Suppressed("web", "env") {
		t.Error("expected web/env to be suppressed")
	}
	if !s.Suppressed("worker", "ports") {
		t.Error("expected worker/ports to be suppressed")
	}
}

func TestSuppressed_NoMatchOnWrongService(t *testing.T) {
	p := writeTempIgnore(t, sampleIgnore)
	s, _ := ignore.Load(p)

	if s.Suppressed("db", "image") {
		t.Error("db/image should not be suppressed")
	}
	if s.Suppressed("worker", "image") {
		t.Error("worker/image should not be suppressed")
	}
}

func TestSuppressed_EmptySet(t *testing.T) {
	s, _ := ignore.Load("")
	if s.Suppressed("web", "image") {
		t.Error("empty set should suppress nothing")
	}
}
