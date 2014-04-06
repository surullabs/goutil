// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved

package concurrent

import (
	"errors"
	. "github.com/surullabs/goutil/testing"
	. "launchpad.net/gocheck"
	"sync"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type WorkQueueSuite struct{}

var _ = Suite(&WorkQueueSuite{})

type taskTest struct {
	task   Task
	result interface{}
	err    interface{}
}

func (t taskTest) Test(ticket Ticket, c *C) {
	response, err := ticket.Wait()
	c.Check(response, DeepEquals, t.result)
	c.Check(err, NilOrErrorMatches, t.err)
}

type queueTest struct {
	name  string
	queue WorkQueue
}

func (q queueTest) Test(c *C) {
	c.Log("Starting ", q.name)
	var wg sync.WaitGroup
	taskList := tasks(&wg, c)
	results := make([]Ticket, len(taskList))
	q.queue.Start()
	defer q.queue.Stop()
	for i, task := range taskList {
		wg.Add(1)
		results[i] = q.queue.Add(task.task)
	}
	for i, task := range taskList {
		task.Test(results[i], c)
	}
	// Wait for the sync group at each turn to ensure we ran everything
	c.Assert(WaitForSyncGroup(&wg, 100*time.Millisecond), IsNil)

	// Now test a double call of start
	q.queue.Start()
	wg.Add(1)
	t := taskList[0]
	t.Test(q.queue.Add(t.task), c)

	// Ensure we stop once
	q.queue.Stop()
}

func tasks(wg *sync.WaitGroup, c *C) []taskTest {
	return []taskTest{
		{
			task: WorkFn(func() {
				wg.Done()
				c.Log("Called WorkFn")
			}),
		},
		{
			task: ErrorWorkFn(func() error {
				wg.Done()
				c.Log("Called ErrorWorkFn")
				return nil
			}),
		},
		{
			task: ResultWorkFn(func() (interface{}, error) {
				wg.Done()
				c.Log("Called ResultWorkFn")
				return nil, nil
			}),
		},
		{
			task: ErrorWorkFn(func() error {
				wg.Done()
				c.Log("Called ErrorWorkFn with error")
				return errors.New("failed")
			}),
			err: "failed",
		},
		{
			task: ResultWorkFn(func() (interface{}, error) {
				wg.Done()
				c.Log("Called ResultWorkFn with error")
				return nil, errors.New("failed result")
			}),
			err: "failed result",
		},
	}
}
func (w *WorkQueueSuite) TestWorkQueue(c *C) {
	for _, test := range []queueTest{
		{"Single Queue Test", NewWorkPool(1)},
		{"Work Pool Test - 2 ", NewWorkPool(2)},
		{"Work Pool Test - 4", NewWorkPool(4)},
	} {
		test.Test(c)
	}
}
