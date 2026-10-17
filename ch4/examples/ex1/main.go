// ex. 1 parts of a function
package main

import "fmt"

func NameIs(name string) { // (1)
	fmt.Printf("Your name is: %s\n", name)
}

func main() {
	NameIs("Gopher") // (2)
}
