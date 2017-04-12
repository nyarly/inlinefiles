package main

import (
	"io"
	"log"
	"regexp"
)

type escaper struct {
	r       io.Reader
	old     []byte
	debug   bool
	already int
}

var doublesRE = regexp.MustCompile(`"`)
var newlsRE = regexp.MustCompile("(?m)\n")

func newEscaper(r io.Reader, d bool) *escaper {
	return &escaper{r, make([]byte, 0), d, -1}
}

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
