package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func fakeXKCD(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	var srvURL string
	comic := func(n int) string {
		return fmt.Sprintf(`{"num":%d,"title":"Comic %d","img":"%s/comics/c%d.png","alt":"alt"}`, n, n, srvURL, n)
	}
	mux.HandleFunc("/info.0.json", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, comic(3)) })
	mux.HandleFunc("/2/info.0.json", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, comic(2)) })
	mux.HandleFunc("/1/info.0.json", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, comic(1)) })
	mux.HandleFunc("/comics/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "PNGDATA "+r.URL.Path) })
	srv := httptest.NewServer(mux)
	srvURL = srv.URL
	return srv
}

func TestDownloadRecent(t *testing.T) {
	srv := fakeXKCD(t)
	defer srv.Close()
	f := Fetcher{Base: srv.URL}
	dir := t.TempDir()

	n, err := f.DownloadRecent(context.Background(), 2, dir)
	if err != nil || n != 2 {
		t.Fatalf("first run wrote %d, %v", n, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 2 || entries[0].Name() != "2-c2.png" || entries[1].Name() != "3-c3.png" {
		t.Fatalf("files = %v", entries)
	}
	n, err = f.DownloadRecent(context.Background(), 3, dir)
	if err != nil || n != 1 {
		t.Fatalf("second run wrote %d (want only #1), %v", n, err)
	}
}
