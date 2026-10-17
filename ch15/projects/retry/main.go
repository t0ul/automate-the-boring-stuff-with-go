// retry runs an operation again when it fails, waiting longer each time, up to a limit.
// Networks fail for a second and come back; a script that gives up at the first error is
// a script you have to babysit. The demo talks to a local server that fails twice first.
//
//	go run ./ch15/projects/retry
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"time"
)

// Policy says how many attempts to make and how long to wait after the first failure;
// each later wait doubles, up to MaxWait.
type Policy struct {
	Attempts int
	Wait     time.Duration
	MaxWait  time.Duration
	Sleep    func(context.Context, time.Duration) error // swapped out in tests
}

// Do runs op until it returns nil, the attempts run out, or the context ends. The error
// returned is the last one op produced, wrapped with the attempt count.
func (p Policy) Do(ctx context.Context, op func() error) error {
	sleep := p.Sleep
	if sleep == nil {
		sleep = sleepCtx
	}
	wait := p.Wait
	var last error
	for attempt := 1; attempt <= p.Attempts; attempt++ {
		last = op()
		if last == nil {
			return nil
		}
		if attempt == p.Attempts {
			break
		}
		slog.Warn("attempt failed, retrying", "attempt", attempt, "wait", wait, "err", last)
		if err := sleep(ctx, wait); err != nil {
			return fmt.Errorf("gave up after %d attempt(s): %w (last error: %v)", attempt, err, last)
		}
		wait = min(wait*2, p.MaxWait)
	}
	return fmt.Errorf("gave up after %d attempt(s): %w", p.Attempts, last)
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// flakyServer answers 503 to the first n requests and 200 afterwards.
func flakyServer(n int32) *httptest.Server {
	var calls atomic.Int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) <= n {
			http.Error(w, "warming up", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ok, here is your data")
	}))
}

func run() error {
	srv := flakyServer(2)
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	policy := Policy{Attempts: 5, Wait: 200 * time.Millisecond, MaxWait: 2 * time.Second}

	var body []byte
	err := policy.Do(context.Background(), func() error {
		resp, err := client.Get(srv.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}
		body = make([]byte, 64)
		n, _ := resp.Body.Read(body)
		body = body[:n]
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("got: %s", body)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
