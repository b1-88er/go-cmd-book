package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "RunAvgFile",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			files:  []string{"testdata/example.csv"},
			expErr: nil,
		},
		{
			name:   "RunAvgMultipleFiles",
			col:    3,
			op:     "avg",
			exp:    "235.92\n",
			files:  []string{"testdata/example.csv", "testdata/example2.csv"},
			expErr: nil,
		},
		{
			name:   "RunFailRead",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{"testdata/example.csv", "testdata/fakefile.csv"},
			expErr: os.ErrNotExist,
		},
		{
			name:   "RunFailColumn",
			col:    0,
			op:     "avg",
			exp:    "",
			files:  []string{"testdata/example.csv"},
			expErr: ErrInvalidColum,
		},
		{
			name:   "RunFailNoFiles",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name:   "RunFailOperation",
			col:    2,
			op:     "invalid",
			exp:    "",
			files:  []string{"testdata/example.csv"},
			expErr: ErrInvalidOperation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var result bytes.Buffer
			err := run(testCase.files, testCase.op, testCase.col, &result)
			assert.ErrorIs(t, err, testCase.expErr)
			assert.Equal(t, testCase.exp, result.String())
		})
	}
}

func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob("testdata/benchmark/*.csv")
	assert.Nil(b, err)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
