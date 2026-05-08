package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysSucceeds(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if err := l.Allow(); err != nil {
		t.Fatalf("expected first Allow to succeed, got %v", err)
	}
}

func TestAllow_SecondCallWithinIntervalBlocked(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(10 * time.Second)
	l.(*struct{ ratelimit.Limiter }) // white-box via exported method only

	// Use the package-level constructor and advance time manually via Reset trick.
	l2 := ratelimit.New(10 * time.Second)
	if err := l2.Allow(); err != nil {
		t.Fatalf("unexpected error on first allow: %v", err)
	}
	// Second call immediately — should be rate limited.
	if err := l2.Allow(); err == nil {
		t.Fatal("expected ErrRateLimited, got nil")
	}
	_ = now
}

func TestAllow_ZeroIntervalNeverLimits(t *testing.T) {
	l := ratelimit.New(0)
	for i := 0; i < 10; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("zero-interval limiter should never block, got %v at iteration %d", err, i)
		}
	}
}

func TestAllow_NegativeIntervalNeverLimits(t *testing.T) {
	l := ratelimit.New(-1 * time.Second)
	if err := l.Allow(); err != nil {
		t.Fatalf("negative interval should never block: %v", err)
	}
	if err := l.Allow(); err != nil {
		t.Fatalf("negative interval should never block on second call: %v", err)
	}
}

func TestReset_AllowsImmediateRescan(t *testing.T) {
	l := ratelimit.New(1 * time.Hour)
	if err := l.Allow(); err != nil {
		t.Fatalf("first allow: %v", err)
	}
	l.Reset()
	if err := l.Allow(); err != nil {
		t.Fatalf("allow after Reset should succeed: %v", err)
	}
}

func TestRemaining_ZeroBeforeFirstScan(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if r := l.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining before first scan, got %v", r)
	}
}

func TestRemaining_PositiveAfterScan(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	_ = l.Allow()
	if r := l.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining after scan, got %v", r)
	}
}

func TestRemaining_ZeroForDisabledLimiter(t *testing.T) {
	l := ratelimit.New(0)
	_ = l.Allow()
	if r := l.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining for disabled limiter, got %v", r)
	}
}

func TestErrRateLimited_IsExported(t *testing.T) {
	if ratelimit.ErrRateLimited == nil {
		t.Fatal("ErrRateLimited must not be nil")
	}
}
