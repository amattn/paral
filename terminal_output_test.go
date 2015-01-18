package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func TestUnerasableOutputter(t *testing.T) {
	buff := new(bytes.Buffer)
	eo := NewErasableOutputter(buff)

	// expected behavior:

	eo.OutputUnerasableString("str\n")
	eo.OutputUnerasableString("str\nstr\n")
	eo.OutputUnerasableString("str\n")

	all, err := ioutil.ReadAll(buff)
	if err != nil {
		t.Fatal(399838273, err)
	}

	allOutput := string(all)

	newline_count := strings.Count(allOutput, "\n")
	assert.Equal(t, newline_count, 4)
}

func TestErasableOutputter(t *testing.T) {
	buff := new(bytes.Buffer)
	eo := NewErasableOutputter(buff)

	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str\nstr\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str\nstr\nstr\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str1\nstr\nstr\nstr\n")
	eo.OutputErasableString("str2\nstr\nstr\nstr\n")
	eo.OutputErasableString("str3\nstr\nstr\nstr\n")
	eo.OutputErasableString("str4\nstr\nstr\n")
	eo.OutputErasableString("str\nstr\n")
	eo.OutputErasableString("str\n")

	all, err := ioutil.ReadAll(buff)
	if err != nil {
		t.Fatal(399838273, err)
	}

	allOutput := string(all)

	newline_count := strings.Count(allOutput, "\n")
	assert.Equal(t, newline_count, 29)

	// we use the terminal metacharacter: "\033[1A" where 1 is the number of lines to go up...
	up1line_count := strings.Count(allOutput, up1line_metacharacter)
	assert.Equal(t, up1line_count, 28)
}

func TestSimpleEraserOutputter(t *testing.T) {
	buff := new(bytes.Buffer)
	eo := NewErasableOutputter(buff)

	eo.OutputErasableString("str11\nstr12\nstr13\nstr14\n")
	eo.EraseLastEraseble()
	eo.OutputUnerasableString("a\nb\nc\nd\ne\nf\ng\n")
	eo.OutputErasableString("str21\nstr22\nstr23\nstr24\n")

	all, err := ioutil.ReadAll(buff)
	if err != nil {
		t.Fatal(399838273, err)
	}

	allOutput := string(all)

	// fmt.Fprintf(os.Stdout, "%v", allOutput)

	newline_count := strings.Count(allOutput, "\n")
	assert.Equal(t, newline_count, 15)
	up1line_count := strings.Count(allOutput, up1line_metacharacter)
	assert.Equal(t, up1line_count, 4)
}

func TestMixedOutputter(t *testing.T) {
	buff := new(bytes.Buffer)
	eo := NewErasableOutputter(buff)

	eo.OutputUnerasableString("a\n")
	eo.OutputUnerasableString("b\nc\n")
	eo.OutputUnerasableString("d\n")
	eo.OutputErasableString("str\n")
	eo.OutputErasableString("str0\n")
	eo.OutputUnerasableString("e\nf\n")
	eo.OutputErasableString("str1\n")
	eo.OutputErasableString("str2\n")
	eo.OutputErasableString("str3\n")
	eo.OutputErasableString("str4\n")
	eo.OutputErasableString("str5\n")
	eo.EraseLastEraseble()
	eo.OutputUnerasableString("g\n")
	eo.OutputUnerasableString("h\n")
	eo.OutputErasableString("str1\nstr\nstr\nstr\n")
	eo.OutputErasableString("str2\nstr\nstr\nstr\n")
	eo.OutputErasableString("str3\nstr\nstr\nstr\n")
	eo.OutputUnerasableString("i\n")
	eo.OutputErasableString("str4\nstr\nstr\nstr\n")
	eo.OutputUnerasableString("j\n")
	eo.OutputErasableString("str5\nstr\nstr\nstr\n")
	eo.OutputErasableString("str6\nstr\nstr\nstr\n")
	eo.OutputErasableString("str71\nstr72\nstr73\nstr74\n")
	eo.EraseLastEraseble()
	eo.OutputUnerasableString("k\n")
	eo.OutputUnerasableString("l\n")
	eo.OutputUnerasableString("m\n")
	eo.OutputUnerasableString("n\n")
	eo.OutputUnerasableString("o\n")
	eo.OutputUnerasableString("p\n")

	all, err := ioutil.ReadAll(buff)
	if err != nil {
		t.Fatal(399838273, err)
	}

	allOutput := string(all)

	// fmt.Fprintf(os.Stdout, "%v", allOutput)

	newline_count := strings.Count(allOutput, "\n")
	assert.Equal(t, newline_count, 51)
	up1line_count := strings.Count(allOutput, up1line_metacharacter)
	assert.Equal(t, up1line_count, 26)
}
