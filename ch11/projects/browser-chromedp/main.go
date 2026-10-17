// browser-chromedp drives a real, headless Chrome: load a page, wait for an element,
// read the title and a heading, click a link, and save a screenshot. Use it when a page
// only makes sense after JavaScript has run; for plain HTML, net/http is enough.
//
//	go run ./ch11/projects/browser-chromedp https://books.toscrape.com/
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// Visit loads url, returns the page title and the text of the first h1, and writes a
// screenshot to shot. Pass an empty shot to skip the screenshot.
func Visit(ctx context.Context, url, shot string) (title, heading string, err error) {
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, chromedp.DefaultExecAllocatorOptions[:]...)
	defer cancel()
	tab, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	if err := chromedp.Do(tab, chromedp.Navigate(url), chromedp.WaitVisible("h1")); err != nil {
		return "", "", fmt.Errorf("load %s: %w", url, err)
	}
	if title, err = chromedp.Run(tab, chromedp.Title()); err != nil {
		return "", "", err
	}
	if heading, err = chromedp.Run(tab, chromedp.Text("h1")); err != nil {
		return "", "", err
	}
	if shot != "" {
		png, err := chromedp.Run(tab, chromedp.FullScreenshot(90))
		if err != nil {
			return title, heading, fmt.Errorf("screenshot: %w", err)
		}
		if err := os.WriteFile(shot, png, 0o644); err != nil {
			return title, heading, err
		}
	}
	return title, heading, nil
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: browser-chromedp <url>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	title, heading, err := Visit(ctx, args[0], "screenshot.png")
	if err != nil {
		return err
	}
	fmt.Printf("title:   %s\nheading: %s\nsaved screenshot.png\n", title, heading)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
