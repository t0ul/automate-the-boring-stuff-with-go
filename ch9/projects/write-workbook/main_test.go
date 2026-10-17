package main

import (
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestBuildInventoryRoundTrip(t *testing.T) {
	f, err := BuildInventory([]Item{{"Pen", 3, 1.5}, {"Pad", 2, 4}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "inv.xlsx")
	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}
	f.Close()

	g, err := excelize.OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if got := g.GetSheetList(); len(got) != 1 || got[0] != "Inventory" {
		t.Fatalf("sheets = %v", got)
	}
	if v, _ := g.GetCellValue("Inventory", "A3"); v != "Pad" {
		t.Errorf("A3 = %q", v)
	}
	if formula, _ := g.GetCellFormula("Inventory", "D4"); formula != "SUM(D2:D3)" {
		t.Errorf("D4 formula = %q", formula)
	}
}
