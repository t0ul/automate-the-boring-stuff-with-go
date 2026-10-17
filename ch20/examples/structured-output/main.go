// structured-output asks the model for JSON in a shape we define, decodes it into a
// struct, and checks every field before believing any of it. When the check fails, the
// model gets one more try with the complaint attached.
//
//	go run ./ch20/examples/structured-output "Picked up 3 kg of apples at the market for 7.50 on Tuesday"
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

// Purchase is the shape we want back. The tags are the contract with the prompt.
type Purchase struct {
	Item     string  `json:"item"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
}

var units = []string{"kg", "g", "l", "ml", "piece", "pack"}

// Validate is the part the model does not get a vote in.
func (p Purchase) Validate() error {
	var problems []string
	if strings.TrimSpace(p.Item) == "" || len(p.Item) > 60 {
		problems = append(problems, "item must be 1 to 60 characters")
	}
	if p.Quantity <= 0 || p.Quantity > 10000 {
		problems = append(problems, "quantity must be between 0 and 10000")
	}
	if !slices.Contains(units, p.Unit) {
		problems = append(problems, "unit must be one of "+strings.Join(units, ", "))
	}
	if p.Price < 0 || p.Price > 100000 {
		problems = append(problems, "price must be between 0 and 100000")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

const system = `Extract a purchase from the user's sentence. Reply with JSON only, exactly this shape:
{"item": "<short noun phrase>", "quantity": <number>, "unit": "<one of: kg, g, l, ml, piece, pack>", "price": <number, total paid>}
Use "piece" when no unit is given. Do not add fields or commentary.`

// Extract asks, validates, and retries once with the validation error fed back.
func Extract(ctx context.Context, c Client, sentence string) (Purchase, error) {
	user := sentence
	var last error
	for attempt := 1; attempt <= 2; attempt++ {
		reply, err := c.Chat(ctx, system, user, true)
		if err != nil {
			return Purchase{}, err
		}
		obj, ok := ExtractJSON(reply)
		if !ok {
			last = fmt.Errorf("no JSON object in reply %q", reply)
		} else {
			var p Purchase
			if err := json.Unmarshal([]byte(obj), &p); err != nil {
				last = fmt.Errorf("bad JSON %q: %w", obj, err)
			} else if err := p.Validate(); err != nil {
				last = fmt.Errorf("invalid: %w (reply was %s)", err, obj)
			} else {
				return p, nil
			}
		}
		user = fmt.Sprintf("%s\n\nYour previous answer was rejected: %v. Reply again with valid JSON only.", sentence, last)
	}
	return Purchase{}, fmt.Errorf("gave up after 2 attempts: %w", last)
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: structured-output <sentence describing a purchase>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	p, err := Extract(ctx, FromEnv(), strings.Join(args, " "))
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", p)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
