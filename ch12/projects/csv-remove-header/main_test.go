package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFolder(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.csv", "h1,h2\n1,2\n3,4\n")
	write("b.csv", "only,header\n")
	write("notes.txt", "not a csv\n")

	done, err := ProcessFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(done) != 2 || done["a.csv"] != 2 || done["b.csv"] != 0 {
		t.Fatalf("done = %v", done)
	}
	got, err := os.ReadFile(filepath.Join(dir, "headerRemoved", "a.csv"))
	if err != nil || string(got) != "1,2\n3,4\n" {
		t.Fatalf("a.csv = %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "headerRemoved", "notes.txt")); err == nil {
		t.Fatal("notes.txt should not have been processed")
	}
}
