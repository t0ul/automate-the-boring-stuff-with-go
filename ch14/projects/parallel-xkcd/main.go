// parallel-xkcd downloads the N most recent xkcd comics with a handful of goroutines at
// once. It is chapter 11's downloader with a worker pool, and it is the pattern for any
// "do this for a thousand URLs" job: a bounded number of workers, a context that stops
// everything on ctrl+c, and the first error reported instead of swallowed.
//
//	go run ./ch14/projects/parallel-xkcd -n 20 -workers 4 -dir xkcd
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
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"
)

var client = &http.Client{Timeout: 30 * time.Second}

// Comic is the subset of xkcd's JSON we use.
type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Img   string `json:"img"`
}

// Fetcher knows where the API lives, so tests can point it at a local server.
type Fetcher struct {
	Base string
}

func (f Fetcher) comic(ctx context.Context, url string) (Comic, error) {
	var c Comic
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return c, err
	}
	req.Header.Set("User-Agent", "automate-the-boring-stuff-with-go/2026")
	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return c, json.NewDecoder(resp.Body).Decode(&c)
}

// Latest returns the newest comic number.
func (f Fetcher) Latest(ctx context.Context) (int, error) {
	c, err := f.comic(ctx, f.Base+"/info.0.json")
	return c.Num, err
}

// Download fetches comic n's metadata and image into dir. It reports whether it wrote a
// file (false when it was already there).
func (f Fetcher) Download(ctx context.Context, n int, dir string) (bool, error) {
	c, err := f.comic(ctx, fmt.Sprintf("%s/%d/info.0.json", f.Base, n))
	if err != nil {
		return false, fmt.Errorf("comic %d: %w", n, err)
	}
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
		return false, fmt.Errorf("image %d: %w", n, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("image %d: %s", n, resp.Status)
	}
	tmp := name + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return false, err
	}
	if err := out.Close(); err != nil {
		return false, err
	}
	return true, os.Rename(tmp, name)
}

// DownloadRecent fetches the n newest comics with at most workers in flight, and returns
// how many files were written. The first error cancels the rest.
func (f Fetcher) DownloadRecent(ctx context.Context, n, workers int, dir string) (int, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}
	latest, err := f.Latest(ctx)
	if err != nil {
		return 0, err
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)
	written := make(chan int, n) // each worker sends 1 when it writes a file
	for num := latest; num > latest-n && num > 0; num-- {
		g.Go(func() error {
			wrote, err := f.Download(ctx, num, dir)
			if err != nil {
				return err
			}
			if wrote {
				written <- 1
				fmt.Printf("#%d\n", num)
			}
			return nil
		})
	}
	err = g.Wait()
	close(written)
	count := 0
	for range written {
		count++
	}
	return count, err
}

func run(args []string) error {
	flags := flag.NewFlagSet("parallel-xkcd", flag.ContinueOnError)
	n := flags.Int("n", 20, "how many recent comics")
	workers := flags.Int("workers", 4, "how many downloads at once")
	dir := flags.String("dir", "xkcd", "folder to save into")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *n < 1 || *workers < 1 {
		return errors.New("-n and -workers must be at least 1")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	start := time.Now()
	count, err := Fetcher{Base: "https://xkcd.com"}.DownloadRecent(ctx, *n, *workers, *dir)
	fmt.Printf("%d new file(s) in %s after %s\n", count, *dir, time.Since(start).Round(time.Millisecond))
	return err
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
