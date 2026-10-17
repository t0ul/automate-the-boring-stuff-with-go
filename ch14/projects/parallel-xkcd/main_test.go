package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDownloadRecentUsesWorkers(t *testing.T) {
	var inFlight, maxInFlight atomic.Int32
	var srvURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/info.0.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"num":12,"title":"latest","img":"%s/comics/12.png"}`, srvURL)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/info.0.json"):
			var n int
			fmt.Sscanf(r.URL.Path, "/%d/info.0.json", &n)
			fmt.Fprintf(w, `{"num":%d,"title":"c","img":"%s/comics/%d.png"}`, n, srvURL, n)
		case strings.HasPrefix(r.URL.Path, "/comics/"):
			cur := inFlight.Add(1)
			for {
				old := maxInFlight.Load()
				if cur <= old || maxInFlight.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			inFlight.Add(-1)
			fmt.Fprint(w, "PNG")
		default:
			http.NotFound(w, r)
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL = srv.URL

	dir := t.TempDir()
	count, err := Fetcher{Base: srv.URL}.DownloadRecent(context.Background(), 8, 3, dir)
	if err != nil || count != 8 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 8 {
		t.Fatalf("%d files written", len(entries))
	}
	if m := maxInFlight.Load(); m < 2 || m > 3 {
		t.Fatalf("max concurrent image downloads = %d, want 2 or 3", m)
	}
}

func TestDownloadRecentStopsOnError(t *testing.T) {
	var srvURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/info.0.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"num":5,"title":"latest","img":"%s/comics/5.png"}`, srvURL)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/3/info.0.json" {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/info.0.json") {
			var n int
			fmt.Sscanf(r.URL.Path, "/%d/info.0.json", &n)
			fmt.Fprintf(w, `{"num":%d,"title":"c","img":"%s/comics/%d.png"}`, n, srvURL, n)
			return
		}
		fmt.Fprint(w, "PNG")
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL = srv.URL

	_, err := Fetcher{Base: srv.URL}.DownloadRecent(context.Background(), 5, 2, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "comic 3") {
		t.Fatalf("expected the comic 3 error, got %v", err)
	}
}
