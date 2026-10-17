// ex. 7 range over a slice
package main

import "fmt"

func main() {
	paperPrices := []float64{1.00, 2.00, 4.00, 8.00, 16.00, 32.00, 64.00, 128.00}
	inflated := make([]float64, len(paperPrices))
	for i, v := range paperPrices {
		inflated[i] = v * 1.08
	}
	fmt.Println(inflated)
}
