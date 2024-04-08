package add

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name string
		a    int
		b    int
		exp  int
	}{
		{
			name: "success",
			a:    1,
			b:    2,
			exp:  3,
		},
		{
			name: "negative",
			a:    -1,
			b:    2,
			exp:  1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := add(testCase.a, testCase.b)
			assert.Equal(t, testCase.exp, got)
		})
	}
}
