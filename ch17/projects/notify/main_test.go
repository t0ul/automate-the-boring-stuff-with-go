package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostWebhook(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("method=%s content-type=%s", r.Method, r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	if err := PostWebhook(context.Background(), srv.URL, "backup done"); err != nil {
		t.Fatal(err)
	}
	if got["text"] != "backup done" || got["content"] != "backup done" {
		t.Fatalf("payload = %v", got)
	}

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "no", 403) }))
	defer bad.Close()
	if err := PostWebhook(context.Background(), bad.URL, "x"); err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("expected a 403 error, got %v", err)
	}
}

func TestDesktopCommand(t *testing.T) {
	if c := DesktopCommand("darwin", "t", "m"); c == nil || c.Args[0] != "osascript" || !strings.Contains(c.Args[2], `"m"`) {
		t.Errorf("darwin = %v", c)
	}
	if c := DesktopCommand("linux", "t", "m"); c == nil || c.Args[0] != "notify-send" {
		t.Errorf("linux = %v", c)
	}
	if c := DesktopCommand("plan9", "t", "m"); c != nil {
		t.Errorf("unknown OS should be nil, got %v", c.Args)
	}
}
