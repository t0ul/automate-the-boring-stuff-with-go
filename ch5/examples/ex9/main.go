// ex. 9 the slices package and the min/max builtins
package main

import (
	"fmt"
	"slices"
)

func main() {
	temps := []int{71, 64, 80, 58, 77}

	fmt.Println(slices.Contains(temps, 80)) // true
	fmt.Println(slices.Index(temps, 58))    // 3
	fmt.Println(slices.Max(temps), slices.Min(temps))

	sorted := slices.Clone(temps) // Sort works in place, so sort a copy
	slices.Sort(sorted)
	fmt.Println(sorted)

	slices.Reverse(sorted)
	fmt.Println(sorted)

	// min and max are builtins for two or more values of the same type
	fmt.Println(min(3, 9, 1), max(2.5, 1.5))
}
