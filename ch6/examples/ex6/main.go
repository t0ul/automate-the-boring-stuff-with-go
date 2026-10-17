// ex. 6 counting things is what maps are for
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "the cat sat on the mat and the cat slept"
	counts := map[string]int{}
	for _, word := range strings.Fields(text) {
		counts[word]++ // a missing key starts at 0, so this just works
	}
	fmt.Println(counts)
	fmt.Println("the:", counts["the"])
}
