// scrape-books reads a catalogue page and prints title, price and star rating for every
// book on it. The default target is books.toscrape.com, a site that exists to be
// scraped. Scraping a site that did not invite it is a different conversation: read its
// robots.txt and terms, go slowly, identify yourself, and prefer an API when there is one.
//
//	go run ./ch11/projects/scrape-books
//	go run ./ch11/projects/scrape-books -url https://books.toscrape.com/catalogue/page-2.html
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{Timeout: 20 * time.Second}

// Book is one row of the catalogue.
type Book struct {
	Title  string
	Price  string
	Rating string
}

// Parse extracts the books from an HTML page in the books.toscrape.com layout.
func Parse(r io.Reader) ([]Book, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}
	var books []Book
	doc.Find("article.product_pod").Each(func(_ int, s *goquery.Selection) {
		title, _ := s.Find("h3 a").Attr("title")
		rating := ""
		if classes, ok := s.Find("p.star-rating").Attr("class"); ok {
			rating = strings.TrimSpace(strings.TrimPrefix(classes, "star-rating"))
		}
		books = append(books, Book{
			Title:  title,
			Price:  strings.TrimSpace(s.Find("p.price_color").Text()),
			Rating: rating,
		})
	})
	return books, nil
}

// Scrape fetches url and parses it.
func Scrape(ctx context.Context, pageURL string) ([]Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "automate-the-boring-stuff-with-go/2026 (+https://www.t0ul.com)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", pageURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s: %s", pageURL, resp.Status)
	}
	return Parse(resp.Body)
}

func run(args []string) error {
	flags := flag.NewFlagSet("scrape-books", flag.ContinueOnError)
	pageURL := flags.String("url", "https://books.toscrape.com/", "catalogue page to read")
	if err := flags.Parse(args); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	books, err := Scrape(ctx, *pageURL)
	if err != nil {
		return err
	}
	for _, b := range books {
		fmt.Printf("%-8s %-6s %s\n", b.Price, b.Rating, b.Title)
	}
	fmt.Fprintf(os.Stderr, "%d book(s) on %s\n", len(books), *pageURL)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
