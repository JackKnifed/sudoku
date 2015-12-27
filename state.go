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

func createBoard(level int) (board map[coord]cell) {
	for i := 0; i < level*level; i++ {
		for j := 0; j < level*level; j++ {
			board[coord{x:i, y:j}].location = coord{x:i, y:j}
			for k := 1; k =< level*level; k++ {
				board[coord{x:i, y:j}].possible = append(board[coord{x:i, y:j}].possible, k) 
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
