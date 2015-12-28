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
// 5) If any x cells only have x possible values, those values are not possible
//  outside of those cells - those values are constrained to those cells.
// 6) If any value only has one possible cell, that is that cell's value.
// 7) If any x values only have x possible cells, those cells only have those
//  possible values - those cells are constrained to those values.
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
		for _, onePossible_ := range each.possible {
			out[onePossible] = append(out[onePossible], id)
		}
	}
	return out
}

// This covers rule 1 from above:
// 1) If all cells are solved, that cluster is solved.
func clusterSolved(cluster []cell, u chan cell) (solved bool) {
	solved = true
	for _, each := range cluster {
		if each.actual == 0 {
			solved = false
		} else if len(each.possible) > 0 {
			solved = false
		}
	}
	return solved
}

// This covers rule 2 from above:
// 2) If any cell is solved, it has no possibles.
func solvedNoPossible(cluster []cell, u chan cell) (changed bool) {
	for _, each := range cluster {
		if each.actual == 0 {
			continue
		}
		if len(each.possible) > 0 {
			changed = true
			u <- cell{
				location: each.location,
				possible: each.possible,
			}
		}
	}
	return changed
}

// Removes known values from other possibles in the same cluster
// Covers rule 3 from above:
// 3) If any cell is solved, that value is not possible in other cells.
func eliminateKnowns(workingCluster []cell, u chan cell) (changed bool) {
	var knownValues []int

	// Loop thru and find all solved values.
	for _, each := range workingCluster {
		if each.actual != 0 {
			knownValues = append(knownValues, each.actual)
		}
	}

	for _, each := range workingCluster {
		if anyInArr(each, knownValues) {
			u <- cell{
				location: each.location,
				possible: subArr(each, knownValues)
			}
			changed = true
		}
	}
	return changed
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

		// should never happen - probably #TODO# to catch this
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

// A helper function to determine the number of valus hit given a specific set
// of cells.
func cellsCost(markedCells map[int]bool, cluster []cell) int {
	var neededValues map[int]bool
	for cellPos, _ := range markedCells {
		for possibleValue, _ := range cluster[cellPos].possible {
			neededValues[possibleValue] = true
		}
	}
	return len(neededValues)
}

func cellLimiterChild(limit int, markedCells map[int]bool, cluster []cell, u chan cell) (changed bool) {
	valueCount := cellsCost(markedCells, cluster)
	// you have overspent - it's a no-go
	if len(markedCells) > limit {
		return false
	}

	// you have room to add more cells (depth first?)
	if len(markedCells) < limit {
		if valueCount < len(markedCells) {
			// #TODO# probably fix this? rework into error?
			panic("less possible values than squares to fill")
		}
		if valueCount > len(markedCells) {
			// you need to try adding each other cell
			for idCell, oneCell := range cluster {
				if oneCell.actual != 0 {
					continue
				}
				if markedCells[idCell] {
					// this cell is already in the map, skip
					continue
				}

				// decend down into looking at that cell
				childMarkedCells := markedCells // #TODO# map copy
				childMarkedCells[idCell] = true
				if cellLimiterChild(limit, childMarkedCells, cluster, u) {
					// if you got true from the child, pass it on
					return true
				}
			}
		}
	}

	// did you fill the cells? if so, mark it
	if valueCount == len(markedCells) {
		// it's a match - so remove values
		for idCell, oneCell := range cluster {
			if oneCell.actual != 0 {
				// already solved
				continue
			}
			if markedCells[idCell] {
				// this cell is a part of the list - no exclusions needed
				continue
			}
			remove := map[int]bool{}
			for potential, _ := range oneCell.possible {
				if markedCells[idCell] {
					// this possibility exists in the map - need to remove it
					remove[potential] = true
				}
			}
			if len(remove) > 0 {
				changed = true
				u <- cell{
					location: oneCell.location,
					possible: remove,
				}
			}
		}
	}
	return changed
}

// This covers rule 5:
// 5) If any x cells have x possible values, those values are not possible
//  outside of those cells - those values are constrained to those cells.
func cellLimiter(cluster []cell, u chan cell) (changed bool) {
	upperBound := len(cluster)
	for _, eachCell := range cluster {
		if eachCell.actual != 0 {
			upperBound--
		}
	}
	for i := 2; i <= upperBound; i++ {
		if cellLimiterChild(i, map[int]bool{}, cluster, u) {
			changed = true
		}
	}
	return changed
}

// This covers rule 6 from above:
// 6) If any value only has one possible cell, that is that cell's value.
func singleCellSolver(index indexedCluster, workingCluster []cell, u chan cell) (changed bool) {
	for val, section := range index {
		if len(section) < 1 {
			// something went terribly wrong here - #TODO# add panic catch?
			panic("Found an unsolved cell with no possible values")
		} else if len(section) == 1 {
			u <- cell{
				location: workingCluster[section[0]].location,
				actual:   val,
			}
			changed = true
		}
	}
	return changed
}

// A helper function to determine the number of cells hit by working across a map
// of values.
func valuesCost(markedVals map[int]bool, index indexedCluster, cluseter []cell) int {
	var neededCells map[int]bool
	for value, _ := range markedVals {
		for oneCell, _ := range index[value] {
			neededCells[oneCell] = true
		}
	}
	return len(neededCells)
}

func valueLimiterChild(limit int, markedValues map[int]bool, index indexedCluster,
	cluster []cell, u chan cell) (changed bool) {
	if len(markedValues) > limit {
		// we have marked more values
		return false
	}
	currentCost := valuesCost(markedValues, index, cluster)
	if currentCost > limit {
		// you're over the budget to spend
		return false
	}
	return changed
}

// THis covers rule 7 from above:
// 7) If any x values have only x possible cells, those cells only have those
//  possible values - those cells are constrained to those values.
/*
func valueLimiter(index indexedCluster, cluster []cell, u chan cell) (changed bool) {

}
*/
