// read-workbook opens an Excel file, lists its sheets, prints the first sheet as a
// table and shows how to read one cell by its name.
//
//	go run ./ch9/projects/read-workbook ./ch9/projects/read-workbook/example.xlsx
//	go run ./ch9/projects/read-workbook -cell C3 ./ch9/projects/read-workbook/example.xlsx
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/xuri/excelize/v2"
)

// Table returns the rows of sheet, padded so every row has as many cells as the widest.
func Table(f *excelize.File, sheet string) ([][]string, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("rows of %s: %w", sheet, err)
	}
	width := 0
	for _, r := range rows {
		width = max(width, len(r))
	}
	for i := range rows {
		for len(rows[i]) < width {
			rows[i] = append(rows[i], "")
		}
	}
	return rows, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("read-workbook", flag.ContinueOnError)
	cell := flags.String("cell", "B2", "one cell to print by name")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: read-workbook [-cell B2] <file.xlsx>")
	}
	f, err := excelize.OpenFile(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	fmt.Println("sheets:", strings.Join(sheets, ", "))
	first := sheets[0]

	rows, err := Table(f, first)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintln(tw, strings.Join(r, "\t"))
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	v, err := f.GetCellValue(first, *cell)
	if err != nil {
		return fmt.Errorf("cell %s: %w", *cell, err)
	}
	fmt.Printf("\n%s!%s = %q\n", first, *cell, v)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
