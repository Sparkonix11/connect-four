package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// Board should be empty
	for r := 0; r < Rows; r++ {
		for c := 0; c < Columns; c++ {
			if board[r][c] != Empty {
				t.Errorf("Expected empty cell at (%d, %d), got %d", r, c, board[r][c])
			}
		}
	}
}

func TestDropDisc(t *testing.T) {
	board := NewBoard()

	// Drop first disc in column 3
	row := board.DropDisc(3, Player1)
	if row != 5 { // Should land at bottom row
		t.Errorf("Expected row 5, got %d", row)
	}
	if board[5][3] != Player1 {
		t.Error("Disc not placed correctly")
	}

	// Drop second disc in same column
	row = board.DropDisc(3, Player2)
	if row != 4 { // Should land one above
		t.Errorf("Expected row 4, got %d", row)
	}

	// Drop disc in invalid column
	row = board.DropDisc(-1, Player1)
	if row != -1 {
		t.Error("Should reject negative column")
	}

	row = board.DropDisc(7, Player1)
	if row != -1 {
		t.Error("Should reject column >= 7")
	}
}

func TestIsColumnFull(t *testing.T) {
	board := NewBoard()

	// Column should not be full initially
	if board.IsColumnFull(0) {
		t.Error("Empty column reported as full")
	}

	// Fill column 0
	for i := 0; i < Rows; i++ {
		board.DropDisc(0, Player1)
	}

	if !board.IsColumnFull(0) {
		t.Error("Full column not reported as full")
	}

	// Try to drop in full column
	row := board.DropDisc(0, Player2)
	if row != -1 {
		t.Error("Should not allow drop in full column")
	}
}

func TestIsBoardFull(t *testing.T) {
	board := NewBoard()

	if board.IsBoardFull() {
		t.Error("Empty board reported as full")
	}

	// Fill entire board
	for c := 0; c < Columns; c++ {
		for r := 0; r < Rows; r++ {
			board.DropDisc(c, Player1)
		}
	}

	if !board.IsBoardFull() {
		t.Error("Full board not reported as full")
	}
}

func TestClone(t *testing.T) {
	board := NewBoard()
	board.DropDisc(3, Player1)
	board.DropDisc(4, Player2)

	clone := board.Clone()

	// Clone should have same values
	if clone[5][3] != Player1 || clone[5][4] != Player2 {
		t.Error("Clone doesn't match original")
	}

	// Modifying clone shouldn't affect original
	clone.DropDisc(3, Player2)
	if board[4][3] != Empty {
		t.Error("Original was modified when clone was modified")
	}
}

func TestToSlice(t *testing.T) {
	board := NewBoard()
	board.DropDisc(0, Player1)
	board.DropDisc(1, Player2)

	slice := board.ToSlice()

	if len(slice) != Rows {
		t.Errorf("Expected %d rows, got %d", Rows, len(slice))
	}
	if len(slice[0]) != Columns {
		t.Errorf("Expected %d columns, got %d", Columns, len(slice[0]))
	}
	if slice[5][0] != 1 {
		t.Error("Player1 disc not in slice")
	}
	if slice[5][1] != 2 {
		t.Error("Player2 disc not in slice")
	}
}
