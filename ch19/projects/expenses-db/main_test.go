package main

import (
	"path/filepath"
	"strings"
	"testing"
)

const sample = "date,category,item,amount\n2026-03-02,food,coffee,4.50\n2026-03-03,books,\"Go book, used\",8.25\n2026-04-01,food,lunch,12.00\n"

func TestImportIsIdempotent(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "e.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	n, err := Import(db, strings.NewReader(sample))
	if err != nil || n != 3 {
		t.Fatalf("first import = %d, %v", n, err)
	}
	n, err = Import(db, strings.NewReader(sample))
	if err != nil || n != 0 {
		t.Fatalf("second import = %d, %v; want 0 new rows", n, err)
	}
	var total float64
	if err := db.QueryRow(`SELECT sum(amount) FROM expenses`).Scan(&total); err != nil || total != 24.75 {
		t.Fatalf("total = %v, %v", total, err)
	}
}

func TestReportAndQuery(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "e.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := Import(db, strings.NewReader(sample)); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Report(db, &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"food", "16.50", "2026-03", "12.75", "2026-04", "12.00"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("report lacks %q:\n%s", want, out.String())
		}
	}
	out.Reset()
	if err := Query(db, &out, `SELECT item FROM expenses WHERE amount > ? ORDER BY item`, 5); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.Contains(got, "Go book, used") || strings.Contains(got, "coffee") {
		t.Fatalf("query output:\n%s", got)
	}
}
