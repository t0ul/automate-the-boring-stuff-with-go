// a function worth testing: turn a title into a URL slug
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Slugify lower-cases, keeps letters and digits, and joins words with single hyphens.
func Slugify(title string) string {
	var b strings.Builder
	lastHyphen := true // so we never start with one
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastHyphen = false
		case !lastHyphen:
			b.WriteByte('-')
			lastHyphen = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

func main() {
	fmt.Println(Slugify("Automate the Boring Stuff with Go, Ch. 15: Testing!"))
}
