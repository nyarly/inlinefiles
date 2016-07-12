package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"

	"golang.org/x/tools/imports"
)

// started from http://stackoverflow.com/questions/17796043/golang-embedding-text-file-into-compiled-executable

const (
	docstring = `Inlines files into a Go source file
Usage: inlinefiles [--package=<name>] [--ext=<suffix>] <source_dir> <output_path>

Options:
  --package Force the name of the package, instead of guessing based on output_path
	--ext Use <suffix> for inlined files instead of ".tmpl"
	`

	header = `// This file was automatically generated based on the contents of *.tmpl
// If you need to update this file, change the contents of those files
// (or add new ones) and run 'go generate'

`
)

func main() {
	parsed, err := docopt.Parse(docstring, nil, true, "", false)
	if err != nil {
		log.Fatal(err)
	}

	sourceDir := parsed[`<source_dir>`].(string)
	targetPath := parsed[`<output_path>`].(string)
	png, ok := parsed[`<name>`]
	if !ok || png == nil {
		absTgt, err := filepath.Abs(targetPath)
		if err != nil {
			log.Fatal(err)
		}
		png = filepath.Base(filepath.Dir(absTgt))
	}
	extg, ok := parsed[`<suffix>`]
	if !ok || extg == nil {
		extg = ".tmpl"
	}
	pn := png.(string)
	ext := extg.(string)

	out := &bytes.Buffer{}
	file, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	fs, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	out.Write([]byte(header))
	out.Write([]byte("package " + pn + "\n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ext) {
			out.Write([]byte(strings.TrimSuffix(f.Name(), ext) + "Tmpl = \""))
			f, err := os.Open(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			r := newEscaper(f)

			io.Copy(out, r)
			out.Write([]byte("\"\n\n"))
		}
	}
	out.Write([]byte(")\n"))
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory")
	}
	fullpath := filepath.Join(cwd, targetPath)

	formattedBytes, err := imports.Process(fullpath, out.Bytes(), nil)
	if err != nil {
		log.Print("Problem formatting to ", fullpath, ": ", err)
		formattedBytes = out.Bytes()
	}
	file.Write(formattedBytes)
}

type escaper struct {
	r   io.Reader
	old []byte
}

var doublesRE = regexp.MustCompile(`"`)
var newlsRE = regexp.MustCompile("(?m)\n")

func (e *escaper) Read(p []byte) (n int, err error) {
	new := make([]byte, len(p)-len(e.old))
	c, err := e.r.Read(new)
	new = append(e.old, new[0:c]...)

	i, n := 0, 0
	for ; i < len(new) && n < len(p); i, n = i+1, n+1 {
		switch new[i] {
		default:
			p[n] = new[i]
		case '"', '\\':
			p[n] = '\\'
			n++
			p[n] = new[i]
		case '\n':
			p[n] = '\\'
			n++
			p[n] = 'n'
		}
	}
	if len(p) < i {
		e.old = new[len(new)-(len(p)-i):]
	} else {
		e.old = new[0:0]
	}

	log.Print(i, "/", n, "\n", len(e.old), ":", string(e.old), "\n", len(p), ":", string(p), "\n\n**************************\n\n")

	return
}

func newEscaper(r io.Reader) *escaper {
	return &escaper{r, make([]byte, 0)}
}
