package main

import (
	"bytes"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"testing"

	"github.com/bmizerany/assert"
)

func TestCounting(t *testing.T) {
	cb := NewCountingBuffer()

	cb.Write([]byte(""))
	assert.Equal(t, uint64(0), cb.TotalIn())
	cb.Write([]byte("hi"))
	assert.Equal(t, uint64(2), cb.TotalIn())
	cb.Write([]byte("hello"))
	assert.Equal(t, uint64(7), cb.TotalIn())

	cb.ClearCount()
	assert.Equal(t, uint64(0), cb.TotalIn())
	cb.WriteString("hi2")
	assert.Equal(t, uint64(3), cb.TotalIn())
	cb.WriteString("hello2")
	assert.Equal(t, uint64(9), cb.TotalIn())
}

// you _should_ be able to make this fail by commenting out the mutex lock/unlock in totalInSafeInc
func TestSimulCounting(t *testing.T) {
	oldgmp := runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(oldgmp)

	cb := NewCountingBuffer()

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			cb.WriteString("1234567")
			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, uint64(70000), cb.TotalIn())
}

func TestReadFrom(t *testing.T) {

	buf := bytes.NewReader([]byte("Hello"))
	cb := NewCountingBuffer()

	n, err := cb.ReadFrom(buf)
	assert.Equal(t, nil, err, 12339872388)

	assert.Equal(t, int64(5), n, 12339872389)
}

func TestReadFromStdOut(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo abc && sleep 0.15 && echo d")
	cberr := NewCountingBuffer()
	cbout := NewCountingBuffer()
	cmd.Stderr = cberr
	cmd.Stdout = cbout

	err := cmd.Run()
	assert.Equal(t, nil, err, 948027283)
	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, uint64(0), cberr.TotalIn())
	assert.Equal(t, uint64(6), cbout.TotalIn())
}
