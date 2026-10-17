package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Server is an outgoing mail server. User and Pass empty means no login (a local relay).
type Server struct {
	Host string
	Port int
	User string
	Pass string
}

// FromEnv reads SMTP_HOST, SMTP_PORT, SMTP_USER and SMTP_PASS.
func FromEnv() (Server, error) {
	s := Server{Host: os.Getenv("SMTP_HOST"), User: os.Getenv("SMTP_USER"), Pass: os.Getenv("SMTP_PASS"), Port: 587}
	if s.Host == "" {
		return s, fmt.Errorf("SMTP_HOST is not set")
	}
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &s.Port); err != nil {
			return s, fmt.Errorf("SMTP_PORT %q: %w", p, err)
		}
	}
	return s, nil
}

// Send delivers msg from one address to the recipients. It upgrades to TLS when the
// server offers STARTTLS, logs in when credentials are set, and refuses to send a
// password over a connection that is not encrypted.
func (s Server) Send(ctx context.Context, from string, to []string, msg []byte) error {
	addr := net.JoinHostPort(s.Host, fmt.Sprint(s.Port))
	conn, err := (&net.Dialer{Timeout: 15 * time.Second}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("connect %s: %w", addr, err)
	}
	c, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return fmt.Errorf("smtp handshake: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: s.Host}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}
	if s.User != "" {
		if _, isTLS := c.TLSConnectionState(); !isTLS {
			return fmt.Errorf("refusing to send the password to %s without TLS", addr)
		}
		if err := c.Auth(smtp.PlainAuth("", s.User, s.Pass, s.Host)); err != nil {
			return fmt.Errorf("login as %s: %w", s.User, err)
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("from %s: %w", from, err)
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("to %s: %w", rcpt, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("server rejected the message: %w", err)
	}
	return c.Quit()
}

// BuildMessage assembles a message with an optional list of attachments.
func BuildMessage(from string, to []string, subject, body string, attachments []string, now time.Time) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMIME-Version: 1.0\r\n",
		from, strings.Join(to, ", "), mime.QEncoding.Encode("utf-8", subject), now.Format(time.RFC1123Z))
	crlfBody := strings.ReplaceAll(body, "\n", "\r\n")

	if len(attachments) == 0 {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		buf.WriteString(crlfBody)
		return buf.Bytes(), nil
	}

	mw := multipart.NewWriter(&buf)
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", mw.Boundary())
	text, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain; charset=utf-8"}})
	if err != nil {
		return nil, err
	}
	io.WriteString(text, crlfBody)
	for _, path := range attachments {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("attachment: %w", err)
		}
		name := filepath.Base(path)
		part, err := mw.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {mime.TypeByExtension(filepath.Ext(name))},
			"Content-Transfer-Encoding": {"base64"},
			"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, name)},
		})
		if err != nil {
			return nil, err
		}
		enc := base64.NewEncoder(base64.StdEncoding, part)
		enc.Write(data)
		enc.Close()
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
