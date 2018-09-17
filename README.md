# goodtimer

Golang timer for humans.

`goodtimer` is a thin wrapper around the standard [time.Timer][1], and it tries to play two roles:

1. A library that helps you use [time.Timer][1] more easily.
2. As well as a demonstration that shows you how to use [time.Timer][1] correctly.


## Installation

```bash
$ go get -u github.com/RussellLuo/goodtimer
```


## Documentation

For usage and examples see the [Godoc][2].


## Why?!

**TL;DR**: The standard [time.Timer][1] is hard to use correctly.

### Timer.Stop

Per the documentation of [Timer.Stop][3], to stop the timer created with NewTimer, you need to check the return value and drain the channel if necessary:

```go
if !t.Stop() {
	<-t.C
}
```

But the draining operation will be blocked if the the program has already received from the Timer's channel before. So someone suggests doing a non-blocking draining:

```go
if !t.Stop() {
	select {
	case <-t.C: // try to drain the channel
	default:
	}
}
```

However, [there is a race condition][4] between draining the channel and sending time into the channel, which may lead to a undrained channel.


### Timer.Reset

To reset a timer, is must have expired or be stopped before. So [Timer.Reset][5] has almost the same issue with [Timer.Stop][3].


### Solutions

Finally, as [Russ Cox][6] suggested ([here][7] and [here][8]), the correct way to use [time.Timer][1] is:

- All the Timer operations ([Timer.Stop][3], [Timer.Reset][5] and receiving from or draining the channel) should be done in the same goroutine.
- The program should manage an extra status showing whether it has received from the Timer's channel or not.


[1]: https://golang.org/pkg/time/#Timer
[2]: https://godoc.org/github.com/RussellLuo/goodtimer 
[3]: https://golang.org/pkg/time/#Timer.Stop
[4]: https://github.com/golang/go/issues/14383#issuecomment-185977844
[5]: https://golang.org/pkg/time/#Timer.Reset
[6]: https://github.com/rsc
[7]: https://github.com/golang/go/issues/11513#issuecomment-157062583
[8]: https://groups.google.com/d/msg/golang-dev/c9UUfASVPoU/tlbK2BpFEwAJ
