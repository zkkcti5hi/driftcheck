// Package ratelimit provides a simple token-bucket style rate limiter used
// to prevent driftcheck scans from running more frequently than a configured
// minimum interval.
//
// Usage:
//
//	l := ratelimit.New(30 * time.Second)
//
//	if err := l.Allow(); err != nil {
//		// scan skipped — too soon since last run
//		fmt.Println("next scan in", l.Remaining())
//		return
//	}
//	// perform scan …
//
// The zero interval (or any negative value) disables rate limiting entirely,
// allowing every call to Allow to succeed.
//
// Limiter is safe for concurrent use.
package ratelimit
