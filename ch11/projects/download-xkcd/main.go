// download-xkcd saves the most recent xkcd comics to a folder. xkcd publishes a JSON
// description of every comic, so there is no HTML to scrape: ask for the metadata, then
// fetch the image it points to. Comics already on disk are skipped, and there is a short
// pause between requests because it is someone else's server.
//
//	go run ./ch11/projects/download-xkcd -n 5 -dir xkcd
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

// Comic is the subset of xkcd's JSON we use.
type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Img   string `json:"img"`
	Alt   string `json:"alt"`
}

// Fetcher knows where the API lives, so tests can point it at a local server.
type Fetcher struct {
	Base  string        // "https://xkcd.com"
	Pause time.Duration // between requests
}

// Latest returns the newest comic's metadata.
func (f Fetcher) Latest(ctx context.Context) (Comic, error) {
	return f.getJSON(ctx, f.Base+"/info.0.json")
}

// Get returns comic number n.
func (f Fetcher) Get(ctx context.Context, n int) (Comic, error) {
	return f.getJSON(ctx, fmt.Sprintf("%s/%d/info.0.json", f.Base, n))
}

func (f Fetcher) getJSON(ctx context.Context, apiURL string) (Comic, error) {
	var c Comic
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return c, err
	}
	req.Header.Set("User-Agent", "automate-the-boring-stuff-with-go/2026")
	resp, err := client.Do(req)
	if err != nil {
		return c, fmt.Errorf("get %s: %w", apiURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c, fmt.Errorf("get %s: %s", apiURL, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, fmt.Errorf("decode %s: %w", apiURL, err)
	}
	return c, nil
}

// SaveImage downloads c's image into dir as <num>-<original name> and reports whether it
// downloaded anything (false when the file already existed).
func (f Fetcher) SaveImage(ctx context.Context, c Comic, dir string) (bool, error) {
	name := filepath.Join(dir, fmt.Sprintf("%d-%s", c.Num, path.Base(c.Img)))
	if _, err := os.Stat(name); err == nil {
		return false, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Img, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("image %d: %w", c.Num, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("image %d: %s", c.Num, resp.Status)
	}
	tmp := name + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return false, fmt.Errorf("write %s: %w", tmp, err)
	}
	if err := out.Close(); err != nil {
		return false, err
	}
	return true, os.Rename(tmp, name)
}

// DownloadRecent saves the n newest comics and returns how many files it wrote.
func (f Fetcher) DownloadRecent(ctx context.Context, n int, dir string) (int, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}
	latest, err := f.Latest(ctx)
	if err != nil {
		return 0, err
	}
	written := 0
	for num := latest.Num; num > latest.Num-n && num > 0; num-- {
		c := latest
		if num != latest.Num {
			time.Sleep(f.Pause)
			if c, err = f.Get(ctx, num); err != nil {
				return written, err
			}
		}
		saved, err := f.SaveImage(ctx, c, dir)
		if err != nil {
			return written, err
		}
		if saved {
			written++
			fmt.Printf("#%d %q\n", c.Num, c.Title)
		} else {
			fmt.Printf("#%d already here\n", c.Num)
		}
	}
	return written, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("download-xkcd", flag.ContinueOnError)
	n := flags.Int("n", 5, "how many recent comics")
	dir := flags.String("dir", "xkcd", "folder to save into")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *n < 1 {
		return errors.New("-n must be at least 1")
	}
	f := Fetcher{Base: "https://xkcd.com", Pause: 500 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	written, err := f.DownloadRecent(ctx, *n, *dir)
	if err != nil {
		return err
	}
	fmt.Printf("%d new file(s) in %s\n", written, *dir)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
