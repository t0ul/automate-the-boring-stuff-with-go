// categorize-expenses fills in the category column of an expenses CSV by asking a local
// model to pick one from a fixed list. The model proposes; the program decides: an
// answer outside the list becomes "other", every answer can be logged, repeated items are
// asked once, and nothing is written until you have looked at the table.
//
//	go run ./ch20/projects/categorize-expenses ./ch20/projects/categorize-expenses/uncategorized.csv
//	go run ./ch20/projects/categorize-expenses -write ./ch20/projects/categorize-expenses/uncategorized.csv
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
)

// Categories is the allowlist. The model never adds to it.
var Categories = []string{"food", "transport", "books", "household", "health", "entertainment", "other"}

// Row is one expense before and after.
type Row struct {
	Date, Item, Amount, Category string
}

// Prompt is what the model sees for one item: the allowed answers and the exact JSON
// shape, and nothing about the program, the file, or what happens next.
func Prompt(item, amount string) (system, user string) {
	system = fmt.Sprintf(`You classify personal expenses. Reply with JSON only: {"category": "<one of: %s>"}`, strings.Join(Categories, ", "))
	user = fmt.Sprintf("Item: %q, amount: %s", item, amount)
	return system, user
}

// Parse turns the model's reply into a category from the allowlist, or "other" with a
// reason. It never returns anything that is not in Categories.
func Parse(reply string) (category, reason string) {
	obj, ok := ExtractJSON(reply)
	if !ok {
		return "other", "no JSON in reply"
	}
	var v struct {
		Category string `json:"category"`
	}
	if err := json.Unmarshal([]byte(obj), &v); err != nil {
		return "other", "bad JSON"
	}
	c := strings.ToLower(strings.TrimSpace(v.Category))
	if !slices.Contains(Categories, c) {
		return "other", fmt.Sprintf("%q is not an allowed category", c)
	}
	return c, ""
}

// Classifier asks the model, with a memory so the same item is asked once.
type Classifier struct {
	Client Client
	cache  map[string]string
	Calls  int
}

// Category returns the category for one item, asking the model only the first time.
func (k *Classifier) Category(ctx context.Context, item, amount string) string {
	if k.cache == nil {
		k.cache = map[string]string{}
	}
	key := strings.ToLower(strings.TrimSpace(item))
	if c, ok := k.cache[key]; ok {
		return c
	}
	system, user := Prompt(item, amount)
	k.Calls++
	reply, err := k.Client.Chat(ctx, system, user, true)
	if err != nil {
		slog.Warn("model call failed", "item", item, "err", err)
		k.cache[key] = "other"
		return "other"
	}
	c, reason := Parse(reply)
	slog.Info("classified", "item", item, "category", c, "reply", strings.TrimSpace(reply), "note", reason)
	k.cache[key] = c
	return c
}

// ReadRows reads date,item,amount[,category] CSV, columns in any order.
func ReadRows(r io.Reader) ([]Row, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, errors.New("no data rows")
	}
	col := map[string]int{}
	for i, h := range records[0] {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	for _, need := range []string{"date", "item", "amount"} {
		if _, ok := col[need]; !ok {
			return nil, fmt.Errorf("missing column %q", need)
		}
	}
	var rows []Row
	for _, rec := range records[1:] {
		row := Row{Date: rec[col["date"]], Item: rec[col["item"]], Amount: rec[col["amount"]]}
		if i, ok := col["category"]; ok && i < len(rec) {
			row.Category = rec[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// Fill classifies every row that has no category yet.
func Fill(ctx context.Context, k *Classifier, rows []Row) {
	for i := range rows {
		if rows[i].Category == "" {
			rows[i].Category = k.Category(ctx, rows[i].Item, rows[i].Amount)
		}
	}
}

func writeRows(w io.Writer, rows []Row) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"date", "item", "amount", "category"})
	for _, r := range rows {
		cw.Write([]string{r.Date, r.Item, r.Amount, r.Category})
	}
	cw.Flush()
	return cw.Error()
}

func run(args []string) error {
	flags := flag.NewFlagSet("categorize-expenses", flag.ContinueOnError)
	write := flags.Bool("write", false, "write <name>-categorized.csv (default: print the table)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: categorize-expenses [-write] <expenses.csv>")
	}
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		return err
	}
	rows, err := ReadRows(f)
	f.Close()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	k := &Classifier{Client: FromEnv()}
	start := time.Now()
	Fill(ctx, k, rows)

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "date\titem\tamount\tcategory")
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Date, r.Item, r.Amount, r.Category)
	}
	tw.Flush()
	fmt.Fprintf(os.Stderr, "%d row(s), %d model call(s), %s\n", len(rows), k.Calls, time.Since(start).Round(100*time.Millisecond))

	if *write {
		out := strings.TrimSuffix(flags.Arg(0), ".csv") + "-categorized.csv"
		of, err := os.Create(out)
		if err != nil {
			return err
		}
		if err := writeRows(of, rows); err != nil {
			of.Close()
			return err
		}
		if err := of.Close(); err != nil {
			return err
		}
		fmt.Println("wrote", out)
	}
	return nil
}

func main() {
	level := slog.LevelWarn
	if os.Getenv("LLM_LOG") != "" {
		level = slog.LevelInfo // LLM_LOG=1 shows every prompt's answer
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
