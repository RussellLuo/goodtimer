package goodtimer_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/goodtimer"
)

func TestGoodTimerReadC(t *testing.T) {
	cases := []struct {
		// in
		timerDuration time.Duration

		// want
		waitingTime  time.Duration
		isZeroResult bool
	}{
		{
			timerDuration: 1 * time.Second,

			waitingTime:  1 * time.Second,
			isZeroResult: false,
		},
		{
			timerDuration: 2 * time.Second,

			waitingTime:  2 * time.Second,
			isZeroResult: false,
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			gt := goodtimer.NewGoodTimer(time.NewTimer(c.timerDuration))

			start := time.Now()
			tv := gt.ReadC()
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

func TestGoodTimerTryReadC(t *testing.T) {
	cases := []struct {
		// in
		timerDuration time.Duration
		timeout       time.Duration

		// want
		waitingTime  time.Duration
		isZeroResult bool
	}{
		{
			timerDuration: 1 * time.Second,
			timeout:       -1 * time.Second,

			waitingTime:  0 * time.Second,
			isZeroResult: true,
		},
		{
			timerDuration: 1 * time.Second,
			timeout:       0 * time.Second,

			waitingTime:  0 * time.Second,
			isZeroResult: true,
		},
		{
			timerDuration: 1 * time.Second,
			timeout:       500 * time.Millisecond,

			waitingTime:  500 * time.Millisecond,
			isZeroResult: true,
		},
		{
			timerDuration: 1 * time.Second,
			timeout:       1100 * time.Millisecond,

			waitingTime:  1 * time.Second,
			isZeroResult: false,
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			gt := goodtimer.NewGoodTimer(time.NewTimer(c.timerDuration))

			start := time.Now()
			tv := gt.TryReadC(c.timeout)
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

func TestGoodTimerReset(t *testing.T) {
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
				gt.ReadC()
			}
			gt.Reset(c.resetDuration)

			start := time.Now()
			tv := gt.ReadC()
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

func TestGoodTimerStop(t *testing.T) {
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
				gt.ReadC()
			}
			gt.Stop()

			tv := gt.TryReadC(c.tryReadCTimeoutAfterStop)
			if tv.IsZero() != c.tryReadCGotZeroResult {
				t.Errorf("isZeroResult: Got (%+v) != Want (%+v)", tv.IsZero(), c.tryReadCGotZeroResult)
			}
		})
	}
}
