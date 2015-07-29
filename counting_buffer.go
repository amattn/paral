package main

import (
	"bytes"
	"io"
	"sync"
)

// This could be split out into its own package at some point.
// I've needed similar functionality in other apps.

// Counts the total number of bytes written since creation or the most recent ClearCount()
// Reset() and Truncate() methods do NOT affect the running counts.
// Currently UnreadByte and UnreadRune do not affect counts, but probably should.
type CountingBuffer struct {
	totalIn  uint64 // do use directly, use accessor methods
	totalOut uint64 // do use directly, use accessor methods, not yet implemented

	bytes.Buffer

	inmutex  *sync.RWMutex
	outmutex *sync.RWMutex
}

func NewCountingBuffer() *CountingBuffer {
	cb := new(CountingBuffer)
	cb.inmutex = new(sync.RWMutex)
	cb.outmutex = new(sync.RWMutex)
	return cb
}

func (cb *CountingBuffer) totalInSafeInc(delta uint64) {
	cb.inmutex.Lock()
	defer cb.inmutex.Unlock()
	cb.totalIn += delta
}

func (cb *CountingBuffer) totalOutSafeInc(delta uint64) {
	cb.outmutex.Lock()
	defer cb.outmutex.Unlock()
	cb.totalOut += delta
}

// Returns the total number of bytes written to the buffer since creation or the most recent ClearCount()
func (cb CountingBuffer) TotalIn() uint64 {
	cb.inmutex.RLock()
	defer cb.inmutex.RUnlock()
	return cb.totalIn
}

// Returns the total number of bytes read from the buffer since creation or the most recent ClearCount()
// Not yet implemented, will always return 0
func (cb CountingBuffer) TotalOut() uint64 {
	cb.outmutex.RLock()
	defer cb.outmutex.RUnlock()
	return cb.totalOut
}

// Clear byte counter.
func (cb *CountingBuffer) ClearCount() {
	cb.inmutex.Lock()
	defer cb.inmutex.Unlock()
	cb.outmutex.Lock()
	defer cb.outmutex.Unlock()
	cb.totalIn = 0
	cb.totalOut = 0
}

func (cb *CountingBuffer) Write(p []byte) (n int, err error) {
	n, err = cb.Buffer.Write(p)
	cb.totalInSafeInc(uint64(n))
	return
}

func (cb *CountingBuffer) WriteString(s string) (n int, err error) {
	n, err = cb.Buffer.WriteString(s)
	cb.totalInSafeInc(uint64(n))
	return
}

// don't be fooled by the name, this writes to the buffer!
func (cb *CountingBuffer) ReadFrom(r io.Reader) (n int64, err error) {

	// we can't use cb.Buffer.ReadFrom(r).  That implementation blocks until EOF.
	// instead we reimplement a simplified version of ReadFrom() here.

	buffer_space := 100 // we assume that most lines don't exceed 80 chars and add some padding.
	buf := make([]byte, buffer_space)

	for {
		// read from r and write into our buf
		m, err := r.Read(buf)
		if m > 0 {
			n += int64(m)

			// copy whatever was written into our buf into cb.Buffer
			mm, err2 := cb.Write(buf[0:m])
			if err2 != nil {
				return int64(mm), err2
			}
		}

		// now that our cb.Buffer has everything, do some cleanup.
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}

		// if m is largish, grow our buffer
		if m > int(0.7*float64(buffer_space)) {
			buffer_space *= 2
			buf = make([]byte, buffer_space)
		}
	}
	return n, nil
}

func (cb *CountingBuffer) WriteByte(c byte) error {
	err := cb.Buffer.WriteByte(c)
	if err == nil {
		cb.totalInSafeInc(1)
	}
	return err
}

func (cb *CountingBuffer) WriteRune(r rune) (n int, err error) {
	n, err = cb.Buffer.WriteRune(r)
	cb.totalInSafeInc(uint64(n))
	return
}
