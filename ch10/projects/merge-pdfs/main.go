// merge-pdfs joins PDF files into one, in the order given, and reports the page count
// of the result. Pass a folder instead of files to merge every PDF in it, sorted by name.
//
//	go run ./ch10/projects/merge-pdfs -out merged.pdf ./ch10/projects/merge-pdfs
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Inputs expands the arguments: a folder becomes its *.pdf files sorted by name; a file
// is taken as given. The output file is never included even if it sits in the folder.
func Inputs(args []string, out string) ([]string, error) {
	var files []string
	outAbs, _ := filepath.Abs(out)
	for _, a := range args {
		info, err := os.Stat(a)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			files = append(files, a)
			continue
		}
		matches, err := filepath.Glob(filepath.Join(a, "*.pdf"))
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			if abs, _ := filepath.Abs(m); abs == outAbs {
				continue
			}
			files = append(files, m)
		}
	}
	if len(files) < 2 {
		return nil, fmt.Errorf("need at least two PDFs, found %d", len(files))
	}
	return files, nil
}

// Merge writes the concatenation of in to out and returns the page count of out.
func Merge(ctx context.Context, in []string, out string) (int, error) {
	if err := api.MergeCreateFile(ctx, in, out, false, nil); err != nil {
		return 0, fmt.Errorf("merge: %w", err)
	}
	n, err := api.PageCountFile(ctx, out)
	if err != nil {
		return 0, fmt.Errorf("count pages of %s: %w", out, err)
	}
	return n, nil
}

func run(args []string) error {
	fs := flag.NewFlagSet("merge-pdfs", flag.ContinueOnError)
	out := fs.String("out", "merged.pdf", "output file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		return errors.New("usage: merge-pdfs [-out merged.pdf] <folder or file.pdf ...>")
	}
	in, err := Inputs(fs.Args(), *out)
	if err != nil {
		return err
	}
	n, err := Merge(context.Background(), in, *out)
	if err != nil {
		return err
	}
	fmt.Printf("merged %s\ninto %s (%d pages)\n", strings.Join(in, ", "), *out, n)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
