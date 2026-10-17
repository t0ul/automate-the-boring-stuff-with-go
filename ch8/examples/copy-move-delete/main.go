// copy, move and delete, the three verbs of file management
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyFile streams src into dst, so a 10 GB file never has to fit in memory.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close() // the Close error matters: it is when the last bytes hit the disk
}

func run() error {
	dir, err := os.MkdirTemp("", "verbs")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	original := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(original, []byte("keep this\n"), 0o644); err != nil {
		return err
	}

	// copy
	backup := filepath.Join(dir, "notes.bak")
	if err := copyFile(original, backup); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	// move (rename works across folders on the same disk; copy then delete otherwise)
	archive := filepath.Join(dir, "archive")
	if err := os.MkdirAll(archive, 0o755); err != nil {
		return err
	}
	moved := filepath.Join(archive, "notes.txt")
	if err := os.Rename(original, moved); err != nil {
		return fmt.Errorf("move: %w", err)
	}

	// delete one file, then a whole tree
	if err := os.Remove(backup); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		fmt.Println("left in the folder:", e.Name())
	}
	return os.RemoveAll(archive)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
