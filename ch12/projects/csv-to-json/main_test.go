package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReadExpensesAnyColumnOrder(t *testing.T) {
	in := "amount,item,date,category\n4.5,coffee,2026-03-02,food\n8.25,\"Go book, used\",2026-03-03,books\n"
	got, err := ReadExpenses(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1].Item != "Go book, used" || got[0].Amount != 4.5 {
		t.Fatalf("got %+v", got)
	}
}

func TestReadExpensesBadAmount(t *testing.T) {
	_, err := ReadExpenses(strings.NewReader("date,category,item,amount\n2026-01-01,x,y,lots\n"))
	if err == nil || !strings.Contains(err.Error(), "row 2") {
		t.Fatalf("expected a row 2 error, got %v", err)
	}
}

func TestSummarizeJSON(t *testing.T) {
	rep := Summarize([]Expense{{"d", "food", "a", 1.5}, {"d", "food", "b", 2}, {"d", "books", "c", 10}})
	data, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	var back map[string]any
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if back["total"].(float64) != 13.5 || back["by_category"].(map[string]any)["food"].(float64) != 3.5 {
		t.Fatalf("json = %s", data)
	}
}
