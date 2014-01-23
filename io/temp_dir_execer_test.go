// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
package io

import (
	"errors"
	"fmt"
	. "github.com/surullabs/goutil/testing"
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type SafeTempDirExecerSuite struct{}

var _ = Suite(&SafeTempDirExecerSuite{})

type safeTempExecerTest struct {
	name          string
	root          string
	prefix        string
	panicWith     interface{}
	errorWith     interface{}
	expectedPanic interface{}
	expectedError interface{}
	fn            func(string) error
}

func (t *safeTempExecerTest) Test(c *C) {
	path, err, panicResult := t.exec()
	c.Check(err, NilOrErrorMatches, t.expectedError)
	c.Check(panicResult, NilOrErrorMatches, t.expectedPanic)
	if exists, existsErr := Exists(path); existsErr != nil {
		c.Fatal(existsErr)
	} else if exists {
		c.Fatalf("Directory %s not deleted", path)
	}
}

func (test *safeTempExecerTest) exec() (dirPath string, err error, panicResult interface{}) {
	defer func() {
		panicResult = recover()
	}()

	s := &SafeTempDirExecer{Location: test.root}
	if test.fn == nil {
		test.fn = func(path string) error {
			dirPath = path
			if test.panicWith != nil {
				panic(fmt.Errorf("%v", test.panicWith))
			}
			if test.errorWith != nil {
				return fmt.Errorf("%v", test.errorWith)
			}
			if isDir, dirErr := IsDir(dirPath); dirErr != nil {
				return dirErr
			} else if !isDir {
				return fmt.Errorf("%s not created", dirPath)
			}
			return nil
		}
	}
	err = s.Exec(test.prefix, test.fn)
	return
}

func (s *SafeTempDirExecerSuite) TestSafeTempDirExecer(c *C) {
	tests := []safeTempExecerTest{
		{
			name: "Default root with empty prefix",
		},
		{
			name:   "Default root with non-empty prefix",
			prefix: "prefix1",
		},
		{
			root:          CreateTestFile(c, "file", []byte("data"), 0600),
			expectedError: "mkdir.*not a directory",
		},
		{
			errorWith:     "This is a test error",
			expectedError: "This is a test error",
		},
		{
			panicWith:     "This is a test panic",
			expectedPanic: "This is a test panic",
		},
	}
	for _, test := range tests {
		test.Test(c)
	}
}

func (s *SafeTempDirExecerSuite) TestRemoveFailure(c *C) {
	oldRemover := remover
	remover = func(d string) error {
		os.RemoveAll(d)
		return errors.New("Remove failed")
	}
	defer func() { remover = oldRemover }()

	tests := []safeTempExecerTest{
		{
			name:          "Remove failed",
			prefix:        "prefix1",
			expectedPanic: "Failed to clear directory: Remove failed",
		},
		{
			errorWith:     "This is a test error",
			expectedPanic: "Failed to clear directory: Remove failed",
		},
		{
			panicWith:     "This is a test panic",
			expectedPanic: "Original Panic: This is a test panic; Failed to clear directory: Remove failed",
		},
	}
	for _, test := range tests {
		test.Test(c)
	}

}
