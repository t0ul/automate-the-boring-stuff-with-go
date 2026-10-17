// ex. 3 iteration order is random, so sort the keys when order matters
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	ages := map[string]int{"Zed": 41, "Ada": 36, "Mo": 29}

	for name, age := range ages { // a different order every run
		fmt.Println(name, age)
	}

	for _, name := range slices.Sorted(maps.Keys(ages)) { // always alphabetical
		fmt.Println(name, ages[name])
	}
}
