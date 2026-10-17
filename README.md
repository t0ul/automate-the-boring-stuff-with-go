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
- [Fin](https://www.t0ul.com/blog/introducing-automate-the-boring-stuff-with-go/)

## 2026 edition

The book is being updated for Go 1.27. The review and the checklist are in [docs/2026-EDITION.md](docs/2026-EDITION.md).

## Credits

Some of the chapter 9–11 projects are adapted from [automateGo](https://github.com/dreddsa5dies/automateGo) by Viktor Solovev (MIT License); see [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md). Everything else is licensed under the GPL-3.0 in [LICENSE](LICENSE).
