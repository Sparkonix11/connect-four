package game

// Board dimensions per SRS FR-GM-001
const (
	Columns = 7
	Rows    = 6
)

// Cell represents the state of a single cell
type Cell int

const (
	Empty   Cell = 0
	Player1 Cell = 1 // Red
	Player2 Cell = 2 // Yellow
)

// Board represents the 7x6 game grid
// Board[row][column] - row 0 is top, row 5 is bottom
type Board [Rows][Columns]Cell

// NewBoard creates an empty game board
func NewBoard() *Board {
	return &Board{}
}

// Clone creates a deep copy of the board for simulation
func (b *Board) Clone() *Board {
	clone := &Board{}
	for r := 0; r < Rows; r++ {
		for c := 0; c < Columns; c++ {
			clone[r][c] = b[r][c]
		}
	}
	return clone
}

// IsColumnFull checks if a column cannot accept more discs
func (b *Board) IsColumnFull(col int) bool {
	if col < 0 || col >= Columns {
		return true
	}
	// Check top row of the column
	return b[0][col] != Empty
}

// IsBoardFull checks if all 42 cells are filled (draw condition)
func (b *Board) IsBoardFull() bool {
	for c := 0; c < Columns; c++ {
		if !b.IsColumnFull(c) {
			return false
		}
	}
	return true
}

// DropDisc drops a disc into the specified column
// Returns the row where the disc landed, or -1 if invalid
func (b *Board) DropDisc(col int, player Cell) int {
	if col < 0 || col >= Columns {
		return -1
	}
	if b.IsColumnFull(col) {
		return -1
	}

	// Find the lowest empty row (gravity simulation)
	for row := Rows - 1; row >= 0; row-- {
		if b[row][col] == Empty {
			b[row][col] = player
			return row
		}
	}
	return -1
}

// GetCell returns the cell value at the specified position
func (b *Board) GetCell(row, col int) Cell {
	if row < 0 || row >= Rows || col < 0 || col >= Columns {
		return Empty
	}
	return b[row][col]
}

// ToSlice converts the board to a 2D slice for JSON serialization
func (b *Board) ToSlice() [][]int {
	result := make([][]int, Rows)
	for r := 0; r < Rows; r++ {
		result[r] = make([]int, Columns)
		for c := 0; c < Columns; c++ {
			result[r][c] = int(b[r][c])
		}
	}
	return result
}

// FromSlice creates a board from a 2D slice
func FromSlice(slice [][]int) *Board {
	board := NewBoard()
	for r := 0; r < Rows && r < len(slice); r++ {
		for c := 0; c < Columns && c < len(slice[r]); c++ {
			board[r][c] = Cell(slice[r][c])
		}
	}
	return board
}

// ValidColumns returns a list of columns that can accept a disc
func (b *Board) ValidColumns() []int {
	var valid []int
	for c := 0; c < Columns; c++ {
		if !b.IsColumnFull(c) {
			valid = append(valid, c)
		}
	}
	return valid
}

// GetDropRow returns the row where a disc would land if dropped in the column
// Returns -1 if the column is full
func (b *Board) GetDropRow(col int) int {
	if col < 0 || col >= Columns || b.IsColumnFull(col) {
		return -1
	}
	for row := Rows - 1; row >= 0; row-- {
		if b[row][col] == Empty {
			return row
		}
	}
	return -1
}
