// csv-to-json turns an expenses CSV into typed JSON: each row becomes a struct, the
// numbers become numbers, and a summary by category comes along for free.
//
//	go run ./ch12/projects/csv-to-json ./ch12/projects/csv-to-json/expenses.csv
package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Expense is one row. The tags choose the JSON names.
type Expense struct {
	Date     string  `json:"date"`
	Category string  `json:"category"`
	Item     string  `json:"item"`
	Amount   float64 `json:"amount"`
}

// Report is what we write out.
type Report struct {
	Expenses   []Expense          `json:"expenses"`
	ByCategory map[string]float64 `json:"by_category"`
	Total      float64            `json:"total"`
}

// ReadExpenses parses CSV with the header date,category,item,amount in any column order.
func ReadExpenses(r io.Reader) ([]Expense, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("empty file")
	}
	col := map[string]int{}
	for i, name := range rows[0] {
		col[strings.ToLower(strings.TrimSpace(name))] = i
	}
	for _, need := range []string{"date", "category", "item", "amount"} {
		if _, ok := col[need]; !ok {
			return nil, fmt.Errorf("missing column %q", need)
		}
	}
	var out []Expense
	for n, row := range rows[1:] {
		amount, err := strconv.ParseFloat(strings.TrimSpace(row[col["amount"]]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: amount %q: %w", n+2, row[col["amount"]], err)
		}
		out = append(out, Expense{
			Date:     row[col["date"]],
			Category: row[col["category"]],
			Item:     row[col["item"]],
			Amount:   amount,
		})
	}
	return out, nil
}

// Summarize builds the report from the rows.
func Summarize(expenses []Expense) Report {
	rep := Report{Expenses: expenses, ByCategory: map[string]float64{}}
	for _, e := range expenses {
		rep.ByCategory[e.Category] += e.Amount
		rep.Total += e.Amount
	}
	return rep
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: csv-to-json <expenses.csv>")
	}
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()
	expenses, err := ReadExpenses(f)
	if err != nil {
		return fmt.Errorf("%s: %w", args[0], err)
	}
	out := strings.TrimSuffix(args[0], ".csv") + ".json"
	data, err := json.MarshalIndent(Summarize(expenses), "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return err
	}
	fmt.Printf("%d expense(s), total %.2f, wrote %s\n", len(expenses), Summarize(expenses).Total, out)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
