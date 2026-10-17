// a struct in, JSON out, and back again
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// The tags say what each field is called in JSON. Without a tag, the field name is used
// as written, capital letter and all. Fields that start with a lower-case letter are
// private and never make it into JSON.
type Book struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Year     int      `json:"year"`
	Pages    int      `json:"pages,omitempty"` // left out of the output when zero
	internal string   // unexported: ignored by encoding/json
}

func main() {
	b := Book{Title: "The Go Programming Language", Authors: []string{"Donovan", "Kernighan"}, Year: 2015}

	out, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))

	in := `{"title":"Automate the Boring Stuff with Go","authors":["Toul Cranmer"],"year":2026,"extra":"ignored"}`
	var got Book
	if err := json.Unmarshal([]byte(in), &got); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", got)
}
