// a database in one file, no server to install: create, insert, query, update, delete
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // registers the "sqlite" driver; pure Go, no C compiler needed
)

type Book struct {
	ID    int
	Title string
	Year  int
}

func run() error {
	dir, err := os.MkdirTemp("", "sqlite-basics")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	db, err := sql.Open("sqlite", filepath.Join(dir, "library.db"))
	if err != nil {
		return err
	}
	defer db.Close()

	// the schema: run it every start; IF NOT EXISTS makes it harmless the second time
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
		id    INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		year  INTEGER NOT NULL
	)`)
	if err != nil {
		return err
	}

	// ? placeholders: the driver passes values separately, so a title with a quote in it
	// is data, never SQL
	for _, b := range []Book{{Title: "The Go Programming Language", Year: 2015}, {Title: "Automate the Boring Stuff with Go", Year: 2026}, {Title: "Let's Go", Year: 2019}} {
		if _, err := db.Exec(`INSERT INTO books (title, year) VALUES (?, ?)`, b.Title, b.Year); err != nil {
			return err
		}
	}

	rows, err := db.Query(`SELECT id, title, year FROM books WHERE year >= ? ORDER BY year`, 2015)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Year); err != nil {
			return err
		}
		fmt.Printf("%d  %-36s %d\n", b.ID, b.Title, b.Year)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// one value: QueryRow
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM books`).Scan(&count); err != nil {
		return err
	}
	fmt.Println("books:", count)

	res, err := db.Exec(`UPDATE books SET year = ? WHERE title = ?`, 2016, "The Go Programming Language")
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	fmt.Println("updated rows:", n)

	if _, err := db.Exec(`DELETE FROM books WHERE year < ?`, 2019); err != nil {
		return err
	}
	db.QueryRow(`SELECT count(*) FROM books`).Scan(&count)
	fmt.Println("books after delete:", count)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
