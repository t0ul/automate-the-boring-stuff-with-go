// two goroutines touching one variable is a data race unless you say who goes first
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// wrong: the ++ is read-modify-write, and the goroutines interleave.
	// Run with `go run -race .` and the race detector will say so.
	unsafe := 0
	for range 1000 {
		wg.Go(func() { unsafe++ })
	}
	wg.Wait()
	fmt.Println("without a mutex:", unsafe, "(often less than 1000)")

	// right: a mutex makes the increment one step at a time
	var mu sync.Mutex
	safe := 0
	for range 1000 {
		wg.Go(func() {
			mu.Lock()
			safe++
			mu.Unlock()
		})
	}
	wg.Wait()
	fmt.Println("with a mutex:   ", safe)
}
