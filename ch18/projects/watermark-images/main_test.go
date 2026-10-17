package main

import (
	"image"
	"image/color"
	"testing"
)

func flat(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestWatermarkChangesOnlyTheCorner(t *testing.T) {
	face, err := LoadFace(24)
	if err != nil {
		t.Fatal(err)
	}
	src := flat(400, 300, color.RGBA{10, 10, 10, 255})
	logo := flat(40, 40, color.RGBA{255, 0, 0, 255})
	out := Watermark(src, "hello", face, logo)

	if out.Bounds() != src.Bounds() {
		t.Fatalf("size changed: %v", out.Bounds())
	}
	if out.RGBAAt(5, 5) != src.RGBAAt(5, 5) {
		t.Error("top-left pixel changed; the stamp should stay in the bottom-right")
	}
	changed := 0
	for y := 200; y < 300; y++ {
		for x := 250; x < 400; x++ {
			if out.RGBAAt(x, y) != src.RGBAAt(x, y) {
				changed++
			}
		}
	}
	if changed < 1600 { // the 40x40 logo alone is 1600 pixels
		t.Fatalf("only %d pixels changed in the bottom-right; expected the logo and the text", changed)
	}
}
