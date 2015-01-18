package main

import (
	"fmt"
	"io"
	"strings"
)

// http://stackoverflow.com/questions/17006262/clearing-output-of-a-terminal-program-in-linux-c
// We currently use "\033[%dA" to move up a line
// We currently use "\r" to move to the beginning of a line
// This isn't really portable...

const (
	up1line_metacharacter = "\033[1A"
)

// kind of like a writer, but not really
type ErasableOutputter struct {
	writer io.Writer

	numberOfProgressLinesWritten int
	lastOutputWidth              int
}

func NewErasableOutputter(writer io.Writer) *ErasableOutputter {
	eo := new(ErasableOutputter)
	eo.writer = writer

	return eo
}

//
func (eo *ErasableOutputter) OutputUnerasableString(str string) {
	fmt.Fprintf(eo.writer, str)
	eo.numberOfProgressLinesWritten = 0
}

func (eo *ErasableOutputter) OutputErasableString(str string) {
	if strings.HasSuffix(str, "\n") == false {
		str = str + "\n"
	}
	eo.EraseLastEraseble()
	fmt.Fprintf(eo.writer, str)
	lines := strings.Split(str, "\n")
	eo.numberOfProgressLinesWritten = len(lines) - 1
	eo.lastOutputWidth = 0
	for _, line := range lines {
		if len(line) > eo.lastOutputWidth {
			eo.lastOutputWidth = len(line)
		}
	}
}

// if the last output was an erasable string, this method will
// "erase" it by moving the cursor back and outputting spaces.
// If the last output was NOT erasable, then this method does nothing.
// This returns the number of lines erased.
func (eo *ErasableOutputter) EraseLastEraseble() int {
	erased := 0
	if eo.numberOfProgressLinesWritten > 0 {
		for i := 0; i < eo.numberOfProgressLinesWritten; i++ {
			fmt.Fprintf(eo.writer, "\r%s", up1line_metacharacter)
			fmt.Fprintf(eo.writer, strings.Repeat(" ", eo.lastOutputWidth))
		}
		fmt.Fprintf(eo.writer, "\r")
		erased = eo.numberOfProgressLinesWritten
		eo.numberOfProgressLinesWritten = 0
		eo.lastOutputWidth = 0
	}
	return erased
}
