// ex. 8 locals are preferred
package main

import "fmt"

var Author = "Gopher"
var Age = 20

func LocalAuthor(name string, age int) {
	n := name // local copies; the globals are untouched
	a := age
	fmt.Printf("Author: %s, Age: %d\n", n, a)
}

func main() {
	LocalAuthor("Gopher2", 40)
	fmt.Printf("Author: %s, Age: %d\n", Author, Age)
}
