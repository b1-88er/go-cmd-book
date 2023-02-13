package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	assert.Nil(t, err)

	result := parseContent(input)

	expected, err := ioutil.ReadFile(goldenFile)

	assert.Nil(t, err)

	assert.Equal(t, string(expected), string(result))
	os.Remove(resultFile)
}

func TestRun(t *testing.T) {
	assert.Nil(t, run(inputFile))
	result, err := ioutil.ReadFile(resultFile)
	assert.Nil(t, err)

	expected, err := ioutil.ReadFile(goldenFile)
	assert.Nil(t, err)

	assert.Equal(t, string(expected), string(result))
	os.Remove(resultFile)
}
