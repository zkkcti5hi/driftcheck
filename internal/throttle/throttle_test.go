package throttle

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestAllow_UnderLimit(t *testing.T) {
	th := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := th.Allow(); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	th := New(2, time.Minute)
	_ = th.Allow()
	_ = th.Allow()
	if err := th.Allow(); err == nil {
		t.Fatal("expected throttle error on third call, got nil")
	}
}

func TestAllow_ErrorMentionsRetry(t *testing.T) {
	th := New(1, time.Minute)
	_ = th.Allow()
	err := th.Allow()
	if err == nil || !strings.Contains(err.Error(), "retry after") {
		t.Fatalf("expected retry message, got: %v", err)
	}
}

func TestAllow_ZeroMaxDisablesThrottle(t *testing.T) {
	th := New(0, time.Minute)
	for i := 0; i < 100; i++ {
		if err := th.Allow(); err != nil {
			t.Fatalf("unexpected error with zero max: %v", err)
		}
	}
}

func TestAllow_WindowReset(t *testing.T) {
	now := time.Now()
	clock := &now
	th := newWithClock(1, 50*time.Millisecond, func() time.Time { return *clock })

	_ = th.Allow()
	if err := th.Allow(); err == nil {
		t.Fatal("expected throttle on second call")
	}

	advanced := now.Add(60 * time.Millisecond)
	clock = &advanced
	if err := th.Allow(); err != nil {
		t.Fatalf("expected allow after window reset: %v", err)
	}
}

func TestRemaining_DecreasesWithCalls(t *testing.T) {
	th := New(5, time.Minute)
	if r := th.Remaining(); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	_ = th.Allow()
	_ = th.Allow()
	if r := th.Remaining(); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestRemaining_DisabledThrottle(t *testing.T) {
	th := New(0, time.Minute)
	if r := th.Remaining(); r != -1 {
		t.Fatalf("expected -1 for disabled throttle, got %d", r)
	}
}

func TestReset_RestoresTokens(t *testing.T) {
	th := New(2, time.Minute)
	_ = th.Allow()
	_ = th.Allow()
	th.Reset()
	if err := th.Allow(); err != nil {
		t.Fatalf("expected allow after reset: %v", err)
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	th := New(50, time.Minute)
	var allowed atomic.Int32
	done := make(chan struct{})
	for i := 0; i < 80; i++ {
		go func() {
			if th.Allow() == nil {
				allowed.Add(1)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 80; i++ {
		<-done
	}
	if allowed.Load() > 50 {
		t.Fatalf("too many allowed: %d (max 50)", allowed.Load())
	}
}
