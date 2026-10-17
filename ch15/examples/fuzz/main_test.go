package main

import (
	"testing"
	"unicode/utf8"
)

// A fuzz test states a property that must hold for ANY input, and lets Go invent inputs.
// Run the seeds only: go test ./ch15/examples/fuzz
// Fuzz for ten seconds:  go test -fuzz=FuzzReverse -fuzztime=10s ./ch15/examples/fuzz
func FuzzReverse(f *testing.F) {
	for _, seed := range []string{"", "a", "Hello, 世界", "!12345"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("not valid UTF-8; reversing bytes of a broken string is undefined")
		}
		twice := Reverse(Reverse(s))
		if twice != s {
			t.Errorf("Reverse(Reverse(%q)) = %q", s, twice)
		}
		if utf8.RuneCountInString(Reverse(s)) != utf8.RuneCountInString(s) {
			t.Errorf("Reverse(%q) changed the character count", s)
		}
	})
}
