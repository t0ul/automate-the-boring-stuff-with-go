package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

// Delivery is what the test server received for one message.
type Delivery struct {
	From string
	To   []string
	Data string
}

// startSMTP runs a minimal SMTP server on a random local port for the duration of the
// test. It speaks enough of the protocol for net/smtp to deliver a message, and it
// deliberately offers neither STARTTLS nor AUTH, like a local relay would.
func startSMTP(t *testing.T) (addr string, got chan Delivery) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	got = make(chan Delivery, 4)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go serve(conn, got)
		}
	}()
	return ln.Addr().String(), got
}

func serve(conn net.Conn, got chan Delivery) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	say := func(s string) { conn.Write([]byte(s + "\r\n")) }
	say("220 smtptest ESMTP")
	var d Delivery
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		cmd := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
			say("250-smtptest")
			say("250 SIZE 10485760")
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			d.From = strings.Trim(line[len("MAIL FROM:"):], "<> ")
			say("250 OK")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			d.To = append(d.To, strings.Trim(line[len("RCPT TO:"):], "<> "))
			say("250 OK")
		case cmd == "DATA":
			say("354 End data with <CR><LF>.<CR><LF>")
			var body strings.Builder
			for {
				l, err := r.ReadString('\n')
				if err != nil {
					return
				}
				if l == ".\r\n" {
					break
				}
				body.WriteString(l)
			}
			d.Data = body.String()
			got <- d
			d = Delivery{}
			say("250 OK queued")
		case cmd == "QUIT":
			say("221 Bye")
			return
		default:
			say("250 OK")
		}
	}
}
