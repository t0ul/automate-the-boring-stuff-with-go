// two ways to build a big string; the benchmark says which is faster and by how much
package main

import (
	"fmt"
	"strings"
)

func joinWithPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p + "\n"
	}
	return s
}

func joinWithBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	parts := []string{"one", "two", "three"}
	fmt.Print(joinWithPlus(parts), joinWithBuilder(parts))
}
