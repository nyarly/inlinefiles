package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscaperReadLength(t *testing.T) {
	assert := assert.New(t)

	testLen := 51200

	needsIt := []byte{byte('"'), byte('\n'), byte('\\')}

	floor := int('\\')
	if int('\n') > floor {
		floor = int('\n')
	}
	if int('"') > floor {
		floor = int('"')
	}
	ceil := 256 - floor

	bs := []byte{}
	for i := 0; i < testLen; i++ {
		bs = append(bs, byte((i%ceil)+floor+1))
	}
	assert.Len(bs, testLen)

	esc := newEscaper(bytes.NewBuffer(append(needsIt, bs...)), false)
	escd, err := ioutil.ReadAll(esc)
	assert.NoError(err)
	assert.Equal([]byte("\\\"\\n\\\\"), escd[:6])
	assert.Equal(bs, escd[6:]) //stripping escapables from prefix
}

func TestBoundaryEscape(t *testing.T) {
	oneByteAtATime(t, "\n", `\n`)
	oneByteAtATime(t, "\\", `\\`)
	oneByteAtATime(t, "\"", `\"`)
	oneByteAtATime(t, "\n\\\"\\\\\"\n", `\n\\\"\\\\\"\n`)
}

func oneByteAtATime(t *testing.T, r, ex string) {
	assert := assert.New(t)

	bs := []byte(r)
	esc := newEscaper(bytes.NewBuffer(bs), false)

	for i, ch := range []byte(ex) {
		into := make([]byte, 1)
		assert.NotPanics(func() {
			n, err := esc.Read(into)
			assert.NoError(err)
			assert.Equal(n, 1)
			assert.Equal(string(ch), string(into), "index: %d", i)
		})
	}

}

func TestEscaperEscapes(t *testing.T) {
	assert := assert.New(t)

	bs := []byte{byte('"'), byte('\n'), byte('\\')}

	esc := newEscaper(bytes.NewBuffer(bs), false)
	escd, err := ioutil.ReadAll(esc)
	assert.NoError(err)
	assert.Equal("\\\"\\n\\\\", string(escd))
}
