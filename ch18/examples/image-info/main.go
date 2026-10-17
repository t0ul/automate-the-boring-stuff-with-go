// open an image file, say what it is, and read one pixel
package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // registering decoders is how image.Decode learns a format
	_ "image/png"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: image-info <file>")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	b := img.Bounds()
	fmt.Printf("format: %s\nsize:   %d x %d\n", format, b.Dx(), b.Dy())

	// colours come back as 16-bit channels; shift by 8 for the familiar 0-255
	r, g, bl, a := img.At(b.Dx()/2, b.Dy()/2).RGBA()
	fmt.Printf("centre pixel: rgba(%d, %d, %d, %d)\n", r>>8, g>>8, bl>>8, a>>8)
}
