// chat-once sends one question to a model running on your own machine and prints the
// answer. Nothing is sent anywhere else.
//
//	ollama pull qwen2.5:1.5b
//	go run ./ch20/examples/chat-once "In one sentence, what is a goroutine?"
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: chat-once <question>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client := FromEnv()
	start := time.Now()
	answer, err := client.Chat(ctx, "You are a concise assistant. Answer in plain text.", strings.Join(args, " "), false)
	if err != nil {
		return err
	}
	fmt.Println(strings.TrimSpace(answer))
	fmt.Fprintf(os.Stderr, "(%s, %s)\n", client.Model, time.Since(start).Round(100*time.Millisecond))
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
