// first match: where does "foo" first appear?
package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`foo.?`)
	fmt.Printf("%q\n", re.FindString("People like to tailgate for the American sport football, not to be confused with football.")) // "foot"
	fmt.Printf("%q\n", re.FindString("this string doesn't have the substring. Bummer..."))                                          // ""
}
