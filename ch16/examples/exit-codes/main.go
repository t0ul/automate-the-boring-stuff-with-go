// exit codes are how your program talks to scripts: 0 is success, anything else is not
package main

import (
	"errors"
	"fmt"
	"os"
)

var errUsage = errors.New("usage: exit-codes <ok|fail|usage>")

func run(args []string) error {
	if len(args) != 1 {
		return errUsage
	}
	switch args[0] {
	case "ok":
		fmt.Println("all good")
		return nil
	case "fail":
		return errors.New("the thing could not be done")
	default:
		return errUsage
	}
}

func main() {
	err := run(os.Args[1:])
	switch {
	case err == nil:
		os.Exit(0) // not needed, main returning does the same; shown for symmetry
	case errors.Is(err, errUsage):
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2) // the convention for "you called me wrong"
	default:
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1) // the convention for "I tried and failed"
	}
}
