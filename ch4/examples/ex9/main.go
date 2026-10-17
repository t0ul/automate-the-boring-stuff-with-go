// ex. 9 errors: return them, wrap them, check them once at the top
package main

import (
	"errors"
	"fmt"
	"os"
)

// ErrEmptyName is a sentinel: a named error value callers can test for.
var ErrEmptyName = errors.New("empty name")

// Hello returns a greeting, or an error when there is nobody to greet.
func Hello(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyName
	}
	return fmt.Sprintf("Hi, %s. Welcome!", name), nil
}

// greetAll stops at the first problem and says which name caused it. The %w verb
// wraps the original error so the caller can still recognise it.
func greetAll(names []string) error {
	for i, name := range names {
		msg, err := Hello(name)
		if err != nil {
			return fmt.Errorf("guest %d: %w", i+1, err)
		}
		fmt.Println(msg)
	}
	return nil
}

// run holds the program's logic and returns its error instead of exiting. main is the
// only place that decides what an error means for the process.
func run() error {
	return greetAll([]string{"Ada", "", "Linus"})
}

func main() {
	if err := run(); err != nil {
		if errors.Is(err, ErrEmptyName) {
			fmt.Fprintln(os.Stderr, "someone forgot to fill in a name:", err)
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		os.Exit(1)
	}
}
