// dues-reminders reads a membership workbook (one row per member, one column per month,
// "paid" where dues came in) and emails everyone who has not paid the latest month.
// It prints what it would send; add -send to actually send through the SMTP server in
// the environment (see send-report).
//
//	go run ./ch17/projects/dues-reminders ./ch17/projects/dues-reminders/members.xlsx
//	go run ./ch17/projects/dues-reminders -send -from treasurer@club.org ./ch17/projects/dues-reminders/members.xlsx
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Member is a row of the sheet.
type Member struct {
	Name  string
	Email string
}

// Unpaid finds the last month column that has any entry and returns the members whose
// cell in that column is not "paid". Columns: Name, Email, then months in order.
func Unpaid(rows [][]string) (month string, members []Member, err error) {
	if len(rows) < 2 || len(rows[0]) < 3 {
		return "", nil, errors.New("need a header row with Name, Email and at least one month")
	}
	header := rows[0]
	latest := -1
	for col := len(header) - 1; col >= 2 && latest < 0; col-- {
		for _, row := range rows[1:] {
			if col < len(row) && strings.TrimSpace(row[col]) != "" {
				latest = col
				break
			}
		}
	}
	if latest < 0 {
		return "", nil, errors.New("no month has any payments recorded")
	}
	for _, row := range rows[1:] {
		if len(row) < 2 || row[0] == "" {
			continue
		}
		paid := latest < len(row) && strings.EqualFold(strings.TrimSpace(row[latest]), "paid")
		if !paid {
			members = append(members, Member{Name: row[0], Email: row[1]})
		}
	}
	return header[latest], members, nil
}

// Reminder is the message text for one member.
func Reminder(m Member, month string) (subject, body string) {
	subject = fmt.Sprintf("%s dues unpaid", month)
	body = fmt.Sprintf("Dear %s,\n\nRecords show that you have not paid dues for %s. Please make this payment as soon as possible.\n\nThank you.", m.Name, month)
	return subject, body
}

func run(args []string) error {
	flags := flag.NewFlagSet("dues-reminders", flag.ContinueOnError)
	send := flags.Bool("send", false, "send the emails (default: print them)")
	from := flags.String("from", "", "sender address (required with -send)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 || (*send && *from == "") {
		return errors.New("usage: dues-reminders [-send -from you@example.com] <members.xlsx>")
	}
	f, err := excelize.OpenFile(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return err
	}
	month, members, err := Unpaid(rows)
	if err != nil {
		return err
	}
	fmt.Printf("%d member(s) unpaid for %s\n", len(members), month)

	var server Server
	if *send {
		if server, err = FromEnv(); err != nil {
			return err
		}
	}
	for _, m := range members {
		subject, body := Reminder(m, month)
		if !*send {
			fmt.Printf("\n--- would send to %s <%s>\nSubject: %s\n%s\n", m.Name, m.Email, subject, body)
			continue
		}
		msg, err := BuildMessage(*from, []string{m.Email}, subject, body, nil, time.Now())
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err = server.Send(ctx, *from, []string{m.Email}, msg)
		cancel()
		if err != nil {
			return fmt.Errorf("%s: %w", m.Email, err)
		}
		fmt.Printf("sent to %s\n", m.Email)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
