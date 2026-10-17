// find-duplicates lists files under a folder that have identical contents, by hashing
// every file with SHA-256. Hashing is the slow part, so a pool of workers does it while
// one goroutine walks the tree and one collects the answers.
//
//	go run ./ch14/projects/find-duplicates ~/Pictures
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"
)

// hashed is one file's result.
type hashed struct {
	path string
	sum  string
	err  error
}

// HashFile returns the hex SHA-256 of a file's contents, streaming so size does not matter.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// FindDuplicates returns groups of paths with identical contents (groups of size >= 2),
// each group sorted, groups sorted by their first path. minSize skips tiny files.
func FindDuplicates(root string, minSize int64, workers int) ([][]string, error) {
	paths := make(chan string)
	results := make(chan hashed)

	// the walker: finds files and sends their paths
	var walkErr error
	go func() {
		defer close(paths)
		walkErr = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // unreadable: skip, do not stop
			}
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil || info.Size() < minSize || !info.Mode().IsRegular() {
				return nil
			}
			paths <- p
			return nil
		})
	}()

	// the workers: hash whatever the walker sends
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for p := range paths {
				sum, err := HashFile(p)
				results <- hashed{p, sum, err}
			}
		})
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	// the collector: this goroutine, grouping by hash
	bySum := map[string][]string{}
	for r := range results {
		if r.err != nil {
			continue // a file that vanished or cannot be read is not a duplicate of anything
		}
		bySum[r.sum] = append(bySum[r.sum], r.path)
	}
	if walkErr != nil {
		return nil, walkErr
	}

	var groups [][]string
	for _, g := range bySum {
		if len(g) >= 2 {
			slices.Sort(g)
			groups = append(groups, g)
		}
	}
	slices.SortFunc(groups, func(a, b []string) int { return slices.Compare(a, b) })
	return groups, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("find-duplicates", flag.ContinueOnError)
	minSize := flags.Int64("min", 1, "ignore files smaller than this many bytes")
	workers := flags.Int("workers", runtime.NumCPU(), "how many files to hash at once")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: find-duplicates [-min bytes] [-workers N] <folder>")
	}
	groups, err := FindDuplicates(flags.Arg(0), *minSize, *workers)
	if err != nil {
		return err
	}
	for i, g := range groups {
		fmt.Printf("group %d (%d copies):\n", i+1, len(g))
		for _, p := range g {
			fmt.Println("  ", p)
		}
	}
	fmt.Fprintf(os.Stderr, "%d group(s) of duplicates\n", len(groups))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
