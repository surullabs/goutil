// Copyright (c) 2014, Surul Software Labs, GmbH
// All rights reserved
//
// This contains utilities for working with Go templates.
package template

import (
	"os"
	"path/filepath"
	"text/template"
)

// A batch templater is used to apply the same data to a batch of
// templates.
type BatchTemplater interface {
	// The keys in the provided map are read as templates from
	// a data source. The result of executing each template
	// with the provided data is written into the data source named by the
	// corresponding value. Any errors will result in the operation
	// being aborted.
	Execute(templates map[string]string, data interface{}) (err error)
}

// A FileBatchTemplater is a BatchTemplater that uses files as a data source and sink
type FileBatchTemplater struct {
	Perm os.FileMode
}

// Creates a new batched templater which will write files out with the specified
// permissions.
func NewFileBatchTemplater(perm os.FileMode) *FileBatchTemplater {
	return &FileBatchTemplater{Perm: perm}
}

// This accepts a glob and a target directory and converts it to a map
// containing entries where the key is a matching file and the value is
// a file in the target directory having the same name.
func GlobBatch(glob, target string) (map[string]string, error) {
	if matches, err := filepath.Glob(glob); err != nil {
		return nil, err
	} else {
		result := make(map[string]string)
		for _, m := range matches {
			result[m] = filepath.Join(target, filepath.Base(m))
		}
		return result, nil
	}
}

func (t *FileBatchTemplater) executeOneFile(src, dest string, data interface{}) (err error) {
	var tpl *template.Template
	if tpl, err = template.ParseFiles(src); err != nil {
		return
	}
	return WriteFile(dest, tpl, data, t.Perm)
}

func (t *FileBatchTemplater) Execute(files map[string]string, data interface{}) (err error) {
	for src, dest := range files {
		if err = t.executeOneFile(src, dest, data); err != nil {
			return
		}
	}
	return
}
