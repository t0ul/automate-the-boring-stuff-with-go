// email-finder reads whatever is on the clipboard, prints every email address it
// finds, and puts the list back on the clipboard, one address per line.
//
//	go run ./ch7/projects/email-finder
package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

// emailPattern is deliberately simple: a word-like local part, an @, a domain with at
// least one dot. It is for finding addresses in prose, not for validating them.
var emailPattern = regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`)

// FindEmails returns the distinct email addresses in text, in order of first appearance.
func FindEmails(text string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range emailPattern.FindAllString(text, -1) {
		m = strings.ToLower(m)
		if seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	return out
}

func run() error {
	text, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}
	emails := FindEmails(text)
	if len(emails) == 0 {
		return errors.New("no email addresses on the clipboard")
	}
	joined := strings.Join(emails, "\n")
	fmt.Println(joined)
	if err := clipboard.WriteAll(joined); err != nil {
		return fmt.Errorf("write clipboard: %w", err)
	}
	fmt.Fprintf(os.Stderr, "copied %d address(es) to the clipboard\n", len(emails))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
