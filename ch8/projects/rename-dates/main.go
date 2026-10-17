// rename-dates finds file names that contain an American date (MM-DD-YYYY) and renames
// them to the European order (DD-MM-YYYY). It only prints the plan unless you pass
// -apply, because renaming is the kind of thing you want to see before you do it.
//
//	go run ./ch8/projects/rename-dates ./reports
//	go run ./ch8/projects/rename-dates -apply ./reports
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// americanDate matches 05-21-2026 and 5-21-2026, with the surrounding text captured so
// the rest of the name survives.
var americanDate = regexp.MustCompile(`^(.*?)(\d{1,2})-(\d{1,2})-(\d{4})(.*)$`)

// EuropeanName returns the file name with its date flipped to DD-MM-YYYY, and whether the
// name contained a plausible American date at all.
func EuropeanName(name string) (string, bool) {
	m := americanDate.FindStringSubmatch(name)
	if m == nil {
		return name, false
	}
	before, month, day, year, after := m[1], m[2], m[3], m[4], m[5]
	if !plausible(month, 1, 12) || !plausible(day, 1, 31) {
		return name, false
	}
	return before + day + "-" + month + "-" + year + after, true
}

func plausible(field string, lo, hi int) bool {
	n := 0
	for _, r := range field {
		n = n*10 + int(r-'0')
	}
	return n >= lo && n <= hi
}

// Plan lists the renames that would happen in folder (non-recursive).
func Plan(folder string) (map[string]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", folder, err)
	}
	plan := map[string]string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if newName, ok := EuropeanName(e.Name()); ok && newName != e.Name() {
			plan[e.Name()] = newName
		}
	}
	return plan, nil
}

func run(args []string) error {
	fs := flag.NewFlagSet("rename-dates", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "actually rename (default: print the plan)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: rename-dates [-apply] <folder>")
	}
	folder := fs.Arg(0)
	plan, err := Plan(folder)
	if err != nil {
		return err
	}
	if len(plan) == 0 {
		fmt.Println("nothing to rename")
		return nil
	}
	for from, to := range plan {
		fmt.Printf("%s -> %s\n", from, to)
		if !*apply {
			continue
		}
		if err := os.Rename(filepath.Join(folder, from), filepath.Join(folder, to)); err != nil {
			return fmt.Errorf("rename %s: %w", from, err)
		}
	}
	if !*apply {
		fmt.Println("(dry run; pass -apply to rename)")
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
