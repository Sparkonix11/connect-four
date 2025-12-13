package bot

import (
	"connect-four/internal/game"
)

// Bot implements the strategic AI opponent
// Decision priority per SRS Appendix B:
// 1. Win if possible (immediate win in next move)
// 2. Block opponent's immediate win
// 3. Create winning opportunities (3-in-a-row with open ends)
// 4. Block opponent's potential winning paths
// 5. Play strategically (center column preference)
type Bot struct{}

// NewBot creates a new bot instance
func NewBot() *Bot {
	return &Bot{}
}

// SelectMove chooses the best column for the bot's next move
func (b *Bot) SelectMove(board *game.Board, botPlayer game.Cell) int {
	opponent := game.Player1
	if botPlayer == game.Player1 {
		opponent = game.Player2
	}

	validColumns := board.ValidColumns()
	if len(validColumns) == 0 {
		return -1
	}

	// Priority 1: Win if possible
	for _, col := range validColumns {
		if b.canWin(board, botPlayer, col) {
			return col
		}
	}

	// Priority 2: Block opponent's immediate win
	for _, col := range validColumns {
		if b.canWin(board, opponent, col) {
			return col
		}
	}

	// Priority 3: Create 3-in-a-row with potential to win
	for _, col := range validColumns {
		if b.createsWinningPath(board, botPlayer, col) {
			// Make sure this move doesn't give opponent a win
			if !b.givesOpponentWin(board, botPlayer, col) {
				return col
			}
		}
	}

	// Priority 4: Block opponent's 3-in-a-row
	for _, col := range validColumns {
		if b.createsWinningPath(board, opponent, col) {
			if !b.givesOpponentWin(board, botPlayer, col) {
				return col
			}
		}
	}

	// Priority 5: Center preference with strategic ordering
	preferredOrder := []int{3, 2, 4, 1, 5, 0, 6} // Center first
	for _, col := range preferredOrder {
		if b.isValidMove(board, col) && !b.givesOpponentWin(board, botPlayer, col) {
			return col
		}
	}

	// Fallback: any valid move
	for _, col := range preferredOrder {
		if b.isValidMove(board, col) {
			return col
		}
	}

	// Last resort
	return validColumns[0]
}

// canWin checks if playing in column will result in immediate win
func (b *Bot) canWin(board *game.Board, player game.Cell, col int) bool {
	testBoard := board.Clone()
	row := testBoard.DropDisc(col, player)
	if row == -1 {
		return false
	}
	return b.checkWinAt(testBoard, row, col, player)
}

// checkWinAt checks if there's a 4-in-a-row at the given position
func (b *Bot) checkWinAt(board *game.Board, row, col int, player game.Cell) bool {
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal down-right
		{1, -1}, // Diagonal down-left
	}

	for _, dir := range directions {
		count := 1

		// Count in positive direction
		for i := 1; i < 4; i++ {
			r, c := row+dir[0]*i, col+dir[1]*i
			if r < 0 || r >= game.Rows || c < 0 || c >= game.Columns {
				break
			}
			if board.GetCell(r, c) != player {
				break
			}
			count++
		}

		// Count in negative direction
		for i := 1; i < 4; i++ {
			r, c := row-dir[0]*i, col-dir[1]*i
			if r < 0 || r >= game.Rows || c < 0 || c >= game.Columns {
				break
			}
			if board.GetCell(r, c) != player {
				break
			}
			count++
		}

		if count >= 4 {
			return true
		}
	}

	return false
}

// createsWinningPath checks if playing creates a 3-in-a-row with potential to extend
func (b *Bot) createsWinningPath(board *game.Board, player game.Cell, col int) bool {
	testBoard := board.Clone()
	row := testBoard.DropDisc(col, player)
	if row == -1 {
		return false
	}

	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal down-right
		{1, -1}, // Diagonal down-left
	}

	for _, dir := range directions {
		count := 1
		openEnds := 0

		// Count and check open end in positive direction
		openInDir := false
		for i := 1; i < 4; i++ {
			r, c := row+dir[0]*i, col+dir[1]*i
			if r < 0 || r >= game.Rows || c < 0 || c >= game.Columns {
				break
			}
			if testBoard.GetCell(r, c) == player {
				count++
			} else if testBoard.GetCell(r, c) == game.Empty {
				openInDir = true
				break
			} else {
				break
			}
		}
		if openInDir {
			openEnds++
		}

		// Count and check open end in negative direction
		openInDir = false
		for i := 1; i < 4; i++ {
			r, c := row-dir[0]*i, col-dir[1]*i
			if r < 0 || r >= game.Rows || c < 0 || c >= game.Columns {
				break
			}
			if testBoard.GetCell(r, c) == player {
				count++
			} else if testBoard.GetCell(r, c) == game.Empty {
				openInDir = true
				break
			} else {
				break
			}
		}
		if openInDir {
			openEnds++
		}

		// 3-in-a-row with at least one open end is a winning path
		if count >= 3 && openEnds >= 1 {
			return true
		}
	}

	return false
}

// givesOpponentWin checks if this move allows opponent to win on next turn
func (b *Bot) givesOpponentWin(board *game.Board, botPlayer game.Cell, col int) bool {
	opponent := game.Player1
	if botPlayer == game.Player1 {
		opponent = game.Player2
	}

	testBoard := board.Clone()
	row := testBoard.DropDisc(col, botPlayer)
	if row == -1 {
		return false
	}

	// Check if opponent can win by playing on top of our move
	aboveRow := row - 1
	if aboveRow >= 0 && testBoard.GetCell(aboveRow, col) == game.Empty {
		// Simulate opponent playing on top
		testBoard2 := testBoard.Clone()
		testBoard2[aboveRow][col] = opponent
		if b.checkWinAt(testBoard2, aboveRow, col, opponent) {
			return true
		}
	}

	return false
}

// isValidMove checks if a column can accept a disc
func (b *Bot) isValidMove(board *game.Board, col int) bool {
	return col >= 0 && col < game.Columns && !board.IsColumnFull(col)
}
