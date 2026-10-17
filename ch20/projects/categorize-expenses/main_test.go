package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestParseNeverLeavesTheAllowlist(t *testing.T) {
	cases := map[string]string{
		`{"category": "food"}`: "food",
		"Sure! Here you go:\n```json\n{\"category\":\"Transport\"}\n```": "transport",
		`{"category": "groceries"}`:                                      "other",
		`{"category": "food and drink"}`:                                 "other",
		`I think it is food.`:                                            "other",
		`{"category": 42}`:                                               "other",
		``:                                                               "other",
	}
	for reply, want := range cases {
		if got, _ := Parse(reply); got != want {
			t.Errorf("Parse(%q) = %q, want %q", reply, got, want)
		}
	}
}

// fakeModel is a real HTTP server speaking the chat API, answering from a table keyed on
// the item named in the user prompt. It stands in for the model, not for the client.
func fakeModel(t *testing.T, answers map[string]string, calls *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls++
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || r.URL.Path != "/v1/chat/completions" {
			http.Error(w, "bad request", 400)
			return
		}
		content := `{"category": "other"}`
		for item, reply := range answers {
			if strings.Contains(req.Messages[1].Content, item) {
				content = reply
			}
		}
		fmt.Fprintf(w, `{"choices":[{"message":{"role":"assistant","content":%q}}]}`, content)
	}))
}

func TestFillCachesAndValidates(t *testing.T) {
	calls := 0
	srv := fakeModel(t, map[string]string{
		"bus pass": `{"category": "transport"}`,
		"coffee":   `{"category": "food"}`,
		"socks":    `{"category": "clothing"}`, // not in the allowlist: becomes other
	}, &calls)
	defer srv.Close()
	k := &Classifier{Client: Client{BaseURL: srv.URL + "/v1", Model: "fake", HTTP: srv.Client()}}
	rows, err := ReadRows(strings.NewReader("date,item,amount\nd,bus pass,32\nd,coffee,5\nd,Coffee,5\nd,socks,9\nd,known,1\n"))
	if err != nil {
		t.Fatal(err)
	}
	rows[4].Category = "books" // already filled: must not be asked
	Fill(context.Background(), k, rows)
	got := []string{rows[0].Category, rows[1].Category, rows[2].Category, rows[3].Category, rows[4].Category}
	want := []string{"transport", "food", "food", "other", "books"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
	if calls != 3 || k.Calls != 3 {
		t.Errorf("model called %d times (client counted %d), want 3: coffee is cached and books was pre-filled", calls, k.Calls)
	}
}

func TestLiveModel(t *testing.T) {
	if os.Getenv("LLM_LIVE") == "" {
		t.Skip("set LLM_LIVE=1 with a local model server running to run this")
	}
	k := &Classifier{Client: FromEnv()}
	if got := k.Category(context.Background(), "monthly bus pass", "32.00"); got != "transport" {
		t.Errorf("live model said %q for a bus pass", got)
	}
}
