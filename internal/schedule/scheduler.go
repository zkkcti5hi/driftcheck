// Package schedule provides periodic drift-check execution.
package schedule

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// RunFunc is the signature of a function that performs a single drift check
// and returns the results. Callers inject this so the scheduler stays
// decoupled from detection logic.
type RunFunc func(ctx context.Context) ([]drift.Result, error)

// OnResultFunc is called after every successful run with the latest results.
type OnResultFunc func(results []drift.Result)

// Scheduler triggers a RunFunc on a fixed interval.
type Scheduler struct {
	interval time.Duration
	run      RunFunc
	onResult OnResultFunc
}

// New creates a Scheduler that calls run every interval and passes results to
// onResult. onResult may be nil if the caller does not need the results.
func New(interval time.Duration, run RunFunc, onResult OnResultFunc) *Scheduler {
	return &Scheduler{
		interval: interval,
		run:      run,
		onResult: onResult,
	}
}

// Start blocks until ctx is cancelled, executing the RunFunc immediately and
// then on every tick of the configured interval.
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("scheduler: starting with interval %s", s.interval)

	// Run once immediately so the first result is available right away.
	s.tick(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			log.Println("scheduler: stopping")
			return
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	results, err := s.run(ctx)
	if err != nil {
		log.Printf("scheduler: run error: %v", err)
		return
	}
	if s.onResult != nil {
		s.onResult(results)
	}
}
