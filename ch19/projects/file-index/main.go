// file-index records every file under a folder (path, size, modified time, SHA-256) in a
// SQLite database, one run at a time, and then answers the questions that take a long
// time to answer by walking: what changed since last time, what is duplicated, what is
// biggest.
//
//	go run ./ch19/projects/file-index index ~/Documents
//	go run ./ch19/projects/file-index changes
//	go run ./ch19/projects/file-index dupes
//	go run ./ch19/projects/file-index largest -n 10
package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS runs (
	id      INTEGER PRIMARY KEY,
	root    TEXT NOT NULL,
	started TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS files (
	run_id INTEGER NOT NULL REFERENCES runs(id),
	path   TEXT NOT NULL,
	size   INTEGER NOT NULL,
	mtime  TEXT NOT NULL,
	sha256 TEXT NOT NULL,
	PRIMARY KEY (run_id, path)
);
CREATE INDEX IF NOT EXISTS files_sha ON files (run_id, sha256);`

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

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Index walks root and records every regular file as a new run. It returns the run id
// and the number of files recorded.
func Index(db *sql.DB, root string, now time.Time) (runID int64, count int, err error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return 0, 0, err
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	res, err := tx.Exec(`INSERT INTO runs (root, started) VALUES (?, ?)`, abs, now.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, 0, err
	}
	runID, _ = res.LastInsertId()
	stmt, err := tx.Prepare(`INSERT INTO files (run_id, path, size, mtime, sha256) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()
	err = filepath.WalkDir(abs, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		sum, err := hashFile(p)
		if err != nil {
			return nil // vanished or unreadable: not this run's problem
		}
		rel, _ := filepath.Rel(abs, p)
		if _, err := stmt.Exec(runID, rel, info.Size(), info.ModTime().UTC().Format(time.RFC3339Nano), sum); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return runID, count, tx.Commit()
}

// lastTwoRuns returns the ids of the latest and the previous run for the same root.
func lastTwoRuns(db *sql.DB) (latest, previous int64, err error) {
	rows, err := db.Query(`SELECT id FROM runs WHERE root = (SELECT root FROM runs ORDER BY id DESC LIMIT 1) ORDER BY id DESC LIMIT 2`)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, 0, err
		}
		ids = append(ids, id)
	}
	if len(ids) < 2 {
		return 0, 0, errors.New("need two runs of the same folder to compare; run index again later")
	}
	return ids[0], ids[1], nil
}

// Changes prints files added, removed or modified between the last two runs.
func Changes(db *sql.DB, w io.Writer) error {
	latest, previous, err := lastTwoRuns(db)
	if err != nil {
		return err
	}
	queries := []struct{ label, sql string }{
		{"added", `SELECT path FROM files WHERE run_id = ?1 AND path NOT IN (SELECT path FROM files WHERE run_id = ?2) ORDER BY path`},
		{"removed", `SELECT path FROM files WHERE run_id = ?2 AND path NOT IN (SELECT path FROM files WHERE run_id = ?1) ORDER BY path`},
		{"modified", `SELECT a.path FROM files a JOIN files b ON a.path = b.path AND b.run_id = ?2 WHERE a.run_id = ?1 AND a.sha256 <> b.sha256 ORDER BY a.path`},
	}
	for _, q := range queries {
		rows, err := db.Query(q.sql, latest, previous)
		if err != nil {
			return err
		}
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err != nil {
				rows.Close()
				return err
			}
			fmt.Fprintf(w, "%-9s %s\n", q.label, p)
		}
		rows.Close()
	}
	return nil
}

// Dupes prints groups of identical files from the latest run.
func Dupes(db *sql.DB, w io.Writer) error {
	rows, err := db.Query(`
		SELECT sha256, size, group_concat(path, char(10))
		FROM files
		WHERE run_id = (SELECT max(id) FROM runs)
		GROUP BY sha256 HAVING count(*) > 1
		ORDER BY size DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		var sum, paths string
		var size int64
		if err := rows.Scan(&sum, &size, &paths); err != nil {
			return err
		}
		n++
		fmt.Fprintf(w, "group %d (%d bytes each, %s...)\n%s\n", n, size, sum[:12], indent(paths))
	}
	fmt.Fprintf(w, "%d group(s) of duplicates\n", n)
	return rows.Err()
}

func indent(lines string) string {
	out := ""
	for _, l := range splitLines(lines) {
		out += "   " + l + "\n"
	}
	return out
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// Largest prints the n biggest files from the latest run.
func Largest(db *sql.DB, w io.Writer, n int) error {
	rows, err := db.Query(`SELECT size, path FROM files WHERE run_id = (SELECT max(id) FROM runs) ORDER BY size DESC LIMIT ?`, n)
	if err != nil {
		return err
	}
	defer rows.Close()
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.AlignRight)
	for rows.Next() {
		var size int64
		var p string
		if err := rows.Scan(&size, &p); err != nil {
			return err
		}
		fmt.Fprintf(tw, "%s\t %s\n", strconv.FormatInt(size, 10), p)
	}
	tw.Flush()
	return rows.Err()
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: file-index [-db file] index <folder> | changes | dupes | largest [-n 10]")
	}
	dbPath := "file-index.db"
	if args[0] == "-db" && len(args) > 2 {
		dbPath, args = args[1], args[2:]
	}
	db, err := Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	switch args[0] {
	case "index":
		if len(args) != 2 {
			return errors.New("index needs a folder")
		}
		id, n, err := Index(db, args[1], time.Now())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "run %d: %d file(s) indexed into %s\n", id, n, dbPath)
		return nil
	case "changes":
		return Changes(db, out)
	case "dupes":
		return Dupes(db, out)
	case "largest":
		flags := flag.NewFlagSet("largest", flag.ContinueOnError)
		n := flags.Int("n", 10, "how many")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		return Largest(db, out, *n)
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
