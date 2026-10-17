// password-checker tells you whether a password meets a simple complexity policy:
// at least 11 characters, with an upper-case letter, a lower-case letter and a digit.
//
//	go run ./ch7/projects/password-checker 'Tr0ub4dor&3xtra'
package main

import (
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"
)

const minLength = 11

var (
	hasUpper = regexp.MustCompile(`[A-Z]`)
	hasLower = regexp.MustCompile(`[a-z]`)
	hasDigit = regexp.MustCompile(`[0-9]`)
)

// Problems returns every rule the password breaks. An empty result means it passes.
func Problems(password string) []string {
	var problems []string
	if utf8.RuneCountInString(password) < minLength {
		problems = append(problems, fmt.Sprintf("shorter than %d characters", minLength))
	}
	if !hasUpper.MatchString(password) {
		problems = append(problems, "no upper-case letter")
	}
	if !hasLower.MatchString(password) {
		problems = append(problems, "no lower-case letter")
	}
	if !hasDigit.MatchString(password) {
		problems = append(problems, "no digit")
	}
	return problems
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: password-checker <password>")
	}
	problems := Problems(args[0])
	if len(problems) == 0 {
		fmt.Println("Password OK")
		return nil
	}
	fmt.Println("Bad Password")
	for _, p := range problems {
		fmt.Println("  -", p)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
