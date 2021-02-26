// Package templatestore provides a convenience for loading inlined templates.
package templatestore

import (
	"bytes"
	htmltmpl "html/template"
	txttmpl "text/template"

	"golang.org/x/tools/godoc/vfs"
)

// LoadText loads a simple template from a possibly inlined VFS
func LoadText(fs vfs.Opener, tName, fName string) (*txttmpl.Template, error) {
	return LoadTextOnto(fs, nil, tName, fName)
}

// LoadTextOnto loads a template from a possibly inlined VFS as an associated template to the parent
func LoadTextOnto(fs vfs.Opener, parent *txttmpl.Template, tName, fName string) (*txttmpl.Template, error) {
	src, err := templateSource(fs, fName)
	if err != nil {
		return nil, err
	}
	var tpl *txttmpl.Template
	if parent == nil {
		tpl = txttmpl.New(tName)
	} else {
		tpl = parent.New(tName)
	}
	return tpl.Parse(src)
}

// LoadHTML loads a simple template from a possibly inlined VFS
func LoadHTML(fs vfs.Opener, tName, fName string) (*htmltmpl.Template, error) {
	return LoadHTMLOnto(fs, nil, tName, fName)
}

// LoadHTMLOnto loads a template from a possibly inlined VFS as an associated template to the parent
func LoadHTMLOnto(fs vfs.Opener, parent *htmltmpl.Template, tName, fName string) (*htmltmpl.Template, error) {
	src, err := templateSource(fs, fName)
	if err != nil {
		return nil, err
	}
	var tpl *htmltmpl.Template
	if parent == nil {
		tpl = htmltmpl.New(tName)
	} else {
		tpl = parent.New(tName)
	}
	return tpl.Parse(src)
}

func templateSource(fs vfs.Opener, fName string) (string, error) {
	tmplFile, err := fs.Open(fName)
	if err != nil {
		return "", err
	}
	tmplB := &bytes.Buffer{}
	_, err = tmplB.ReadFrom(tmplFile)
	return tmplB.String(), err
}
