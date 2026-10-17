// a transaction makes many writes one unit: all of them land, or none, and it is faster
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func run() error {
	dir, err := os.MkdirTemp("", "transactions")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	db, err := sql.Open("sqlite", filepath.Join(dir, "t.db"))
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE readings (sensor TEXT, value REAL)`); err != nil {
		return err
	}

	// one statement at a time: every Exec is its own transaction and its own disk sync
	start := time.Now()
	for i := range 500 {
		if _, err := db.Exec(`INSERT INTO readings VALUES (?, ?)`, "a", float64(i)); err != nil {
			return err
		}
	}
	fmt.Println("500 separate inserts:", time.Since(start).Round(time.Millisecond))

	// the same work inside one transaction
	start = time.Now()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO readings VALUES (?, ?)`)
	if err != nil {
		return err
	}
	for i := range 500 {
		if _, err := stmt.Exec("b", float64(i)); err != nil {
			tx.Rollback()
			return err
		}
	}
	stmt.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Println("500 inserts in one transaction:", time.Since(start).Round(time.Millisecond))

	// all or nothing: a failure in the middle rolls everything back
	err = func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback() // a no-op after Commit; the safety net otherwise
		if _, err := tx.Exec(`INSERT INTO readings VALUES ('c', 1)`); err != nil {
			return err
		}
		if true {
			return errors.New("sensor c reported nonsense, abandoning the batch")
		}
		return tx.Commit()
	}()
	fmt.Println("batch c:", err)
	var n int
	db.QueryRow(`SELECT count(*) FROM readings WHERE sensor = 'c'`).Scan(&n)
	fmt.Println("rows from sensor c:", n)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
