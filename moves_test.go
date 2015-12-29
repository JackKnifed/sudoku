package sudoku

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func loadCluster(inputFile string) ([][]cell, error) {
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	var v [][]cell
	err = Unmarshal(b, v)
	if err != nil {
		return err
	}
	return v
}

func loadBools(inputFile string) ([]bool, error) {
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	var v []bool
	err = Unmarshal(b, v)
	if err != nil {
		return err
	}
	return v
}

func loadIndex(inputFile string) ([]bool, error) {
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	var v []indexedCluster
	err = Unmarshal(b, v)
	if err != nil {
		return err
	}
	return v
}

func updateCatcher(in, out chan cell) {
	defer close(out)
	more := true
	var updates Cluster
	var val cell
	for more {
		val, more = <-in
		if more {
			updates = append(updates, val)
		}
	}
	sort.Sort(updates)
	for _, each := range updates {
		out <- each
	}
	return
}

func TestIndexCluster(t *testing.T) {
	input, err := loadCluster("testdata/moves_input.json")
	if err != nil {
		t.FatalF("expected input could not be loaded - %v", err)
	}

	expected, err := loadIndex("testdata/moves_index.json")
	if err != nil {
		t.FatalF("expected cound not be loaded - %v", err)
	}

	if len(input) != len(expected) {
		t.FatalF("input and expected do not have the same number of tests")
	}

	for i := 0; i < len(input); i++ {
		output := indexCluster(input[i])
		assert.Equal(t, expected, output, "test %d - expected index and returned index differ", i)
	}
}

func TestClusterSolved(t *testing.T) {
	input, err := loadCluster("testdata/moves_input.json")
	if err != nil {
		t.FatalF("expected input could not be loaded - %v", err)
	}

	expected, err := loadBools("testdata/moves_one.json")
	if err != nil {
		t.FatalF("expected cound not be loaded - %v", err)
	}

	if len(input) != len(expected) {
		t.FatalF("input and expected do not have the same number of tests")
	}

	for i := 0; i < len(input); i++ {
		var update chan cell
		var more bool
		result := clusterSolved(input[i], update)
		select {
		case _, more = <-update:
		default:
			more = true
		}

		assert.Equal(t, expected[i], result, "test %d - return value differs", i)
		assert.Equal(t, expected[i], more, "test %d - channel status differs", i)
	}
}

func TestSolvedNoPossible(t *testing.T) {
	input, err := loadCluster("testdata/moves_input.json")
	if err != nil {
		t.FatalF("input could not be loaded - %v", err)
	}

	expectedResponse, err := loadCluster("testdata/moves_two_response.json")
	if err != nil {
		t.FatalF("expected cound not be loaded - %v", err)
	}

	expected, err := loadCluster("testdata/moves_two_updates.json")
	if err != nil {
		t.FatalF("expected cound not be loaded - %v", err)
	}

	if len(input) != len(expected) {
		t.FatalF("input and expected count differs")
	}

	for i := 0; i < len(input); i++ {
		var updates, back chan cell
		var more bool
		go updateCatcher(updates, back)
		result := clusterSolved(input[i], update)

		assert.Equal(t, expectedResponse[i], result, "expected and actual return differs")
		close(updates)
		sort.Sort(expected[i])

		for _, update := range expected[i] {
			assert.Equal(t, expected[i], <-back, "test %d - return value differs", i)
		}

		junk, more := <-update
		if more {
			t.Fail()
			t.Logf("test %d - channel status differs", i)
			for more {
				t.Log(junk)
				junk, more = <-update
			}
		}
	}
}
