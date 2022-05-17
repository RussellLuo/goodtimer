package goodtimer_test

import (
	"context"
	"testing"
	"time"

	"github.com/RussellLuo/goodtimer"
)

func TestGoodTimer_ReadC(t *testing.T) {
	cases := []struct {
		name string

		// in
		timerDuration time.Duration
		timeout       time.Duration

		// want
		waitingTime  time.Duration
		isZeroResult bool
	}{
		{
			name: "timer 1s wait 1s",

			timerDuration: 1 * time.Second,

			waitingTime:  1 * time.Second,
			isZeroResult: false,
		},
		{
			name: "timer 2s wait 1s",

			timerDuration: 2 * time.Second,

			waitingTime:  2 * time.Second,
			isZeroResult: false,
		},
		{
			name: "timer 1s timeout -1s wait 0s",

			timerDuration: 1 * time.Second,
			timeout:       -1 * time.Second,

			waitingTime:  0 * time.Second,
			isZeroResult: true,
		},
		{
			name: "timer 1s timeout 500ms wait 500ms",

			timerDuration: 1 * time.Second,
			timeout:       500 * time.Millisecond,

			waitingTime:  500 * time.Millisecond,
			isZeroResult: true,
		},
		{
			name: "timer 1s timeout 1.1s wait 1s",

			timerDuration: 1 * time.Second,
			timeout:       1100 * time.Millisecond,

			waitingTime:  1 * time.Second,
			isZeroResult: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gt := goodtimer.NewGoodTimer(time.NewTimer(c.timerDuration))

			ctx := context.Background()
			if c.timeout != 0 {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, c.timeout)
				defer cancel()
			}

			start := time.Now()
			tv := gt.ReadC(ctx)
			elapsed := time.Now().Sub(start)

			err := 10 * time.Millisecond
			if int(elapsed/err) != int(c.waitingTime/err) {
				t.Errorf("waitingTime: Got (%+v0ms) != Want (%+v0ms)", int(elapsed/err), int(c.waitingTime/err))
			}
			if tv.IsZero() != c.isZeroResult {
				t.Errorf("isZeroResult: Got (%+v) != Want (%+v)", tv.IsZero(), c.isZeroResult)
			}
		})
	}
}

func TestGoodTimer_Reset(t *testing.T) {
	cases := []struct {
		// in
		timerDuration    time.Duration
		readCBeforeReset bool
		resetDuration    time.Duration

		// want
		waitingTime  time.Duration
		isZeroResult bool
	}{
		{
			timerDuration:    1 * time.Second,
			readCBeforeReset: true,
			resetDuration:    1500 * time.Millisecond,

			waitingTime:  1500 * time.Millisecond,
			isZeroResult: false,
		},
		{
			timerDuration:    1 * time.Second,
			readCBeforeReset: false,
			resetDuration:    1500 * time.Millisecond,

			waitingTime:  1500 * time.Millisecond,
			isZeroResult: false,
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			gt := goodtimer.NewGoodTimer(time.NewTimer(c.timerDuration))

			if c.readCBeforeReset {
				gt.ReadC(context.Background())
			}
			gt.Reset(c.resetDuration)

			start := time.Now()
			tv := gt.ReadC(context.Background())
			elapsed := time.Now().Sub(start)

			err := 10 * time.Millisecond
			if int(elapsed/err) != int(c.waitingTime/err) {
				t.Errorf("waitingTime: Got (%+v0ms) != Want (%+v0ms)", int(elapsed/err), int(c.waitingTime/err))
			}
			if tv.IsZero() != c.isZeroResult {
				t.Errorf("isZeroResult: Got (%+v) != Want (%+v)", tv.IsZero(), c.isZeroResult)
			}
		})
	}
}

func TestGoodTimer_Stop(t *testing.T) {
	cases := []struct {
		// in
		timerDuration            time.Duration
		readCBeforeStop          bool
		tryReadCTimeoutAfterStop time.Duration

		// want
		tryReadCGotZeroResult bool
	}{
		{
			timerDuration:            1 * time.Second,
			readCBeforeStop:          true,
			tryReadCTimeoutAfterStop: 1 * time.Second,

			tryReadCGotZeroResult: true,
		},
		{
			timerDuration:            1 * time.Second,
			readCBeforeStop:          false,
			tryReadCTimeoutAfterStop: 1 * time.Second,

			tryReadCGotZeroResult: true,
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			gt := goodtimer.NewGoodTimer(time.NewTimer(c.timerDuration))

			if c.readCBeforeStop {
				gt.ReadC(context.Background())
			}
			gt.Stop()

			ctx, cancel := context.WithTimeout(context.Background(), c.tryReadCTimeoutAfterStop)
			defer cancel()

			tv := gt.ReadC(ctx)
			if tv.IsZero() != c.tryReadCGotZeroResult {
				t.Errorf("isZeroResult: Got (%+v) != Want (%+v)", tv.IsZero(), c.tryReadCGotZeroResult)
			}
		})
	}
}
