// ex10
package main

import (
	"fmt"
	"strconv"
)

func main() {
	strVar := "100"
	intVar, err := strconv.Atoi(strVar)
	if err != nil {
		fmt.Println("not a number:", err)
		return
	}
	// %T prints the type of a value
	fmt.Printf("%v is a %T\n", intVar, intVar)
}
