package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.875, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}

	testCases := []struct {
		name string
		op   statsFunc
		exp  []float64
	}{
		{"sum", sum, []float64{300, 86.515, -30, 436}},
		{"avg", avg, []float64{37.5, 6.655, -15, 72.66666666666667}},
	}
	for _, testCase := range testCases {
		for i, testExp := range testCase.exp {
			name := fmt.Sprintf("%s_%d", testCase.name, i)
			t.Run(name, func(t *testing.T) {
				result := testCase.op(data[i])
				assert.Equal(t, testExp, result)
			})
		}
	}
}

func TestCSV2Float(t *testing.T) {
	csvData := `ip addr,reqs,response time
192.168.0.199,2056,236
192.168.0.88,899,220
192.168.0.199,3054,226
192.168.0.100,4133,218
192.168.0.199,950,238
`

	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{
			name: "col2",
			col:  2,
			exp:  []float64{2056, 899, 3054, 4133, 950}, expErr: nil,
			r: bytes.NewBufferString(csvData),
		},
		{
			name: "col3",
			col:  3,
			exp:  []float64{236, 220, 226, 218, 238}, expErr: nil,
			r: bytes.NewBufferString(csvData),
		},
		{
			name:   "FileRead",
			col:    1,
			exp:    nil,
			expErr: iotest.ErrTimeout,
			r:      iotest.TimeoutReader(bytes.NewReader([]byte{0})),
		},
		{
			name:   "FileNotNumber",
			col:    1,
			exp:    nil,
			expErr: ErrNotNumber,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "FailedInvalidCol",
			col:    4,
			exp:    nil,
			expErr: ErrInvalidColum,
			r:      bytes.NewBufferString(csvData),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := csv2float(testCase.r, testCase.col)
			assert.ErrorIs(t, err, testCase.expErr)
			assert.Equal(t, testCase.exp, result)
		})
	}
}
