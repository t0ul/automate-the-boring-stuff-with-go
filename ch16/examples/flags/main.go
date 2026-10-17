// the flag package: typed options, defaults, help text, and what is left over
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	name := flag.String("name", "world", "who to greet")
	times := flag.Int("times", 1, "how many times")
	loud := flag.Bool("loud", false, "shout")
	pause := flag.Duration("pause", 0, "wait between greetings, e.g. 500ms")
	flag.Parse() // reads os.Args; -h or -help prints the usage for free

	greeting := "hello, " + *name
	if *loud {
		greeting = "HELLO, " + *name + "!"
	}
	for range *times {
		fmt.Println(greeting)
		time.Sleep(*pause)
	}
	// whatever was not a flag: positional arguments
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "also received:", flag.Args())
	}
}
