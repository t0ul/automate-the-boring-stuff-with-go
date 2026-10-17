package main

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunOnce(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	var out strings.Builder
	code, _, err := RunOnce(context.Background(), time.Second, &out, "sh", "-c", "echo hi; exit 3")
	if err != nil || code != 3 || !strings.Contains(out.String(), "hi") {
		t.Fatalf("code=%d err=%v out=%q", code, err, out.String())
	}

	code, took, err := RunOnce(context.Background(), 200*time.Millisecond, &out, "sleep", "5")
	if err == nil || code != -1 || took > 3*time.Second {
		t.Fatalf("deadline not enforced: code=%d took=%s err=%v", code, took, err)
	}

	if _, _, err := RunOnce(context.Background(), time.Second, &out, "no-such-program-xyz"); err == nil {
		t.Fatal("expected an error for a missing program")
	}
}
