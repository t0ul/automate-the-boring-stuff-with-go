// ex. 5 returning values
package main

import "fmt"

func EvenOdd(numbers []int) ([]int, []int) { // (1)
	var odd, even []int
	for _, n := range numbers {
		if n%2 == 0 {
			even = append(even, n)
		} else {
			odd = append(odd, n)
		}
	}
	return odd, even // (2)
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	odd, even := EvenOdd(numbers) // (3)
	fmt.Println("Odd: ", odd)
	fmt.Println("Even:", even)
}
