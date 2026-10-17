// read-pdf-text prints the page count of a PDF and the plain text of every page.
// PDFs store drawing instructions, not paragraphs, so extracted text is only as tidy
// as the document that produced it.
//
//	go run ./ch10/projects/read-pdf-text ./ch10/projects/read-pdf-text/sample.pdf
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PageTexts returns the plain text of each page, in order.
func PageTexts(path string) ([]string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var pages []string
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			pages = append(pages, "")
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", i, err)
		}
		pages = append(pages, strings.TrimSpace(text))
	}
	return pages, nil
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: read-pdf-text <file.pdf>")
	}
	pages, err := PageTexts(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s: %d page(s)\n", args[0], len(pages))
	for i, text := range pages {
		fmt.Printf("\n--- page %d ---\n%s\n", i+1, text)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
