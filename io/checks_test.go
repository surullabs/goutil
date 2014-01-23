// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
package io

import (
	"errors"
	. "github.com/surullabs/goutil/testing"
	. "launchpad.net/gocheck"
	"os"
)

type ChecksSuite struct{}

var _ = Suite(&ChecksSuite{})

func (s *ChecksSuite) TestExists(c *C) {
	defer func() { stat = os.Stat }()

	file := CreateTestFile(c, "file", []byte("data"), 0700)
	if exists, err := Exists(file); err != nil {
		c.Fatal(err)
	} else if !exists {
		c.Fatal("Should exist")
	}

	os.RemoveAll(file)
	if exists, err := Exists(file); err != nil {
		c.Fatal(err)
	} else if exists {
		c.Fatal("Should not exist")
	}

	stat = func(s string) (os.FileInfo, error) { return nil, errors.New("Stat failed") }
	_, err := Exists(file)
	c.Check(err, ErrorMatches, "Stat failed")
}

func (s *ChecksSuite) TestIsDir(c *C) {
	defer func() { stat = os.Stat }()

	file := CreateTestFile(c, "file", []byte("data"), 0700)
	if isDir, err := IsDir(file); err != nil {
		c.Fatal(err)
	} else if isDir {
		c.Fatal("Should not be dir")
	}

	os.RemoveAll(file)
	if isDir, err := IsDir(file); err != nil {
		if !os.IsNotExist(err) {
			c.Fatal(err)
		}
	} else if isDir {
		c.Fatal("Should not exist")
	}

	dir := c.MkDir()
	if isDir, err := IsDir(dir); err != nil {
		c.Fatal(err)
	} else if !isDir {
		c.Fatal("Should be directory")
	}

	stat = func(s string) (os.FileInfo, error) { return nil, errors.New("Stat failed") }
	_, err := IsDir(file)
	c.Check(err, ErrorMatches, "Stat failed")
}
