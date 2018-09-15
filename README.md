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

However, [there is a race condition between draining the channel and sending time into the channel][4].

Finally, the rules of thumb are:

- [Timer.Stop][3] should not be done concurrent to other receives from the Timer's channel.
- The programm should manage an extra status showing whether it has received from the Timer's channel or not.


### Timer.Reset

To reset a timer, is must have expired or be stopped before. So [Timer.Reset][5] has almost the same issue with [Timer.Stop][3].

Some issues and the corresponding suggested solutions:

- [time: Timer.C can still trigger even after Timer.Reset is called][6] ([suggested solution][7])
- [Timer.Reset][8] ([suggested solution][9])


[1]: https://golang.org/pkg/time/#Timer
[2]: https://godoc.org/github.com/RussellLuo/goodtimer 
[3]: https://golang.org/pkg/time/#Timer.Stop
[4]: https://github.com/golang/go/issues/14383#issuecomment-185977844
[5]: https://golang.org/pkg/time/#Timer.Reset
[6]: https://github.com/golang/go/issues/11513
[7]: https://github.com/golang/go/issues/11513#issuecomment-157062583
[8]: https://groups.google.com/forum/#!topic/golang-dev/c9UUfASVPoU
[9]: https://groups.google.com/d/msg/golang-dev/c9UUfASVPoU/tlbK2BpFEwAJ
