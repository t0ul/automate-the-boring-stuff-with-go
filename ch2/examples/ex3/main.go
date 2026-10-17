// ex3
package main

import "fmt"

// int -> uint -> float64 if you're coding along
func pay(money, bills float64) float64 {
	return money - bills
}

func main() {
	// a literal like 5000.0 is a float64 unless you say otherwise
	var paycheck, bills float64 = 4000, 5000.0
	fmt.Printf("bank account balance: %.2f\n", pay(paycheck, bills))
}
