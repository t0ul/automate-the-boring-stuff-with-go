// log/slog: log lines with fields a program (or a grep) can read later
package main

import (
	"errors"
	"log/slog"
	"os"
	"time"
)

func main() {
	// the default logger writes readable text to stderr
	slog.Info("backup started", "folder", "/Users/you/Documents", "files", 1284)
	slog.Warn("skipping unreadable file", "path", "/Users/you/Documents/locked.xlsx")

	// a JSON logger to a file is what you want for a scheduled job
	f, err := os.CreateTemp("", "run-*.log")
	if err != nil {
		slog.Error("cannot open log", "err", err)
		os.Exit(1)
	}
	defer os.Remove(f.Name())
	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("job", "nightly-backup", "run_id", time.Now().Format("20060102-150405"))

	logger.Debug("walking tree")
	logger.Info("done", "took", 4200*time.Millisecond, "bytes", 93_123_440)
	logger.Error("upload failed", "err", errors.New("connection reset"), "attempt", 3)

	f.Close()
	data, _ := os.ReadFile(f.Name())
	os.Stdout.Write(data)
}
