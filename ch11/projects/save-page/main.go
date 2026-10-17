// save-page downloads a URL to a file without ever holding the whole thing in memory:
// the response body streams straight into the file with io.Copy.
//
//	go run ./ch11/projects/save-page -out go-spec.html https://go.dev/ref/spec
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"io"
)

var client = &http.Client{Timeout: 2 * time.Minute}

// Download streams url into path and returns the number of bytes written.
func Download(ctx context.Context, rawURL, path string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "automate-the-boring-stuff-with-go/2026")
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get %s: %s", rawURL, resp.Status)
	}
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", path, err)
	}
	n, copyErr := io.Copy(f, resp.Body)
	if err := f.Close(); copyErr == nil && err != nil {
		copyErr = err
	}
	if copyErr != nil {
		return n, fmt.Errorf("save %s: %w", path, copyErr)
	}
	return n, nil
}

// DefaultName picks a file name from the last part of the URL path, or index.html.
func DefaultName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "download"
	}
	base := path.Base(u.Path)
	if base == "." || base == "/" || base == "" {
		return "index.html"
	}
	return base
}

func run(args []string) error {
	fs := flag.NewFlagSet("save-page", flag.ContinueOnError)
	out := fs.String("out", "", "file to write (default: last part of the URL)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: save-page [-out file] <url>")
	}
	rawURL := fs.Arg(0)
	if *out == "" {
		*out = DefaultName(rawURL)
	}
	n, err := Download(context.Background(), rawURL, *out)
	if err != nil {
		return err
	}
	fmt.Printf("saved %d bytes to %s\n", n, *out)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
