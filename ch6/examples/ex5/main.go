// ex. 5 a map of structs, keyed by name
package main

import (
	"fmt"
	"maps"
	"slices"
)

type City struct {
	Population int
	SizeSqMi   float64
}

func main() {
	cities := map[string]City{
		"New York": {Population: 8000000, SizeSqMi: 302},
		"Houston":  {Population: 2000000, SizeSqMi: 665},
		"Austin":   {Population: 980000, SizeSqMi: 326},
	}
	for _, name := range slices.Sorted(maps.Keys(cities)) {
		c := cities[name]
		fmt.Printf("%-9s %8d people, %.0f/sq mi\n", name, c.Population, float64(c.Population)/c.SizeSqMi)
	}
}
