package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode"
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

func newWrapper(raw string, num int) *commandWrapper {
	wrapper := new(commandWrapper)
	wrapper.raw = raw
	wrapper.num = num
	wrapper.cberr = NewCountingBuffer()
	wrapper.cbout = NewCountingBuffer()
	return wrapper
}
func (wrapper commandWrapper) outputFileName() string {
	filename := filenamification(wrapper.raw)
	if len(filename) > 50 {
		filename = filename[:50]
	}
	return fmt.Sprintf("%v%d_%s%v", output_file_prefix, wrapper.num, filename, output_file_suffix)
}
func (wrapper commandWrapper) errorFileName() string {
	filename := filenamification(wrapper.raw)
	if len(filename) > 50 {
		filename = filename[:50]
	}
	return fmt.Sprintf("%v%d_%s%v", error_file_prefix, wrapper.num, filename, error_file_suffix)
}

func (wrapper commandWrapper) duration() time.Duration {
	if wrapper.end.IsZero() {
		return time.Now().Sub(wrapper.start)
	}
	return wrapper.end.Sub(wrapper.start)
}

// we use sh -c to run the command.  It means this program doesn't have
// to worry about command line argument parsing.
func (wrapper *commandWrapper) run() (err error) {

	// fmt.Println(wrapper.num, "running:", wrapper.raw)

	cmd := exec.Command("sh", "-c", wrapper.raw)
	cmd.Stderr = wrapper.cberr
	cmd.Stdout = wrapper.cbout
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	var wg sync.WaitGroup

	handleOutput := func(filename string, writer io.WriterTo, stopChan chan bool) {
		defer wg.Done()
		fo, err := os.Create(filename)
		if err == nil {
			defer fo.Close()
			for {
				select {
				case <-stopChan:
					// one last write to clear out anything in the buffer.
					writer.WriteTo(fo)
					return
				default:
					_, err := writer.WriteTo(fo)
					if err != nil {
						return
					}
				}
				time.Sleep(3 * time.Millisecond)
			}
		}
	}

	outStopCh := make(chan bool)
	errStopCh := make(chan bool)
	if capture_output_to_file {
		wg.Add(2)
		go handleOutput(wrapper.outputFileName(), wrapper.cbout, outStopCh)
		go handleOutput(wrapper.errorFileName(), wrapper.cberr, errStopCh)
	}
	// fe, err := os.Create(wrapper.errorFileName())
	// if err == nil {
	// 	defer fe.Close()
	// 	wrapper.cberr.WriteTo(fe)
	// }

	wrapper.cmd = cmd

	wrapper.start = time.Now()
	err = cmd.Start()
	if err != nil {
		log.Println("start failure", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("run failure", err)
		exiterr, isExitErr := err.(*exec.ExitError)
		if isExitErr {
			wrapper.pstate = exiterr.ProcessState
		}
	}
	wrapper.end = time.Now()

	// fmt.Sprintln(wrapper.num, "total time elapsed:", wrapper.end.Sub(wrapper.start))

	outStopCh <- true
	errStopCh <- true
	wg.Wait()
	return
}

// Given a string, return a string that can be used as a filename
// 1. replace whitespace with underbars
// 2. replace ![a-Z0-9-_.] with dashes
// 3. trim any leaading or trailing dots, underbars or dashes
// 4. replace any groups of dashes or underbars with a single dash or underbar
func filenamification(input string) string {
	// step 1
	mapping1 := func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}
	// step 2
	mapping2 := func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r == '-' || r == '_' || r == '.' {
			return r
		}
		return '-'
	}
	step1 := strings.Map(mapping1, input)
	step2 := strings.Map(mapping2, step1)
	// step 3
	step3 := strings.Trim(step2, "-_.")

	// step 4
	step4 := step3

OuterLoop:
	for {
		// log.Println("step4", step4)
		switch {
		case strings.Contains(step4, "--"):
			step4 = strings.Replace(step4, "--", "-", -1)
		case strings.Contains(step4, "__"):
			step4 = strings.Replace(step4, "__", "_", -1)
		case strings.Contains(step4, ".."):
			step4 = strings.Replace(step4, "..", ".", -1)
		case strings.Contains(step4, "_-"):
			step4 = strings.Replace(step4, "_-", "_", -1)
		case strings.Contains(step4, "-_"):
			step4 = strings.Replace(step4, "-_", "_", -1)
		default:
			break OuterLoop
		}
	}

	return step4
}
