// Package ratelimit white-box tests using injectable clock.
package ratelimit

import (
	"testing"
	"time"
)

func newWithClock(interval time.Duration, clock func() time.Time) *Limiter {
	return &Limiter{
		interval: interval,
		now:      clock,
	}
}

func TestAllow_AllowedAfterIntervalElapsed(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	clock := func() time.Time { return current }

	l := newWithClock(10*time.Second, clock)

	if err := l.Allow(); err != nil {
		t.Fatalf("first allow: %v", err)
	}

	// Advance by less than interval — should be blocked.
	current = base.Add(5 * time.Second)
	if err := l.Allow(); err == nil {
		t.Fatal("expected rate limit within interval")
	}

	// Advance past interval — should succeed.
	current = base.Add(10 * time.Second)
	if err := l.Allow(); err != nil {
		t.Fatalf("expected allow after interval: %v", err)
	}
}

func TestRemaining_AccurateWithClock(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	current := base
	clock := func() time.Time { return current }

	l := newWithClock(30*time.Second, clock)
	_ = l.Allow()

	current = base.Add(10 * time.Second)
	got := l.Remaining()
	want := 20 * time.Second
	if got != want {
		t.Fatalf("Remaining() = %v, want %v", got, want)
	}
}

func TestRemaining_ZeroWhenExactlyElapsed(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	current := base
	clock := func() time.Time { return current }

	l := newWithClock(15*time.Second, clock)
	_ = l.Allow()

	current = base.Add(15 * time.Second)
	if r := l.Remaining(); r != 0 {
		t.Fatalf("expected 0 when exactly elapsed, got %v", r)
	}
}
