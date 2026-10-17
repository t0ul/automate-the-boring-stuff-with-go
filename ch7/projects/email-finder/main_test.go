package main

import (
	"reflect"
	"testing"
)

func TestFindEmails(t *testing.T) {
	text := `Contact Ada <ada@example.com> or ops@mail.example.org.
Duplicate: ADA@example.com. Not an address: foo@bar, @nope.com, x@y.z`
	got := FindEmails(text)
	want := []string{"ada@example.com", "ops@mail.example.org"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindEmails = %v, want %v", got, want)
	}
}

func TestFindEmailsNone(t *testing.T) {
	if got := FindEmails("nothing here"); len(got) != 0 {
		t.Fatalf("expected none, got %v", got)
	}
}
