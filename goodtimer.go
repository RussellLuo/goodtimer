package goodtimer

import (
	"time"
)

// GoodTimer wraps the standard time.Timer to provide more user-friendly interfaces.
// NOTE: All the the functions of GoodTimer *should* be used in the same goroutine.
type GoodTimer struct {
	t    *time.Timer // The actual timer
	read bool        // Whether t.C has already been read from
}

// NewGoodTimer creates an instance of GoodTimer.
func NewGoodTimer(t *time.Timer) *GoodTimer {
	return &GoodTimer{t: t}
}

// ReadC read and return the time value from the wrapped timer's channel C.
// It returns a zero time value if the wrapped timer's channel C has already been read from.
func (gt *GoodTimer) ReadC() time.Time {
	if gt.read {
		return time.Time{}
	}
	tv := <-gt.t.C
	gt.read = true
	return tv
}

// TryReadC is a non-blocking version of ReadC, and it will wait for at most the duration d.
// Like ReadC, it returns a zero time value if the wrapped timer's channel C has already been read from.
func (gt *GoodTimer) TryReadC(timeout time.Duration) time.Time {
	if gt.read {
		return time.Time{}
	}
	select {
	case tv := <-gt.t.C:
		gt.read = true
		return tv
	case <-time.After(timeout):
		return time.Time{}
	}
}

// Reset changes the timer to expire after duration d.
//
// Implementation reference:
//     - https://github.com/golang/go/issues/11513#issuecomment-157062583
//     - https://groups.google.com/d/msg/golang-dev/c9UUfASVPoU/tlbK2BpFEwAJ
func (gt *GoodTimer) Reset(d time.Duration) {
	gt.Stop()
	gt.t.Reset(d)
	gt.read = false
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already expired or been stopped.
//
// Implementation reference:
//     - https://golang.org/pkg/time/#Timer.Stop
func (gt *GoodTimer) Stop() bool {
	stopped := gt.t.Stop()
	if !stopped && !gt.read {
		// Drain the gt.t.C if it has not been read from already
		<-gt.t.C
	}
	return stopped
}
