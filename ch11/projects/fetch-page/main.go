// fetch-page makes one HTTP GET the right way: a client with a timeout, a context, a
// User-Agent that says who you are, a status check, and a body that is always closed.
//
//	go run ./ch11/projects/fetch-page https://www.t0ul.com/
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const userAgent = "automate-the-boring-stuff-with-go/2026 (+https://www.t0ul.com)"

// Client is the one HTTP client every project in this chapter shares the shape of.
var Client = &http.Client{Timeout: 15 * time.Second}

// Fetch GETs url and returns the body, up to limit bytes, plus the content type.
func Fetch(ctx context.Context, pageURL string, limit int64) (body []byte, contentType string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := Client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("get %s: %w", pageURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("get %s: %s", pageURL, resp.Status)
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return nil, "", fmt.Errorf("read %s: %w", pageURL, err)
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: fetch-page <url>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	body, ct, err := Fetch(ctx, args[0], 1<<20)
	if err != nil {
		return err
	}
	fmt.Printf("content-type: %s\nbytes: %d\n\n", ct, len(body))
	if len(body) > 500 {
		body = body[:500]
	}
	fmt.Printf("%s\n", body)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
