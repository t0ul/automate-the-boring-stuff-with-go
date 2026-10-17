package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIndexChangesDupesLargest(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.txt", "same")
	write("sub/b.txt", "same")
	write("c.txt", "a much longer file than the others")
	write("gone.txt", "will be deleted")

	db, err := Open(filepath.Join(t.TempDir(), "idx.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if _, n, err := Index(db, root, time.Now()); err != nil || n != 4 {
		t.Fatalf("first index = %d, %v", n, err)
	}
	var out strings.Builder
	if err := Changes(db, &out); err == nil {
		t.Fatal("changes with one run should fail")
	}

	write("c.txt", "modified!")
	write("new.txt", "brand new")
	os.Remove(filepath.Join(root, "gone.txt"))
	if _, n, err := Index(db, root, time.Now()); err != nil || n != 4 {
		t.Fatalf("second index = %d, %v", n, err)
	}

	out.Reset()
	if err := Changes(db, &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"added     new.txt", "removed   gone.txt", "modified  c.txt"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("changes lacks %q:\n%s", want, out.String())
		}
	}

	out.Reset()
	if err := Dupes(db, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "1 group(s)") || !strings.Contains(out.String(), "a.txt") || !strings.Contains(out.String(), filepath.Join("sub", "b.txt")) {
		t.Errorf("dupes:\n%s", out.String())
	}

	out.Reset()
	if err := Largest(db, &out, 1); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "c.txt") {
		t.Errorf("largest:\n%s", out.String())
	}
}
