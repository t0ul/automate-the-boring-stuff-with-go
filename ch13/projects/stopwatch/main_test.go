package main

import (
	"strings"
	"testing"
	"time"
)

func TestRunWithFakeClock(t *testing.T) {
	// a clock that advances a fixed amount every time it is read
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	steps := []time.Duration{0, 1500 * time.Millisecond, 700 * time.Millisecond, 2 * time.Second}
	i := 0
	clock := func() time.Time {
		now = now.Add(steps[i])
		i++
		return now
	}
	var out strings.Builder
	laps := Run(strings.NewReader("\n\nq\n"), &out, clock)
	if len(laps) != 2 {
		t.Fatalf("got %d laps: %+v", len(laps), laps)
	}
	if laps[1].Split != 700*time.Millisecond || laps[1].Total != 2200*time.Millisecond {
		t.Fatalf("lap 2 = %+v", laps[1])
	}
	if !strings.Contains(out.String(), "Lap #2:   0.70s (  2.20s total)") {
		t.Fatalf("output = %q", out.String())
	}
}
