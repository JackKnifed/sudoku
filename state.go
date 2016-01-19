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
	excluded []int
}

type cluster []cell

type board struct {
	size     int
	clusters []cluster
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

func clusterFilter(update <-chan coord, in <-chan board, out [][]chan<- cluster, status [][]<-chan struct{}) {

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
				position := getPos(changed, i, curBoard.size)
				if _, open := <-status[i][position]; !open {
					// skip this update if the cluster is already solved
					continue
				}
				curCluster, err := clusterPicker(curBoard, i, changed)
				if err != nil {
					panic(err) // #TODO# replace this panic
				}
				out[i][position] <- curCluster
			}
		}
	}
}

// takes an input and sends it out as it can
// exits when in is closed
// closes out on exit
func clusterSticky(in <-chan cluster, out chan<- cluster) {
	var more bool
	defer close(out)

	// preload the locCluster - don't do anything until you have that
	locCluster := <-in
	for {
		select {
		case locCluster, more = <-in:
			// if you get an update from upstream, do that first
			if !more {
				return
			}
		// pass the update downstream then block till you get another update from upstream
		case out <- locCluster:
			locCluster = <-in
		}
	}
}

// like a buffered channel, but no limit to the buffer size
// closes out on exit
// exits when in is closed
func updateBuffer(in <-chan cell, out chan<- cell) {
	var updates []interface{}
	var singleUpdate interface{}
	var open bool
	defer close(out)

	for {
		if len(updates) < 1 {
			// if you currently don't have anything to pass, WAIT FOR SOMETHING
			singleUpdate, open = <-in
			if !open {
				return
			}
			updates = append(updates, singleUpdate)
			continue
		}
		select {
		case singleUpdate, open = <-in:
			if !open {
				return
			}
			updates = append(updates, singleUpdate)
		case out <- updates[0]:
			updates = updates[1:]
		}
	}
}

// boardCache serves a given (or newer) update out as many times as requested
// closes `out` on exit
// exits when in `is` closed
func boardCache(in chan board, out chan<- board) {
	currentBoard := <-in
	var done bool
	defer close(out)

	for {
		select {
		case currentBoard, done = <-in:
			if done {
				return
			}
		case in <- currentBoard:
		case out <- currentBoard:
		}
	}
}

// looks for a send from all status channels
// if recieved or closed on all status chan, send on `idle`
// if all status is closed, close idle and exits
func idleCheck(status [][]<-chan interface{}, idle chan<- interface{}) {
	var more, isIdle, isSolved bool
	defer close(idle)

	for {
	start:
		isSolved, isIdle = true, true
		for _, middle := range status {
			for _, inner := range status {
				select {
				case _, more = <-inner:
					if more {
						isSolved = false
					}
				default:
					isSolved = false
					isIdle = false
					continue start
				}
			}
		}
		if isSolved {
			return
		} else if isIdle {
			idle <- nil
		}
	}
}

// takes a given cluster, and runs it through all of the moves
// exits when one of the conditions is met, or when the update channel is closed
// closes the status channel on exit
func clusterWorker(in <-chan cluster, status chan<- struct{}, updates chan<- cell, problems chan<- error) {
	defer close(status)

	var more bool
	var newCluster cluster
	var index indexedCluster
	var changes, newChanges []cell

	for {
		select {
		case newCluster, more = <-in:
			if !more {
				// if the channel is closed, exit
				return
			}
			if clusterSolved(cluster) {
				// if the cell is solved, exit
				return
			}

			newChanges = solvedNoPossible(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = eliminateKnowns(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = singleValueSolver(cluster.cells)
			changes = append(changes, newChanges)

			newChanges = cellLimiter(cluster.cells)
			changes = append(changes, newChanges)

			index = indexCluster(cluster.cells)

			newChanges = singleCellSolver(index, cluster.cells)
			changes = append(changes, newChanges)

			newChanges = valueLimiter(index, cluster.cells)
			changes = append(changes, newChanges)

			// feed all those changes into the update queue
			for len(changes) > 1 {
				updates <- changes[0]
				changes = changes[1:]
			}
		case status <- nil:
			// report idle only if there is nothing to do - order matters
		}
	}
}

// processes updates
// priority is problems, updates, then status checks
func updateProcessor(curBoard chan board, status <-chan struct{}, updates <-chan cell, posChange chan<- coord, problems <-chan error) {
	defer close(curBoard)
	defer close(posChange)

	var newBoard board
	var err error

	for {
		select {
		case err = <-problems:
			// any errors get top priority
		case cellChange = <-updates:
			// any udpates are handled before idle checks
			newBoard, err = changeBoard(<-curBoard, cellChange)
		case _, solved = <-status:
			// status check if neither of the above happens - or solved check
		}
	}
}

func changeBoard(in board, u cell) (out board, err error) {
	target := in.clusters[u.location.x][u.location.y]
	if cell.actual != 0 {
		if target.actual != 0
		for i = 1; i <= board.level; i++ {
			cell.excluded = append(cell.excluded, i)
		}
		cell.excluded = dedupArr(cell.excluded)
		in.clusters[cell.location.x][cell.location.y] = cell
	} else {
		in.clusters[cell.location.x][cell.location.y].excluded = addArr(board.clusters[cell.location.x][cell.location.y].excluded, cell.excluded)
	}


}
