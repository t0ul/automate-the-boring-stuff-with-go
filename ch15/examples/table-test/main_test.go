package main

import "testing"

// A table-driven test: one slice of cases, one loop, one failure message that names the
// case. Adding a case is adding a line.
func TestSlugify(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"punctuation", "Ch. 15: Testing!", "ch-15-testing"},
		{"runs of separators", "a  --  b", "a-b"},
		{"leading and trailing", " -x- ", "x"},
		{"unicode letters", "Café Olé", "café-olé"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Slugify(c.in); got != c.want {
				t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
