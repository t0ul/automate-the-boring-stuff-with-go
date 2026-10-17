// watermark-images stamps text, a logo, or both onto the bottom-right corner of every
// image in a folder, writing copies to a watermarked/ sub-folder. The originals are not
// touched.
//
//	go run ./ch18/projects/watermark-images -text "© Toul Cranmer" -logo ./ch18/projects/watermark-images/sample/logo.png ./ch18/projects/watermark-images/sample
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// LoadFace builds a font face from the Go font that ships with x/image, so there is no
// font file to find on the reader's machine.
func LoadFace(size float64) (font.Face, error) {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72})
}

// StampText draws text in the bottom-right corner: a dark shadow, then light text on top,
// so it reads on both bright and dark photos.
func StampText(dst *image.RGBA, text string, face font.Face, margin int) {
	d := &font.Drawer{Dst: dst, Face: face}
	width := d.MeasureString(text).Ceil()
	b := dst.Bounds()
	x := b.Max.X - margin - width
	y := b.Max.Y - margin
	d.Src = image.NewUniform(color.RGBA{0, 0, 0, 160})
	d.Dot = fixed.P(x+2, y+2)
	d.DrawString(text)
	d.Src = image.NewUniform(color.RGBA{255, 255, 255, 220})
	d.Dot = fixed.P(x, y)
	d.DrawString(text)
}

// Overlay composites logo onto the bottom-right corner, above any text, respecting the
// logo's own transparency.
func Overlay(dst *image.RGBA, logo image.Image, margin, liftBy int) {
	lb := logo.Bounds()
	b := dst.Bounds()
	at := image.Rect(b.Max.X-margin-lb.Dx(), b.Max.Y-margin-liftBy-lb.Dy(), b.Max.X-margin, b.Max.Y-margin-liftBy)
	draw.Draw(dst, at, logo, lb.Min, draw.Over)
}

// Watermark returns a stamped copy of src. Either text or logo may be empty/nil.
func Watermark(src image.Image, text string, face font.Face, logo image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	margin := max(12, b.Dx()/50)
	lift := 0
	if text != "" {
		StampText(dst, text, face, margin)
		lift = face.Metrics().Height.Ceil() + margin/2
	}
	if logo != nil {
		Overlay(dst, logo, margin, lift)
	}
	return dst
}

func decode(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("watermark-images", flag.ContinueOnError)
	text := flags.String("text", "", "text to stamp")
	logoPath := flags.String("logo", "", "PNG logo to stamp")
	fontSize := flags.Float64("font-size", 28, "text size in points")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 || (*text == "" && *logoPath == "") {
		return errors.New("usage: watermark-images [-text T] [-logo file.png] <folder>")
	}
	folder := flags.Arg(0)
	face, err := LoadFace(*fontSize)
	if err != nil {
		return err
	}
	var logo image.Image
	if *logoPath != "" {
		if logo, err = decode(*logoPath); err != nil {
			return err
		}
	}
	outDir := filepath.Join(folder, "watermarked")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	n := 0
	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if e.IsDir() || (ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			continue
		}
		src := filepath.Join(folder, e.Name())
		if abs, _ := filepath.Abs(src); *logoPath != "" && abs == mustAbs(*logoPath) {
			continue // do not stamp the logo onto itself
		}
		img, err := decode(src)
		if err != nil {
			return err
		}
		out := filepath.Join(outDir, strings.TrimSuffix(e.Name(), ext)+".png")
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		if err := png.Encode(f, Watermark(img, *text, face, logo)); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		fmt.Println("wrote", out)
		n++
	}
	fmt.Printf("%d image(s) watermarked\n", n)
	return nil
}

func mustAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
