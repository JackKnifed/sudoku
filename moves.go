package sudoku

// Sudoku has one rule - no value can be repeated horizontally, vertically,
// or in the same section. This gives rise to a few simple rules - and those
// rules are applied across each cluster - without even knowning the orientation
// of the cluster.
//
// In order to make solving possible however, that one rule is enforced with
// the following rules:
//
// 1) If all cells are solved, that cluster is solved.
// 2) If any cell is solved, it has no possibles.
// 3) If any cell is solved, that value is not possible in other cells.
// 4) If any cell only has one possible value, that is that cell's value.
// 5)If any x cells have x possible values, other cells in the cluster cannot be
//  those values - those values are constrained to those cells.
// 6) If any value only has one possible cell, that is that cell's value.
// 7) If any x values have x possible cells, other values are not possible
//  in those cells - those cells are constrained to those values.
//
// Additional Helper functions are included and explained later.

// indexedCLuster is a datatype for an index for the values of the cluster.
// Each possible value is a key in the map. THe value of each key is an array of
// possible locations for that value - the index and order are not defined,
// instead the values of the array are the indexes of possible cells in for that
// value.
type intArray []int
type indexedCluster map[int]intArray

func indexCluster(in []cell) (out indexedCluster) {
	for id, each := range in {
		for onePossible, _ := range each.possible {
			out[onePossible] = append(out[onePossible], id)
		}
	}
	return out
}

// This covers the 4th rule from above:
// 4) If any cell only has one possible value, that is that cell's value.
func singleValueSolver(cluster []cell, u chan cell) (changed bool) {
	for _, workingCell := range cluster {
		// if the cell is already solved, skip this.
		if workingCell.actual != 0 {
			continue
		}

		// if there is more than one possible value for this cell, you should be good
		if len(workingCell.possible) > 1 {
			continue
		}

		// should never happen
		if len(workingCell.possible) < 1 {
			panic("Found an unsolved cell with no possible values")
		}

		changed = true

		// this empty for loop will set value to each key in the map
		// since there's one key, it gives me that one key
		var key int
		for key, _ = range workingCell.possible {
		}

		// send back an update for this cell
		u <- cell{
			location: workingCell.location,
			actual:   key,
		}
	}
	return changed
}

// Removes known values from other possibles in the same cluster
func eliminatePossibles(workingCluster []cell, u chan cell) bool {
	var solved map[int]bool
	var remove map[int]bool
	var changed bool

	for _, each := range workingCluster {
		if each.actual != 0 {
			solved[each.actual] = true
		}
	}

	for _, each := range workingCluster {
		for i, potential := range each.possible {
			if potential && solved[i] {
				remove[i] = true
				changed = true
			}
		}
		if len(remove) > 0 {
			// send back removal of possibles & reset
			u <- cell{
				location: each.location,
				possible: remove,
			}
			remove = map[int]bool{}
		}
	}
	return changed
}

func confirmIndexed(index indexedCluster, workingCluster []cell, u chan cell) bool {
	var changed bool
	for val, section := range index {
		// if len(section.targets) < 1 {
		if len(section) < 1 {
			// something went terribly wrong here
			// } else if len(section.targets) == 1 {
		} else if len(section) == 1 {
			u <- cell{
				// location: workingCluster[section.targets[0]].location,
				location: workingCluster[section[0]].location,
				actual:   val,
			}
			changed = true
		}
	}
	return changed
}

func additionalCost(addVal int, markedVals map[int]bool, index indexedCluster) int {
	// skip this one if it's already marked
	if markedVals[addVal] {
		return -1
	}
	var newCols map[int]bool
	// for target, _ := range index[addVal].targets {
	for target, _ := range index[addVal] {
		newCols[target] = true
	}
	for preIndexed, _ := range markedVals {
		// for target, _ := range index[preIndexed].targets {
		for target, _ := range index[preIndexed] {
			delete(newCols, target)
		}
	}
	return len(newCols)
}

// Given certain premarked values, searches an index to find the next value to
// mark while staying under the given budget. If something is found to match
// the required values under the squares budget, changes are made.
func findPairsChild(markedVals map[int]bool, budget, required int, index indexedCluster, workingCluster []cell, u chan cell) (changed bool) {
	for possibleVal, locations := range index {
		if markedVals[possibleVal] {
			// this cell is already on the list, so skip it
			continue
		}
		// if len(locations.targets) > budget {
		if len(locations) > budget {
			// adding this cell will never fit in the budget, so skip it
			continue
		}
		// you can't add anything, so return false
		if budget < 1 {
			return false
		}
		if deduction := additionalCost(possibleVal, markedVals, index); deduction < budget {
			markedValsCopy := markedVals
			markedValsCopy[possibleVal] = true
			// if you already have the hit, mark it
			if required >= len(markedValsCopy) {
				if processSquares() {

				}

			}

			// if you need to go down recursively, do it
			if findPairsChild(makredValsCopy, budget-deduction, required, index, workingCluster) {
				return true
			}
		}
	}
}

// func findPairs(index indexedCluster, workingCluster []cell, u chan cell) (changed bool) {
// 	searchCount := 2

// 	localIndex := index
// 	for searchCount < len(index) {
// 		for val, valIndex := range localIndex {
// 			if len(valIndex.targets) <= searchCount {
// 				localIndex := index
// 				// this could be a first in a pair
// 				// probalby use something recursive to find the second match? idk

// 			}
// 		}
// 	}
// }

// func cheapestAddition(markedVals map[int]bool, index indexedCluster)(int, int){
// 	var markedCells map[int]bool
// 	for k, v :=
// }
