package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func TestOutputName(t *testing.T) {
	if got := OutputName("a/contract.pdf"); got != "a/contract_stamped.pdf" {
		t.Fatalf("OutputName = %q", got)
	}
}

func TestStamp(t *testing.T) {
	out := filepath.Join(t.TempDir(), "stamped.pdf")
	if err := Stamp(context.Background(), "contract.pdf", out, "DRAFT"); err != nil {
		t.Fatal(err)
	}
	if err := api.ValidateFile(context.Background(), out, nil, nil); err != nil {
		t.Fatalf("stamped file does not validate: %v", err)
	}
	n, err := api.PageCountFile(context.Background(), out)
	if err != nil || n != 1 {
		t.Fatalf("pages = %d, %v", n, err)
	}
}
