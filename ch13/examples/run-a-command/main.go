// run another program: capture its output, check its exit code, give it a deadline
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	// a program's output as a string
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("go says: %s", out)

	// a failing program: the error carries the exit code
	err = exec.Command("go", "run", "./does-not-exist").Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		fmt.Println("exit code:", exitErr.ExitCode())
	}

	// a deadline: the context kills the command if it runs too long
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "help")
	cmd.Stdout = nil // discard; set to os.Stdout to stream it to your terminal
	if err := cmd.Run(); err != nil {
		fmt.Println("go help:", err)
	} else {
		fmt.Println("go help finished within the deadline")
	}
}
