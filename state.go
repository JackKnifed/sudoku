package sudoku

// Holds the global Sudoku board state

// Communication is handled by four major mechanics:
// * A channel to distribute the current board state
// * A channel to notify threads of updates
// * A channel to update known & possible values

// coord contains x and y elements for a given position
type coord struct {
	x int
	y int
}

// A cell holds all of the required knowledge for a cell.
// It contains it's own address on the board (guarenteed unique)
// if actual is set to 0, actual value is unknown
type cell struct {
	location coord
	actual   int
	possible []int
}

type ClusterUpdate struct {
	cells []cell,
	version int
}


/*
type Cluster interface {
	Len() int
	Swap(i, j int)
	Less(i, j int)
}
*/

func (c Cluster) Len() int {
	return len(c)
}
func (c Cluster) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Cluster) Less(i, j int) bool {
	if c[i].x < c[j].x {
		return true
	}
	if c[i].x == c[j].x && c[i].y < c[j].y {
		return true
	}
	return false
}

func createBoard(level int) (board map[coord]cell) {
	for i := 0; i < level*level; i++ {
		for j := 0; j < level*level; j++ {
			board[coord{x: i, y: j}].location = coord{x: i, y: j}
			for k := 1; k <= level*level; k++ {
				board[coord{x: i, y: j}].possible = append(board[coord{x: i, y: j}].possible, k)
			}
		}
	}
	return board
}

// The boardQueue is run in a go thread
// It serves a given value to any requestor whenever asked, or recieves updates to the value to serve.
func boardQueue(up, down chan map[coord]cell) {
	currentBoard := <-up
	var done bool
	defer close(down)

	for {
		select {
		case currentBoard, done = <-up:
			if done {
				return
			}
		case down <- currentBoard:
		}
	}
}

// The updateQueue holds an update and pushes it down stream.
// It will not push the same update down more than once.
func updateQueue(up, down chan coord) {
	var push bool
	var done bool
	var val coord
	defer close(down)

	for {
		if push {
			// Otherwise wait to either push it down, or updates
			select {
			case val, done = <-up:
				if done {
					return
				}
				push = true
			case down <- val:
				push = false
			}
		} else {
			// If it's been pushed, wait for upstream.
			val, done = <-up
			if done {
				return
			}
			push = true
		}
	}
}

func broadcastUpdate(up chan coord, childs []chan coord) {
	var val coord
	var done bool

	// clean up downstream channels
	defer func() {
		for _, each := range childs {
			close(each)
		}
	}()

	for {
		val, done = <-up
		if done {
			return
		}
		for _, each := range childs {
			each <- val
		}
	}
}

func clusterThread(updates chan ClusterUpdate, cluster chan ClusterUpdate, status chan int) {
	var currentVersion int
	defer close(status)

	for {
		select {
		case status <- currentVersion:
		case newCluster := <-cluster:
			currentVersion = newCluster.version
			var changes []cell
			if !clusterSolved(cluster) {
				return
			}

			newChanges := solvedNoPossible(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = eliminateKnowns(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = singleValueSolver(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = cellLimiter(cluster.cells)
			changes = append(changes, newChanges)

			index := indexCluster(cluster.cells)

			newChanges = singleCellSolver(index, cluster.cells)
			changes = append(changes, newChanges)

			newChanges = valueLimiter(index, cluster.cells)
			changes = append(changes, newChanges)

			updates <- ClusterUpdate{
				version: currentVersion,
				cells: changes,
			}
		}
	}
}