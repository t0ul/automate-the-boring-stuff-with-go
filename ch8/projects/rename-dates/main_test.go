package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEuropeanName(t *testing.T) {
	cases := []struct {
		in, want string
		ok       bool
	}{
		{"invoice_05-21-2026.pdf", "invoice_21-05-2026.pdf", true},
		{"5-1-2026 notes.txt", "1-5-2026 notes.txt", true},
		{"photo_13-40-2026.jpg", "photo_13-40-2026.jpg", false},
		{"no-date-here.txt", "no-date-here.txt", false},
		{"2026-05-21.txt", "2026-05-21.txt", false},
	}
	for _, c := range cases {
		got, ok := EuropeanName(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("EuropeanName(%q) = %q, %v; want %q, %v", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestPlan(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"a_01-02-2026.txt", "b.txt", "c_12-31-2025.csv"} {
		if err := os.WriteFile(filepath.Join(dir, n), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan, err := Plan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != 2 || plan["a_01-02-2026.txt"] != "a_02-01-2026.txt" || plan["c_12-31-2025.csv"] != "c_31-12-2025.csv" {
		t.Fatalf("plan = %v", plan)
	}
}
