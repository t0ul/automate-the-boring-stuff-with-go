// ex. 3 the same thing, refactored
package main

import "fmt"

func greetReaders(names []string) {
	for _, name := range names {
		fmt.Printf("Hello %s\n", name)
	}
}

func main() {
	readers := []string{
		"Gopher", "Gopher2", "Gopher3", "Gopher4", "Gopher5",
		"Gopher6", "Gopher7", "Gopher8", "Gopher9", "Gopher10",
	}
	greetReaders(readers)
}
