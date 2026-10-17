// countdown counts down from N seconds, then gets your attention: it plays a sound where
// the operating system has a command-line player, and rings the terminal bell everywhere.
//
//	go run ./ch13/projects/countdown -seconds 5
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// AlarmCommand returns the command that makes a noise on this operating system, or nil
// when there is no safe default. The caller falls back to the terminal bell.
func AlarmCommand(goos string) *exec.Cmd {
	switch goos {
	case "darwin":
		return exec.Command("afplay", "/System/Library/Sounds/Glass.aiff")
	case "linux":
		return exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga")
	case "windows":
		return exec.Command("powershell", "-c", "[console]::beep(880,600)")
	}
	return nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("countdown", flag.ContinueOnError)
	seconds := flags.Int("seconds", 10, "how long to count down")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *seconds < 1 {
		return fmt.Errorf("-seconds must be at least 1")
	}
	for left := *seconds; left > 0; left-- {
		fmt.Printf("\r%3d ", left) // \r returns to the start of the line, so the number updates in place
		time.Sleep(time.Second)
	}
	fmt.Println("\rtime's up!")

	fmt.Print("\a") // the terminal bell
	if cmd := AlarmCommand(runtime.GOOS); cmd != nil {
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "(no sound:", err, ")")
		}
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
