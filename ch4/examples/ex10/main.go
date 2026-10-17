// ex. 10 type parameters, for when the same function should work on several types
package main

import "fmt"

// Number is a constraint: any of these types is allowed.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum works for a slice of any Number. T is a type parameter.
func Sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// Map applies f to every element and returns the results. The two type parameters
// are independent: a []string can become a []int.
func Map[T, U any](values []T, f func(T) U) []U {
	out := make([]U, 0, len(values))
	for _, v := range values {
		out = append(out, f(v))
	}
	return out
}

func main() {
	fmt.Println(Sum([]int{1, 2, 3}))
	fmt.Println(Sum([]float64{1.5, 2.25}))
	lengths := Map([]string{"go", "gopher", "generics"}, func(s string) int { return len(s) })
	fmt.Println(lengths)
}
