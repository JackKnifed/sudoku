package sudoku

import (
	"fmt"
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
		}, {
			[]int{},
			[]int{},
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

func ExampleDedupArr() {
	input := []int{2, 1, 8, 4, 16, 6, 16, 1, 4}
	fmt.Println(dedupArr(input))

	// Output:
	// [1 2 4 6 8 16]
}

func TestInArr(t *testing.T) {
	var tests = []struct {
		in    []int
		check int
		out   bool
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			3, true,
		}, {
			[]int{1, 1, 3, 4, 4, 4},
			7, false,
		},
	}

	for _, testRun := range tests {
		result := inArr(testRun.in, testRun.check)
		assert.Equal(t, testRun.out, result, "wrong response - expected %d", testRun.check)
	}
}

func TestAnyInArr(t *testing.T) {
	var tests = []struct {
		in    []int
		check []int
		out   bool
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			true,
		}, {
			[]int{1, 1, 3, 4, 4, 4},
			[]int{1},
			true,
		}, {
			[]int{1, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			true,
		}, {
			[]int{1, 1, 3, 4, 4, 4},
			[]int{7},
			false,
		}, {
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			true,
		}, {
			[]int{},
			[]int{1, 2, 3, 4, 5, 6},
			false,
		}, {
			[]int{1, 2, 3, 4, 5, 6},
			[]int{},
			false,
		},
	}

	for id, testRun := range tests {
		result := anyInArr(testRun.in, testRun.check)
		assert.Equal(t, testRun.out, result, "wrong response - %d test", id)
	}
}

func TestAddArr(t *testing.T) {
	var tests = []struct {
		a   []int
		b   []int
		out []int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		}, {
			[]int{1, 2, 3},
			[]int{4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		}, {
			[]int{1, 2, 3, 4, 6},
			[]int{1, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		}, {
			[]int{},
			[]int{1, 4, 5, 6},
			[]int{1, 4, 5, 6},
		}, {
			[]int{1, 2, 3, 4, 5, 6},
			[]int{},
			[]int{1, 2, 3, 4, 5, 6},
		},
	}

	for id, testRun := range tests {
		result := addArr(testRun.a, testRun.b)
		assert.Equal(t, testRun.out, result, "wrong response test %d", id)
	}
}

func TestSubArr(t *testing.T) {
	var tests = []struct {
		a   []int
		b   []int
		out []int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			[]int{},
		}, {
			[]int{1, 2, 3},
			[]int{4, 5, 6},
			[]int{1, 2, 3},
		}, {
			[]int{1, 2, 3, 4, 6},
			[]int{1, 4, 5, 6},
			[]int{2, 3},
		}, {
			[]int{},
			[]int{1, 4, 5, 6},
			[]int{},
		}, {
			[]int{1, 2, 3, 4, 5, 6},
			[]int{},
			[]int{1, 2, 3, 4, 5, 6},
		},
	}

	for id, testRun := range tests {
		result := subArr(testRun.a, testRun.b)
		assert.Equal(t, testRun.out, result, "wrong response test %d", id)
	}
}

func TestAndArr(t *testing.T) {
	var tests = []struct {
		a   []int
		b   []int
		out []int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		}, {
			[]int{1, 2, 3},
			[]int{4, 5, 6},
			[]int{},
		}, {
			[]int{1, 2, 3, 7},
			[]int{4, 5, 6, 7},
			[]int{7},
		}, {
			[]int{1, 2, 3, 4, 6},
			[]int{1, 4, 5, 6},
			[]int{1, 4, 6},
		}, {
			[]int{},
			[]int{1, 4, 5, 6},
			[]int{},
		}, {
			[]int{1, 2, 3, 4, 5, 6},
			[]int{},
			[]int{},
		},
	}

	for id, testRun := range tests {
		result := andArr(testRun.a, testRun.b)
		assert.Equal(t, testRun.out, result, "wrong response test %d", id)
	}
}
