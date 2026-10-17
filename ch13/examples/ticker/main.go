// do something every so often: Sleep for a one-off pause, Ticker for a rhythm
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("pausing for half a second")
	time.Sleep(500 * time.Millisecond)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // always stop a ticker you are done with
	deadline := time.After(1100 * time.Millisecond)

	for {
		select { // whichever happens first
		case t := <-ticker.C:
			fmt.Println("tick", t.Format("15:04:05.000"))
		case <-deadline:
			fmt.Println("done")
			return
		}
	}
}
