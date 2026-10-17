// ex. 2 else
package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	n := rand.IntN(100)
	if n%2 == 0 {
		// runs when the condition is true
		fmt.Printf("Number is even: %d\n", n)
	} else {
		// runs when it is not
		fmt.Printf("Number is odd: %d\n", n)
	}
}
