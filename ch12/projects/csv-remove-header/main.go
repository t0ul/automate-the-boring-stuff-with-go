// csv-remove-header strips the first row from every CSV file in a folder and writes the
// result to a headerRemoved/ sub-folder, leaving the originals alone. It is the job you
// get when a tool insists on a header row and the file you were given has one you do not
// want merged into the data.
//
//	go run ./ch12/projects/csv-remove-header ./ch12/projects/csv-remove-header/sample
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StripHeader reads src as CSV and writes every row but the first to dst. It returns the
// number of rows written.
func StripHeader(src, dst string) (int, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	rows, err := csv.NewReader(in).ReadAll()
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", src, err)
	}
	if len(rows) == 0 {
		return 0, nil
	}
	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	w := csv.NewWriter(out)
	if err := w.WriteAll(rows[1:]); err != nil {
		out.Close()
		return 0, fmt.Errorf("write %s: %w", dst, err)
	}
	return len(rows) - 1, out.Close()
}

// ProcessFolder strips every *.csv directly inside folder into folder/headerRemoved.
func ProcessFolder(folder string) (map[string]int, error) {
	outDir := filepath.Join(folder, "headerRemoved")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	matches, err := filepath.Glob(filepath.Join(folder, "*.csv"))
	if err != nil {
		return nil, err
	}
	done := map[string]int{}
	for _, src := range matches {
		name := filepath.Base(src)
		if !strings.HasSuffix(strings.ToLower(name), ".csv") {
			continue
		}
		n, err := StripHeader(src, filepath.Join(outDir, name))
		if err != nil {
			return done, err
		}
		done[name] = n
	}
	return done, nil
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: csv-remove-header <folder>")
	}
	done, err := ProcessFolder(args[0])
	if err != nil {
		return err
	}
	for name, n := range done {
		fmt.Printf("%s: %d data row(s) kept\n", name, n)
	}
	fmt.Printf("wrote %d file(s) to %s\n", len(done), filepath.Join(args[0], "headerRemoved"))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
