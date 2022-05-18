package goodtimer_test

import (
	"context"
	"fmt"
	"time"

	"github.com/RussellLuo/goodtimer"
)

func Example_readC() {
	gt := goodtimer.New(2 * time.Second)

	// Wait until we read from the wrapped timer's channel C.
	ctx := context.Background()
	if tv := gt.ReadC(ctx); !tv.IsZero() {
		fmt.Println("The timer fires")
	}

	// Output:
	// The timer fires
}

func Example_readCWithTimeout() {
	gt := goodtimer.New(2 * time.Second)

	// Wait up to 1 second to read from the wrapped timer's channel C.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if tv := gt.ReadC(ctx); tv.IsZero() {
		fmt.Println("Timed out before the timer fires")
	}

	// Output:
	// Timed out before the timer fires
}

func Example_stop() {
	gt := goodtimer.New(2 * time.Second)

	// Any operations in the current goroutine.

	if gt.Stop() {
		fmt.Println("The timer is stopped before firing")
	}

	// Output:
	// The timer is stopped before firing
}

func Example_reset() {
	gt := goodtimer.New(2 * time.Second)

	// Any operations in the current goroutine.

	gt.Reset(2 * time.Second)

	// Now you can use the timer gt again.
}
