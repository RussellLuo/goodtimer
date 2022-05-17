package goodtimer

import (
	"context"
	"time"
)

// GoodTimer wraps the standard time.Timer to provide more user-friendly interfaces.
//
// **NOTE**: All the the functions of GoodTimer *should* be used in the same goroutine.
type GoodTimer struct {
	t    *time.Timer // The actual timer
	read bool        // Whether t.C has already been read from
}

// NewGoodTimer creates an instance of GoodTimer.
func NewGoodTimer(t *time.Timer) *GoodTimer {
	return &GoodTimer{t: t}
}

// ReadC waits until it can read from the wrapped timer's channel C, or ctx's deadline, if any, expires.
// It returns the time value received from the channel C, a zero time value if the channel C has already been read from or if ctx's deadline expires.
func (gt *GoodTimer) ReadC(ctx context.Context) time.Time {
	if gt.read {
		return time.Time{}
	}

	select {
	case tv := <-gt.t.C:
		gt.read = true
		return tv
	case <-ctx.Done():
		return time.Time{}
	}
}

// Reset changes the timer to expire after duration d.
func (gt *GoodTimer) Reset(d time.Duration) {
	gt.Stop()
	gt.t.Reset(d)
	gt.read = false
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already expired or been stopped.
func (gt *GoodTimer) Stop() bool {
	stopped := gt.t.Stop()
	if !stopped && !gt.read {
		// Drain the gt.t.C if it has not been read from already
		<-gt.t.C
	}
	return stopped
}
