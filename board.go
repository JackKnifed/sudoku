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

func (b board) Nil() bool { return len(b.cells) == 0 }

// Row positions are the first coordinate of [][]cells.
// For example, in a 9x9 grid, the rows would be numbered as such
// -------------------------
// | 0 0 0 | 0 0 0 | 0 0 0 |
// | 1 1 1 | 1 1 1 | 1 1 1 |
// | 2 2 2 | 2 2 2 | 2 2 2 |
// |-------+-------+--------
// | 3 3 3 | 3 3 3 | 3 3 3 |
// | 4 4 4 | 4 4 4 | 4 4 4 |
// | 5 5 5 | 5 5 5 | 5 5 5 |
// |-------+-------+--------
// | 6 6 6 | 6 6 6 | 6 6 6 |
// | 7 7 7 | 7 7 7 | 7 7 7 |
// | 8 8 8 | 8 8 8 | 8 8 8 |
// -------------------------
func (b board) Row(id uint) []*cell {
	if b.Nil() {
		return []*cell{}
	}
	return b.cells[id]
}

// Col positions are the second coordinate of [][]cells
// For example, in a 9x9 grid, the rows would be numbered as such
// -------------------------
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// |-------+-------+--------
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// |-------+-------+--------
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// | 0 1 2 | 3 4 5 | 6 7 8 |
// -------------------------
func (b board) Col(x, y uint) (col []*cell) {
	_ = x
	if b.Nil() {
		return []*cell{}
	}
	index := 0
	for _, row := range b.cells {
		row[index] = row[y]
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
// y  -------------------------------------
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
func (b board) Block(x, y uint) (block []*cell) {
	if b.Nil() {
		return []*cell{}
	}

	xBlock := x / b.blockSize.x
	yBlock := y / b.blockSize.y
	blockID := yBlock*b.blockAcross.y + xBlock

	rowStart := (blockID / b.blockAcross.x) * b.blockSize.y
	colStart := (blockID % b.blockAcross.y) * b.blockSize.x

	for _, rowSlice := range b.cells[rowStart : rowStart+b.blockSize.y] {
		block = append(block, rowSlice[colStart:colStart+b.blockSize.x]...)
	}
	return
}
