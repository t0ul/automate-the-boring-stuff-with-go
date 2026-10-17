// now, formatting, parsing, and arithmetic with time
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("now:", now)

	// Formatting uses a reference time, not %Y-%m-%d: Mon Jan 2 15:04:05 MST 2006.
	// You write that moment in the layout you want, and Go copies the shape.
	fmt.Println(now.Format("2006-01-02"))
	fmt.Println(now.Format("Monday, 2 January 2006 at 3:04 PM"))
	fmt.Println(now.Format(time.RFC3339)) // the layout to use when another program will read it

	// Parsing is the same layout idea in reverse
	due, err := time.Parse("2006-01-02 15:04", "2026-12-24 18:30")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// A Duration is a count of nanoseconds with a readable String
	until := time.Until(due)
	fmt.Printf("%s until %s (%.0f days)\n", until.Round(time.Minute), due.Format("Jan 2"), until.Hours()/24)

	// Add, Sub, Before, After, Truncate
	fmt.Println(due.Add(48 * time.Hour).Weekday())
	fmt.Println(due.Sub(now) > 0, due.After(now))
	fmt.Println(now.Truncate(time.Hour).Format(time.Kitchen))
}
