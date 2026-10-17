// errgroup: run tasks concurrently, stop all of them at the first error, keep that error
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func work(ctx context.Context, name string, fail bool) error {
	for i := range 5 {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: cancelled at step %d\n", name, i)
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		if fail && i == 1 {
			return errors.New(name + ": disk full")
		}
	}
	fmt.Printf("%s: finished\n", name)
	return nil
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(2) // at most two tasks at once
	g.Go(func() error { return work(ctx, "alpha", false) })
	g.Go(func() error { return work(ctx, "beta", true) })
	g.Go(func() error { return work(ctx, "gamma", false) })

	if err := g.Wait(); err != nil {
		fmt.Println("first error:", err)
	}
}
