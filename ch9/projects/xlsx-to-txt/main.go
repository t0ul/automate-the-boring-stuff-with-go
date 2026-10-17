// xlsx-to-txt is the reverse of txt-to-xlsx: every column of the first sheet becomes a
// text file named after the column letter, one cell per line.
//
//	go run ./ch9/projects/xlsx-to-txt ./ch9/projects/xlsx-to-txt/table.xlsx
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Columns reads the first sheet column by column, trimming trailing empty cells.
func Columns(f *excelize.File) ([][]string, error) {
	cols, err := f.GetCols(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}
	for i := range cols {
		for len(cols[i]) > 0 && cols[i][len(cols[i])-1] == "" {
			cols[i] = cols[i][:len(cols[i])-1]
		}
	}
	return cols, nil
}

// WriteColumns writes each column to dir/column_<letter>.txt and returns the paths.
func WriteColumns(dir string, cols [][]string) ([]string, error) {
	var written []string
	for i, col := range cols {
		letter, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return written, err
		}
		path := filepath.Join(dir, "column_"+letter+".txt")
		if err := os.WriteFile(path, []byte(strings.Join(col, "\n")+"\n"), 0o644); err != nil {
			return written, err
		}
		written = append(written, path)
	}
	return written, nil
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: xlsx-to-txt <file.xlsx>")
	}
	f, err := excelize.OpenFile(args[0])
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	cols, err := Columns(f)
	if err != nil {
		return err
	}
	paths, err := WriteColumns(".", cols)
	if err != nil {
		return err
	}
	for _, p := range paths {
		fmt.Println("wrote", p)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
