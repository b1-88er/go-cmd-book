package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	assert.Equal(t, count(b, false, false), 4)
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("one\ntwo\nthree")
	assert.Equal(t, count(b, false, true), 3)
}

func TestCountByes(t *testing.T) {
	b := bytes.NewBufferString("12345")
	assert.Equal(t, count(b, true, false), 5)
}
