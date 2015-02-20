package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"
)

type Scheduler struct {
	maxSimul int

	numberSubmitted int
	todo            []*commandWrapper
	inFlight        map[int]*commandWrapper // key is commandWrapper.num
	finished        []*commandWrapper

	sync.WaitGroup

	ticker    *time.Ticker
	outputter *ErasableOutputter
}

func NewScheduler(maxSimul int) *Scheduler {
	return &Scheduler{
		maxSimul:  maxSimul,
		outputter: NewErasableOutputter(os.Stdout),
	}
}

func (s *Scheduler) AddCommands(cmds ...*commandWrapper) {
	s.numberSubmitted += len(cmds)
	s.todo = append(s.todo, cmds...)
}

func (s *Scheduler) numInFlight() int {
	return len(s.inFlight)
}

func (s *Scheduler) Start() {
	// log.Println("Trace Start()")
	s.Add(1)
	go s.start()

	go func() {
		s.ticker = time.NewTicker(347 * time.Millisecond)

		for _ = range s.ticker.C {
			s.outputter.OutputErasableString(s.Progress())
		}
	}()
}

func (s *Scheduler) start() {
	defer s.Done()

	// log.Println("Trace start()")
	if len(s.todo) <= 0 {
		log.Println("Trace no commands to process")
		return
	}
	s.inFlight = make(map[int]*commandWrapper)
	s.finished = make([]*commandWrapper, 0, len(s.todo))

	for len(s.todo) > 0 {
		// log.Println("Trace start() outer for loop")

		// are we at max inflight capacity?  just wait
		if s.numInFlight() >= s.maxSimul {
			time.Sleep(1 * time.Millisecond)
		} else {
			nextWrapper := s.todo[0]
			s.todo = s.todo[1:]
			s.Add(1)
			s.inFlight[nextWrapper.num] = nextWrapper

			go func(wrapper *commandWrapper) {
				wrapper.run()
				s.finished = append(s.finished, wrapper)
				exit := "exit status 0"
				if wrapper.pstate != nil {
					exit = wrapper.pstate.String()
				}
				finished_str := fmt.Sprintf("Command %d finished, %v, time elapsed:%v stdout bytes:%d stderr bytes:%d\n", wrapper.num, exit, ShortString(wrapper.duration()), wrapper.cbout.TotalIn(), wrapper.cberr.TotalIn())
				s.outputter.EraseLastEraseble()
				s.outputter.OutputUnerasableString(finished_str)
				s.outputter.OutputErasableString(s.Progress())
				delete(s.inFlight, wrapper.num)
				s.Done()
			}(nextWrapper)
		}
	}
}

func (s *Scheduler) Wait() {
	s.WaitGroup.Wait()
	s.ticker.Stop()
	s.outputter.OutputErasableString(s.Progress())
}

func (s *Scheduler) Progress() string {
	// return ""

	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)
	w.Init(buf, 5, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "\tnot yet started\tin progress\tfinished\tsubmitted\t\n")
	fmt.Fprintf(w, "count:\t%d\t%d\t%d\t%d\t\n", len(s.todo), len(s.inFlight), len(s.finished), s.numberSubmitted)
	submitted := float32(s.numberSubmitted)
	fmt.Fprintf(w,
		"percent:\t%3.1f\t %3.1f\t   %3.1f\t\t\n",
		100*float32(len(s.todo))/submitted,
		100*float32(len(s.inFlight))/submitted,
		100*float32(len(s.finished))/submitted,
	)

	keys := make([]int, 0, len(s.inFlight))
	for num, _ := range s.inFlight {
		keys = append(keys, num)
	}

	sort.Ints(keys)

	for _, num := range keys {
		cmd := s.inFlight[num]
		if cmd != nil {

			raw := cmd.raw
			if len(raw) > 10 {
				raw = raw[0:9]
			}
			fmt.Fprintf(w, "Command %d running %+v time elapsed:%v stdout bytes:%d stderr bytes:%d\n", num, raw, ShortString(cmd.duration()), cmd.cbout.TotalIn(), cmd.cberr.TotalIn())
		}
	}

	w.Flush()
	b, err := ioutil.ReadAll(buf)
	if err != nil {
		panic(err)
	}
	return string(b)
}
