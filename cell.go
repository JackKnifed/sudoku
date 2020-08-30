package sudoku

// Holds the global Sudoku board state

// Communication is handled by four major mechanics:
// * A channel to distribute the current board state
// * A channel to notify threads of updates
// * A channel to update known & possible values

const (
	InvalidValue = "invalid value for cell"
)

func uintToBshift(i uint) uint {
	return 1 << i
}

// coord contains x and y elements for a given position - 0 indexed
type coord struct {
	x uint
	y uint
}

// A cell holds all of the required knowledge for a cell.
// It contains it's own address on the board (guaranteed unique)
// if actual is set to 0, actual value is unknown
type cell struct {
	// memberBoard points back to the board
	board *board

	location coord // 0 indexed coordinates of the cell within the board

	provided bool
	solved   uint
	excluded uint
}

func (c *cell) ExcludePossible(val uint) error {
	if c.solved != 0 {
		return sudokuError{errType: ErrCellAlreadSolved}
	}
	// set the bit to excluded
	c.excluded = c.excluded | uintToBshift(val)
	return nil
}

func (c cell) IsPossible(val uint) bool {
	// c.excluded are the values that have been excluded
	// ^c.excluded are the values still possible
	return ^c.excluded&uintToBshift(val) > 0
}

func (c *cell) SetStartValue(val uint) error {
	err := c.SetValue(val)
	c.provided = true
	return err
}

func (c *cell) SetValue(val uint) error {
	if c.solved != 0 || !c.IsPossible(val) {
		return sudokuError{
			errType: ErrInvalidValue,
			args:    []uinterface{}{val},
		}
	}
	c.solved = val
	c.excluded = ^uintToBshift(val)
	return nil
}

func (c cell) BlockNum(width uint) uint {
	xBlock := c.location.x / c.board.blockSize.x
	yBlock := c.location.y / c.board.blockSize.y
	return yBlock*c.board.blockAcross.y + xBlock
}
