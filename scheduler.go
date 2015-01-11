package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"text/tabwriter"
	"time"
)

type Scheduler struct {
	maxSimul int

	numberSubmitted int
	todo            []*commandWrapper
	inProgress      map[int]*commandWrapper // key is commandWrapper.num
	finished        []*commandWrapper

	sync.WaitGroup

	ticker *time.Ticker
}

func NewScheduler(maxSimul int) *Scheduler {
	return &Scheduler{
		maxSimul: maxSimul,
	}
}

func (s *Scheduler) AddCommands(cmds ...*commandWrapper) {
	s.numberSubmitted += len(cmds)
	s.todo = append(s.todo, cmds...)
}

func (s *Scheduler) inFlight() int {
	return len(s.inProgress)
}

func (s *Scheduler) Start() {
	// log.Println("Trace Start()")
	s.Add(1)
	go s.start()

	go func() {
		s.ticker = time.NewTicker(200 * time.Millisecond)

		for now := range s.ticker.C {
			OutputProgess(now, s)
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
	s.inProgress = make(map[int]*commandWrapper)
	s.finished = make([]*commandWrapper, 0, len(s.todo))

	for len(s.todo) > 0 {
		// log.Println("Trace start() outer for loop")

		// are we at max inflight capacity?  just wait
		if s.inFlight() >= s.maxSimul {
			time.Sleep(1 * time.Millisecond)
		} else {
			nextWrapper := s.todo[0]
			s.todo = s.todo[1:]
			s.Add(1)
			s.inProgress[nextWrapper.num] = nextWrapper

			go func(wrapper *commandWrapper) {
				wrapper.run()

				s.finished = append(s.finished, wrapper)
				exit := "exit status 0"
				if wrapper.pstate != nil {
					exit = wrapper.pstate.String()
				}
				OutputString(s, fmt.Sprintf("Command %d finished, %v, time elapsed: %v\n", wrapper.num, exit, wrapper.duration()))
				delete(s.inProgress, wrapper.num)
				s.Done()
			}(nextWrapper)
		}
	}
}

func (s *Scheduler) Wait() {
	s.WaitGroup.Wait()
	s.ticker.Stop()
	OutputProgess(time.Now(), s)
}

func (s *Scheduler) Progress(t time.Time) string {

	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)
	w.Init(buf, 5, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "\tnot yet started\tin progress\tfinished\tsubmitted\t\n")
	fmt.Fprintf(w, "count:\t%d\t%d\t%d\t%d\t\n", len(s.todo), len(s.inProgress), len(s.finished), s.numberSubmitted)
	submitted := float32(s.numberSubmitted)
	fmt.Fprintf(w,
		"percent:\t%3.1f\t %3.1f\t   %3.1f\t\t\n",
		100*float32(len(s.todo))/submitted,
		100*float32(len(s.inProgress))/submitted,
		100*float32(len(s.finished))/submitted,
	)
	w.Flush()
	b, err := ioutil.ReadAll(buf)
	if err != nil {
		panic(err)
	}
	return string(b)
}
