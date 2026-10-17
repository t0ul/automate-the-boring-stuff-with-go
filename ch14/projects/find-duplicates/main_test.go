package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicates(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.txt", "same content")
	write("sub/b.txt", "same content")
	write("sub/deep/c.txt", "same content")
	write("d.txt", "different")
	write("e.txt", "other twin")
	write("f.txt", "other twin")
	write("tiny.txt", "")

	groups, err := FindDuplicates(dir, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 || len(groups[0]) != 3 || len(groups[1]) != 2 {
		t.Fatalf("groups = %v", groups)
	}
	if filepath.Base(groups[0][0]) != "a.txt" || filepath.Base(groups[1][0]) != "e.txt" {
		t.Fatalf("ordering = %v", groups)
	}
}

func TestHashFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "x")
	if err := os.WriteFile(p, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	sum, err := HashFile(p)
	if err != nil || sum != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("sum = %s, %v", sum, err)
	}
}
