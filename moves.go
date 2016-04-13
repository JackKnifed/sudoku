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
// 2) If any cell is solved, it has all exclusions.
// 3) If any cell is solved, that value is excluded in other cells.
// 4) If any cell has all but one excluded value, that is that cell's value.
// 5) If any x cells have the same x values, the missing values are elsewhere
//  excluded - those values are constrained to those cells.
// 6) If any value is present in only one cell, that is that cell's value.
// 7) If any x values are possible in x cells, all other values are excluded in
//  those cells.
//
// Additional Helper functions are included first.

// indexedCLuster is a datatype for an index for the values of the cluster.
// Each possible value is a key in the map. THe value of each key is an array of
// possible locations for that value - the index and order are not defined,
// instead the values of the array are the indexes of possible cells in for that
// value.
type intArray []int
type indexedCluster map[int]intArray

var fullArray intArray

func init() {
	for i, _ := range in {
		fullArray = append(fullArray, i)
	}
}

// indexCluster takes a cluster of excluded values, and returns an index of
//  the possible locations for each value.
func indexCluster(in []cell) (out indexedCluster) {

	// add the fullArray to every value location
	for i, _ := range in {
		copy(out[i], fullArray)
	}

	// I can simply delete known values from the array
	for id, each := range out {
		if (each.actual) != 0 {
			delete(out, id)
		}
		for _, oneExcluded := range each.excluded {
			out[id] = subArr(out[id], []int{oneExcluded})
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
// 2) If any cell is solved, it has all exclusions.
func solvedNoPossible(cluster []cell) (changes []cell) {
	for id, each := range cluster {
		if each.actual != 0 && each.excluded < len(cluster) {
			newExclusion := subArr(fullArray, each.excluded)
			changes = append(changes, cell{location: each.location, excluded: newExc})
		}
	}
	return
}

// Removes known values from other possibles in the same cluster
// Covers rule 3 from above:
// 3) If any cell is solved, that value is excluded in other cells.
func eliminateKnowns(cluster []cell) (changes []cell) {
	var knownValues []int

	// Loop thru and find all solved values.
	for _, each := range cluster {
		if each.actual != 0 {
			knownValues = append(knownValues, each.actual)
		}
	}

	for _, each := range cluster {
		if !allInArr(each.excluded, knownValues){
			newExclusion := subArr(knownValues, each.excluded)
			changes = append(changes, cell{location: each.location, excluded: newExc})
		}
	}
	return
}

// This covers the 4th rule from above:
// 4) If any cell has all but one excluded value, that is that cell's value.
func singleValueSolver(cluster []cell) (changes []cell) {
	for _, each := range cluster {
		// skip this cell if it's already solved
		if each.actual != 0 {
			continue
		}

		// should never happen - probably #TODO# to catch this
		if len(each.excluded) >= len(cluster) {
			panic("Found an unsolved cell with all values excluded")
		}

		if len(each.excluded) == len(cluster) - 1 {
			// send back an update for this cell
			solvedValue := subArr(fullArray, each.excluded)
			changes = append(changes, cell{location: each.location,
				actual: solvedValue[0], possible: fullArray})
		}
	}
	return
}

// A helper function to determine the values certain cells hit
// ##TODO## review
func valuesPainted(markedCells []int, cluster []cell) (neededValues []int) {
	var first, second []int
	first = cluster[markedCells[0]].excluded

	switch {
	case len(markedCells) < 1:
		// this will only happen when passed one cell
		return []int[]
	case len(markedCells) > 1 :
		remainder = valuesPainted(markedCells[1:], cluster)
	default:
		remainder = cluster[markedCells[1]].excluded
	}

	neededValues = subArr(fullArray, andArr(first, second))
	return
}
func valuesPainted(markedCells []int, cluster []cell) (neededValues []int) {
	for _, oneCell := range markedCells {
		neededValues = addArr(neededValues, cluster[oneCell].possible)
	}
	return neededValues
}

// A helper function to determine the number of valus hit given a specific set
// of cells.
// ##TODO## review
func cellsCost(markedCells []int, cluster []cell) int {
	return len(valuesPainted(markedCells, cluster))
}

// ##TODO## - review
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
// 5) If any x cells have the same x values, the missing values are elsewhere
//  excluded - those values are constrained to those cells.
// ##TODO##
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
// 6) If any value is possible in only one cell, that is that cell's value.
// ##TODO##
func singleCellSolver(index indexedCluster, cluster []cell) (changes []cell) {
	for val, section := range index {
		if len(section) < 1 {
			// something went terribly wrong here - #TODO# add panic catch?
			panic("Found an unsolved cell with no possible values")
		} else if len(section) == 1 {
			changes = append(changes, cell{
				location: cluster[section[0]].location,
				actual:   val,
				possible: fullArray,
			})
		}
	}
	return
}

// A helper function to determine what cells are painted by given values
// ##TODO## review
func cellsPainted(markedVals []int, index indexedCluster) (neededCells []int) {
	for _, value := range markedVals {
		neededCells = addArr(neededCells, index[value])
	}
	return
}

// A helper function to determine the number of cells hit by working across a map
// of values.
// ##TODO## review
func valuesCost(markedVals []int, index indexedCluster) int {
	return len(cellsPainted(markedVals, index))
}

// ##TODO## review
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

// This covers rule 7 from above:
// 7) If any x values are possible in x cells, all other values are excluded in
//  those cells.
// ##TODO##
func valueLimiter(index indexedCluster, cluster []cell) (changes []cell) {
	upperBound := len(index)
	for i := 2; i <= upperBound; i++ {
		changes = append(changes, valueLimiterCHild(i, []int{}, index, cluster))
	}
	return changes
}
