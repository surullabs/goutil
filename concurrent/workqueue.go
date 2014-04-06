// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved

package concurrent

import (
	"fmt"
	"github.com/surullabs/fault"
	"sync"
)

var _ = fmt.Print

type result struct {
	value interface{}
	err   error
}

type functionTask func() (interface{}, error)

func (f functionTask) newTicket() *workTask {
	return &workTask{fn: f, signal: make(chan result, 1)}
}

func ResultWorkFn(fn func() (interface{}, error)) Task {
	return functionTask(fn)
}

func WorkFn(fn func()) Task {
	return ResultWorkFn(func() (interface{}, error) {
		fn()
		return nil, nil
	})
}

func ErrorWorkFn(fn func() error) Task {
	return ResultWorkFn(func() (interface{}, error) {
		return nil, fn()
	})

}

type workTask struct {
	fn     functionTask
	signal chan result
}

func (w *workTask) do() {
	var res result
	res.value, res.err = w.fn()
	w.signal <- res
	close(w.signal)
}

func (w *workTask) Wait() (interface{}, error) {
	r := <-w.signal
	return r.value, r.err
}

type Ticket interface {
	Wait() (interface{}, error)
	do()
}

type Task interface {
	newTicket() *workTask
}

type WorkQueue interface {
	Start()
	Add(Task) Ticket
	Stop()
}

type worker struct {
	work chan *workTask
}

func (s *worker) start() {
	if s.work != nil {
		return
	}
	s.work = make(chan *workTask, 100)
	var wg sync.WaitGroup
	wg.Add(1)
	go s.loop(&wg)
	wg.Wait()
}

func (s *worker) stop() {
	if s.work != nil {
		close(s.work)
		s.work = nil
	}
}

func (s *worker) loop(wg *sync.WaitGroup) {
	// Signal that the goroutine has been started
	wg.Done()
	for {
		select {
		case t, done := <-s.work:
			if t != nil {
				t.do()
			}
			if !done {
				return
			}
		}
	}
}

func (s *worker) add(t Task) Ticket {
	ticket := t.newTicket()
	s.work <- ticket
	return ticket
}

type workPool struct {
	workers []worker
	m       sync.Mutex
	index   int
}

func (w *workPool) Start() {
	w.m.Lock()
	defer w.m.Unlock()

	for i := range w.workers {
		w.workers[i].start()
	}
}

func (w *workPool) Add(t Task) Ticket {
	// TODO: Move away from the Mutex since it is inefficient for single queues.
	w.m.Lock()
	defer w.m.Unlock()

	ticket := w.workers[w.index].add(t)
	w.index = (w.index + 1) % len(w.workers)
	return ticket
}

func (w *workPool) Stop() {
	w.m.Lock()
	defer w.m.Unlock()
	for i := range w.workers {
		w.workers[i].stop()
	}
}

func NewWorkPool(size int) WorkQueue {
	return &workPool{workers: make([]worker, size)}
}

type TicketBatch struct {
	tickets []Ticket
}

func (t *TicketBatch) Add(ticket Ticket) {
	if t.tickets == nil {
		t.tickets = make([]Ticket, 0)
	}
	t.tickets = append(t.tickets, ticket)
}

func (t *TicketBatch) Wait() error {
	if t.tickets == nil {
		return nil
	}
	errs := &fault.ErrorChain{}
	for _, ticket := range t.tickets {
		_, err := ticket.Wait()
		errs.Chain(err)
	}
	return errs.AsError()
}
