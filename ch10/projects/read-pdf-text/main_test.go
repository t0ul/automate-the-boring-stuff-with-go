package main

import (
	"strings"
	"testing"
)

func TestPageTexts(t *testing.T) {
	pages, err := PageTexts("sample.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(pages))
	}
	if !strings.Contains(pages[0], "Chapter 10") || !strings.Contains(pages[1], "Page two") {
		t.Fatalf("unexpected text: %q / %q", pages[0], pages[1])
	}
}
