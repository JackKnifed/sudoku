package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDedupArr(t *testing.T) {
	var tests = []struct {
		in  []int
		out []int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		}, {
			[]int{1, 1, 3, 4, 4, 4},
			[]int{1, 3, 4},
		}, {
			[]int{2, 1, 8, 4, 16, 6},
			[]int{1, 2, 4, 6, 8, 16},
		}, {
			[]int{2, 1, 8, 4, 16, 6, 16, 1, 4},
			[]int{1, 2, 4, 6, 8, 16},
		},
	}

	for _, testRun := range tests {
		var input, output []int
		input = testRun.in[:]
		output = dedupArr(input)

		assert.Equal(t, testRun.in, input, "input changed")
		assert.Equal(t, testRun.out, output, "output was incorrect")
	}
}
