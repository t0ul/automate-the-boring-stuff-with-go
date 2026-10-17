// multiplication-table writes an N by N multiplication table to a workbook, with the
// first row and column in bold.
//
//	go run ./ch9/projects/multiplication-table -n 12
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// Fill writes the table for 1..n into sheet, leaving A1 empty like a paper table.
func Fill(f *excelize.File, sheet string, n int) error {
	bold, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return err
	}
	for i := 1; i <= n; i++ {
		rowHeader, err := excelize.CoordinatesToCellName(1, i+1)
		if err != nil {
			return err
		}
		colHeader, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return err
		}
		for _, cell := range []string{rowHeader, colHeader} {
			if err := f.SetCellValue(sheet, cell, i); err != nil {
				return err
			}
			if err := f.SetCellStyle(sheet, cell, cell, bold); err != nil {
				return err
			}
		}
		for j := 1; j <= n; j++ {
			cell, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheet, cell, i*j); err != nil {
				return err
			}
		}
	}
	return nil
}

func run(args []string) error {
	fs := flag.NewFlagSet("multiplication-table", flag.ContinueOnError)
	n := fs.Int("n", 10, "table size")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *n < 1 || *n > 1000 {
		return fmt.Errorf("n must be between 1 and 1000, got %d", *n)
	}
	f := excelize.NewFile()
	defer f.Close()
	if err := Fill(f, "Sheet1", *n); err != nil {
		return err
	}
	out := fmt.Sprintf("multiplication_table_%d.xlsx", *n)
	if err := f.SaveAs(out); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Println("wrote", out)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
