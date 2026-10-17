package main

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

var sheet = [][]string{
	{"Name", "Email", "Jan", "Feb", "Mar"},
	{"Ada", "ada@x.io", "paid", "paid", "paid"},
	{"Mo", "mo@x.io", "paid", "paid", ""},
	{"Zed", "zed@x.io", "paid", "", ""},
	{"", "", "", "", ""},
}

func TestUnpaid(t *testing.T) {
	month, members, err := Unpaid(sheet)
	if err != nil {
		t.Fatal(err)
	}
	if month != "Mar" || len(members) != 2 || members[0].Name != "Mo" || members[1].Email != "zed@x.io" {
		t.Fatalf("month=%s members=%+v", month, members)
	}
	if _, _, err := Unpaid([][]string{{"Name", "Email", "Jan"}, {"A", "a@x", ""}}); err == nil {
		t.Fatal("expected an error when no month has payments")
	}
}

func TestRemindersActuallySend(t *testing.T) {
	addr, got := startSMTP(t)
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	server := Server{Host: host, Port: port}

	_, members, _ := Unpaid(sheet)
	for _, m := range members {
		subject, body := Reminder(m, "Mar")
		msg, err := BuildMessage("treasurer@x.io", []string{m.Email}, subject, body, nil, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		if err := server.Send(context.Background(), "treasurer@x.io", []string{m.Email}, msg); err != nil {
			t.Fatal(err)
		}
	}
	for range members {
		select {
		case d := <-got:
			if !strings.Contains(d.Data, "Subject: Mar dues unpaid") || !strings.Contains(d.Data, "not paid dues for Mar") {
				t.Errorf("unexpected message:\n%s", d.Data)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("a reminder never arrived")
		}
	}
}
