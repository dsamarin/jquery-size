package main

import (
	"io"
	"sync/atomic"
)

// Counter is a Writer that loses all data written and records the bytes written.
type Counter struct {
	io.Writer
	count uint64
}

// NewCounter creates a new Writer that loses all data written to it
// and saves and increments only by the length of the data written.
func NewCounter() *Counter {
	return &Counter{}
}

func (counter *Counter) Write(buffer []byte) (int, error) {
	size := len(buffer)
	atomic.AddUint64(&counter.count, uint64(size))
	return size, nil
}

// Count returns the number of bytes written.
func (counter *Counter) Count() uint64 {
	return atomic.LoadUint64(&counter.count)
}
