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

// ErrorIsNil checker.

type errorIsNilChecker struct {
	*CheckerInfo
}

// The IsNil checker tests whether the obtained value is nil.
//
// For example:
//
//    c.Assert(err, IsNil)
//
var ErrorIsNil Checker = &errorIsNilChecker{
	&CheckerInfo{Name: "ErrorIsNil", Params: []string{"value"}},
}

func (checker *errorIsNilChecker) Check(params []interface{}, names []string) (result bool, err string) {
	switch v := params[0].(type) {
	case nil:
		return true, ""
	case error:
		return false, v.Error()
	default:
		return false, fmt.Sprintf("expected error, found %T:%v", v, v)
	}
}
