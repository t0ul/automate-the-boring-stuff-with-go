// location: where in the string is the match?
package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`treasure`)
	fmt.Println(re.FindStringIndex("Some maps lead to riches and some don't. But rotate this one, squint, and X marks the treasure.")) // [start end]
	fmt.Println(re.FindStringIndex("there be none here, try another map") == nil)                                                      // true
}
