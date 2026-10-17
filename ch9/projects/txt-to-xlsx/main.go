// txt-to-xlsx takes text files and writes one workbook where each file is a column and
// each line is a cell.
//
//	go run ./ch9/projects/txt-to-xlsx ch9/projects/txt-to-xlsx/first.txt ch9/projects/txt-to-xlsx/second.txt ch9/projects/txt-to-xlsx/third.txt
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadLines returns the lines of a file without their line endings.
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

// ColumnsToWorkbook puts each slice of lines into its own column of a new workbook, with
// the given names as a header row.
func ColumnsToWorkbook(names []string, columns [][]string) (*excelize.File, error) {
	f := excelize.NewFile()
	for c := range columns {
		col, err := excelize.ColumnNumberToName(c + 1)
		if err != nil {
			return nil, err
		}
		values := append([]any{names[c]}, toAny(columns[c])...)
		if err := f.SetSheetCol("Sheet1", col+"1", &values); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func toAny(s []string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: txt-to-xlsx <file.txt> [more.txt ...]")
	}
	var names []string
	var columns [][]string
	for _, path := range args {
		lines, err := ReadLines(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		names = append(names, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		columns = append(columns, lines)
	}
	f, err := ColumnsToWorkbook(names, columns)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.SaveAs("combined.xlsx"); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Printf("wrote combined.xlsx with %d column(s)\n", len(columns))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
