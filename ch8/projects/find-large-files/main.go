// find-large-files walks a folder and lists every file at or above a size, largest
// first. Sizes read the way people write them: 500KB, 100MB, 2GB.
//
//	go run ./ch8/projects/find-large-files -min 100MB ~/Downloads
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Entry is one file that met the threshold.
type Entry struct {
	Path string
	Size int64
}

var units = []struct {
	suffix string
	mult   int64
}{
	{"GB", 1 << 30}, {"MB", 1 << 20}, {"KB", 1 << 10}, {"B", 1},
}

// ParseSize turns "100MB", "2gb", "512" (bytes) into a byte count.
func ParseSize(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	for _, u := range units {
		if num, ok := strings.CutSuffix(s, u.suffix); ok {
			n, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
			if err != nil || n < 0 {
				return 0, fmt.Errorf("bad size %q", s)
			}
			return int64(n * float64(u.mult)), nil
		}
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("bad size %q", s)
	}
	return n, nil
}

// HumanSize formats bytes with one decimal and the largest fitting unit.
func HumanSize(n int64) string {
	for _, u := range units {
		if n >= u.mult && u.mult > 1 {
			return fmt.Sprintf("%.1f %s", float64(n)/float64(u.mult), u.suffix)
		}
	}
	return fmt.Sprintf("%d B", n)
}

// FindLarge returns the files under root whose size is at least min, largest first.
// Folders it cannot read are skipped, not fatal.
func FindLarge(root string, minSize int64) ([]Entry, error) {
	var found []Entry
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Size() >= minSize {
			found = append(found, Entry{Path: path, Size: info.Size()})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", root, err)
	}
	slices.SortFunc(found, func(a, b Entry) int { return int(b.Size - a.Size) })
	return found, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("find-large-files", flag.ContinueOnError)
	minFlag := flags.String("min", "100MB", "smallest size to report")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: find-large-files [-min 100MB] <folder>")
	}
	minSize, err := ParseSize(*minFlag)
	if err != nil {
		return err
	}
	entries, err := FindLarge(flags.Arg(0), minSize)
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Printf("%10s  %s\n", HumanSize(e.Size), e.Path)
	}
	fmt.Fprintf(os.Stderr, "%d file(s) at or above %s\n", len(entries), HumanSize(minSize))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
