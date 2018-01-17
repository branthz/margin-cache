// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// timer manage for expired keys

package hashmap

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type timer struct {
	i int // heap index
	// Timer wakes up at when, and then at when+period, ... (period > 0 only)
	// each time calling f(arg, now) in the timer goroutine, so f must be
	// a well-behaved function and not block.
	when   int64
	period int64
	f      func(interface{}, string)
	arg    interface{}
	key    string
}

type timerst struct {
	lock         sync.RWMutex
	created      bool
	sleeping     bool
	rescheduling bool
	sleepUntil   int64
	waitnote     note
	waitc        chan int
	t            []*timer
}

var timers *timerst
var timerOnce sync.Once
func setup() {
	timerOnce.Do(func(){
	timers = new(timerst)
	timers.waitc = make(chan int)
	timers.waitnote.slp = make(chan int64)
	timers.waitnote.recv = make(chan int)
	//go timers.waitnote.notewait()
	})
}

func goready(n int) {
	timers.waitc <- n
}

func gopark() {
	<-timers.waitc
}

// stopTimer removes t from the timer heap if it is there.
// It returns true if t was removed, false if t wasn't even there.
func stopTimer(t *timer) bool {
	return deltimer(t)
}

func addtimer(t *timer) {
	timers.lock.Lock()
	addtimerLocked(t)
	timers.lock.Unlock()
}

// Add a timer to the heap and start or kick timerproc if the new timer is
// earlier than any of the others.
// Timers are locked.
func addtimerLocked(t *timer) {
	// when must never be negative; otherwise timerproc will overflow
	// during its delta calculation and never expire other runtime timers.
	if t.when < 0 {
		return
		//t.when = 1<<63 - 1
	}
	t.i = len(timers.t)
	timers.t = append(timers.t, t)
	siftupTimer(t.i)
	if t.i == 0 {
		// siftup moved to top: new earliest deadline.
		if timers.sleeping {
			timers.waitnote.notewakeup()
			timers.sleeping = false
		}
		if timers.rescheduling {
			timers.rescheduling = false
			goready(0)
		}
	}
	if !timers.created {
		timers.created = true
		go timerproc()
	}
}

// Delete timer t from the heap.
// Do not need to update the timerproc: if it wakes up early, no big deal.
func deltimer(t *timer) bool {
	// Dereference t so that any panic happens before the lock is held.
	// Discard result, because t might be moving in the heap.
	_ = t.i

	timers.lock.Lock()
	// t may not be registered anymore and may have
	// a bogus i (typically 0, if generated by Go).
	// Verify it before proceeding.
	i := t.i
	last := len(timers.t) - 1
	if i < 0 || i > last || timers.t[i] != t {
		timers.lock.Unlock()
		return false
	}
	if i != last {
		timers.t[i] = timers.t[last]
		timers.t[i].i = i
	}
	timers.t[last] = nil
	timers.t = timers.t[:last]
	if i != last {
		siftupTimer(i)
		siftdownTimer(i)
	}
	timers.lock.Unlock()
	return true
}

// Timerproc runs the time-driven events.
// It sleeps until the next event in the timers heap.
// If addtimer inserts a new earlier event, it wakes timerproc early.
func timerproc() {
	for {
		timers.lock.Lock()
		timers.sleeping = false
		now := time.Now().UnixNano()
		delta := int64(-1)
		for {
			if len(timers.t) == 0 {
				delta = -1
				break
			}
			t := timers.t[0]
			delta = t.when - now
			if delta > 0 {
				break
			}
			if t.period > 0 {
				// leave in heap but adjust next time to fire
				t.when += t.period * (1 + -delta/t.period)
				siftdownTimer(0)
			} else {
				// remove from heap
				last := len(timers.t) - 1
				if last > 0 {
					timers.t[0] = timers.t[last]
					timers.t[0].i = 0
				}
				timers.t[last] = nil
				timers.t = timers.t[:last]
				if last > 0 {
					siftdownTimer(0)
				}
				t.i = -1 // mark as removed
			}
			f := t.f
			arg := t.arg
			key := t.key
			timers.lock.Unlock()
			f(arg, key)
			timers.lock.Lock()
		}
		if delta < 0 {
			// No timers left - put goroutine to sleep.
			timers.rescheduling = true
			timers.lock.Unlock()
			gopark()
			continue
		}
		// At least one timer pending. Sleep until then.
		timers.sleeping = true
		timers.sleepUntil = now + delta
		//noteclear(&timers.waitnote)
		timers.lock.Unlock()
		timers.waitnote.notetsleepg(delta)
	}
}

type note struct {
	key  uintptr
	slp  chan int64
	recv chan int
}

const noteFree = uintptr(0)
const noteSleep = uintptr(1)

func (n *note) noteclear() {
	n.key = noteFree
}

func (n *note) notewakeup() {
	if atomic.CompareAndSwapUintptr(&n.key, noteSleep, noteFree) {
		n.recv <- 0
	}
}

func (n *note) notewait() {
	var t int64
	for {
		t = <-n.slp
		//fmt.Printf("note wait:----%d\n", t/1000000)
		time.Sleep(time.Duration(t))
		if atomic.CompareAndSwapUintptr(&n.key, noteSleep, noteFree) {
			//fmt.Printf("note timer will end sleeping\n")
			n.recv <- 1
		}
		//fmt.Printf("note wait end,next turn\n")
	}
}

func (n *note) notetsleepg(t int64) {
	if t < 1000 {
		return
	}

	atomic.StoreUintptr(&n.key, noteSleep)
	select {
	case <- n.recv:
		break
	case <- time.After(time.Duration(t)):
		break
	}
}

// Heap maintenance algorithms.
func siftupTimer(i int) {
	t := timers.t
	when := t[i].when
	tmp := t[i]
	for i > 0 {
		p := (i - 1) / 4 // parent
		if when >= t[p].when {
			break
		}
		t[i] = t[p]
		t[i].i = i
		t[p] = tmp
		t[p].i = p
		i = p
	}
}

func siftdownTimer(i int) {
	t := timers.t
	n := len(t)
	when := t[i].when
	tmp := t[i]
	for {
		c := i*4 + 1 // left child
		c3 := c + 2  // mid child
		if c >= n {
			break
		}
		w := t[c].when
		if c+1 < n && t[c+1].when < w {
			w = t[c+1].when
			c++
		}
		if c3 < n {
			w3 := t[c3].when
			if c3+1 < n && t[c3+1].when < w3 {
				w3 = t[c3+1].when
				c3++
			}
			if w3 < w {
				w = w3
				c = c3
			}
		}
		if w >= when {
			break
		}
		t[i] = t[c]
		t[i].i = i
		t[c] = tmp
		t[c].i = c
		i = c
	}
}

var seed = rand.New(rand.NewSource(time.Now().Unix()))

func getRand() int64 {
	return seed.Int63n(200)
}

/*
func main() {
	setup()
	for i := 1; i < 1000; i++ {
		t := new(timer)
		t.when = time.Now().UnixNano() + 1e9*getRand()
		t.f = timerAction
		t.arg = "hello"
		t.seq = uint32(i)
		addtimer(t)
	}
	stop := make(chan bool)
	<-stop
}*/
