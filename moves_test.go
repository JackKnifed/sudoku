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

func Test
