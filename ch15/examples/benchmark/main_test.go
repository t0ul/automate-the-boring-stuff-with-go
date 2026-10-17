package main

import (
	"fmt"
	"testing"
)

// go test -bench=. -benchmem ./ch15/examples/benchmark
var parts = func() []string {
	p := make([]string, 2000)
	for i := range p {
		p[i] = fmt.Sprintf("line number %d of the report", i)
	}
	return p
}()

func BenchmarkJoinWithPlus(b *testing.B) {
	for b.Loop() {
		joinWithPlus(parts)
	}
}

func BenchmarkJoinWithBuilder(b *testing.B) {
	for b.Loop() {
		joinWithBuilder(parts)
	}
}

func TestBothAgree(t *testing.T) {
	if joinWithPlus(parts) != joinWithBuilder(parts) {
		t.Fatal("the two implementations differ")
	}
}
