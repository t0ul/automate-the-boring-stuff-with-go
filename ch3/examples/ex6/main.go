// ex 6. If you're a reader, then a thousand thanks
package main

import "fmt"

func main() {
	fmt.Println("Enter your name")
	var n string
	fmt.Scan(&n)
	switch {
	case n == "Toul":
		fmt.Printf("The author is %s\n", n)
	default:
		// range over an integer counts 0, 1, 2 ... 999
		for i := range 1000 {
			fmt.Printf("i: %d Thank you %s for reading\n", i, n)
		}
	}
}
