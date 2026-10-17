// goroutines run functions at the same time; a WaitGroup waits for them to finish
package main

import (
	"fmt"
	"sync"
	"time"
)

func fetch(name string, took time.Duration) {
	time.Sleep(took) // stands in for a network call
	fmt.Printf("%s done after %s\n", name, took)
}

func main() {
	start := time.Now()

	// one after another: about 600 ms
	fetch("a", 300*time.Millisecond)
	fetch("b", 300*time.Millisecond)
	fmt.Println("sequential:", time.Since(start).Round(100*time.Millisecond))

	// at the same time: about 300 ms
	start = time.Now()
	var wg sync.WaitGroup
	for _, name := range []string{"c", "d", "e"} {
		wg.Go(func() { fetch(name, 300*time.Millisecond) })
	}
	wg.Wait() // blocks until every goroutine started with wg.Go has returned
	fmt.Println("concurrent:", time.Since(start).Round(100*time.Millisecond))
}
