// walk a whole tree and count files by extension
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	counts := map[string]int{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // stop on the first folder we cannot read
		}
		if d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir // do not descend into git's internals
		}
		if !d.IsDir() {
			counts[filepath.Ext(path)]++
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(counts)
}
