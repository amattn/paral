package main

// This is a lousy, proof of concept implementation.
// Don't reference this code unless you have no alternative.

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func OutputString(s *Scheduler, str string) {
	EraseProgess(s)
	fmt.Fprintf(os.Stdout, str)
	lastOutputWasProgress = false
	OutputProgess(time.Now(), s)
}

var lastOutputWasProgress bool
var lastOutputWidth int

func EraseProgess(s *Scheduler) string {
	prog := s.Progress(time.Now())
	if lastOutputWasProgress {
		lines := strings.Count(prog, "\n")
		fmt.Fprintf(os.Stdout, "\r\033[%dA", lines)
		fmt.Fprintf(os.Stdout, strings.Repeat(" ", lastOutputWidth))
		fmt.Fprintf(os.Stdout, "\r")
	}
	return prog
}
func OutputProgess(t time.Time, s *Scheduler) {
	prog := EraseProgess(s)
	fmt.Fprintf(os.Stdout, prog)
	lastOutputWasProgress = true
	lines := strings.Split(prog, "\n")
	lastOutputWidth = 0
	for _, line := range lines {
		if len(line) > lastOutputWidth {
			lastOutputWidth = len(line)
		}
	}
}
