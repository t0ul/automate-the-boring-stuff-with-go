// ex. 1 a map: keys of one type, values of another
package main

import "fmt"

func main() {
	nyc := map[string]int{
		"population": 8000000,
		"size":       302,
	}
	htx := map[string]int{
		"population": 2000000,
		"size":       665,
	}
	fmt.Println(htx["population"], nyc["population"])
}
