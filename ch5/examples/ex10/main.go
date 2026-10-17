// ex. 10 a sub-slice shares memory with its parent
package main

import (
	"fmt"
	"slices"
)

func main() {
	week := []string{"Mon", "Tue", "Wed", "Thu", "Fri"}

	midweek := week[1:4] // looks separate, but points into the same array
	midweek[0] = "TUESDAY"
	fmt.Println(week) // week changed too

	independent := slices.Clone(week[1:4]) // a real copy
	independent[0] = "tue"
	fmt.Println(week) // unchanged this time
	fmt.Println(independent)
}
