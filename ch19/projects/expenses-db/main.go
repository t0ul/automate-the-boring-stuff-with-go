// expenses-db keeps the expenses from chapter 12 in a SQLite file and answers questions
// about them with SQL. Import a CSV as often as you like (duplicates are skipped),
// print a report, or run any query.
//
//	go run ./ch19/projects/expenses-db import ./ch12/projects/csv-to-json/expenses.csv
//	go run ./ch19/projects/expenses-db report
//	go run ./ch19/projects/expenses-db query "SELECT item, amount FROM expenses ORDER BY amount DESC LIMIT 3"
package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS expenses (
	id       INTEGER PRIMARY KEY,
	date     TEXT NOT NULL,
	category TEXT NOT NULL,
	item     TEXT NOT NULL,
	amount   REAL NOT NULL,
	UNIQUE (date, category, item, amount)
);`

// Open opens (or creates) the database and makes sure the schema exists.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// Import loads CSV rows with the columns date, category, item, amount (any order) and
// returns how many were new. Rows already present are skipped, so re-importing is safe.
func Import(db *sql.DB, r io.Reader) (int, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, errors.New("empty file")
	}
	col := map[string]int{}
	for i, name := range rows[0] {
		col[strings.ToLower(strings.TrimSpace(name))] = i
	}
	for _, need := range []string{"date", "category", "item", "amount"} {
		if _, ok := col[need]; !ok {
			return 0, fmt.Errorf("missing column %q", need)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO expenses (date, category, item, amount) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	added := 0
	for n, row := range rows[1:] {
		amount, err := strconv.ParseFloat(strings.TrimSpace(row[col["amount"]]), 64)
		if err != nil {
			return 0, fmt.Errorf("row %d: amount %q: %w", n+2, row[col["amount"]], err)
		}
		res, err := stmt.Exec(row[col["date"]], row[col["category"]], row[col["item"]], amount)
		if err != nil {
			return 0, err
		}
		if k, _ := res.RowsAffected(); k > 0 {
			added++
		}
	}
	return added, tx.Commit()
}

// Query runs any SELECT and prints the result as a table, whatever its columns are.
func Query(db *sql.DB, w io.Writer, query string, args ...any) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(cols, "\t"))
	values := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		cells := make([]string, len(cols))
		for i, v := range values {
			switch x := v.(type) {
			case []byte:
				cells[i] = string(x)
			case float64:
				cells[i] = strconv.FormatFloat(x, 'f', 2, 64)
			default:
				cells[i] = fmt.Sprint(x)
			}
		}
		fmt.Fprintln(tw, strings.Join(cells, "\t"))
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return tw.Flush()
}

// Report prints totals by category and by month.
func Report(db *sql.DB, w io.Writer) error {
	fmt.Fprintln(w, "by category")
	if err := Query(db, w, `SELECT category, count(*) AS items, round(sum(amount), 2) AS total FROM expenses GROUP BY category ORDER BY total DESC`); err != nil {
		return err
	}
	fmt.Fprintln(w, "\nby month")
	return Query(db, w, `SELECT substr(date, 1, 7) AS month, round(sum(amount), 2) AS total FROM expenses GROUP BY month ORDER BY month`)
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: expenses-db [-db file] import <csv> | report | query <sql>")
	}
	dbPath := "expenses.db"
	if args[0] == "-db" && len(args) > 2 {
		dbPath, args = args[1], args[2:]
	}
	db, err := Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	switch args[0] {
	case "import":
		if len(args) != 2 {
			return errors.New("import needs a CSV file")
		}
		f, err := os.Open(args[1])
		if err != nil {
			return err
		}
		defer f.Close()
		n, err := Import(db, f)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "imported %d new row(s) into %s\n", n, dbPath)
		return nil
	case "report":
		return Report(db, out)
	case "query":
		if len(args) != 2 {
			return errors.New("query needs one SQL string")
		}
		return Query(db, out, args[1])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
