// todo is a complete command-line tool: subcommands, flags per subcommand, a help text
// embedded in the binary, a version stamped in at build time, proper exit codes, and
// state in a JSON file in the right place for the operating system.
//
//	go run ./ch16/projects/todo add buy milk
//	go run ./ch16/projects/todo list
//	go run ./ch16/projects/todo done 1
package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed usage.txt
var usage string

// version is overwritten at build time:
//
//	go build -ldflags "-X main.version=1.0.0" -o todo ./ch16/projects/todo
var version = "dev"

var errUsage = errors.New("usage")

// run executes one command line against the store at path and writes to out.
func run(args []string, path string, out io.Writer, now func() time.Time) error {
	if len(args) == 0 {
		return errUsage
	}
	cmd, rest := args[0], args[1:]
	if cmd == "version" {
		fmt.Fprintln(out, "todo", version)
		return nil
	}
	store, err := Load(path)
	if err != nil {
		return err
	}
	switch cmd {
	case "add":
		text := strings.TrimSpace(strings.Join(rest, " "))
		if text == "" {
			return fmt.Errorf("add: nothing to add: %w", errUsage)
		}
		n := store.Add(text, now())
		fmt.Fprintf(out, "added #%d: %s\n", n, text)
		return store.Save()
	case "list":
		flags := flag.NewFlagSet("list", flag.ContinueOnError)
		flags.SetOutput(out)
		all := flags.Bool("all", false, "include done items")
		if err := flags.Parse(rest); err != nil {
			return err
		}
		shown := 0
		for i, it := range store.Items {
			if it.Done != nil && !*all {
				continue
			}
			mark := " "
			if it.Done != nil {
				mark = "x"
			}
			fmt.Fprintf(out, "%3d [%s] %s\n", i+1, mark, it.Text)
			shown++
		}
		if shown == 0 {
			fmt.Fprintln(out, "nothing to do")
		}
		return nil
	case "done", "rm":
		if len(rest) != 1 {
			return fmt.Errorf("%s: needs one item number: %w", cmd, errUsage)
		}
		n, err := strconv.Atoi(rest[0])
		if err != nil {
			return fmt.Errorf("%s: %q is not a number: %w", cmd, rest[0], errUsage)
		}
		if cmd == "done" {
			err = store.MarkDone(n, now())
		} else {
			err = store.Remove(n)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s #%d\n", cmd, n)
		return store.Save()
	default:
		return fmt.Errorf("unknown command %q: %w", cmd, errUsage)
	}
}

func main() {
	path, err := DefaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	err = run(os.Args[1:], path, os.Stdout, time.Now)
	switch {
	case err == nil:
	case errors.Is(err, errUsage):
		if err != errUsage { // a wrapped usage error says what was wrong; the bare one just wants the help text
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	default:
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
