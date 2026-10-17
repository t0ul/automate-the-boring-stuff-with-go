// a message with an attachment is a multipart message: text part, then the file in base64
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
)

// WithAttachment builds a multipart/mixed message with one text part and one file.
func WithAttachment(from, to, subject, body, filename string, file []byte) ([]byte, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fmt.Fprintf(&buf, "From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\n", from, to, subject)
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", mw.Boundary())

	text, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain; charset=utf-8"}})
	if err != nil {
		return nil, err
	}
	io.WriteString(text, strings.ReplaceAll(body, "\n", "\r\n"))

	att, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"application/octet-stream"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filename)},
	})
	if err != nil {
		return nil, err
	}
	enc := base64.NewEncoder(base64.StdEncoding, att)
	enc.Write(file)
	enc.Close()

	if err := mw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	msg, err := WithAttachment("me@example.com", "you@example.com", "The report", "Attached as promised.", "report.csv", []byte("date,amount\n2026-10-17,42\n"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(string(msg))
}
