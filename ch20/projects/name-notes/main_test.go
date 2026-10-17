package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAccept(t *testing.T) {
	good := map[string]string{
		`{"name": "budget-meeting-minutes"}`:             "budget-meeting-minutes",
		"```json\n{\"name\":\"Lemon-Cake-Recipe\"}\n```": "lemon-cake-recipe",
		`{"name": "budget_meeting_minutes"}`:             "budget-meeting-minutes",
		`{"name": "Grandma's lemon cake"}`:               "grandmas-lemon-cake",
		`{"name": "  packing list:  Lisbon  "}`:          "packing-list-lisbon",
	}
	for reply, want := range good {
		if got, err := Accept(reply); err != nil || got != want {
			t.Errorf("Accept(%q) = %q, %v; want %q", reply, got, err, want)
		}
	}
	bad := []string{
		`{"name": "minutes"}`,                 // one word
		`{"name": "a-b-c-d-e-f-g"}`,           // too many words
		`{"name": "notes.backup"}`,            // tidies to one word
		`{"name": "../elsewhere"}`,            // tidies to one word; the slashes never survive
		`{"name": ""}`,                        // empty
		`the note is about a budget meeting`,  // no JSON
		`{"title": "budget-meeting-minutes"}`, // wrong key
	}
	for _, reply := range bad {
		if got, err := Accept(reply); err == nil {
			t.Errorf("Accept(%q) = %q, want an error", reply, got)
		}
	}
}

func fakeModel(t *testing.T, answers map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		content := `{"name": "untitled-note"}`
		for key, reply := range answers {
			if strings.Contains(req.Messages[1].Content, key) {
				content = reply
			}
		}
		fmt.Fprintf(w, `{"choices":[{"message":{"role":"assistant","content":%q}}]}`, content)
	}))
}

func TestPlanHandlesCollisionsAndRejects(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.txt", "cake recipe one")
	write("b.txt", "cake recipe two")
	write("c.md", "gibberish")
	write("lemon-cake.txt", "already named")
	write("photo.jpg", "not a note")

	srv := fakeModel(t, map[string]string{
		"cake recipe":   `{"name": "lemon-cake"}`,
		"gibberish":     `{"name": "???"}`,
		"already named": `{"name": "lemon-cake"}`,
	})
	defer srv.Close()
	client := Client{BaseURL: srv.URL + "/v1", Model: "fake", HTTP: srv.Client()}

	plan, err := Plan(context.Background(), client, dir)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Rename{}
	for _, r := range plan {
		got[r.From] = r
	}
	if len(plan) != 4 {
		t.Fatalf("plan has %d entries, want 4 (photo.jpg skipped): %+v", len(plan), plan)
	}
	if got["a.txt"].To != "lemon-cake-2.txt" || got["b.txt"].To != "lemon-cake-3.txt" {
		t.Errorf("collision handling: a=%q b=%q", got["a.txt"].To, got["b.txt"].To)
	}
	if got["c.md"].To != "" || !strings.Contains(got["c.md"].Note, "does not match") {
		t.Errorf("rejected name should be kept with a reason, got %+v", got["c.md"])
	}
	if got["lemon-cake.txt"].To != "lemon-cake-4.txt" {
		t.Errorf("existing file name must count as taken, got %+v", got["lemon-cake.txt"])
	}
}
