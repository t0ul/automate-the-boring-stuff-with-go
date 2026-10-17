// ex. 2 add, look up, test for presence, delete
package main

import "fmt"

func main() {
	stock := map[string]int{"apples": 12, "pears": 3}

	stock["plums"] = 20 // add
	stock["pears"]++    // update

	// a missing key gives the zero value, which is why the comma-ok form exists
	fmt.Println(stock["bananas"])
	n, ok := stock["bananas"]
	fmt.Println(n, ok)

	delete(stock, "apples")
	fmt.Println(len(stock), stock)
}
