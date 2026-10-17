// ex. 3 else if
package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	n := rand.IntN(100)
	if n%2 == 0 {
		fmt.Printf("Number is even: %d\n", n)
	} else if n%2 == 1 {
		fmt.Printf("Number is odd: %d\n", n)
	}
}
