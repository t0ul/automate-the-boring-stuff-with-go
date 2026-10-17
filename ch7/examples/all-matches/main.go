// all matches: every place the pattern appears
package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`a.`)
	fmt.Printf("%q\n", re.FindAllString("An apple in the big apple is known to cause lots of apathy for some reason", -1))
	fmt.Printf("%q\n", re.FindAllString("paranormal", 2)) // at most 2 matches
	fmt.Printf("%q\n", re.FindAllString("graal", -1))
	fmt.Printf("%q\n", re.FindAllString("none", -1)) // []
}
