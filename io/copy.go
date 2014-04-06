// Copyright 2014, Surul Software Labs GmbH
// All rights reserved.

package io

import (
	"github.com/surullabs/fault"
	"io"
	"os"
)

var check fault.FaultCheck = fault.NewChecker()

func CopyFile(src, dst string) (err error) {
	defer check.Recover(&err)
	srcfi := check.Return(os.Stat(src)).(os.FileInfo)
	check.Truef(srcfi.Mode().IsRegular(), "Unsupported file type for copy: %s %v", srcfi.Name(), srcfi.Mode())
	if dstfi, dErr := os.Stat(dst); dErr != nil {
		check.Truef(os.IsNotExist(dErr), "%v", dErr)
	} else {
		check.Truef(dstfi.Mode().IsRegular(), "Unsupported file type for copy: %s %v", dstfi.Name(), dstfi.Mode())
		check.Truef(os.SameFile(srcfi, dstfi), "Destination %s exists", dst)
		return
	}

	in := check.Return(os.Open(src)).(*os.File)
	defer in.Close()
	out := check.Return(os.Create(dst)).(*os.File)
	defer out.Close()
	check.Return(io.Copy(out, in))
	check.Error(out.Sync())
	return
}
