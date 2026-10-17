package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.UserAgent(), "automate-the-boring-stuff-with-go") {
			t.Errorf("User-Agent = %q", r.UserAgent())
		}
		switch r.URL.Path {
		case "/ok":
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("hello, world"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	body, ct, err := Fetch(context.Background(), srv.URL+"/ok", 1024)
	if err != nil || string(body) != "hello, world" || ct != "text/plain" {
		t.Fatalf("Fetch = %q, %q, %v", body, ct, err)
	}
	if _, _, err := Fetch(context.Background(), srv.URL+"/missing", 1024); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected a 404 error, got %v", err)
	}
	body, _, err = Fetch(context.Background(), srv.URL+"/ok", 5)
	if err != nil || string(body) != "hello" {
		t.Fatalf("limit not applied: %q, %v", body, err)
	}
}
