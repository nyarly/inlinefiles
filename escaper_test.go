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

func TestEscaperEscapes(t *testing.T) {
	assert := assert.New(t)

	bs := []byte{byte('"'), byte('\n'), byte('\\')}

	esc := newEscaper(bytes.NewBuffer(bs), false)
	escd, err := ioutil.ReadAll(esc)
	assert.NoError(err)
	assert.Equal("\\\"\\n\\\\", string(escd))
}
