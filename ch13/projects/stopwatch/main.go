// stopwatch: press Enter for a lap, ctrl+c (or type q) to stop. Prints each lap's time
// and the running total, then a summary.
//
//	go run ./ch13/projects/stopwatch
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Lap is one press of Enter.
type Lap struct {
	Number int
	Split  time.Duration // since the previous lap
	Total  time.Duration // since the start
}

// Format renders one lap line the way the program prints it.
func (l Lap) Format() string {
	return fmt.Sprintf("Lap #%d: %6.2fs (%6.2fs total)", l.Number, l.Split.Seconds(), l.Total.Seconds())
}

// Run reads lap triggers from in until it closes or a line starting with q arrives, and
// writes lap lines to out. It is written against io.Reader so a test can drive it.
func Run(in io.Reader, out io.Writer, clock func() time.Time) []Lap {
	start := clock()
	last := start
	var laps []Lap
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(sc.Text())), "q") {
			break
		}
		now := clock()
		lap := Lap{Number: len(laps) + 1, Split: now.Sub(last), Total: now.Sub(start)}
		laps = append(laps, lap)
		fmt.Fprintln(out, lap.Format())
		last = now
	}
	return laps
}

func main() {
	fmt.Println("Press Enter for a lap, q then Enter to stop. Starting now.")
	laps := Run(os.Stdin, os.Stdout, time.Now)
	if len(laps) == 0 {
		fmt.Println("no laps")
		return
	}
	fastest := laps[0]
	for _, l := range laps[1:] {
		if l.Split < fastest.Split {
			fastest = l
		}
	}
	fmt.Printf("%d laps, fastest lap #%d at %.2fs\n", len(laps), fastest.Number, fastest.Split.Seconds())
}
