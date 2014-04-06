// Copyright 2014, Surul Software Labs GmbH
// All rights reserved.

package template

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"text/template"
)

// Utility to write a template out to a file
func WriteFile(file string, tpl *template.Template, data interface{}, mode os.FileMode) (err error) {
	var f *os.File
	if f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode); err != nil {
		return
	} else {
		defer f.Close()
		if err = tpl.Execute(f, data); err != nil {
			return
		}
	}

	// Workaround for https://code.google.com/p/go/issues/detail?id=6288
	// This is a hack!
	var bytes []byte
	if bytes, err = ioutil.ReadFile(file); err != nil {
		return
	}

	var re *regexp.Regexp
	if re, err = regexp.Compile(`<no value[^>]*>`); err != nil {
		panic(err)
	}

	if re.FindIndex(bytes) != nil {
		return errors.New("Missing variables present")
	}
	return
}
