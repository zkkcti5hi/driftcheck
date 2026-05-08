package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/schedule"
)

func TestScheduler_CallsRunImmediately(t *testing.T) {
	var callCount int32

	runFn := func(_ context.Context) ([]drift.Result, error) {
		atomic.AddInt32(&callCount, 1)
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := schedule.New(10*time.Second, runFn, nil)

	done := make(chan struct{})
	go func() {
		s.Start(ctx)
		close(done)
	}()

	// Give the goroutine a moment to execute the immediate tick.
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	if got := atomic.LoadInt32(&callCount); got < 1 {
		t.Fatalf("expected at least 1 call, got %d", got)
	}
}

func TestScheduler_TicksOnInterval(t *testing.T) {
	var callCount int32

	runFn := func(_ context.Context) ([]drift.Result, error) {
		atomic.AddInt32(&callCount, 1)
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := schedule.New(30*time.Millisecond, runFn, nil)

	done := make(chan struct{})
	go func() {
		s.Start(ctx)
		close(done)
	}()

	time.Sleep(120 * time.Millisecond)
	cancel()
	<-done

	// Expect at least 3 calls (immediate + ~2–3 ticks in 120 ms).
	if got := atomic.LoadInt32(&callCount); got < 3 {
		t.Fatalf("expected >= 3 calls, got %d", got)
	}
}

func TestScheduler_RunErrorDoesNotStop(t *testing.T) {
	var callCount int32

	runFn := func(_ context.Context) ([]drift.Result, error) {
		atomic.AddInt32(&callCount, 1)
		return nil, errors.New("boom")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := schedule.New(30*time.Millisecond, runFn, nil)

	done := make(chan struct{})
	go func() {
		s.Start(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	if got := atomic.LoadInt32(&callCount); got < 2 {
		t.Fatalf("expected scheduler to keep running after error, got %d calls", got)
	}
}

func TestScheduler_OnResultCalledWithResults(t *testing.T) {
	expected := []drift.Result{{Service: "web", Drifted: true}}

	runFn := func(_ context.Context) ([]drift.Result, error) {
		return expected, nil
	}

	received := make(chan []drift.Result, 1)
	onResult := func(r []drift.Result) {
		select {
		case received <- r:
		default:
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := schedule.New(10*time.Second, runFn, onResult)
	go s.Start(ctx)

	select {
	case r := <-received:
		if len(r) != 1 || r[0].Service != "web" {
			t.Fatalf("unexpected results: %+v", r)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for onResult callback")
	}
}
