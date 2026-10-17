package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestFit(t *testing.T) {
	cases := []struct{ w, h, max, ww, hh int }{
		{800, 600, 300, 300, 225},
		{300, 900, 300, 100, 300},
		{200, 100, 300, 200, 100},
		{1000, 1000, 100, 100, 100},
	}
	for _, c := range cases {
		if w, h := Fit(c.w, c.h, c.max); w != c.ww || h != c.hh {
			t.Errorf("Fit(%d,%d,%d) = %dx%d, want %dx%d", c.w, c.h, c.max, w, h, c.ww, c.hh)
		}
	}
}

func writePNG(t *testing.T, path string, w, h int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestProcess(t *testing.T) {
	dir := t.TempDir()
	writePNG(t, filepath.Join(dir, "wide.png"), 800, 600)
	writePNG(t, filepath.Join(dir, "tall.png"), 300, 900)
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("skip me"), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := Process(dir, 300, 2)
	if err != nil || n != 2 {
		t.Fatalf("Process = %d, %v", n, err)
	}
	for name, want := range map[string][2]int{"wide.jpg": {300, 225}, "tall.jpg": {100, 300}} {
		f, err := os.Open(filepath.Join(dir, "thumbs", name))
		if err != nil {
			t.Fatal(err)
		}
		cfg, _, err := image.DecodeConfig(f)
		f.Close()
		if err != nil || cfg.Width != want[0] || cfg.Height != want[1] {
			t.Errorf("%s = %dx%d, %v; want %v", name, cfg.Width, cfg.Height, err, want)
		}
	}
}
