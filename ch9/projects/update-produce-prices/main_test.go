package main

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestUpdatePrices(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	rows := [][]any{
		{"PRODUCE", "COST PER POUND", "POUNDS SOLD", "TOTAL"},
		{"Garlic", 1.19, 10, 11.9},
		{"Okra", 2.26, 4, 9.04},
		{"Garlic", 1.19, 2, 2.38},
	}
	for i, r := range rows {
		if err := f.SetSheetRow("Sheet1", excelCell(i+1), &r); err != nil {
			t.Fatal(err)
		}
	}
	n, err := UpdatePrices(f, "Sheet1", map[string]float64{"Garlic": 3.07})
	if err != nil || n != 2 {
		t.Fatalf("UpdatePrices = %d, %v", n, err)
	}
	if v, _ := f.GetCellValue("Sheet1", "D2"); v != "30.7" {
		t.Errorf("D2 = %q, want 30.7", v)
	}
	if v, _ := f.GetCellValue("Sheet1", "B3"); v != "2.26" {
		t.Errorf("Okra price changed: %q", v)
	}
}

func excelCell(row int) string { return "A" + itoa(row) }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func TestFixtureHasTargets(t *testing.T) {
	f, err := excelize.OpenFile("produceSales.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	n, err := UpdatePrices(f, "Sheet1", corrections)
	if err != nil || n == 0 {
		t.Fatalf("fixture produced %d updates, %v", n, err)
	}
}
