package main

import (
	"encoding/json"
	"testing"
)

func TestGet(t *testing.T) {
	var v any
	if err := json.Unmarshal([]byte(`{"a":{"b":[10,{"c":"deep"}]},"n":3}`), &v); err != nil {
		t.Fatal(err)
	}
	cases := map[string]any{"n": 3.0, "a.b.0": 10.0, "a.b.1.c": "deep"}
	for path, want := range cases {
		got, err := Get(v, path)
		if err != nil || got != want {
			t.Errorf("Get(%q) = %v, %v; want %v", path, got, err, want)
		}
	}
	for _, bad := range []string{"x", "a.b.7", "a.b.1.c.d", "a.b.one"} {
		if _, err := Get(v, bad); err == nil {
			t.Errorf("Get(%q) should fail", bad)
		}
	}
}

func TestRender(t *testing.T) {
	if s, _ := Render("plain"); s != "plain" {
		t.Errorf("string rendered as %q", s)
	}
	s, _ := Render(map[string]any{"k": 1.0})
	if s != "{\n  \"k\": 1\n}" {
		t.Errorf("object rendered as %q", s)
	}
}
