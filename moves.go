package sudoku

type indexedCluster map[int]struct{
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
	for val, section := range index{
		if len(section.targets) < 1{
			// something went terribly wrong here
		} else if len(section.targets) == 1 {
			u <- cell{
				location: workingCluster[section.targets[0]].location,
				actual: val,
			}
			changed = true
		}
	}
	return changed
}

func findPairs(index indexedCluster, workingCluster []cell, u chan cell) (changed bool) {
	searchCount := 2

	for searchCount < len(index) {
		for val, valIndex := range localIndex {
			if len(valIndex) <= searchCount {
				localIndex := index
				// this could be a first in a pair
				// probalby use something recursive to find the second match? idk

			}
		}
	}
}