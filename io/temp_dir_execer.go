// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
// Contains io utilities. This is currently used as a dumping ground for all
// Surul IO Utilities.
package io

import (
	"fmt"
	"io/ioutil"
	"os"
)

// A TempDirExecer must create a temporary directory and then run the provided
// function with the directory as an argument. It must return any error that the
// function returns. The directory must be deleted by the TempDirExecer after
// the function has completed.
type TempDirExecer interface {
	// Create a temporary directory with 'prefix' as a name prefix and run 'fn' in
	// this directory.
	Exec(prefix string, fn func(string) error) error
}

var remover func(string) error = os.RemoveAll

// A SafeTempDirExecer exposes the TempDirExecer interface and guarantees cleanup of the
// temporary directory. If cleanup fails Exec will panic.
type SafeTempDirExecer struct {
	// Directory in which the temporary directory will be created. If empty the default
	// temporary directory is used.
	Location string
}

// Create a temporary directory in s.Location using 'prefix'. 'fn' is
// then run with the full directory path as the argument. If any errors occur when deleting
// the directory this will panic.
func (s *SafeTempDirExecer) Exec(prefix string, fn func(string) error) (err error) {
	var rootTempDir string
	if rootTempDir, err = ioutil.TempDir(s.Location, prefix); err != nil {
		return
	}
	defer func() {
		failure := recover()
		// Now delete the temporary directory
		if removeErr := remover(rootTempDir); removeErr != nil {
			// There are 3 possible errors here. The panic, the function error
			// and this remove error. Under normal circumstances it out be ok
			// to return an error, but here we could be leaking passwords
			// through the template script and so this will panic
			if failure != nil {
				failure = fmt.Errorf("Original Panic: %v; Failed to clear directory: %v", failure, removeErr)
			} else {
				failure = fmt.Errorf("Failed to clear directory: %v", removeErr)
			}
		}
		if failure != nil {
			panic(failure)
		}
	}()

	err = fn(rootTempDir)
	return
}
