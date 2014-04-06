// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved

package testing

import (
	"errors"
	"sync"
	"time"
)

func WaitForSyncGroup(wg *sync.WaitGroup, d time.Duration) error {
	c := make(chan int)
	defer func() { close(c) }()

	go func() {
		wg.Wait()
		c <- 1
	}()

	select {
	case <-c:
		return nil
	case <-time.After(d):
		return errors.New("wait group did not finish")
	}
}
