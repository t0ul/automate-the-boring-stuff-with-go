package main

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

func TestBuildQuiz(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	quiz := BuildQuiz(r, capitals)
	if len(quiz) != len(capitals) {
		t.Fatalf("got %d questions, want %d", len(quiz), len(capitals))
	}
	seen := map[string]bool{}
	for _, q := range quiz {
		if seen[q.State] {
			t.Errorf("state %s asked twice", q.State)
		}
		seen[q.State] = true
		if q.Options[q.Answer] != capitals[q.State] {
			t.Errorf("%s: answer index points at %q, want %q", q.State, q.Options[q.Answer], capitals[q.State])
		}
		opts := q.Options[:]
		if len(slices.Compact(slices.Sorted(slices.Values(opts)))) != 4 {
			t.Errorf("%s: options not distinct: %v", q.State, q.Options)
		}
	}
}

func TestBuildQuizIsReproducible(t *testing.T) {
	a := BuildQuiz(rand.New(rand.NewPCG(7, 7)), capitals)
	b := BuildQuiz(rand.New(rand.NewPCG(7, 7)), capitals)
	if !slices.Equal(a, b) {
		t.Fatal("same seed produced different quizzes")
	}
}

func TestRender(t *testing.T) {
	quiz := BuildQuiz(rand.New(rand.NewPCG(3, 3)), capitals)
	paper, key := Render(2, quiz)
	if !strings.Contains(paper, "Form 2") || strings.Count(paper, "What is the capital of") != 50 {
		t.Fatalf("paper malformed:\n%s", paper[:200])
	}
	if strings.Count(key, "\n") < 50 {
		t.Fatalf("key too short:\n%s", key)
	}
}
