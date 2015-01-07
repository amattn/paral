package main

import (
	"fmt"
	"time"
)

type command struct {
	raw    string
	num    int
	output []byte
}

func (cmd command) outputFileName() string {
	return fmt.Sprintf("%v_%d.out", output_file_prefix, cmd.num)
}

func (cmd command) run() {
	fmt.Println(cmd.num, "running:", cmd.raw)
	time.Sleep(time.Duration(cmd.num*400) * time.Millisecond)
	fmt.Println(cmd.num, "finishing:", cmd.raw)
}
