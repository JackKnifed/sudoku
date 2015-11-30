package sudoku

// look at the possible values for a cell - if there's only one, mark that cell as solved
func confirmCell(workingCell cell, c chan cell) bool {
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
			c <- cell{
				location: workingCell.location,
				actual:   firstFind,
			}
			return true
		}
	}
	return false
}

// Removes known values from other possibles in the same cluster
func eliminatePossibles(workingCluster []cell, c chan cell) bool {
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
			c <- cell{
				location: each.location,
				possible: remove,
			}
			remove = map[int]bool{}
		}
	}
	return changed
}
