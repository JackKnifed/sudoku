package sudoku

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func loadCellCluster(inputFile string) ([][]cell, error) {
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	var v interface{}
	err = Unmarshal(b, v)
	if err != nil {
		return err
	}
	return v
}

func TestIndexCluster(t *testing.T) {
	input, err := loadCellCluster("testdata/moves_input.json")
	if err != nil {
		t.FatalF("expected input could not be loaded - %v", err)
	}

	expected, err := loadCellCluster("testdata/moves_index.json")
	if err != nil {
		t.FatalF("expected results cound not be loaded - %v", err)
	}

	if len(input) != len(expected) {
		t.FatalF("input and expected do not have the same number of tests")
	}

	for i := 0; i < len(input); i++ {
		output := indexCluster(input[i])
		sort.Sort(expected)
		sort.Sort(output)
		assert.Equal(t, expected, output, "expected index and returned index differ")
	}
}
