package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestNextZipName(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "docs")
	if err := os.Mkdir(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := NextZipName(folder)
	if err != nil || got != folder+"_1.zip" {
		t.Fatalf("first name = %q, %v", got, err)
	}
	if err := os.WriteFile(folder+"_1.zip", nil, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = NextZipName(folder)
	if err != nil || got != folder+"_2.zip" {
		t.Fatalf("second name = %q, %v", got, err)
	}
}

func TestBackupToZip(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "docs")
	for _, p := range []string{"a.txt", "sub/b.txt", "sub/deeper/c.txt"} {
		full := filepath.Join(folder, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("content of "+p), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	zipPath := filepath.Join(dir, "docs_1.zip")
	n, err := BackupToZip(folder, zipPath)
	if err != nil || n != 3 {
		t.Fatalf("BackupToZip = %d, %v", n, err)
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	slices.Sort(names)
	want := []string{"docs/a.txt", "docs/sub/b.txt", "docs/sub/deeper/c.txt"}
	if !slices.Equal(names, want) {
		t.Fatalf("entries = %v, want %v", names, want)
	}
}
