// ex. 4 a struct groups named fields of different types
package main

import "fmt"

// City is a type we define. Capitalised field names are visible to other packages.
type City struct {
	Name       string
	Population int
	SizeSqMi   float64
}

func (c City) Density() float64 {
	return float64(c.Population) / c.SizeSqMi
}

func main() {
	nyc := City{Name: "New York", Population: 8000000, SizeSqMi: 302}
	htx := City{Name: "Houston", Population: 2000000, SizeSqMi: 665}

	fmt.Printf("%s: %.0f people per square mile\n", nyc.Name, nyc.Density())
	fmt.Printf("%+v\n", htx) // %+v prints field names too
}
