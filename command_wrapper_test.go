package main

import "testing"

func TestFilenamification(t *testing.T) {
	inputs := []string{
		"",
		"sleep 1",
		"sleep 1 && echo 2",
		"1234567890.-abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ.-_",
		"sleep 1 && echo \"Hello, World!\"",
		"compiler -t build --flag some/thing/path jfk2 'isndf' ()",
	}
	expecteds := []string{
		"",
		"sleep_1",
		"sleep_1_echo_2",
		"1234567890.-abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"sleep_1_echo_Hello_World",
		"compiler_t_build_flag_some-thing-path_jfk2_isndf",
	}

	if len(inputs) != len(expecteds) {
		t.Fatal("Invalid test setup: len(inputs) != len(expecteds)", len(inputs), len(expecteds))
	}

	for i, input := range inputs {

		candidate := filenamification(input)
		expected := expecteds[i]
		if candidate != expected {
			t.Error("candidate != expected", i, candidate, expected)
		}
	}
}
