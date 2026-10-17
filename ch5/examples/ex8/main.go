// ex. 8 the same loop, the long way
package main

import "fmt"

func main() {
	paperPrices := []float64{1.00, 2.00, 4.00, 8.00, 16.00, 32.00, 64.00, 128.00}
	inflated := make([]float64, len(paperPrices))
	for i := 0; i < len(paperPrices); i++ {
		inflated[i] = paperPrices[i] * 1.08
	}
	fmt.Println(inflated)
}
