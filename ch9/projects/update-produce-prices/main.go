// update-produce-prices fixes wrong unit prices in a sales workbook and recomputes the
// TOTAL column for every row it touched. The corrections live in one map, so the next
// time prices change you edit three lines instead of three hundred cells.
//
//	go run ./ch9/projects/update-produce-prices ./ch9/projects/update-produce-prices/produceSales.xlsx
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

var corrections = map[string]float64{
	"Garlic": 3.07,
	"Celery": 1.19,
	"Lemon":  1.27,
}

// UpdatePrices rewrites COST PER POUND (column B) and TOTAL (column D) for every row
// whose PRODUCE (column A) has a corrected price. It returns how many rows changed.
func UpdatePrices(f *excelize.File, sheet string, prices map[string]float64) (int, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return 0, err
	}
	changed := 0
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		price, ok := prices[row[0]]
		if !ok {
			continue
		}
		pounds, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return changed, fmt.Errorf("row %d: pounds %q: %w", i+1, row[2], err)
		}
		n := i + 1
		if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", n), price); err != nil {
			return changed, err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", n), round2(price*pounds)); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func round2(x float64) float64 { return float64(int(x*100+0.5)) / 100 }

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: update-produce-prices <produceSales.xlsx>")
	}
	f, err := excelize.OpenFile(args[0])
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	sheet := f.GetSheetList()[0]
	n, err := UpdatePrices(f, sheet, corrections)
	if err != nil {
		return err
	}
	out := "updatedProduceSales.xlsx"
	if err := f.SaveAs(out); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Printf("updated %d row(s), wrote %s\n", n, out)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
