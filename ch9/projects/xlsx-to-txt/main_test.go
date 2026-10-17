package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestColumnsAndWrite(t *testing.T) {
	f, err := excelize.OpenFile("table.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cols, err := Columns(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(cols) != 3 || cols[0][0] != "Fruit" || cols[2][4] != "summer" {
		t.Fatalf("columns = %v", cols)
	}
	dir := t.TempDir()
	paths, err := WriteColumns(dir, cols)
	if err != nil || len(paths) != 3 {
		t.Fatalf("WriteColumns = %v, %v", paths, err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "column_B.txt"))
	if err != nil || string(b) != "Colour\nred\ngreen\npurple\nbrown\n" {
		t.Fatalf("column_B.txt = %q, %v", b, err)
	}
}
