// write-workbook creates an Excel file from scratch: a named sheet, a bold header,
// a few rows, sensible column widths and a SUM formula.
//
//	go run ./ch9/projects/write-workbook
package main

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// Item is one line of inventory.
type Item struct {
	Name     string
	Quantity int
	Price    float64
}

// BuildInventory fills a new workbook with the items and returns it unsaved.
func BuildInventory(items []Item) (*excelize.File, error) {
	f := excelize.NewFile()
	const sheet = "Inventory"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return nil, err
	}
	header := []any{"Item", "Quantity", "Unit price", "Value"}
	if err := f.SetSheetRow(sheet, "A1", &header); err != nil {
		return nil, err
	}
	bold, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return nil, err
	}
	if err := f.SetCellStyle(sheet, "A1", "D1", bold); err != nil {
		return nil, err
	}
	for i, it := range items {
		row := i + 2
		values := []any{it.Name, it.Quantity, it.Price}
		if err := f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &values); err != nil {
			return nil, err
		}
		if err := f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("B%d*C%d", row, row)); err != nil {
			return nil, err
		}
	}
	last := len(items) + 1
	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", last+1), "Total"); err != nil {
		return nil, err
	}
	if err := f.SetCellFormula(sheet, fmt.Sprintf("D%d", last+1), fmt.Sprintf("SUM(D2:D%d)", last)); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(sheet, "A", "A", 18); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(sheet, "B", "D", 12); err != nil {
		return nil, err
	}
	return f, nil
}

func run() error {
	items := []Item{
		{"Keyboard", 12, 39.50},
		{"Mouse", 30, 14.99},
		{"Monitor", 5, 219.00},
		{"USB-C cable", 80, 6.25},
	}
	f, err := BuildInventory(items)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.SaveAs("inventory.xlsx"); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Println("wrote inventory.xlsx")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
