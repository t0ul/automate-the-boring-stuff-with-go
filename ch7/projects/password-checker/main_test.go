package main

import "testing"

func TestProblems(t *testing.T) {
	cases := []struct {
		password string
		want     int
	}{
		{"CorrectHorse9", 0},
		{"short1A", 1},
		{"alllowercase11", 1},
		{"ALLUPPERCASE11", 1},
		{"NoDigitsHereAtAll", 1},
		{"", 4},
		{"Ünïcödé-Länge-1", 0},
	}
	for _, c := range cases {
		if got := Problems(c.password); len(got) != c.want {
			t.Errorf("Problems(%q) = %v (%d), want %d problems", c.password, got, len(got), c.want)
		}
	}
}
