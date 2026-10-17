package main

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestTableOnFixture(t *testing.T) {
	f, err := excelize.OpenFile("example.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := Table(f, "Sheet1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 8 || rows[0][1] != "Item" || rows[1][2] != "73" {
		t.Fatalf("unexpected table: %v", rows)
	}
}

func TestTablePadsRaggedRows(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	_ = f.SetCellValue("Sheet1", "A1", "a")
	_ = f.SetCellValue("Sheet1", "C2", "c")
	rows, err := Table(f, "Sheet1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || len(rows[0]) != 3 || len(rows[1]) != 3 {
		t.Fatalf("rows not padded: %v", rows)
	}
}
