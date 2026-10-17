// quiz-builder writes multiple-choice quizzes about U.S. state capitals, each with the
// questions and the answer options in a different random order, plus an answer key.
// Change the capitals table and you have a quiz generator for anything.
//
//	go run ./ch7/projects/quiz-builder -n 3 -out out
package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var capitals = map[string]string{
	"Alabama": "Montgomery", "Alaska": "Juneau", "Arizona": "Phoenix", "Arkansas": "Little Rock",
	"California": "Sacramento", "Colorado": "Denver", "Connecticut": "Hartford", "Delaware": "Dover",
	"Florida": "Tallahassee", "Georgia": "Atlanta", "Hawaii": "Honolulu", "Idaho": "Boise",
	"Illinois": "Springfield", "Indiana": "Indianapolis", "Iowa": "Des Moines", "Kansas": "Topeka",
	"Kentucky": "Frankfort", "Louisiana": "Baton Rouge", "Maine": "Augusta", "Maryland": "Annapolis",
	"Massachusetts": "Boston", "Michigan": "Lansing", "Minnesota": "Saint Paul", "Mississippi": "Jackson",
	"Missouri": "Jefferson City", "Montana": "Helena", "Nebraska": "Lincoln", "Nevada": "Carson City",
	"New Hampshire": "Concord", "New Jersey": "Trenton", "New Mexico": "Santa Fe", "New York": "Albany",
	"North Carolina": "Raleigh", "North Dakota": "Bismarck", "Ohio": "Columbus", "Oklahoma": "Oklahoma City",
	"Oregon": "Salem", "Pennsylvania": "Harrisburg", "Rhode Island": "Providence", "South Carolina": "Columbia",
	"South Dakota": "Pierre", "Tennessee": "Nashville", "Texas": "Austin", "Utah": "Salt Lake City",
	"Vermont": "Montpelier", "Virginia": "Richmond", "Washington": "Olympia", "West Virginia": "Charleston",
	"Wisconsin": "Madison", "Wyoming": "Cheyenne",
}

// Question is one quiz item: the prompt, four options and the index of the right one.
type Question struct {
	State   string
	Options [4]string
	Answer  int
}

// BuildQuiz returns every entry of answers as a question, in random order, with the
// correct capital and three wrong ones shuffled among the options.
func BuildQuiz(r *rand.Rand, answers map[string]string) []Question {
	states := slices.Sorted(mapsKeys(answers))
	r.Shuffle(len(states), func(i, j int) { states[i], states[j] = states[j], states[i] })

	allCapitals := make([]string, 0, len(answers))
	for _, c := range answers {
		allCapitals = append(allCapitals, c)
	}
	slices.Sort(allCapitals)

	quiz := make([]Question, 0, len(states))
	for _, state := range states {
		correct := answers[state]
		wrong := make([]string, 0, 3)
		for len(wrong) < 3 {
			c := allCapitals[r.IntN(len(allCapitals))]
			if c != correct && !slices.Contains(wrong, c) {
				wrong = append(wrong, c)
			}
		}
		q := Question{State: state, Answer: r.IntN(4)}
		w := 0
		for i := range q.Options {
			if i == q.Answer {
				q.Options[i] = correct
			} else {
				q.Options[i] = wrong[w]
				w++
			}
		}
		quiz = append(quiz, q)
	}
	return quiz
}

func mapsKeys(m map[string]string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Render returns the quiz text and the answer key text for quiz number n.
func Render(n int, quiz []Question) (paper, key string) {
	var p, k strings.Builder
	fmt.Fprintf(&p, "Name:\n\nDate:\n\nPeriod:\n\n%40s State Capitals Quiz (Form %d)\n\n", "", n)
	fmt.Fprintf(&k, "State Capitals Quiz (Form %d): answer key\n\n", n)
	for i, q := range quiz {
		fmt.Fprintf(&p, "%d. What is the capital of %s?\n", i+1, q.State)
		for j, opt := range q.Options {
			fmt.Fprintf(&p, "    %c. %s\n", 'A'+j, opt)
		}
		p.WriteString("\n")
		fmt.Fprintf(&k, "%d. %c\n", i+1, 'A'+q.Answer)
	}
	return p.String(), k.String()
}

func run(args []string) error {
	flags := flag.NewFlagSet("quiz-builder", flag.ContinueOnError)
	count := flags.Int("n", 3, "how many quiz forms to write")
	outDir := flags.String("out", "out", "directory for the quiz files")
	seed := flags.Uint64("seed", 0, "random seed (0 = different every run)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", *outDir, err)
	}
	if *seed == 0 {
		*seed = rand.Uint64()
	}
	r := rand.New(rand.NewPCG(*seed, *seed))
	for n := 1; n <= *count; n++ {
		paper, key := Render(n, BuildQuiz(r, capitals))
		quizPath := filepath.Join(*outDir, fmt.Sprintf("capitals_quiz_%d.txt", n))
		keyPath := filepath.Join(*outDir, fmt.Sprintf("capitals_quiz_%d_answers.txt", n))
		if err := os.WriteFile(quizPath, []byte(paper), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", quizPath, err)
		}
		if err := os.WriteFile(keyPath, []byte(key), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", keyPath, err)
		}
		fmt.Println("wrote", quizPath, "and", keyPath)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
