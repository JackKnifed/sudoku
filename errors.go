package sudoku

import "fmt"

type sudokuError struct {
	errType uint8
	args    []interface{}
}

func (e sudokuError) Error() string {
	if e.errType == 0 || e.errType > uint8(len(errFmtStrings)) {
		return fmt.Sprintf(errFmtStrings[0], e.errType)
	}
	return fmt.Sprintf(errFmtStrings[e.errType], e.args)
}

const (
	ErrInvalidValue = 1 << iota
	ErrCellAlreadSolved
)

var errFmtStrings = []string{
	"err %d undefined",
	"invalid value %s for cell",
	"cell is already solved",
}
