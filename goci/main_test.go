package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		stderr string
		expErr *stepErr
	}{
		{
			name:   "success",
			proj:   "testdata/tool",
			out:    "Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS",
			stderr: "",
			expErr: nil,
		},
		{
			name:   "validation error",
			proj:   "testdata/toolErr",
			out:    "",
			expErr: &stepErr{step: "go build"},
		},
		{
			name:   "format error",
			proj:   "testdata/toolFmtErr",
			out:    "Go build: SUCCESS\nGo test: SUCCESS\n",
			expErr: &stepErr{step: "go fmt"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := bytes.Buffer{}
			err := run(testCase.proj, &out)
			assert.Equal(t, testCase.out, out.String())

			// both assertions do the same thing
			if err != nil {
				assert.ErrorIs(t, testCase.expErr, err)

			}
			if expErr, ok := (err).(*stepErr); ok {
				assert.Equal(t, expErr.step, testCase.expErr.step)
			}

		})
	}
}
