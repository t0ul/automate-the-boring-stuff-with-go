// genfixtures writes every input file the chapter 9 and 10 projects read. The data is
// invented and deterministic, so anyone can regenerate it and diff the result.
//
//	go run ./tools/genfixtures          # everything
//	go run ./tools/genfixtures ch18     # only the steps whose name starts with ch18
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"ch9/projects/read-workbook/example.xlsx", exampleWorkbook},
		{"ch9/projects/census-by-county/census.xlsx", censusWorkbook},
		{"ch9/projects/update-produce-prices/produceSales.xlsx", produceWorkbook},
		{"ch9/projects/txt-to-xlsx/*.txt", textFiles},
		{"ch9/projects/xlsx-to-txt/table.xlsx", tableWorkbook},
		{"ch10/projects/*/*.pdf", pdfs},
		{"ch17/projects/dues-reminders/members.xlsx", membersWorkbook},
		{"ch18/projects/thumbnails/sample/*", sampleImages},
		{"ch18/projects/watermark-images/sample/*", watermarkImages},
	}
	only := ""
	if len(os.Args) > 1 {
		only = os.Args[1]
	}
	for _, s := range steps {
		if !strings.HasPrefix(s.name, only) {
			continue
		}
		if err := s.fn(); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
		fmt.Println("wrote", s.name)
	}
	return nil
}

func save(f *excelize.File, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return f.Close()
}

func exampleWorkbook() error {
	f := excelize.NewFile()
	rows := [][]any{
		{"Date", "Item", "Quantity"},
		{"2026-03-02", "Apples", 73},
		{"2026-03-02", "Cherries", 85},
		{"2026-03-03", "Pears", 14},
		{"2026-03-04", "Oranges", 52},
		{"2026-03-04", "Apples", 152},
		{"2026-03-05", "Bananas", 23},
		{"2026-03-06", "Strawberries", 98},
	}
	for i, row := range rows {
		if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+1), &row); err != nil {
			return err
		}
	}
	if _, err := f.NewSheet("Sheet2"); err != nil {
		return err
	}
	if err := f.SetCellValue("Sheet2", "A1", "Hello from the second sheet"); err != nil {
		return err
	}
	return save(f, "ch9/projects/read-workbook/example.xlsx")
}

func censusWorkbook() error {
	r := rand.New(rand.NewPCG(2026, 9))
	counties := map[string][]string{
		"Vermont":  {"Addison", "Bennington", "Chittenden", "Windsor"},
		"Idaho":    {"Ada", "Bonneville", "Kootenai"},
		"Delaware": {"Kent", "New Castle", "Sussex"},
	}
	f := excelize.NewFile()
	header := []any{"CensusTract", "State", "County", "POP2010"}
	if err := f.SetSheetRow("Sheet1", "A1", &header); err != nil {
		return err
	}
	row, tract := 2, 100
	for _, state := range []string{"Delaware", "Idaho", "Vermont"} {
		for _, county := range counties[state] {
			for n := 0; n < 3+r.IntN(4); n++ {
				tract++
				values := []any{fmt.Sprintf("%06d", tract), state, county, 1500 + r.IntN(6000)}
				if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", row), &values); err != nil {
					return err
				}
				row++
			}
		}
	}
	return save(f, "ch9/projects/census-by-county/census.xlsx")
}

func produceWorkbook() error {
	r := rand.New(rand.NewPCG(2026, 10))
	produce := []struct {
		name  string
		price float64
	}{
		{"Potatoes", 0.86}, {"Okra", 2.26}, {"Fava beans", 2.69}, {"Watermelon", 0.66},
		{"Garlic", 1.19}, {"Parsnips", 2.27}, {"Asparagus", 2.49}, {"Avocados", 3.23},
		{"Celery", 1.19}, {"Lemon", 1.27}, {"Carrots", 1.07}, {"Oranges", 1.10},
	}
	f := excelize.NewFile()
	header := []any{"PRODUCE", "COST PER POUND", "POUNDS SOLD", "TOTAL"}
	if err := f.SetSheetRow("Sheet1", "A1", &header); err != nil {
		return err
	}
	for row := 2; row <= 40; row++ {
		p := produce[r.IntN(len(produce))]
		pounds := 5 + r.IntN(60)
		values := []any{p.name, p.price, pounds, round2(p.price * float64(pounds))}
		if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", row), &values); err != nil {
			return err
		}
	}
	return save(f, "ch9/projects/update-produce-prices/produceSales.xlsx")
}

func round2(x float64) float64 { return float64(int(x*100+0.5)) / 100 }

func textFiles() error {
	files := map[string]string{
		"first.txt":  "Monday\nTuesday\nWednesday\nThursday\nFriday\n",
		"second.txt": "coffee\ntea\nwater\n",
		"third.txt":  "one\ntwo\nthree\nfour\nfive\nsix\nseven\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join("ch9/projects/txt-to-xlsx", name), []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func tableWorkbook() error {
	f := excelize.NewFile()
	cols := [][]any{
		{"Fruit", "apple", "pear", "plum", "fig"},
		{"Colour", "red", "green", "purple", "brown"},
		{"Season", "autumn", "autumn", "summer", "summer"},
	}
	for c, col := range cols {
		name, err := excelize.ColumnNumberToName(c + 1)
		if err != nil {
			return err
		}
		if err := f.SetSheetCol("Sheet1", name+"1", &col); err != nil {
			return err
		}
	}
	return save(f, "ch9/projects/xlsx-to-txt/table.xlsx")
}

func pdfs() error {
	docs := map[string][]string{
		"ch10/projects/read-pdf-text/sample.pdf": {
			"Automate the Boring Stuff with Go\nChapter 10: PDFs\n\nThis is page one of the sample document.",
			"Page two.\nThe text on this page was extracted by a Go program.",
		},
		"ch10/projects/merge-pdfs/part1.pdf":        {"Part one, page one.", "Part one, page two."},
		"ch10/projects/merge-pdfs/part2.pdf":        {"Part two, the only page."},
		"ch10/projects/merge-pdfs/part3.pdf":        {"Part three, page one.", "Part three, page two.", "Part three, page three."},
		"ch10/projects/watermark-pdfs/contract.pdf": {"Service agreement\n\nSection 1. Both parties agree to keep this document confidential."},
	}
	for path, pages := range docs {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, MiniPDF(pages), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// savePNG and saveJPEG write an image, creating the folder if needed.
func savePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func saveJPEG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 85}); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// gradient paints a horizontal blue-to-orange gradient with a soft circle in it.
func gradient(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := float64(x) / float64(w)
			r := uint8(30 + 200*t)
			g := uint8(90 + 60*t)
			b := uint8(200 - 150*t)
			d := math.Hypot(float64(x-w/2), float64(y-h/2))
			if d < float64(h)/4 {
				r, g, b = 250, 240, 200
			}
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// stripes paints diagonal stripes, a shape that shows resampling quality.
func stripes(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x+y)/24%2 == 0 {
				img.Set(x, y, color.RGBA{40, 120, 80, 255})
			} else {
				img.Set(x, y, color.RGBA{230, 230, 220, 255})
			}
		}
	}
	return img
}

func sampleImages() error {
	if err := savePNG("ch18/projects/thumbnails/sample/sunset.png", gradient(800, 600)); err != nil {
		return err
	}
	if err := saveJPEG("ch18/projects/thumbnails/sample/stripes.jpg", stripes(640, 480)); err != nil {
		return err
	}
	return savePNG("ch18/projects/thumbnails/sample/tall.png", gradient(300, 900))
}

func watermarkImages() error {
	if err := savePNG("ch18/projects/watermark-images/sample/photo.png", gradient(640, 480)); err != nil {
		return err
	}
	// a small round logo with a transparent background
	logo := image.NewRGBA(image.Rect(0, 0, 96, 96))
	for y := 0; y < 96; y++ {
		for x := 0; x < 96; x++ {
			d := math.Hypot(float64(x-48), float64(y-48))
			switch {
			case d < 40:
				logo.Set(x, y, color.RGBA{20, 20, 20, 230})
			case d < 46:
				logo.Set(x, y, color.RGBA{255, 255, 255, 230})
			}
		}
	}
	return savePNG("ch18/projects/watermark-images/sample/logo.png", logo)
}

func membersWorkbook() error {
	f := excelize.NewFile()
	rows := [][]any{
		{"Name", "Email", "Jan", "Feb", "Mar"},
		{"Ada Lovelace", "ada@example.com", "paid", "paid", "paid"},
		{"Grace Hopper", "grace@example.com", "paid", "paid", ""},
		{"Linus Torvalds", "linus@example.com", "paid", "", ""},
		{"Ken Thompson", "ken@example.com", "paid", "paid", "paid"},
		{"Rob Pike", "rob@example.com", "", "paid", ""},
	}
	for i, row := range rows {
		if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+1), &row); err != nil {
			return err
		}
	}
	return save(f, "ch17/projects/dues-reminders/members.xlsx")
}
