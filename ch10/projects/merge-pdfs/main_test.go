package main

import (
	"context"
	"path/filepath"
	"testing"
)

func TestInputsFromFolder(t *testing.T) {
	in, err := Inputs([]string{"."}, "merged.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if len(in) != 3 || filepath.Base(in[0]) != "part1.pdf" || filepath.Base(in[2]) != "part3.pdf" {
		t.Fatalf("Inputs = %v", in)
	}
}

func TestMerge(t *testing.T) {
	out := filepath.Join(t.TempDir(), "merged.pdf")
	n, err := Merge(context.Background(), []string{"part1.pdf", "part2.pdf", "part3.pdf"}, out)
	if err != nil {
		t.Fatal(err)
	}
	if n != 6 {
		t.Fatalf("merged page count = %d, want 6", n)
	}
}
