package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type commandWrapper struct {
	raw    string
	num    int
	output []byte

	cmd    *exec.Cmd
	err    error            // any error produced by cmd.Start() or cmd.Wait()
	pstate *os.ProcessState // if cmd.Wait() returns an ExitErr, we pull out the process state

	cberr *CountingBuffer
	cbout *CountingBuffer

	start time.Time
	end   time.Time
}

func (wrapper commandWrapper) outputFileName() string {
	return fmt.Sprintf("%v_%d.out", output_file_prefix, wrapper.num)
}

func (wrapper commandWrapper) duration() time.Duration {
	return wrapper.end.Sub(wrapper.start)
}

// we use sh -c to run the command.  It means this program doesn't have
// to worry about command line argument parsing.
func (wrapper *commandWrapper) run() (err error) {

	// fmt.Println(wrapper.num, "running:", wrapper.raw)

	cmd := exec.Command("sh", "-c", wrapper.raw)
	wrapper.cberr = NewCountingBuffer()
	wrapper.cbout = NewCountingBuffer()
	cmd.Stderr = wrapper.cberr
	cmd.Stdout = wrapper.cbout
	wrapper.cmd = cmd

	wrapper.start = time.Now()
	err = cmd.Start()
	if err != nil {
		log.Println("start failure", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		// log.Println("run failure", err)
		exiterr, isExitErr := err.(*exec.ExitError)
		if isExitErr {
			wrapper.pstate = exiterr.ProcessState
		}
	}
	wrapper.end = time.Now()

	// fmt.Sprintln(wrapper.num, "total time elapsed:", wrapper.end.Sub(wrapper.start))
	return
}
