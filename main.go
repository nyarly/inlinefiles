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
	"text/template"

	"github.com/docopt/docopt-go"

	"golang.org/x/tools/imports"
)

//go:generate inlinefiles --package=main . templates.go

const (
	docstring = `Inlines files into a Go source file
Usage: inlinefiles [options] <source_dir> <output_path>

Options:
	-d --debug        Debug output
	--package=<name>  Force the name of the package, instead of guessing based on output_path
	--ext=<suffix>    Use <suffix> for inlined files. Equivalent to --glob='*<suffix>'
	--glob=<pattern>  Use <pattern> to restrict files included. Default: '*'
	--vfs=<name>      Put the templates into a mapfs in a variable called <vfs_name>
	`
)

type rootCtx struct {
	PackageName string
	Templates   []templateCtx
	MapFSName   string
}

type templateCtx struct {
	SourceFile   string
	SourceReader io.Reader
}

func (c templateCtx) ConstantName() string {
	ext := filepath.Ext(c.SourceFile)
	return strings.TrimSuffix(c.SourceFile, ext) + "Tmpl"
}

func (c templateCtx) Contents() (string, error) {
	b, e := ioutil.ReadAll(c.SourceReader)
	return string(b), e
}

func main() {
	parsed, err := docopt.Parse(docstring, nil, true, "", false)
	if err != nil {
		log.Fatal(err)
	}

	debug := parsed[`--debug`].(bool)

	var tmpl *template.Template
	ctx := rootCtx{}

	sourceDir := parsed[`<source_dir>`].(string)
	targetPath := parsed[`<output_path>`].(string)
	png, ok := parsed[`--package`]
	if !ok || png == nil {
		absTgt, err := filepath.Abs(targetPath)
		if err != nil {
			log.Fatal(err)
		}
		png = filepath.Base(filepath.Dir(absTgt))
	}
	ctx.PackageName = png.(string)

	glob := `*`
	if extg := parsed[`--ext`]; extg != nil {
		glob = `*` + extg.(string)
	}
	if globg := parsed[`--glob`]; globg != nil {
		glob = globg.(string)
	}

	mfg, useMapFS := parsed[`--vfs`]
	if useMapFS && mfg != nil {
		ctx.MapFSName = mfg.(string)
		tmpl = template.Must(template.New("root").Parse(useMapFSTmpl))
	} else {
		tmpl = template.Must(template.New("root").Parse(useConstantsTmpl))
	}

	out := &bytes.Buffer{}
	file, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			if f.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		match, err := filepath.Match(glob, f.Name())
		if err != nil {
			return err
		}

		if match {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			sourcePath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}

			ctx.Templates = append(ctx.Templates, templateCtx{
				SourceFile:   sourcePath,
				SourceReader: newEscaper(f, debug),
			})
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(out, ctx)

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
	r       io.Reader
	old     []byte
	debug   bool
	already int
}

func newEscaper(r io.Reader, d bool) *escaper {
	return &escaper{r, make([]byte, 0), d, -1}
}

var doublesRE = regexp.MustCompile(`"`)
var newlsRE = regexp.MustCompile("(?m)\n")

func (e *escaper) Read(p []byte) (n int, err error) {
	var new []byte
	if len(p) > len(e.old) {
		new = make([]byte, len(p)-len(e.old))
		var c int
		c, err = e.r.Read(new)

		if err != nil {
			if e.debug {
				log.Print(err == io.EOF, err)
			}
			if err != io.EOF {
				return 0, err
			}
		}
		new = append(e.old, new[0:c]...)
	} else {
		new = e.old
	}

	i, n := 0, 0
	for ; i < len(new) && n < len(p); i, n = i+1, n+1 {
		switch new[i] {
		default:
			p[n] = new[i]
		case '"', '\\':
			if e.already != i {
				p[n] = '\\'
				e.already = i
				i--
			} else {
				p[n] = new[i]
			}
		case '\n':
			p[n] = '\\'
			new[i] = 'n'
			i--
		}
	}

	e.old = new[i:]
	if e.already == i {
		e.already = 0
	} else {
		e.already = -1
	}

	if e.debug {
		log.Print(i, "/", n, " - ", err)
	}

	if len(e.old) > 0 && err == io.EOF {
		err = nil
	}

	if i == 0 && n == 0 {
		err = io.EOF
	}

	if e.debug {
		log.Print(i, "/", n, " - ", err, "\n", len(e.old), ":", string(e.old), "\n", len(p), ":", string(p), "\n\n**************************\n\n")
	}
	return
}
