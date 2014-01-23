// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
package io

import (
	"os"
	"time"
)

var stat func(string) (os.FileInfo, error) = os.Stat

func Exists(path string) (bool, error) {
	if _, err := stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func IsDir(path string) (bool, error) {
	if s, err := stat(path); err != nil {
		return false, err
	} else {
		return s.IsDir(), nil
	}

}

// This will wait for the path to exist for a Duration of timeout,
// polling every tick seconds. If the file exists within the timeout
// it will return nil. If not an error will be returned. You can check
// if it is just due to the file not existing by calling os.IsNotExist
func WaitTillExists(path string, tick, timeout time.Duration) (err error) {
	end := time.Now().Add(timeout)

	if _, err = stat(path); err == nil {
		return
	}

	for time.Now().Before(end) && err != nil {
		time.Sleep(tick)
		_, err = stat(path)
	}

	return
}
