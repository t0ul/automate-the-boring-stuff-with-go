// shrink an image to half size with a good resampling filter
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: resize <in> <out.png>")
		os.Exit(2)
	}
	in, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer in.Close()
	src, _, err := image.Decode(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx()/2, b.Dy()/2))
	// CatmullRom is slow and sharp; ApproxBiLinear is fast and soft. Thumbnails want sharp.
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)

	out, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer out.Close()
	if err := png.Encode(out, dst); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%dx%d -> %dx%d, wrote %s\n", b.Dx(), b.Dy(), dst.Bounds().Dx(), dst.Bounds().Dy(), os.Args[2])
}
