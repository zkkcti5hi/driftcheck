// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently drift scans are allowed to run programmatically.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when a scan is attempted before the minimum
// interval between scans has elapsed.
var ErrRateLimited = errors.New("rate limited: minimum scan interval not elapsed")

// Limiter enforces a minimum interval between successive drift scans.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Limiter that allows at most one scan per interval.
// A zero or negative interval disables rate limiting.
func New(interval time.Duration) *Limiter {
	return &Limiter{
		interval: interval,
		now:      time.Now,
	}
}

// Allow returns nil if a scan may proceed, or ErrRateLimited if the minimum
// interval since the last allowed scan has not yet elapsed.
// When nil is returned the internal timestamp is updated.
func (l *Limiter) Allow() error {
	if l.interval <= 0 {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if !l.last.IsZero() && now.Sub(l.last) < l.interval {
		return ErrRateLimited
	}
	l.last = now
	return nil
}

// Reset clears the last-scan timestamp, allowing the next call to Allow to
// succeed regardless of when it occurs.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = time.Time{}
}

// Remaining returns how long the caller must wait before Allow will succeed.
// Returns 0 if a scan is currently permitted.
func (l *Limiter) Remaining() time.Duration {
	if l.interval <= 0 {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.last.IsZero() {
		return 0
	}
	elapsed := l.now().Sub(l.last)
	if elapsed >= l.interval {
		return 0
	}
	return l.interval - elapsed
}
