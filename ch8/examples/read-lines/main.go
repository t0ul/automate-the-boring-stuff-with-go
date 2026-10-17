// read a file line by line, and write a file in one call
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.MkdirTemp("", "read-lines")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir) // clean up the temp folder when main returns

	path := filepath.Join(dir, "todo.txt")
	// os.WriteFile creates or replaces the file in one call; 0o644 is "owner can write, all can read"
	if err := os.WriteFile(path, []byte("buy milk\nwalk dog\n\nfix bike\n"), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer f.Close()

	// bufio.Scanner hands you one line at a time, however big the file is
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		n++
		fmt.Printf("%d. %s\n", n, line)
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
