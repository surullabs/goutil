// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved.

// This provides various testing utilities that are compatible with
// the gocheck package.
package testing

import (
	"fmt"
	. "launchpad.net/gocheck"
)

type nilErrorMatcher struct {
	*CheckerInfo
}

var NilOrErrorMatches = &nilErrorMatcher{
	&CheckerInfo{Name: "NilOrErrorMatches", Params: []string{"obtained", "regex"}},
}

func (n *nilErrorMatcher) Check(params []interface{}, names []string) (result bool, error string) {
	if params[0] == nil && params[1] != nil {
		return false, fmt.Sprintf("Obtained error is nil when expecting non-nil error %v", params[1])
	} else if params[0] != nil && params[1] == nil {
		return false, fmt.Sprintf("Obtained error %v when expecting no error", params[0])
	} else if params[0] == nil && params[1] == nil {
		return true, ""
	}
	return ErrorMatches.Check(params, names)
}
