// notify gets a message to you: a chat webhook (Slack, Discord, Mattermost and friends
// all accept a small JSON POST), a desktop notification, or plain stdout. It is the
// last line of every scheduled job: "done, 1284 files, 3 errors".
//
//	go run ./ch17/projects/notify -desktop "backup finished"
//	NOTIFY_WEBHOOK=https://hooks.slack.com/services/... go run ./ch17/projects/notify "backup finished"
//	go run ./ch13/projects/every -interval 1h -- sh -c 'go run ./ch11/projects/download-xkcd -n 1 && go run ./ch17/projects/notify "xkcd fetched"'
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

// Payload builds the JSON body a chat webhook expects. Slack-style hooks read "text",
// Discord reads "content"; sending both keeps one program working with either.
func Payload(message string) []byte {
	b, _ := json.Marshal(map[string]string{"text": message, "content": message})
	return b
}

// PostWebhook delivers the message and treats any non-2xx reply as a failure.
func PostWebhook(ctx context.Context, url, message string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(Payload(message)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("webhook: %s", resp.Status)
	}
	return nil
}

// DesktopCommand returns the operating system's notification command, or nil.
func DesktopCommand(goos, title, message string) *exec.Cmd {
	switch goos {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		return exec.Command("osascript", "-e", script)
	case "linux":
		return exec.Command("notify-send", title, message)
	case "windows":
		ps := fmt.Sprintf(`[void][System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms'); `+
			`$n = New-Object System.Windows.Forms.NotifyIcon; $n.Icon = [System.Drawing.SystemIcons]::Information; `+
			`$n.Visible = $true; $n.ShowBalloonTip(5000, %q, %q, 'Info')`, title, message)
		return exec.Command("powershell", "-NoProfile", "-Command", ps)
	}
	return nil
}

func run(args []string) error {
	flags := flag.NewFlagSet("notify", flag.ContinueOnError)
	webhook := flags.String("webhook", os.Getenv("NOTIFY_WEBHOOK"), "chat webhook URL (default $NOTIFY_WEBHOOK)")
	desktop := flags.Bool("desktop", false, "also show a desktop notification")
	title := flags.String("title", "notify", "title for the desktop notification")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() == 0 {
		return errors.New("usage: notify [-webhook URL] [-desktop] [-title T] <message>")
	}
	message := strings.Join(flags.Args(), " ")

	var errs []error
	delivered := 0
	if *webhook != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := PostWebhook(ctx, *webhook, message); err != nil {
			errs = append(errs, err)
		} else {
			delivered++
		}
	}
	if *desktop {
		if cmd := DesktopCommand(runtime.GOOS, *title, message); cmd != nil {
			if err := cmd.Run(); err != nil {
				errs = append(errs, fmt.Errorf("desktop: %w", err))
			} else {
				delivered++
			}
		}
	}
	fmt.Println(message) // stdout always, so a log captures it too
	if delivered == 0 && len(errs) > 0 {
		return errors.Join(errs...)
	}
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, "warning:", e)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
