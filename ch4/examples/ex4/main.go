// ex. 4 arguments of different types
package main

import "fmt"

func HappyBirthday(name string, age int) {
	fmt.Printf("Happy birthday %s, you're turning %d\n", name, age)
}

func main() {
	HappyBirthday("Go", 20)
}
