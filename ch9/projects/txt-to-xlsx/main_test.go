package main

import "testing"

func TestColumnsToWorkbook(t *testing.T) {
	f, err := ColumnsToWorkbook([]string{"days", "drinks"}, [][]string{{"Mon", "Tue", "Wed"}, {"tea"}})
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cases := map[string]string{"A1": "days", "A4": "Wed", "B1": "drinks", "B2": "tea", "B3": ""}
	for cell, want := range cases {
		if got, _ := f.GetCellValue("Sheet1", cell); got != want {
			t.Errorf("%s = %q, want %q", cell, got, want)
		}
	}
}

func TestReadLinesFixture(t *testing.T) {
	lines, err := ReadLines("first.txt")
	if err != nil || len(lines) != 5 || lines[0] != "Monday" {
		t.Fatalf("ReadLines = %v, %v", lines, err)
	}
}
