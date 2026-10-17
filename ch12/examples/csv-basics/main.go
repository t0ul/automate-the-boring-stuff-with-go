// read a CSV file with the standard library, then write one
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	data := "name,city,age\nAda,London,36\nLinus,Helsinki,54\n\"Grace, Rear Admiral\",Arlington,85\n"

	// csv.Reader handles quotes and commas inside quotes for you
	r := csv.NewReader(strings.NewReader(data))
	rows, err := r.ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	header, body := rows[0], rows[1:]
	fmt.Println("columns:", header)
	for _, row := range body {
		fmt.Printf("%s lives in %s\n", row[0], row[1])
	}

	// csv.Writer quotes whatever needs quoting
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"name", "note"})
	w.Write([]string{"Ada", "wrote the first program, in 1843"})
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
