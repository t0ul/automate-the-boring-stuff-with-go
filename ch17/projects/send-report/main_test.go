package main

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBuildMessagePlain(t *testing.T) {
	msg, err := BuildMessage("a@x.io", []string{"b@x.io"}, "Hi there", "line one\nline two", nil, time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	s := string(msg)
	for _, want := range []string{"From: a@x.io\r\n", "To: b@x.io\r\n", "Subject: Hi there\r\n", "Content-Type: text/plain; charset=utf-8\r\n\r\nline one\r\nline two"} {
		if !strings.Contains(s, want) {
			t.Errorf("message lacks %q:\n%s", want, s)
		}
	}
}

func TestSendWithAttachment(t *testing.T) {
	addr, got := startSMTP(t)
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)

	file := filepath.Join(t.TempDir(), "report.csv")
	if err := os.WriteFile(file, []byte("date,amount\n2026-10-17,42\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	msg, err := BuildMessage("me@x.io", []string{"you@x.io", "them@x.io"}, "Report", "Attached.", []string{file}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := (Server{Host: host, Port: port}).Send(context.Background(), "me@x.io", []string{"you@x.io", "them@x.io"}, msg); err != nil {
		t.Fatal(err)
	}
	select {
	case d := <-got:
		if d.From != "me@x.io" || len(d.To) != 2 {
			t.Errorf("envelope = %+v", d)
		}
		if !strings.Contains(d.Data, `filename="report.csv"`) || !strings.Contains(d.Data, "ZGF0ZSxhbW91bnQK") { // base64 of "date,amount\n"
			t.Errorf("attachment missing from delivered data:\n%s", d.Data)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server received nothing")
	}
}

func TestSendRefusesPasswordWithoutTLS(t *testing.T) {
	addr, _ := startSMTP(t)
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	err := (Server{Host: host, Port: port, User: "me", Pass: "secret"}).Send(context.Background(), "me@x.io", []string{"you@x.io"}, []byte("Subject: x\r\n\r\nhi"))
	if err == nil || !strings.Contains(err.Error(), "without TLS") {
		t.Fatalf("expected a refusal, got %v", err)
	}
}
