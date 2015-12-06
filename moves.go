package sudoku

type indexedCluster map[int]struct {
	targets []int
}

func indexCluster(in []cell) (out indexedCluster) {
	for id, each := range in {
		for onePossible, _ := range each.possible {
			tmp := out[onePossible]
			tmp.targets = append(tmp.targets, id)
			out[onePossible] = tmp
		}
	}
	return out
}

// look at the possible values for a cell - if there's only one, mark that cell as solved
func confirmCell(workingCell cell, u chan cell) bool {
	var firstFind int

	for loc, each := range workingCell.possible {
		if each {
			if firstFind != 0 {
				// I found two possible values for this cell
				return false
			}
			firstFind = loc
		}
		if firstFind != 0 {
			// send back an update and return true - you changed something
			u <- cell{
				location: workingCell.location,
				actual:   firstFind,
			}
			return true
		}
	}
	return false
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
		if len(section.targets) < 1 {
			// something went terribly wrong here
		} else if len(section.targets) == 1 {
			u <- cell{
				location: workingCluster[section.targets[0]].location,
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
	for target, _ := range index[addVal].targets {
		newCols[target] = true
	}
	for preIndexed, _ := range markedVals {
		for target, _ := range index[preIndexed].targets {
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
				if processSquares(){

				}


			}

			// if you need to go down recursively, do it
			if findPairsChild(makredValsCopy, budget - deduction, required, index, workingCluster) {
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
