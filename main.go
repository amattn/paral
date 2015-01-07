package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
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
	fmt.Println(">", time.Now())

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
	cmds := make([]*command, 0, len(args))
	for i, a := range flag.Args() {
		cmd := command{
			raw: a,
			num: i,
		}
		fmt.Printf("%3d, %v\n", i, cmd.raw)
		cmds = append(cmds, &cmd)
	}

	// Scheduler

	s := MakeScheduler(max_simul)
	s.Start()

	for _, cmd := range cmds {
		s.intake <- cmd
	}
	s.Close()

	s.Wait()
	log.Println("Finished", s)
}

type Scheduler struct {
	maxSimul int
	intake   chan *command

	inFlight int

	sync.WaitGroup
}

func MakeScheduler(maxSimul int) Scheduler {
	return Scheduler{
		maxSimul: maxSimul,
		intake:   make(chan *command),
	}
}

func (s *Scheduler) Start() {
	go func() {
		s.Add(1)
		defer s.Done()

		for cmd := range s.intake {
			log.Println(s.inFlight, s.maxSimul)
			for s.inFlight >= s.maxSimul {
				time.Sleep(10 * time.Millisecond)
			}
			s.Add(1)
			s.inFlight++
			go func(cmd *command) {
				cmd.run()
				s.inFlight--
				s.Done()
			}(cmd)
		}
	}()
}

func (s *Scheduler) Close() {
	close(s.intake)
}
