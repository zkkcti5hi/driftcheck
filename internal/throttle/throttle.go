// Package throttle provides a token-bucket style throttle for limiting
// the rate at which drift-check notifications or webhook calls are dispatched.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Throttle controls how many events may pass through per window.
type Throttle struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	tokens   int
	windowAt time.Time
	now      func() time.Time
}

// New creates a Throttle that allows at most max events per window duration.
// A zero or negative max disables throttling (all calls are allowed).
func New(max int, window time.Duration) *Throttle {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, now func() time.Time) *Throttle {
	t := &Throttle{
		max:    max,
		window: window,
		now:    now,
	}
	t.windowAt = now()
	t.tokens = max
	return t
}

// Allow reports whether the caller is permitted to proceed.
// It returns an error describing the wait time when throttled.
func (t *Throttle) Allow() error {
	if t.max <= 0 || t.window <= 0 {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if now.Sub(t.windowAt) >= t.window {
		t.windowAt = now
		t.tokens = t.max
	}

	if t.tokens <= 0 {
		retry := t.windowAt.Add(t.window).Sub(now)
		return fmt.Errorf("throttled: retry after %s", retry.Round(time.Millisecond))
	}

	t.tokens--
	return nil
}

// Remaining returns the number of tokens left in the current window.
func (t *Throttle) Remaining() int {
	if t.max <= 0 || t.window <= 0 {
		return -1
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.now().Sub(t.windowAt) >= t.window {
		return t.max
	}
	return t.tokens
}

// Reset clears the current window, restoring all tokens immediately.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.windowAt = t.now()
	t.tokens = t.max
}
