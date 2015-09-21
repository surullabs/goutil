// Copyright 2014, Surul Software Labs GmbH
// All rights reserved.
package concurrent

import (
	"sync"

	"github.com/surullabs/fault"
)

type ErrorCollector struct {
	m      sync.Mutex
	errors fault.ErrorChain
}

func (c *ErrorCollector) Append(err error) {
	if err == nil {
		return
	}
	c.m.Lock()
	defer c.m.Unlock()
	c.errors.Append(err)
}

func (c *ErrorCollector) Error() error {
	c.m.Lock()
	defer c.m.Unlock()
	return (&c.errors).AsError()
}
