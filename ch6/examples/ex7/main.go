// ex. 7 a slice of structs is a table
package main

import (
	"fmt"
	"slices"
)

type Expense struct {
	Day    string
	Item   string
	Amount float64
}

func main() {
	expenses := []Expense{
		{"Mon", "coffee", 4.50},
		{"Mon", "lunch", 12.00},
		{"Tue", "coffee", 4.50},
		{"Tue", "book", 28.99},
	}

	total := 0.0
	byDay := map[string]float64{}
	for _, e := range expenses {
		total += e.Amount
		byDay[e.Day] += e.Amount
	}
	fmt.Printf("total %.2f\n", total)
	for _, day := range []string{"Mon", "Tue"} {
		fmt.Printf("%s %.2f\n", day, byDay[day])
	}

	// sort the table by amount, biggest first
	slices.SortFunc(expenses, func(a, b Expense) int {
		switch {
		case a.Amount > b.Amount:
			return -1
		case a.Amount < b.Amount:
			return 1
		}
		return 0
	})
	fmt.Println(expenses[0])
}
