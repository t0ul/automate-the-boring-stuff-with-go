# Automate the Boring Stuff with Go: 2026 edition

The book was written for Go 1.19 (chapters published November 2022 to June 2023). This
file is the review that started the 2026 edition and the checklist that tracks it. The
chapter text lives on [t0ul.com](https://www.t0ul.com/blog/introducing-automate-the-boring-stuff-with-go/);
this repository holds the code.

## What the review found (2026-10-17)

Toolchain and build

- Go is 1.27 (go.dev/dl). Chapter 1 installed go1.19.3 and explained GOPATH; chapter 2
  ran `go get` outside a module. Neither works the way the text says any more.
- The repository had no `go.mod`, so every project with a third-party import failed with
  "no required module provides package" on a current toolchain.
- No tests, no CI. `go vet` flagged an unused `fmt.Errorf` and two `Printf` calls with no
  verb. `io/ioutil` (deprecated since Go 1.16) was still in three chapter 11 projects.
  Four network projects called `http.Get` with no timeout.

Libraries the book taught

| Taught | Status | Replacement |
|---|---|---|
| github.com/tealeg/xlsx | archived (August 2025) | github.com/xuri/excelize/v2 |
| github.com/jung-kurt/gofpdf | archived (2021); fork go-pdf/fpdf archived (March 2025) | github.com/pdfcpu/pdfcpu |
| github.com/ledongthuc/pdf | active | kept, for text extraction |
| github.com/tebeka/selenium | slow (last push January 2025), needs a chromedriver binary | github.com/chromedp/chromedp |
| opesun/goquery, jackdanger/collectlinks (text only) | unmaintained since 2016 | github.com/PuerkitoBio/goquery |
| github.com/atotto/clipboard | active | kept |

Content

- Chapters 6 (maps, 329 words), 8 (files, 596 words) and 10 (PDFs, 333 words) are thin;
  chapters 8 to 11 are code with headings and almost no prose.
- Structs are never taught. Generics, the `slices` and `maps` packages, `min`/`max`,
  range over int, error wrapping (`%w`, `errors.Is`, `errors.As`, `errors.Join`),
  `log/slog`, goroutines and `go test` do not appear.
- Topics the Python original covers that the Go book never reached: CSV and JSON, time
  and scheduling, launching programs, email and notifications, images, debugging, testing,
  command-line tools, building and sharing binaries.

Provenance

- Nineteen projects in chapters 7 to 11 and sixteen fixture files were derived from or
  copied out of dreddsa5dies/automateGo (MIT). The notice only named chapters 9 to 11.
  The 2026 edition rewrites every one of them clean-room and regenerates the fixtures, then
  removes the notice. A provenance check under `tools/` keeps it that way.

Hygiene

- The `chN.md` chapter copies had drifted hundreds of lines from the published posts.
  Deleted; the site is the single source of truth.
- Program outputs (`save.pdf`, `saveTable.xlsx`, `saveFile*.txt`, `save.txt`,
  `updateproduceSales.xlsx`) were committed as if they were inputs. Deleted and ignored.

## The edition, in order

Code (this repository)

- [x] A. Root module, CI (build, vet, test, gofmt, govulncheck), vet fixes, chapter copies
  and outputs removed, this file.
- [x] B. Clean-room rewrite of the derived projects on excelize, pdfcpu, ledongthuc/pdf,
  goquery and chromedp; own fixtures from `tools/genfixtures`; tests beside every project;
  `tools/provenance` in CI; notice removed.

Text (t0ul.com)

- [ ] C. Chapters 1, 2, 4, 5, 6, 7, 8, 9, 10, 11 updated for Go 1.27 and the new code.
- [ ] D. New chapters: 12 CSV and JSON, 13 time, scheduling and launching programs,
  14 concurrency for automation, 15 testing and debugging, 16 command-line tools and
  shipping binaries.
- [ ] E. Introduction updated with what changed; cover and listing say 2026 edition.

## Conventions for 2026 code

- One module at the root; every example and project is its own `package main` directory
  and runs with `go run ./chN/<dir>`.
- `main` calls `run() error` and exits non-zero on error. Errors are returned and wrapped
  with `%w`; `log.Fatal` appears only in `main`.
- Network calls use an `http.Client` with a timeout and a `context`.
- Inputs are committed next to the project; outputs are written next to the project and
  ignored by git.
- Where a project has pure logic, a `_test.go` sits beside it.
