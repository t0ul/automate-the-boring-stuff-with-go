package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const page = `<html><body>
<article class="product_pod">
  <p class="star-rating Three"></p>
  <h3><a href="a.html" title="A Light in the Attic">A Light in the ...</a></h3>
  <p class="price_color">£51.77</p>
</article>
<article class="product_pod">
  <p class="star-rating One"></p>
  <h3><a href="b.html" title="Tipping the Velvet">Tipping the Velvet</a></h3>
  <p class="price_color">£53.74</p>
</article>
</body></html>`

func TestParse(t *testing.T) {
	books, err := Parse(strings.NewReader(page))
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 || books[0] != (Book{"A Light in the Attic", "£51.77", "Three"}) || books[1].Rating != "One" {
		t.Fatalf("books = %+v", books)
	}
}

func TestScrape(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(page))
	}))
	defer srv.Close()
	books, err := Scrape(context.Background(), srv.URL)
	if err != nil || len(books) != 2 {
		t.Fatalf("Scrape = %+v, %v", books, err)
	}
}
