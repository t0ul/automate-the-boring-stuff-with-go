package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultName(t *testing.T) {
	cases := map[string]string{
		"https://go.dev/ref/spec":    "spec",
		"https://go.dev/":            "index.html",
		"https://go.dev":             "index.html",
		"https://x.y/a/b/report.pdf": "report.pdf",
	}
	for in, want := range cases {
		if got := DefaultName(in); got != want {
			t.Errorf("DefaultName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDownload(t *testing.T) {
	payload := bytes.Repeat([]byte("0123456789"), 100_000) // 1 MB
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer srv.Close()
	path := filepath.Join(t.TempDir(), "big.bin")
	n, err := Download(context.Background(), srv.URL+"/big.bin", path)
	if err != nil || n != int64(len(payload)) {
		t.Fatalf("Download = %d, %v", n, err)
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, payload) {
		t.Fatalf("file content differs (%d bytes), %v", len(got), err)
	}
}
