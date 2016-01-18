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
		for _, onePossible := range each.possible {
			out[onePossible] = append(out[onePossible], id)
		}
	}
	return out
}

// This covers rule 1 from above:
// 1) If all cells are solved, that cluster is solved.
func clusterSolved(cluster []cell) (solved bool) {
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
func solvedNoPossible(cluster []cell) (changes []cell) {
	for _, each := range cluster {
		if each.actual == 0 {
			continue
		}
		if len(each.possible) > 0 {
			changes = append(changes, cell{location: each.location, possible: each.possible})
		}
	}
	return
}

// Removes known values from other possibles in the same cluster
// Covers rule 3 from above:
// 3) If any cell is solved, that value is not possible in other cells.
func eliminateKnowns(cluster []cell) (changes []cell) {
	var knownValues []int

	// Loop thru and find all solved values.
	for _, each := range cluster {
		if each.actual != 0 {
			knownValues = append(knownValues, each.actual)
		}
	}

	for _, each := range cluster {
		if len(andArr(each.possible, knownValues)) > 0 {
			changes = append(changes, cell{
				location: each.location,
				possible: subArr(each.possible, knownValues),
			})
		}
	}
	return
}

// This covers the 4th rule from above:
// 4) If any cell only has one possible value, that is that cell's value.
func singleValueSolver(cluster []cell) (changes []cell) {
	for _, each := range cluster {
		if each.actual == 0 {
			continue
		}

		// should never happen - probably #TODO# to catch this
		if len(each.possible) < 1 {
			panic("Found an unsolved cell with no possible values")
		}

		if len(each.possible) == 1 {
			// send back an update for this cell
			changes = append(changes, cell{location: each.location,
				actual: each.possible[0], possible: []int{}})
		}

	}
	return
}

// A helper function to determine the values certain cells hit
func valuesPainted(markedCells []int, cluster []cell) (neededValues []int) {
	for _, oneCell := range markedCells {
		neededValues = addArr(neededValues, cluster[oneCell].possible)
	}
	return neededValues
}

// A helper function to determine the number of valus hit given a specific set
// of cells.
func cellsCost(markedCells []int, cluster []cell) int {
	return len(valuesPainted(markedCells, cluster))
}

func cellLimiterChild(limit int, markedCells []int, cluster []cell) (changes []cell) {
	valueCount := cellsCost(markedCells, cluster)
	switch {

	case valueCount < len(markedCells):
		// #TODO# probably fix this? rework into error?
		panic("less possible values than squares to fill")
	// you have overspent - it's a no-go
	case len(markedCells) > limit:
		return []cell{}

	// did you fill the cells? if so, mark it
	case valueCount == len(markedCells):
		valuesCovered := valuesPainted(markedCells, cluster)
		// it's a match - so remove values
		for id, each := range cluster {
			if each.actual != 0 {
				// already solved
				continue
			}
			if inArr(markedCells, id) {
				// this cell is a part of the list - no exclusions needed
				continue
			}

			toRemove := andArr(each.possible, valuesCovered)
			if len(toRemove) > 0 {
				changes = append(changes, cell{
					location: each.location,
					possible: toRemove,
				})
			}
		}
	// you have room to add more cells (depth first?)
	case len(markedCells) < limit:
		// you need to try adding each other cell
		for idCell, oneCell := range cluster {
			if oneCell.actual != 0 {
				continue
			}
			if inArr(markedCells, idCell) {
				// this cell is already in the map, skip
				continue
			}

			// decend down into looking at that cell
			changes = append(changes,
				cellLimiterChild(limit, append(markedCells, idCell), cluster)...)
		}
	}
	return changes
}

// This covers rule 5:
// 5) If any x cells have x possible values, those values are not possible
//  outside of those cells - those values are constrained to those cells.
func cellLimiter(cluster []cell) (changes []cell) {
	upperBound := len(cluster)
	for _, eachCell := range cluster {
		if eachCell.actual != 0 {
			upperBound--
		}
	}
	for i := 2; i <= upperBound; i++ {
		changes = append(changes, cellLimiterChild(i, []int{}, cluster)...)
	}
	return
}

// This covers rule 6 from above:
// 6) If any value only has one possible cell, that is that cell's value.
func singleCellSolver(index indexedCluster, cluster []cell) (changes []cell) {
	for val, section := range index {
		if len(section) < 1 {
			// something went terribly wrong here - #TODO# add panic catch?
			panic("Found an unsolved cell with no possible values")
		} else if len(section) == 1 {
			changes = append(changes, cell{
				location: cluster[section[0]].location,
				actual:   val,
				possible: []int{},
			})
		}
	}
	return
}

// A helper function to determine what cells are painted by given values
func cellsPainted(markedVals []int, index indexedCluster) (neededCells []int) {
	for _, value := range markedVals {
		neededCells = addArr(neededCells, index[value])
	}
	return
}

// A helper function to determine the number of cells hit by working across a map
// of values.
func valuesCost(markedVals []int, index indexedCluster) int {
	return len(cellsPainted(markedVals, index))
}

func valueLimiterChild(limit int, markedValues []int, index indexedCluster,
	cluster []cell) (changes bool) {
	cellCount := valuesCost(markedValues, index)
	switch {
	case cellCount < len(markedValues):
		panic("less cells available than the values that need to go in them")
	case len(markedValues) > limit:
		// we have marked more values
		return false
	case cellCount == len(markedValues):
		// you have exactly as many values as cells
		cellsCovered := cellsPainted(markedValues, index)
		for id, each := range cellsCovered {
			if toRemove := subArr(each.possible, markedValues); len(toRemove) > 0 {
				changes = append(changes, cell{
					location: each.location,
					possible: toRemove,
				})
			}
		}
	case len(markedValues) < limit:
		// you can mark another value and see where that gets you
		for value, _ := range index {
			if inArr(markedValues, value) {
				// this value is already marked
				continue
			}
			// decend down into looking at that value
			changes = append(changes, valueLimiterChild(limit, append(markedValues, value),
				index, cluster, u))
		}
	}
	return
}

// THis covers rule 7 from above:
// 7) If any x values have only x possible cells, those cells only have those
//  possible values - those cells are constrained to those values.
func valueLimiter(index indexedCluster, cluster []cell) (changes []cell) {
	upperBound := len(index)
	for i := 2; i <= upperBound; i++ {
		changes = append(changes, valueLimiterCHild(i, []int{}, index, cluster))
	}
	return changes
}
