package goodtimer_test

import (
	"fmt"
	"time"

	"github.com/RussellLuo/goodtimer"
)

func Example_blockingRead() {
	// import "fmt"
	// import "time"
	// import "github.com/RussellLuo/goodtimer"

	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Read from the wrapped timer's channel C
	if tv := gt.ReadC(); !tv.IsZero() {
		fmt.Println("The timer fires")
	}
}

func Example_nonBlockingRead() {
	// import "fmt"
	// import "time"
	// import "github.com/RussellLuo/goodtimer"

	t := time.NewTimer(2 * time.Second)
	gt := goodtimer.NewGoodTimer(t)

	// Read from the wrapped timer's channel C, in a non-blocking way
	if tv := gt.TryReadC(1 * time.Second); tv.IsZero() {
		fmt.Println("Timed out before the timer fires")
	}
}
