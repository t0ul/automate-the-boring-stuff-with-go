package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestParseSize(t *testing.T) {
	cases := map[string]int64{"512": 512, "1KB": 1024, "1.5MB": 1536 * 1024, "2gb": 2 << 30, " 3 kb ": 3072}
	for in, want := range cases {
		got, err := ParseSize(in)
		if err != nil || got != want {
			t.Errorf("ParseSize(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "abc", "-1MB", "1TB"} {
		if _, err := ParseSize(bad); err == nil {
			t.Errorf("ParseSize(%q) should fail", bad)
		}
	}
}

func TestHumanSize(t *testing.T) {
	if got := HumanSize(1536); got != "1.5 KB" {
		t.Errorf("HumanSize(1536) = %q", got)
	}
	if got := HumanSize(10); got != "10 B" {
		t.Errorf("HumanSize(10) = %q", got)
	}
}

func TestFindLarge(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, size int) {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, bytes.Repeat([]byte("x"), size), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("small.txt", 10)
	write("big.bin", 5000)
	write("nested/bigger.bin", 9000)
	got, err := FindLarge(dir, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || filepath.Base(got[0].Path) != "bigger.bin" || filepath.Base(got[1].Path) != "big.bin" {
		t.Fatalf("FindLarge = %+v", got)
	}
}
