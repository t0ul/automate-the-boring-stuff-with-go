// backup-to-zip copies a whole folder into a numbered zip file next to it. Run it again
// and you get the next number, so you keep a history of backups.
//
//	go run ./ch8/projects/backup-to-zip ./ch8/projects/backup-to-zip/sample
package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// NextZipName returns "<folder>_1.zip", or _2, _3, ... whichever does not exist yet.
func NextZipName(folder string) (string, error) {
	folder = filepath.Clean(folder)
	for n := 1; ; n++ {
		name := fmt.Sprintf("%s_%d.zip", folder, n)
		_, err := os.Stat(name)
		if errors.Is(err, fs.ErrNotExist) {
			return name, nil
		}
		if err != nil {
			return "", fmt.Errorf("check %s: %w", name, err)
		}
	}
}

// BackupToZip writes every regular file under folder into zipPath, with paths relative
// to the folder's parent so the archive unpacks into a folder of the same name.
func BackupToZip(folder, zipPath string) (count int, err error) {
	out, err := os.Create(zipPath)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", zipPath, err)
	}
	defer func() {
		if cerr := out.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	zw := zip.NewWriter(out)
	base := filepath.Base(filepath.Clean(folder))

	err = filepath.WalkDir(folder, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}
		if err := addFile(zw, path, filepath.ToSlash(filepath.Join(base, rel))); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return count, fmt.Errorf("walk %s: %w", folder, err)
	}
	if err := zw.Close(); err != nil {
		return count, fmt.Errorf("finish zip: %w", err)
	}
	return count, nil
}

func addFile(zw *zip.Writer, path, nameInZip string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := zw.Create(nameInZip)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: backup-to-zip <folder>")
	}
	folder := strings.TrimSuffix(args[0], string(filepath.Separator))
	info, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a folder", folder)
	}
	zipPath, err := NextZipName(folder)
	if err != nil {
		return err
	}
	n, err := BackupToZip(folder, zipPath)
	if err != nil {
		return err
	}
	fmt.Printf("backed up %d file(s) to %s\n", n, zipPath)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
