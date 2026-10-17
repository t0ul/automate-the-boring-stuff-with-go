// watermark-pdfs stamps a word diagonally across every page of a PDF, the way "DRAFT"
// or "CONFIDENTIAL" appears on documents that leave the building.
//
//	go run ./ch10/projects/watermark-pdfs -text DRAFT ./ch10/projects/watermark-pdfs/contract.pdf
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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Stamp writes a copy of in to out with text across every page.
func Stamp(ctx context.Context, in, out, text string) error {
	desc := "font:Helvetica, points:56, rot:45, op:0.25, fillcolor:#b00020, scale:0.8 rel"
	wm, err := api.TextWatermark(ctx, text, desc, true, false, types.POINTS, nil)
	if err != nil {
		return fmt.Errorf("describe watermark: %w", err)
	}
	if err := api.AddWatermarksFile(ctx, in, out, nil, wm, nil); err != nil {
		return fmt.Errorf("stamp %s: %w", in, err)
	}
	return nil
}

// OutputName puts "_stamped" before the extension: contract.pdf -> contract_stamped.pdf.
func OutputName(in string) string {
	ext := filepath.Ext(in)
	return strings.TrimSuffix(in, ext) + "_stamped" + ext
}

func run(args []string) error {
	flags := flag.NewFlagSet("watermark-pdfs", flag.ContinueOnError)
	text := flags.String("text", "DRAFT", "the word to stamp")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() == 0 {
		return errors.New("usage: watermark-pdfs [-text DRAFT] <file.pdf ...>")
	}
	for _, in := range flags.Args() {
		out := OutputName(in)
		if err := Stamp(context.Background(), in, out, *text); err != nil {
			return err
		}
		fmt.Println("wrote", out)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
