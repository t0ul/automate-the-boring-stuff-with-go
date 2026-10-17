// ex. 1 slicing an array
package main

import "fmt"

func main() {
	evens := [5]int{2, 4, 6, 8, 10}
	s := evens[1:4] // from index 1 up to, but not including, index 4
	fmt.Println(s)
}
