# Automate the Boring Stuff with Go

Source code for [Automate the Boring Stuff with Go](https://www.t0ul.com/blog/introducing-automate-the-boring-stuff-with-go/), a beginner-friendly guide to programming in Go. The chapters are on t0ul.com; this repository holds every example and project.

## Running the code

Install the current Go release from [go.dev/dl](https://go.dev/dl/), clone this repository, and run any example or project by its directory:

```bash
git clone https://github.com/t0ul/automate-the-boring-stuff-with-go.git
cd automate-the-boring-stuff-with-go
go run ./ch3/projects/GuessTheNumber
```

There is one Go module at the root, so `go build ./...` builds everything and `go vet ./...` checks it. CI runs build, vet, test, gofmt and govulncheck on every change.

## Chapters 

- [Introduction](https://www.t0ul.com/blog/automate-go-ch1-introduction/)
- [Setting up Developer Environment](https://www.t0ul.com/blog/automate-go-ch1-introduction/)
- [The Basics of GO](https://www.t0ul.com/blog/automate-go-ch2-go-basics/)
- [Logic Flow](https://www.t0ul.com/blog/automate-go-ch3-control-flow/)
- [Functions](https://www.t0ul.com/blog/automate-go-ch4-functions/)
- [Slices](https://www.t0ul.com/blog/automate-go-ch5-slices/)
- [Maps](https://www.t0ul.com/blog/automate-go-ch6-syntax/)
- [String Manipulation](https://www.t0ul.com/blog/automate-go-ch7-strings/)
- [Files](https://www.t0ul.com/blog/automate-go-ch8-files/)
- [Excel Files](https://www.t0ul.com/blog/automate-go-ch9-excel-files/)
- [PDFs](https://www.t0ul.com/blog/automate-go-ch10-pdfs/)
- [Web Requests](https://www.t0ul.com/blog/automate-go-ch11-web-requests/)
- [CSV and JSON](https://www.t0ul.com/blog/automate-go-ch12-csv-and-json/)
- [Time, Scheduling and Running Programs](https://www.t0ul.com/blog/automate-go-ch13-time-and-scheduling/)
- [Concurrency](https://www.t0ul.com/blog/automate-go-ch14-concurrency/)
- [Testing and Debugging](https://www.t0ul.com/blog/automate-go-ch15-testing-and-debugging/)
- [Command-Line Tools and Shipping Binaries](https://www.t0ul.com/blog/automate-go-ch16-command-line-tools/)
- [Fin](https://www.t0ul.com/blog/introducing-automate-the-boring-stuff-with-go/)

## 2026 edition

The book is being updated for Go 1.27. The review and the checklist are in [docs/2026-EDITION.md](docs/2026-EDITION.md).

## Projects by chapter

| Chapter | Projects |
|---|---|
| 3 Control flow | `ch3/projects/GuessTheNumber`, `ch3/projects/RockPaperScissors` |
| 7 Strings | `ch7/projects/PlainSearch`, `Replacement`, `email-finder`, `password-checker`, `quiz-builder` |
| 8 Files | `ch8/projects/backup-to-zip`, `find-large-files`, `rename-dates` |
| 9 Excel | `ch9/projects/read-workbook`, `census-by-county`, `write-workbook`, `update-produce-prices`, `multiplication-table`, `txt-to-xlsx`, `xlsx-to-txt` (excelize) |
| 10 PDFs | `ch10/projects/read-pdf-text`, `merge-pdfs`, `watermark-pdfs` (pdfcpu, ledongthuc/pdf) |
| 11 Web | `ch11/projects/fetch-page`, `search-term`, `save-page`, `download-xkcd`, `scrape-books` (goquery), `browser-chromedp` (headless Chrome) |
| 12 CSV and JSON | `ch12/projects/csv-remove-header`, `csv-to-json`, `json-get` |
| 13 Time and scheduling | `ch13/projects/stopwatch`, `countdown`, `every` |
| 14 Concurrency | `ch14/projects/parallel-xkcd`, `find-duplicates` |
| 15 Testing and debugging | `ch15/projects/retry` (plus table, fuzz and benchmark examples) |
| 16 Command-line tools | `ch16/projects/todo` (subcommands, embed, version via ldflags, `build.sh`) |

Every project is `go run ./<dir>`; inputs sit next to the code and outputs are written next to the code and ignored by git. `go run ./tools/genfixtures` regenerates the workbook and PDF inputs; `go run ./tools/provenance` checks that no file is copied from anywhere else.

## License

GPL-3.0, see [LICENSE](LICENSE).
