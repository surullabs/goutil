// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved.
//
// This provides various testing utilities that are compatible with
// the gocheck package.
package testing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"os/exec"
	"path/filepath"
)

// This creates a file named fileName in a temporary directory
// specifically created for this file. The contents are then
// written to the file. This file will be deleted at the end
// of a test. If creation of the file fails, the test will be
// aborted. The method will return the full path to the created
// file.
func CreateTestFile(c *C, fileName string, contents []byte, perm os.FileMode) string {
	topDir := c.MkDir()
	filePath := filepath.Join(topDir, "file")
	if err := ioutil.WriteFile(filePath, contents, perm); err != nil {
		c.Fatal(err)
	}
	return filePath
}

type fileMatcher struct {
	*CheckerInfo
}

var FileMatches = &fileMatcher{
	&CheckerInfo{Name: "FileMatches", Params: []string{"obtained", "expected"}},
}

func (n *fileMatcher) Check(params []interface{}, names []string) (result bool, errorStr string) {
	data1, err1 := ioutil.ReadFile(params[0].(string))
	if err1 != nil {
		return false, err1.Error()
	}
	data2, err2 := ioutil.ReadFile(params[1].(string))
	if err2 != nil {
		return false, err2.Error()
	}
	if bytes.Compare(data1, data2) != 0 {
		return false, fmt.Sprintf("Contents of %v and %v do not match", params[0], params[1])
	}
	return true, ""
}

type directoryMatcher struct {
	*CheckerInfo
}

var DirectoryMatches = &directoryMatcher{
	&CheckerInfo{Name: "DirectoryMatches", Params: []string{"obtained", "expected"}},
}

func bothDirOrMissing(dir1, dir2 string) (bool, string) {
	stat1, err1 := os.Stat(dir1)
	stat2, err2 := os.Stat(dir2)
	if err1 != nil || err2 != nil {
		if os.IsNotExist(err1) && os.IsNotExist(err2) {
			return true, ""
		}
		if err1 != nil {
			return false, err1.Error()
		}
		return false, err2.Error()
	}
	if !stat1.IsDir() {
		return false, fmt.Sprintf("%s is not a directory", dir1)
	}
	if !stat2.IsDir() {
		return false, fmt.Sprintf("%s is not a directory", dir2)
	}
	return true, ""
}

func (n *directoryMatcher) Check(params []interface{}, names []string) (result bool, errorStr string) {
	dir1, dir2 := params[0].(string), params[1].(string)
	if bothDir, err := bothDirOrMissing(dir1, dir2); !bothDir {
		return bothDir, err
	}

	walkMatcher := func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir1, path)
		if err != nil {
			return err
		}
		path2 := filepath.Join(dir2, rel)

		if bothDir, _ := bothDirOrMissing(path, path2); bothDir {
			return nil
		}

		if matches, errStr := FileMatches.Check([]interface{}{path, path2}, names); !matches {
			return fmt.Errorf("%s", errStr)
		}
		return nil
	}
	if err := filepath.Walk(dir1, walkMatcher); err != nil {
		return false, err.Error()
	}
	return true, ""
}

type directoryContains struct {
	*CheckerInfo
}

var HasFilesNamed = &directoryContains{
	&CheckerInfo{Name: "HasFilesNamed", Params: []string{"obtained", "expected"}},
}

func (n *directoryContains) Check(params []interface{}, names []string) (result bool, errorStr string) {
	dirName := params[0].(string)
	fileNames := make(map[string]struct{})
	exists := struct{}{}
	if filesInDir, err := filepath.Glob(filepath.Join(dirName, "*")); err != nil {
		return false, err.Error()
	} else {
		for _, file := range filesInDir {
			fileNames[filepath.Base(file)] = exists
		}
	}
	filesToCheck := params[1].([]string)
	for _, file := range filesToCheck {
		if _, exists := fileNames[file]; !exists {
			return false, fmt.Sprintf("%s does not exist", file)
		}
	}
	return true, ""
}

// Function to copy a directory tree in UNIX-ey OSes. This
// just runs cp -r
func UnsafeUnixCopyTree(src, dest string) (err error) {
	if output, err := exec.Command("cp", "-r", src, dest).CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return
}
