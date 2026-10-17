// three string tools you will use constantly: Cut, Builder, and runes
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	// strings.Cut splits at the first separator and tells you whether it was there
	key, value, found := strings.Cut("timeout=30s", "=")
	fmt.Println(key, value, found)

	// strings.Builder assembles a big string without copying it on every +=
	var b strings.Builder
	for i := range 3 {
		fmt.Fprintf(&b, "line %d\n", i+1)
	}
	fmt.Print(b.String())

	// a string is bytes; a rune is one character. They differ as soon as you leave ASCII.
	s := "héllo"
	fmt.Println(len(s), utf8.RuneCountInString(s)) // 6 bytes, 5 characters
	for i, r := range s {                          // range walks runes, not bytes
		fmt.Printf("%d:%c ", i, r)
	}
	fmt.Println()
}
