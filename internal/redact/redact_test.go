package redact_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/redact"
)

func TestNew_DefaultKeys(t *testing.T) {
	r := redact.New(nil)
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestIsSensitive_MatchesSubstring(t *testing.T) {
	r := redact.New(nil)
	cases := []struct {
		name      string
		want      bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"SECRET_TOKEN", true},
		{"AUTH_HEADER", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"APP_NAME", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := r.IsSensitive(tc.name)
			if got != tc.want {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestEnv_RedactsSensitivePairs(t *testing.T) {
	r := redact.New(nil)
	input := []string{
		"APP_NAME=myapp",
		"DB_PASSWORD=supersecret",
		"PORT=8080",
		"API_KEY=abc123",
	}
	got := r.Env(input)

	if got[0] != "APP_NAME=myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", got[0])
	}
	if got[1] != "DB_PASSWORD=***REDACTED***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", got[1])
	}
	if got[2] != "PORT=8080" {
		t.Errorf("expected PORT unchanged, got %q", got[2])
	}
	if got[3] != "API_KEY=***REDACTED***" {
		t.Errorf("expected API_KEY redacted, got %q", got[3])
	}
}

func TestEnv_PairWithoutEquals(t *testing.T) {
	r := redact.New(nil)
	input := []string{"MALFORMED"}
	got := r.Env(input)
	if got[0] != "MALFORMED" {
		t.Errorf("expected malformed pair unchanged, got %q", got[0])
	}
}

func TestEnv_CustomKeys(t *testing.T) {
	r := redact.New([]string{"MYAPP_INTERNAL"})
	input := []string{
		"MYAPP_INTERNAL_TOKEN=xyz",
		"DB_PASSWORD=should_not_be_redacted",
	}
	got := r.Env(input)
	if got[0] != "MYAPP_INTERNAL_TOKEN=***REDACTED***" {
		t.Errorf("expected custom key redacted, got %q", got[0])
	}
	if got[1] != "DB_PASSWORD=should_not_be_redacted" {
		t.Errorf("expected DB_PASSWORD unchanged with custom keys, got %q", got[1])
	}
}

func TestEnv_EmptySlice(t *testing.T) {
	r := redact.New(nil)
	got := r.Env([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
