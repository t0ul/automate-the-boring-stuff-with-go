// Ex. 1 Tell me if a number is even or odd
package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	// pick a random whole number from 0 to 99
	n := rand.IntN(100)
	// 'if' condition
	if n%2 == 0 {
		// code block that runs only when the condition is true
		fmt.Printf("Number is even: %d\n", n)
		return // leave main, so the line below does not run
	}
	// reached only when the condition was false
	fmt.Printf("Number is odd: %d\n", n)
}
