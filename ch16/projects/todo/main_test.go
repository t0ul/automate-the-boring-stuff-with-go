package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixedClock() time.Time { return time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC) }

func TestAddListDoneRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "todo.json") // folder does not exist yet
	do := func(args ...string) string {
		var out strings.Builder
		if err := run(args, path, &out, fixedClock); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return out.String()
	}
	do("add", "buy", "milk")
	do("add", "walk the dog")
	if got := do("list"); !strings.Contains(got, "  1 [ ] buy milk") || !strings.Contains(got, "  2 [ ] walk the dog") {
		t.Fatalf("list = %q", got)
	}
	do("done", "1")
	if got := do("list"); strings.Contains(got, "buy milk") {
		t.Fatalf("done item still listed: %q", got)
	}
	if got := do("list", "-all"); !strings.Contains(got, "  1 [x] buy milk") {
		t.Fatalf("list -all = %q", got)
	}
	do("rm", "1")
	if got := do("list"); !strings.Contains(got, "  1 [ ] walk the dog") {
		t.Fatalf("after rm, list = %q", got)
	}

	// the state survived every call because each one reloaded the file
	store, err := Load(path)
	if err != nil || len(store.Items) != 1 || store.Items[0].Text != "walk the dog" {
		t.Fatalf("store = %+v, %v", store, err)
	}
}

func TestUsageErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todo.json")
	for _, args := range [][]string{{}, {"add"}, {"done"}, {"done", "x"}, {"fly"}} {
		var out strings.Builder
		err := run(args, path, &out, fixedClock)
		if !errors.Is(err, errUsage) {
			t.Errorf("run(%v) = %v, want a usage error", args, err)
		}
	}
	var out strings.Builder
	if err := run([]string{"done", "9"}, path, &out, fixedClock); err == nil || errors.Is(err, errUsage) {
		t.Errorf("done 9 on an empty list = %v, want a plain error", err)
	}
}

func TestVersion(t *testing.T) {
	var out strings.Builder
	if err := run([]string{"version"}, "", &out, fixedClock); err != nil || !strings.HasPrefix(out.String(), "todo dev") {
		t.Fatalf("version = %q, %v", out.String(), err)
	}
}
