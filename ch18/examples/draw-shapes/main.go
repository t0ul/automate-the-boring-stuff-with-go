// make an image from nothing: a background, a rectangle, a circle, saved twice
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))

	// fill the whole thing with one colour
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{245, 240, 225, 255}), image.Point{}, draw.Src)

	// a rectangle is a Draw onto a sub-rectangle
	draw.Draw(img, image.Rect(40, 40, 200, 160), image.NewUniform(color.RGBA{200, 60, 40, 255}), image.Point{}, draw.Src)

	// a circle is a loop over pixels inside its bounding box
	cx, cy, r := 290, 170, 80
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				img.Set(x, y, color.RGBA{40, 90, 200, 255})
			}
		}
	}

	if err := save("shapes.png", func(f *os.File) error { return png.Encode(f, img) }); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := save("shapes.jpg", func(f *os.File) error { return jpeg.Encode(f, img, &jpeg.Options{Quality: 90}) }); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("wrote shapes.png and shapes.jpg")
}

func save(name string, encode func(*os.File) error) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := encode(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
