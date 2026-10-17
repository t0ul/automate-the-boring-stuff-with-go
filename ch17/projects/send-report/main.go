// send-report emails a message with attachments through any SMTP server. Server
// settings come from the environment (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS), never
// from the code or the command line, where they would end up in your shell history.
//
//	export SMTP_HOST=smtp.fastmail.com SMTP_PORT=587 SMTP_USER=me@example.com SMTP_PASS=app-password
//	go run ./ch17/projects/send-report -from me@example.com -to boss@example.com -subject "October report" -attach report.xlsx -body "Attached."
//	go run ./ch17/projects/send-report -dry-run -from a@b.c -to d@e.f -subject hi -body hello
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func run(args []string) error {
	flags := flag.NewFlagSet("send-report", flag.ContinueOnError)
	from := flags.String("from", "", "sender address")
	to := flags.String("to", "", "recipient addresses, comma separated")
	subject := flags.String("subject", "", "subject line")
	body := flags.String("body", "", "message text")
	attach := flags.String("attach", "", "files to attach, comma separated")
	dryRun := flags.Bool("dry-run", false, "print the message instead of sending it")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *from == "" || *to == "" || *subject == "" {
		return errors.New("usage: send-report -from A -to B[,C] -subject S [-body text] [-attach f1,f2] [-dry-run]")
	}
	recipients := strings.Split(*to, ",")
	var attachments []string
	if *attach != "" {
		attachments = strings.Split(*attach, ",")
	}
	msg, err := BuildMessage(*from, recipients, *subject, *body, attachments, time.Now())
	if err != nil {
		return err
	}
	if *dryRun {
		os.Stdout.Write(msg)
		return nil
	}
	server, err := FromEnv()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := server.Send(ctx, *from, recipients, msg); err != nil {
		return err
	}
	fmt.Printf("sent %q to %s via %s (%d bytes)\n", *subject, *to, server.Host, len(msg))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
