package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchURL(t *testing.T) {
	got := SearchURL("https://example.org/search.json", "go & tea", 3)
	want := "https://example.org/search.json?fields=title%2Cauthor_name%2Cfirst_publish_year&limit=3&q=go+%26+tea"
	if got != want {
		t.Fatalf("SearchURL = %s\nwant %s", got, want)
	}
}

func TestSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "gophers" {
			t.Errorf("q = %q", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"numFound":1,"docs":[{"title":"The Go Programming Language","author_name":["Alan Donovan","Brian Kernighan"],"first_publish_year":2015}]}`))
	}))
	defer srv.Close()
	results, err := Search(context.Background(), SearchURL(srv.URL, "gophers", 1))
	if err != nil || len(results) != 1 || results[0].PublishYear != 2015 || len(results[0].Authors) != 2 {
		t.Fatalf("Search = %+v, %v", results, err)
	}
}
