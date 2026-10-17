// an email is text: a few header lines, a blank line, the body
package main

import (
	"fmt"
	"strings"
	"time"
)

// PlainMessage builds a plain-text email in the format every mail server expects
// (RFC 5322). Lines end in \r\n, headers come first, one blank line separates the body.
func PlainMessage(from, to, subject, body string, now time.Time) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "Date: %s\r\n", now.Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(strings.ReplaceAll(body, "\n", "\r\n"))
	return []byte(b.String())
}

func main() {
	msg := PlainMessage("me@example.com", "you@example.com", "Lunch?", "Noon at the usual place.\nBring the report.", time.Now())
	fmt.Print(string(msg))
}
