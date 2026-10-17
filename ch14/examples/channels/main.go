// channels carry values between goroutines; this is the worker-pool shape
package main

import (
	"fmt"
	"sync"
)

func main() {
	jobs := make(chan int)
	results := make(chan string)

	// three workers read jobs until the channel is closed
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Go(func() {
			for n := range jobs {
				results <- fmt.Sprintf("worker %d squared %d = %d", w, n, n*n)
			}
		})
	}

	// one goroutine feeds the jobs and closes the channel when there are no more
	go func() {
		for n := 1; n <= 6; n++ {
			jobs <- n
		}
		close(jobs)
	}()

	// close results once every worker is finished, so the range below ends
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r)
	}
}
