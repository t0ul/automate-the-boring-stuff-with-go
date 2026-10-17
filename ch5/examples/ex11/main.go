// ex. 11 range over an integer
package main

import "fmt"

func main() {
	for i := range 3 {
		fmt.Println("round", i+1)
	}
}
