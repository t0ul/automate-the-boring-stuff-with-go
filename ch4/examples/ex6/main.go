// ex. 6 returning values of different types
package main

import "fmt"

func nameAndAge(name string, age int) (string, int) { // (1)
	return name, age
}

func main() {
	name, age := nameAndAge("Gopher", 20)
	fmt.Printf("name: %s, age: %d\n", name, age)
}
