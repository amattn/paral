package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// https://groups.google.com/forum/#!msg/golang-nuts/a9PitPAHSSU/ziQw1-QHw3EJ
const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

var show_h bool
var show_help bool
var show_version bool
var output_file_prefix string
var max_simul int

func init() {
	flag.BoolVar(&show_h, "h", false, "show help message and exit(0)")
	flag.BoolVar(&show_help, "help", false, "show help message and exit(0)")
	flag.BoolVar(&show_version, "version", false, "show version info and exit(0)")
	flag.StringVar(&output_file_prefix, "out", "out", "prefix name of output file which contains stdout and stderr of processes")
	flag.IntVar(&max_simul, "n", runtime.NumCPU(), "maximum number of simultaneous processes.  Defaults to NumCPU().  To run all processes simultaneously, set to 0.")
}

func main() {
	fmt.Println("paral")
	fmt.Printf("> paral version %v (build %d)\n", Version(), BuildNumber())
	fmt.Printf("> Go Version: %v %v/%v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("> NumCPU(): %v\n", runtime.NumCPU())
	fmt.Printf("> GOMAXPROCS(): %v\n", runtime.GOMAXPROCS(0))
	start := time.Now()
	fmt.Println(">", start)

	// command line flags:
	flag.Parse()

	if max_simul <= 0 {
		fmt.Println(">", "Maximum Simulataneous processes:", "unlimited")
		max_simul = MaxInt
	} else {
		fmt.Println(">", "Maximum Simulataneous processes:", max_simul)
	}

	if show_version {
		os.Exit(0)
	}

	if show_h || show_help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("Commands:")
	args := flag.Args()
	cmds := make([]*commandWrapper, 0, len(args))
	for i, a := range flag.Args() {
		cmd := commandWrapper{
			raw: a,
			num: i,
		}
		fmt.Printf("%3d, %v\n", i, cmd.raw)
		cmds = append(cmds, &cmd)
	}

	// Scheduler

	s := NewScheduler(max_simul)
	s.AddCommands(cmds...)
	s.Start()
	s.Wait()
	// log.Println("TRACE Finished", s)

	// Log some output

	// calculate how much time would theoretically have been spent if we did things in serial:
	var total_process_time time.Duration
	for _, cmd := range s.finished {
		total_process_time += cmd.end.Sub(cmd.start)
	}
	log.Println("Sum of all process time:", total_process_time)

	// total time taken:
	total_duration := time.Now().Sub(start)
	log.Println("     Total time elapsed:", total_duration)

	// show some gain/loss
	ratio := float64(total_process_time) / float64(time.Now().Sub(start))
	if ratio >= 1 {
		log.Printf("      effeciency ratio (gain): %0.2f", ratio)
	} else if ratio == 1 {
		log.Printf("             effeciency ratio: %0.2f", ratio)
	} else {
		log.Printf("        effeciency ratio: %0.2f", (ratio))
	}
}
