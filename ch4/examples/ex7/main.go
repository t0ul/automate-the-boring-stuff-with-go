// ex. 7 globals are tricky
package main

import "fmt"

// package-level ("global") variables
var Author = "Gopher"
var Age = 20

func LocalAuthor(name string, age int) {
	Author = name // this changes the global for everyone
	Age = age
	fmt.Printf("Author: %s, Age: %d\n", name, age)
}

func main() {
	LocalAuthor("Gopher2", 40)
	fmt.Printf("Author: %s, Age: %d\n", Author, Age) // not what the globals said at the top
}
