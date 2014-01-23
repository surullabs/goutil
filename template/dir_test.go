// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
package template

import (
	. "github.com/surullabs/goutil/testing"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type FileTemplaterSuite struct{}

var _ = Suite(&FileTemplaterSuite{})

func (s *FileTemplaterSuite) TestFileTemplater(c *C) {
	tests := []struct {
		glob        string
		expectedErr interface{}
		vars        map[string]string
		targetDir   string
		golden      string
	}{
		{
			glob:      "testdata/*.tpl",
			vars:      map[string]string{"var1": "value1", "var2": "value2"},
			targetDir: c.MkDir(),
			golden:    "testdata/expected",
		},
		{
			glob:        "testdata/*.bad",
			expectedErr: ".*unexpected unclosed action.*",
		},
		{
			glob:        "testdata/*.tpl",
			vars:        map[string]string{"var1": "value1", "var2": "value2"},
			targetDir:   c.MkDir() + "/NonExistentDirectory/",
			expectedErr: ".*no such file or directory.*",
		},
	}

	for _, test := range tests {
		matches, _ := GlobBatch(test.glob, test.targetDir)
		err := NewFileBatchTemplater(0600).Execute(matches, test.vars)
		c.Assert(err, NilOrErrorMatches, test.expectedErr)
		if err == nil {
			c.Assert(test.targetDir, DirectoryMatches, test.golden)
		}
	}
}
