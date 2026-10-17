// json-get pulls one value out of a JSON document by path, the way you would with jq,
// and pretty-prints the whole document when you give it no path. It reads a file or
// standard input, so it chains with the other tools in this book.
//
//	go run ./ch12/projects/json-get ./ch9/projects/census-by-county/census_summary.json
//	curl -s https://xkcd.com/info.0.json | go run ./ch12/projects/json-get -path title
//	go run ./ch12/projects/json-get -path comics.1.title ./ch12/projects/json-get/sample.json
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Get walks v along a dotted path. Object keys are names; array elements are indexes.
func Get(v any, path string) (any, error) {
	if path == "" {
		return v, nil
	}
	for _, part := range strings.Split(path, ".") {
		switch node := v.(type) {
		case map[string]any:
			next, ok := node[part]
			if !ok {
				return nil, fmt.Errorf("no key %q", part)
			}
			v = next
		case []any:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(node) {
				return nil, fmt.Errorf("no element %q in an array of %d", part, len(node))
			}
			v = node[i]
		default:
			return nil, fmt.Errorf("cannot descend into %T with %q", v, part)
		}
	}
	return v, nil
}

// Render prints strings bare and everything else as indented JSON, which is what you
// want when the result goes straight into another command.
func Render(v any) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	out, err := json.MarshalIndent(v, "", "  ")
	return string(out), err
}

func run(args []string, stdin io.Reader) error {
	flags := flag.NewFlagSet("json-get", flag.ContinueOnError)
	path := flags.String("path", "", "dotted path, e.g. comics.1.title (empty = whole document)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var in io.Reader = stdin
	switch flags.NArg() {
	case 0:
	case 1:
		f, err := os.Open(flags.Arg(0))
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	default:
		return errors.New("usage: json-get [-path a.b.0] [file.json]")
	}
	var v any
	if err := json.NewDecoder(in).Decode(&v); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}
	got, err := Get(v, *path)
	if err != nil {
		return err
	}
	s, err := Render(got)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
