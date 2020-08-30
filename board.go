package sudoku

type board struct {
	cells [][]*cell // x coord then y coord

	blockSize   coord // # of cells in each block in each direction
	blockAcross coord // # of blocks on the board in each direction
}

func (b board) Init() {
	b.cells = make([][]*cell, b.blockSize.x*b.blockAcross.x)
	for i := range b.cells {
		b.cells[i] = make([]*cell, b.blockSize.y*b.blockAcross.y)
		for j := range b.cells[i] {
			b.cells[i][j] = &cell{
				board: &b,
				location: coord{
					x: uint(i),
					y: uint(j),
				},
			}
		}
	}
}

func (b board) Width() uint  { return b.blockSize.x * b.blockAcross.x }
func (b board) Height() uint { return b.blockSize.y * b.blockAcross.y }

func (b board) Nil() bool { return len(b.cells) == 0 }

// Row positions are the first coordinate of [][]cells.
// For example, in a 9x9 grid, the rows would be numbered as such
//  x: 0 1 2   3 5 6   7 8 9
// y -------------------------
// 1 | 0 0 0 | 0 0 0 | 0 0 0 |
// 2 | 1 1 1 | 1 1 1 | 1 1 1 |
// 3 | 2 2 2 | 2 2 2 | 2 2 2 |
//   |-------+-------+--------
// 4 | 3 3 3 | 3 3 3 | 3 3 3 |
// 5 | 4 4 4 | 4 4 4 | 4 4 4 |
// 6 | 5 5 5 | 5 5 5 | 5 5 5 |
//   |-------+-------+--------
// 7 | 6 6 6 | 6 6 6 | 6 6 6 |
// 8 | 7 7 7 | 7 7 7 | 7 7 7 |
// 9 | 8 8 8 | 8 8 8 | 8 8 8 |
//   -------------------------
func (b board) Row(id uint) []*cell {
	if b.Nil() {
		return []*cell{}
	}
	return b.cells[id]
}

// Col positions are the second coordinate of [][]cells
// For example, in a 9x9 grid, the rows would be numbered as such
//  x: 0 1 2   3 5 6   7 8 9
// y -------------------------
// 1 | 0 1 2 | 3 4 5 | 6 7 8 |
// 2 | 0 1 2 | 3 4 5 | 6 7 8 |
// 3 | 0 1 2 | 3 4 5 | 6 7 8 |
//   |-------+-------+--------
// 4 | 0 1 2 | 3 4 5 | 6 7 8 |
// 5 | 0 1 2 | 3 4 5 | 6 7 8 |
// 6 | 0 1 2 | 3 4 5 | 6 7 8 |
//   |-------+-------+--------
// 7 | 0 1 2 | 3 4 5 | 6 7 8 |
// 8 | 0 1 2 | 3 4 5 | 6 7 8 |
// 9 | 0 1 2 | 3 4 5 | 6 7 8 |
//   -------------------------
func (b board) col(id uint) (col []*cell) {
	if b.Nil() {
		return []*cell{}
	}
	index := 0
	for _, row := range b.cells {
		row[index] = row[id]
	}
	return
}

// Block positions are left to right blocks of [][]cells
// For example, in a 9x9 grid, the rows would be numbered as such - the expected case
// -------------------------
// | 0 0 0 | 1 1 1 | 2 2 2 |
// | 0 0 0 | 1 1 1 | 2 2 2 |
// | 0 0 0 | 1 1 1 | 2 2 2 |
// |-------+-------+--------
// | 3 3 3 | 4 4 4 | 5 5 5 |
// | 3 3 3 | 4 4 4 | 5 5 5 |
// | 3 3 3 | 4 4 4 | 5 5 5 |
// |-------+-------+--------
// | 6 6 6 | 7 7 7 | 8 8 8 |
// | 6 6 6 | 7 7 7 | 8 8 8 |
// | 6 6 6 | 7 7 7 | 8 8 8 |
// -------------------------
// More uinteresting, take a nonsense example (that still must work)
// blockAcross := coord{x:3, y:4}
// blockSize := coord{x: 5, y:6}
//   x: 0 1 2 3 4   5 6 7 8 9   a b c d e
//  y -------------------------------------
//  0 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//  1 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//  2 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//  3 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//  4 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//  5 | 0 0 0 0 0 | 1 1 1 1 1 | 2 2 2 2 2 |
//    |-----------+-----------+-----------|
//  6 | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//  7 | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//  8 | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//  9 | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//  a | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//  b | 3 3 3 3 3 | 4 4 4 4 4 | 5 5 5 5 5 |
//    |-----------+-----------+-----------|
//  c | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
//  d | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
//  e | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
//  f | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
// 10 | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
// 11 | 6 6 6 6 6 | 7 7 7 7 7 | 8 8 8 8 8 |
//    |-----------+-----------+-----------|
// 12 | 9 9 9 9 9 | a a a a a | b b b b b |
// 13 | 9 9 9 9 9 | a a a a a | b b b b b |
// 14 | 9 9 9 9 9 | a a a a a | b b b b b |
// 15 | 9 9 9 9 9 | a a a a a | b b b b b |
// 16 | 9 9 9 9 9 | a a a a a | b b b b b |
// 17 | 9 9 9 9 9 | a a a a a | b b b b b |
//    |-----------------------------------|
func (b board) Block(id uint) (block []*cell) {
	if b.Nil() {
		return []*cell{}
	}

	rowStart := (id / b.blockAcross.x) * b.blockSize.y
	colStart := (id % b.blockAcross.y) * b.blockSize.x

	for _, rowSlice := range b.cells[rowStart : rowStart+b.blockSize.y] {
		block = append(block, rowSlice[colStart:colStart+b.blockSize.x]...)
	}
	return
}

// Returns a slice of all positions next to the current one
// For example, in a 9x9 grid, if input [3,6] (marked below by @)
//  x: 0 1 2   3 5 6   7 8 9
// y -------------------------
// 1 | 0 1 2 | 3 4 5 | 6 7 8 |
// 2 | 0 1 2 | 3 4 5 | 6 7 8 |
// 3 | 0 1 2 | 3 4 5 | 6 7 8 |
//   |-------+-------+--------
// 4 | 0 1 2 | 3 4 5 | 6 7 8 |
// 5 | 0 1 R | 3 R 5 | 6 7 8 |
// 6 | 0 1 2 | @ 4 5 | 6 7 8 |
//   |-------+-------+--------
// 7 | 0 1 R | 3 R 5 | 6 7 8 |
// 8 | 0 1 2 | 3 4 5 | 6 7 8 |
// 9 | 0 1 2 | 3 4 5 | 6 7 8 |
//   -------------------------
// This returns cells marked by X [2,5], [4,5], [2,7], [4,7]
// Also note, if this is on an edge or corner, it will return less values.
func (b board) Corners(pos coord) (corners []*cell) {
	appendCell := func(x, y uint) {
		corners = append(corners, b.cells[x][y])
	}
	walkY := func(x, y uint) {
		if y > 0 {
			appendCell(x, y-1)
		} else if y < b.Width()-1 {
			appendCell(x, y+1)
		}
	}
	walkX := func(x, y uint) {
		if x > 0 {
			walkY(x-1, y)
		} else if x < b.Height()-1 {
			walkY(x+1, y)
		}
	}
	walkX(pos.x, pos.y)
	return
}
