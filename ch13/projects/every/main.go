// every runs a command on a fixed interval and logs each run: start, exit code, how long
// it took. It is the smallest useful scheduler, for when cron feels like too much and
// "leave a terminal open" is acceptable. Everything after -- is the command.
//
//	go run ./ch13/projects/every -interval 1h -- go run ./ch11/projects/download-xkcd -n 1
//	go run ./ch13/projects/every -interval 2s -times 3 -- date
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

// RunOnce runs the command with a timeout and returns its exit code and duration. Output
// is streamed to out so a long job is visible while it runs.
func RunOnce(ctx context.Context, timeout time.Duration, out io.Writer, name string, args ...string) (code int, took time.Duration, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = out
	cmd.Stderr = out
	start := time.Now()
	err = cmd.Run()
	took = time.Since(start)
	var exitErr *exec.ExitError
	switch {
	case ctx.Err() != nil:
		return -1, took, fmt.Errorf("killed after %s: %w", timeout, ctx.Err()) // the deadline hit
	case err == nil:
		return 0, took, nil
	case errors.As(err, &exitErr):
		return exitErr.ExitCode(), took, nil // the command ran and failed; that is a result, not our error
	default:
		return -1, took, err // could not start at all
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("every", flag.ContinueOnError)
	interval := flags.Duration("interval", time.Hour, "time between runs, e.g. 30m, 2h")
	timeout := flags.Duration("timeout", 10*time.Minute, "kill a run that takes longer than this")
	times := flags.Int("times", 0, "stop after this many runs (0 = forever)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() == 0 {
		return errors.New("usage: every [-interval 1h] [-timeout 10m] [-times N] -- <command> [args]")
	}
	name, cmdArgs := flags.Arg(0), flags.Args()[1:]

	// ctrl+c cancels the context, which stops the loop and kills a running command
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for n := 1; ; n++ {
		fmt.Printf("[%s] run %d: %s %v\n", time.Now().Format(time.Kitchen), n, name, cmdArgs)
		code, took, err := RunOnce(ctx, *timeout, os.Stdout, name, cmdArgs...)
		if err != nil {
			fmt.Printf("[%s] run %d could not complete after %s: %v\n", time.Now().Format(time.Kitchen), n, took.Round(time.Millisecond), err)
		} else {
			fmt.Printf("[%s] run %d exit %d after %s\n", time.Now().Format(time.Kitchen), n, code, took.Round(time.Millisecond))
		}
		if *times > 0 && n >= *times {
			return nil
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			fmt.Println("stopped")
			return nil
		}
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
