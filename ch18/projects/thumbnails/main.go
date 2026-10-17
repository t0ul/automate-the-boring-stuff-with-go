// thumbnails makes a small copy of every image in a folder, longest side -size pixels,
// into a thumbs/ sub-folder. Resampling is CPU work, so the files are processed by a
// pool of workers (chapter 14).
//
//	go run ./ch18/projects/thumbnails -size 300 ./ch18/projects/thumbnails/sample
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
)

// Fit returns the size that fits w x h inside a max x max box, keeping the aspect ratio.
// Images already smaller than the box are left alone.
func Fit(w, h, max int) (int, int) {
	if w <= max && h <= max {
		return w, h
	}
	if w >= h {
		return max, h * max / w
	}
	return w * max / h, max
}

// Thumbnail returns src scaled to fit in a max x max box.
func Thumbnail(src image.Image, max int) image.Image {
	b := src.Bounds()
	w, h := Fit(b.Dx(), b.Dy(), max)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Src, nil)
	return dst
}

// One reads src, writes a JPEG thumbnail to dst, and returns the thumbnail's size.
func One(src, dst string, max int) (int, int, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, 0, err
	}
	defer in.Close()
	img, _, err := image.Decode(in)
	if err != nil {
		return 0, 0, fmt.Errorf("decode %s: %w", src, err)
	}
	thumb := Thumbnail(img, max)
	out, err := os.Create(dst)
	if err != nil {
		return 0, 0, err
	}
	if err := jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85}); err != nil {
		out.Close()
		return 0, 0, err
	}
	b := thumb.Bounds()
	return b.Dx(), b.Dy(), out.Close()
}

// Process thumbnails every .png, .jpg and .jpeg directly inside folder into
// folder/thumbs, with at most workers running at once. It returns how many it wrote.
func Process(folder string, max, workers int) (int, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return 0, err
	}
	outDir := filepath.Join(folder, "thumbs")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return 0, err
	}
	var g errgroup.Group
	g.SetLimit(workers)
	done := make(chan string, len(entries))
	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if e.IsDir() || (ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			continue
		}
		src := filepath.Join(folder, e.Name())
		dst := filepath.Join(outDir, strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))+".jpg")
		g.Go(func() error {
			w, h, err := One(src, dst, max)
			if err != nil {
				return err
			}
			done <- fmt.Sprintf("%s -> %dx%d", e.Name(), w, h)
			return nil
		})
	}
	err = g.Wait()
	close(done)
	n := 0
	for line := range done {
		fmt.Println(line)
		n++
	}
	return n, err
}

func run(args []string) error {
	flags := flag.NewFlagSet("thumbnails", flag.ContinueOnError)
	size := flags.Int("size", 300, "longest side of a thumbnail, in pixels")
	workers := flags.Int("workers", runtime.NumCPU(), "images to process at once")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: thumbnails [-size 300] [-workers N] <folder>")
	}
	n, err := Process(flags.Arg(0), *size, *workers)
	if err != nil {
		return err
	}
	fmt.Printf("%d thumbnail(s) in %s\n", n, filepath.Join(flags.Arg(0), "thumbs"))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
