// search-term builds a URL with query parameters from whatever you type, calls a JSON
// API and prints the results. The API is Open Library's book search, which needs no key.
//
//	go run ./ch11/projects/search-term go programming language
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

// Result is the part of one Open Library hit we care about.
type Result struct {
	Title       string   `json:"title"`
	Authors     []string `json:"author_name"`
	PublishYear int      `json:"first_publish_year"`
}

// SearchURL builds <base>?q=<terms>&limit=<n>, escaping the terms properly.
func SearchURL(base, terms string, limit int) string {
	q := url.Values{}
	q.Set("q", terms)
	q.Set("limit", fmt.Sprint(limit))
	q.Set("fields", "title,author_name,first_publish_year")
	return base + "?" + q.Encode()
}

// Search calls the API and decodes its "docs" array.
func Search(ctx context.Context, searchURL string) ([]Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "automate-the-boring-stuff-with-go/2026")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search: %s", resp.Status)
	}
	var payload struct {
		Docs []Result `json:"docs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return payload.Docs, nil
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: search-term <words to search for>")
	}
	terms := strings.Join(args, " ")
	u := SearchURL("https://openlibrary.org/search.json", terms, 5)
	fmt.Println("GET", u)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	results, err := Search(ctx, u)
	if err != nil {
		return err
	}
	for i, r := range results {
		fmt.Printf("%d. %s (%d) by %s\n", i+1, r.Title, r.PublishYear, strings.Join(r.Authors, ", "))
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
