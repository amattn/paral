package main

import (
	"bytes"
	"io"
)

// This could be split out into its own package at some point.
// I've needed similar functionality in other apps.

// Counts the total number of bytes written since creation or the most recent ClearCount()
// Reset() and Truncate() methods do NOT affect the running counts.
type CountingBuffer struct {
	totalIn  uint64
	totalOut uint64 // not yet implemented

	bytes.Buffer
}

func NewCountingBuffer() *CountingBuffer {
	return new(CountingBuffer)
}

func (cb *CountingBuffer) totalInSafeInc(i uint64) {

}

// Returns the total number of bytes written to the buffer since creation or the most recent ClearCount()
func (cb CountingBuffer) TotalIn() uint64 {
	return cb.totalIn
}

// Clear byte counter.
func (cb *CountingBuffer) ClearCount() {
	cb.totalIn = 0
}

func (cb *CountingBuffer) Write(p []byte) (n int, err error) {
	n, err = cb.Buffer.Write(p)
	if err != nil {
		cb.totalInSafeInc(uint64(n))
	}
	return
}

func (cb *CountingBuffer) WriteString(s string) (n int, err error) {
	n, err = cb.Buffer.WriteString(s)
	if err != nil {
		cb.totalInSafeInc(uint64(n))
	}
	return
}

// don't be fooled by the name, this writes to the buffer!
func (cb *CountingBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = cb.Buffer.ReadFrom(r)
	if err != nil {
		cb.totalInSafeInc(uint64(n))
	}
	return
}

func (cb *CountingBuffer) WriteByte(c byte) error {
	err := cb.Buffer.WriteByte(c)
	if err != nil {
		cb.totalInSafeInc(1)
	}
	return err
}

func (cb *CountingBuffer) WriteRune(r rune) (n int, err error) {
	n, err = cb.Buffer.WriteRune(r)
	if err != nil {
		cb.totalInSafeInc(uint64(n))
	}
	return
}
