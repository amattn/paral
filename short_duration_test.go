package main

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestSeconds(t *testing.T) {
	inputs := []time.Duration{
		142000000000 * time.Nanosecond, // 142 seconds
		14200000000 * time.Nanosecond,  // 14.2 seconds
		1420000000 * time.Nanosecond,   // 1.42 seconds
		142000000 * time.Nanosecond,
		14200000 * time.Nanosecond,
		1420000 * time.Nanosecond, // 5.42 milliseconds
		142000 * time.Nanosecond,
		14200 * time.Nanosecond,
		1420 * time.Nanosecond, // 1.42 microseconds
		142 * time.Nanosecond,
		14 * time.Nanosecond,
		1 * time.Nanosecond, // 1 nanosecond
		0 * time.Second,
		-142000000000 * time.Nanosecond, // 142 seconds
		-14200000000 * time.Nanosecond,  // 14.2 seconds
		-1420000000 * time.Nanosecond,   // 1.42 seconds
		-142000000 * time.Nanosecond,
		-14200000 * time.Nanosecond,
		-1420000 * time.Nanosecond, // 5.42 milliseconds
		-142000 * time.Nanosecond,
		-14200 * time.Nanosecond,
		-1420 * time.Nanosecond, // 1.42 microseconds
		-142 * time.Nanosecond,
		-14 * time.Nanosecond,
		-1 * time.Nanosecond, // 1 nanosecond
	}

	expecteds := []string{
		"2m22s",
		"14.2s",
		"1.42s",
		"142ms",
		"14.2ms",
		"1.42ms",
		"142µs",
		"14.2µs",
		"1.42µs",
		"142ns",
		"14ns",
		"1ns",
		"0",
		"-2m22s",
		"-14.2s",
		"-1.42s",
		"-142ms",
		"-14.2ms",
		"-1.42ms",
		"-142µs",
		"-14.2µs",
		"-1.42µs",
		"-142ns",
		"-14ns",
		"-1ns",
	}

	assert.Equal(t, len(inputs), len(expecteds), "37473948793 Bad test harness")

	for i, input := range inputs {
		exp := expecteds[i]
		candidate := ShortString(input)
		assert.Equal(t, exp, candidate, "23987428u4 ", i)
	}
}

func TestMoreThanSecond(t *testing.T) {
	inputs := []time.Duration{
		1 * time.Second, // 142 seconds
		14 * time.Second,
		142 * time.Second,
		1420 * time.Second,
		14200 * time.Second,
		142000 * time.Second,
		1380000 * time.Second,
		1420000 * time.Second,
		-1 * time.Second, // 142 seconds
		-14 * time.Second,
		-142 * time.Second,
		-1420 * time.Second,
		-14200 * time.Second,
		-142000 * time.Second,
		-1380000 * time.Second,
		-1420000 * time.Second,
	}
	expecteds := []string{
		"1.00s",
		"14.0s",
		"2m22s",
		"23m40s",
		"3h56m",
		"39h26m",
		"16.0d",
		"16.4d",
		"-1.00s",
		"-14.0s",
		"-2m22s",
		"-23m40s",
		"-3h56m",
		"-39h26m",
		"-16.0d",
		"-16.4d",
	}
	assert.Equal(t, len(inputs), len(expecteds), "98723984 Bad test harness")

	for i, input := range inputs {
		exp := expecteds[i]
		candidate := ShortString(input)
		assert.Equal(t, exp, candidate, "98723985 ", i)
	}
}
