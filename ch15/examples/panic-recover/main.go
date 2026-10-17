// what a panic is, what it prints, and the one place recover belongs
package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

// process stands in for work on one item that might blow up on bad data.
func process(item []int) int {
	return item[0] * 2 // panics when item is empty: index out of range
}

// safely runs fn and turns a panic into an error, so one bad item does not take down
// the loop that is processing a thousand of them. This is the legitimate use of recover.
func safely(fn func() int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v\n%s", r, debug.Stack())
		}
	}()
	return fn(), nil
}

func main() {
	items := [][]int{{1}, {}, {3}}
	for i, item := range items {
		result, err := safely(func() int { return process(item) })
		if err != nil {
			fmt.Fprintf(os.Stderr, "item %d failed: %v\n", i, err)
			continue
		}
		fmt.Printf("item %d -> %d\n", i, result)
	}
	fmt.Println("the loop survived")
}
