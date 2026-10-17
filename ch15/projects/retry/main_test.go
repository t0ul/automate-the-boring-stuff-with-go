package main

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// recordingSleep never sleeps; it writes down what it was asked to wait.
func recordingSleep(waits *[]time.Duration) func(context.Context, time.Duration) error {
	return func(_ context.Context, d time.Duration) error {
		*waits = append(*waits, d)
		return nil
	}
}

func TestDoSucceedsAfterFailures(t *testing.T) {
	var waits []time.Duration
	p := Policy{Attempts: 5, Wait: 100 * time.Millisecond, MaxWait: 250 * time.Millisecond, Sleep: recordingSleep(&waits)}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		if calls < 4 {
			return errors.New("nope")
		}
		return nil
	})
	if err != nil || calls != 4 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
	want := []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 250 * time.Millisecond}
	if len(waits) != len(want) {
		t.Fatalf("waits = %v", waits)
	}
	for i := range want {
		if waits[i] != want[i] {
			t.Fatalf("waits = %v, want %v", waits, want)
		}
	}
}

func TestDoGivesUp(t *testing.T) {
	var waits []time.Duration
	p := Policy{Attempts: 3, Wait: time.Millisecond, MaxWait: time.Second, Sleep: recordingSleep(&waits)}
	err := p.Do(context.Background(), func() error { return errors.New("still broken") })
	if err == nil || !strings.Contains(err.Error(), "after 3 attempt(s)") || !errors.Is(err, err) {
		t.Fatalf("err = %v", err)
	}
	if len(waits) != 2 {
		t.Fatalf("slept %d times, want 2", len(waits))
	}
}

func TestDoStopsWhenContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := Policy{Attempts: 10, Wait: time.Hour, MaxWait: time.Hour} // real sleep, which must not happen
	calls := 0
	err := p.Do(ctx, func() error { calls++; return errors.New("x") })
	if calls != 1 || !errors.Is(err, context.Canceled) {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}
