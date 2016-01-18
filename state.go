package sudoku

// Holds the global Sudoku board state

// Communication is handled by four major mechanics:
// * A channel to distribute the current board state
// * A channel to notify threads of updates
// * A channel to update known & possible values

import "math"

const (
	boardRow    = 0
	boardCol    = 1
	boardSquare = 2
)

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

type cluster []cell

type board struct {
	size     int
	clusters []cluster
}

type ClusterUpdate struct {
	cells   []cell
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

func createBoard(size int) board {
	var newBoard board
	newBoard.size = size
	for i := 0; i < size*size; i++ {
		for j := 0; j < size*size; j++ {
			board[coord{x: i, y: j}].location = coord{x: i, y: j}
			for k := 1; k <= size*size; k++ {
				board[coord{x: i, y: j}].possible = append(board[coord{x: i, y: j}].possible, k)
			}
		}
	}
	return newBoard
}

func getPos(position coord, orientation, size int) (int, error) {
	if position.x >= size*size {
		return -1, errors.New("x position is larger than the board")
	}
	if position.y >= size*size {
		return -1, errors.New("y position is larger than the board")
	}
	switch orientation {
	case boardRow:
		return position.x
	case boardCol:
		return position.y
	case boardSquare:
		return ((position.y / size) * size) + (position.x / size)
	default:
		return -1, errors.New("bad position")
	}
}

func clusterPicker(in board, orient int, position coord) (cluster, error) {
	if position.x >= len(in) {
		return cluster{}, errors.New("x coord out of range")
	} else if position.y >= len(in) {
		return cluster{}, errors.New("y coord out of range")
	}
	switch orient {
	case boardRow:
		return in[position.y], nil
	case boardCol:
		var result cluster
		for _, each := range in {
			result = append(result, each[position.x])
		}
		return outCluster, nil
	case boardSquare:
		var result cluster
		boardSize := int(math.Sqrt(len(board)))
		startX := (position.x / boardSize) * boardSize
		endX := ((position.x / boardSize) + 1) * boardSize
		startY := (position.y / boardSize) * boardSize
		endY := ((position.y / boardSize) + 1) * boardSize
		for x := startX; x < endX; x++ {
			for y := startY; y < endY; y++ {
				result = append(result, in[x][y])
			}
		}
		return result, nil
	default:
		return cluster{}, errors.New("bad orientation")
	}
}

func boardFilter(update <-chan coord, in <-chan board, out [][]chan<- cluster, status [][]chan<- struct{}) {

	var toWork cluster
	defer closeArrArrChan(out)
	// don't actually do anything until you have a board state to work with
	<-in
	for {
		select {
		case changed, more := <-update:
			if !more {
				return
			}
			for i = boardRow; i <= boardSquare; i++ {
				curBoard := <-in
				if _, open := <-status[i][getPos(changed, i, curBoard.size)]; !open {
					// skip this update if the cluster is already solved
					continue
				}
				toWork, err := clusterPicker(curBoard, i, changed)
				if err != nil {
					panic(err) // #TODO# replace this panic
				}
				out[i][getPos(changed, i, curBoard.size)] <- toWork
			}
		}
	}
}

// The boardQueue is run in a go thread
// It serves a given value to any requestor whenever asked, or recieves updates to the value to serve.
func boardCache(in <-chan board, out chan<- board) {
	currentBoard := <-in
	var done bool
	defer close(out)

	for {
		select {
		case currentBoard, done = <-in:
			if done {
				return
			}
		case out <- currentBoard:
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

func clusterFilter(in <-chan board)

func cluster(updates chan ClusterUpdate, cluster chan ClusterUpdate, status chan int) {
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
				cells:   changes,
			}
		}
	}
}
