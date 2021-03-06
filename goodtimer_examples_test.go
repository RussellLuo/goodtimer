package goodtimer_test

import (
	"fmt"
	"time"

	"github.com/RussellLuo/goodtimer"
)

func Example_blockingRead() {
	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Read from the wrapped timer's channel C.
	if tv := gt.ReadC(); !tv.IsZero() {
		fmt.Println("The timer fires")
	}

	// Output:
	// The timer fires
}

func Example_nonBlockingRead() {
	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Read from the wrapped timer's channel C, in a non-blocking way.
	if tv := gt.TryReadC(1 * time.Second); tv.IsZero() {
		fmt.Println("Timed out before the timer firing")
	}

	// Output:
	// Timed out before the timer firing
}

func Example_stop() {
	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Any operations in the current goroutine.

	if gt.Stop() {
		fmt.Println("The timer is stopped before firing")
	}

	// Output:
	// The timer is stopped before firing
}

func Example_reset() {
	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Any operations in the current goroutine.

	gt.Reset(2 * time.Second)

	// Now you can use the timer gt again.
}
