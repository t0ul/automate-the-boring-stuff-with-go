// census-by-county reads a census workbook (one row per census tract) and totals the
// population and the tract count per county, then writes the result as JSON so other
// programs can use it.
//
//	go run ./ch9/projects/census-by-county ./ch9/projects/census-by-county/census.xlsx
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Tract is one row of the workbook.
type Tract struct {
	State  string
	County string
	Pop    int
}

// CountyStats is what we know about one county after reading every tract.
type CountyStats struct {
	Tracts int `json:"tracts"`
	Pop    int `json:"pop"`
}

// ReadTracts loads the first sheet. It expects the columns CensusTract, State, County,
// POP2010 in that order, with a header row.
func ReadTracts(path string) ([]Tract, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}
	var tracts []Tract
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}
		pop, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, fmt.Errorf("row %d: population %q: %w", i+1, row[3], err)
		}
		tracts = append(tracts, Tract{State: row[1], County: row[2], Pop: pop})
	}
	return tracts, nil
}

// Aggregate returns state -> county -> stats.
func Aggregate(tracts []Tract) map[string]map[string]CountyStats {
	out := map[string]map[string]CountyStats{}
	for _, t := range tracts {
		if out[t.State] == nil {
			out[t.State] = map[string]CountyStats{}
		}
		s := out[t.State][t.County]
		s.Tracts++
		s.Pop += t.Pop
		out[t.State][t.County] = s
	}
	return out
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: census-by-county <census.xlsx>")
	}
	tracts, err := ReadTracts(args[0])
	if err != nil {
		return err
	}
	stats := Aggregate(tracts)

	for _, state := range slices.Sorted(keys(stats)) {
		total := 0
		for _, c := range stats[state] {
			total += c.Pop
		}
		fmt.Printf("%s: %d people in %d counties\n", state, total, len(stats[state]))
		for _, county := range slices.Sorted(keys(stats[state])) {
			c := stats[state][county]
			fmt.Printf("  %-12s %6d people, %d tracts\n", county, c.Pop, c.Tracts)
		}
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("census_summary.json", data, 0o644); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}
	fmt.Println("wrote census_summary.json")
	return nil
}

func keys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
