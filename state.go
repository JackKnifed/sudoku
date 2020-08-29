package sudoku

// Holds the global Sudoku board state

// Communication is handled by four major mechanics:
// * A channel to distribute the current board state
// * A channel to notify threads of updates
// * A channel to update known & possible values

const (
	InvalidValue = "invalid value for cell"
)

func intToBshift(i uint8) uint16 {
	return 1 << i
}

// coord contains x and y elements for a given position - 0 indexed
type coord struct {
	x uint8
	y uint8
}

// A cell holds all of the required knowledge for a cell.
// It contains it's own address on the board (guaranteed unique)
// if actual is set to 0, actual value is unknown
type cell struct {
	// memberBoard points back to the board
	board *board

	location coord // 0 indexed coordinates of the cell within the board

	provided bool
	solved   uint8
	possible uint16
}

func (c cell) row() []cell {
	if c.board == nil {
		return []cell{}
	}
	return c.board.cells[c.location.x]
}

func (c cell) col() (col []cell) {
	if c.board == nil {
		return []cell{}
	}
	for _, row := range c.board.cells {
		col = append(col, row[c.location.y])
	}
	return
}

func (c cell) blockStartAndEnd() (x1, x2, y1, y2 uint8) {
	startAndEnd := func(coord uint8) (uint8, uint8) {
		if coord >= c.board.width { // outside board bounds, so just cap it
			coord = c.board.width - 1
		}
		// the following three lines are laid out so it's clearer what is going on
		coord /= 3 // compress it down to each block with integer division (block 0, 1, 2)
		coord *= 3 // expand that block back out (start 0, 3, or 6)
		return coord + 1, coord + 3
	}

	x1, x2 = startAndEnd(c.location.x)
	y1, y2 = startAndEnd(c.location.y)
	return
}

func (c cell) block() (block []cell) {
	if c.board == nil {
		return []cell{}
	}

	x1, x2, y1, y2 := c.blockStartAndEnd()
	for _, partRow := range c.board.cells[x1:x2] {
		for _, cellInPartRow := range partRow[y1:y2] {
			block = append(block, cellInPartRow)
		}
	}
	return
}

func (c *cell) ExcludePossible(val uint8) error {
	if c.solved != 0 {
		return sudokuError{errType: ErrCellAlreadSolved}
	}
	c.possible = c.possible &^ intToBshift(val)
	return nil
}

func (c cell) IsPossible(val uint8) bool {
	return c.possible&intToBshift(val) > 0
}

func (c *cell) SetStartValue(val uint8) error {
	err := c.SetValue(val)
	c.provided = true
	return err
}

func (c *cell) SetValue(val uint8) error {
	if c.solved != 0 || !c.IsPossible(val) {
		return sudokuError{
			errType: ErrInvalidValue,
			args:    []interface{}{val},
		}
	}
	c.solved = val
	c.possible = intToBshift(val)
	return nil
}

type board struct {
	width uint8    // how wide/high the board is
	cells [][]cell // x coord then y coord
}
