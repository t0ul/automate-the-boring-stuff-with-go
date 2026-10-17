// name-notes suggests a descriptive file name for every text note in a folder by asking
// a local model what each note is about. The model only ever supplies a candidate name;
// a strict pattern decides whether it is acceptable, collisions are resolved by the
// program, and nothing is renamed without -apply.
//
//	go run ./ch20/projects/name-notes ./ch20/projects/name-notes/sample
//	go run ./ch20/projects/name-notes -apply ./ch20/projects/name-notes/sample
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// validName is the only shape a new name may take: two to six lowercase words joined by
// single hyphens. No path separators, no dots, no spaces, nothing surprising.
var validName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+){1,5}$`)

const system = `You name text files. Given the contents of a note, reply with JSON only:
{"name": "<2 to 6 lowercase words joined by hyphens that say what the note is about>"}`

// Suggest asks the model for a name and returns it only if it passes the pattern.
func Suggest(ctx context.Context, c Client, contents string) (string, error) {
	// a note can be long; the first few hundred characters say what it is about
	excerpt := contents
	if len(excerpt) > 600 {
		excerpt = excerpt[:600]
	}
	reply, err := c.Chat(ctx, system, "Note contents:\n"+excerpt, true)
	if err != nil {
		return "", err
	}
	return Accept(reply)
}

// Accept extracts the name from the reply, tidies it the way a person would (lower
// case, hyphens for spaces and underscores, punctuation dropped) and then checks it
// against the pattern. The tidying is deterministic code; the pattern is still the gate.
func Accept(reply string) (string, error) {
	obj, ok := ExtractJSON(reply)
	if !ok {
		return "", fmt.Errorf("no JSON in reply %q", reply)
	}
	var v struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(obj), &v); err != nil {
		return "", fmt.Errorf("bad JSON %q", obj)
	}
	name := Tidy(v.Name)
	if !validName.MatchString(name) {
		return "", fmt.Errorf("suggested name %q does not match %s", v.Name, validName)
	}
	return name, nil
}

// Tidy lower-cases, turns runs of spaces, underscores and hyphens into one hyphen, and
// drops every other character that is not a letter or digit.
func Tidy(s string) string {
	var b strings.Builder
	lastHyphen := true
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastHyphen = false
		case r == ' ' || r == '_' || r == '-':
			if !lastHyphen {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// Rename is one planned change.
type Rename struct {
	From, To string
	Note     string // why there is no To, when there is not
}

// Plan proposes a new name for every .txt and .md file in folder. Names the model gets
// wrong are left as they are, with the reason; duplicate suggestions get a -2, -3 suffix.
func Plan(ctx context.Context, c Client, folder string) ([]Rename, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	taken := map[string]bool{}
	for _, e := range entries {
		taken[e.Name()] = true
	}
	var plan []Rename
	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if e.IsDir() || (ext != ".txt" && ext != ".md") {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(folder, e.Name()))
		if err != nil {
			return nil, err
		}
		name, err := Suggest(ctx, c, string(contents))
		if err != nil {
			plan = append(plan, Rename{From: e.Name(), Note: err.Error()})
			continue
		}
		candidate := name + ext
		for n := 2; taken[candidate]; n++ {
			candidate = fmt.Sprintf("%s-%d%s", name, n, ext)
		}
		taken[candidate] = true
		plan = append(plan, Rename{From: e.Name(), To: candidate})
	}
	return plan, nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("name-notes", flag.ContinueOnError)
	apply := flags.Bool("apply", false, "rename the files (default: print the plan)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: name-notes [-apply] <folder>")
	}
	folder := flags.Arg(0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	plan, err := Plan(ctx, FromEnv(), folder)
	if err != nil {
		return err
	}
	for _, r := range plan {
		if r.To == "" {
			fmt.Printf("%-14s kept (%s)\n", r.From, r.Note)
			continue
		}
		fmt.Printf("%-14s -> %s\n", r.From, r.To)
		if *apply {
			if err := os.Rename(filepath.Join(folder, r.From), filepath.Join(folder, r.To)); err != nil {
				return err
			}
		}
	}
	if !*apply {
		fmt.Println("(dry run; pass -apply to rename)")
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
