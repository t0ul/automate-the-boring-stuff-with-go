package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestVisit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>Hello page</title></head><body><h1 id="h">Rendered by Chrome</h1>
<script>document.getElementById("h").textContent += " and JavaScript"</script></body></html>`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	shot := filepath.Join(t.TempDir(), "shot.png")
	title, heading, err := Visit(ctx, srv.URL, shot)
	var execErr *exec.Error
	if errors.As(err, &execErr) {
		t.Skipf("no Chrome on this machine: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if title != "Hello page" || heading != "Rendered by Chrome and JavaScript" {
		t.Fatalf("title=%q heading=%q", title, heading)
	}
	if info, err := os.Stat(shot); err != nil || info.Size() == 0 {
		t.Fatalf("screenshot missing: %v", err)
	}
}
