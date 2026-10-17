package main

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestFill(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	if err := Fill(f, "Sheet1", 6); err != nil {
		t.Fatal(err)
	}
	checks := map[string]string{"A1": "", "B1": "1", "A7": "6", "G7": "36", "D3": "6"}
	for cell, want := range checks {
		if got, _ := f.GetCellValue("Sheet1", cell); got != want {
			t.Errorf("%s = %q, want %q", cell, got, want)
		}
	}
}
