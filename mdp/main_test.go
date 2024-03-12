package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	assert.Nil(t, err)

	result, err := parseContent(input, "")

	assert.Nil(t, err)

	expected, err := os.ReadFile(goldenFile)

	assert.Nil(t, err)

	assert.Equal(t, string(expected), string(result))
	os.Remove(resultFile)
}

func TestRun(t *testing.T) {
	var mockStdout bytes.Buffer
	assert.Nil(t, run(inputFile, "", &mockStdout, true))
	resultFile := strings.TrimSpace(mockStdout.String())
	result, err := os.ReadFile(resultFile)
	assert.Nil(t, err)

	expected, err := os.ReadFile(goldenFile)
	assert.Nil(t, err)

	assert.Equal(t, string(expected), string(result))
	os.Remove(resultFile)
}
