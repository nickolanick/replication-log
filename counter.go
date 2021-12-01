package main

import (
	"sync/atomic"
)

type Counter struct {
	// total ordering
	counter int32
}

func (c *Counter) get() int {
	c.counter = atomic.AddInt32(&c.counter, 1)

	return int(c.counter)
}

var counter = Counter{0}
