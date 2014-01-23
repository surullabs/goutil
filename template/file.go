// Copyright 2014, Surul Software Labs GmbH
// All rights reserved.
//
package template

import (
	"os"
	"text/template"
)

// Utility to write a template out to a file
func WriteFile(file string, tpl *template.Template, data interface{}, mode os.FileMode) error {
	if f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode); err != nil {
		return err
	} else {
		defer f.Close()
		return tpl.Execute(f, data)
	}
}
