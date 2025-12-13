package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// GameStatus represents the current state of a game
type GameStatus string

const (
	GameStatusWaiting      GameStatus = "waiting"
	GameStatusInProgress   GameStatus = "in_progress"
	GameStatusDisconnected GameStatus = "disconnected"
	GameStatusFinished     GameStatus = "finished"
)

// GameResult represents the outcome of a finished game
type GameResult string

const (
	ResultPlayer1Win GameResult = "player1"
	ResultPlayer2Win GameResult = "player2"
	ResultDraw       GameResult = "draw"
	ResultForfeit    GameResult = "forfeit"
)

// PlayerInfo holds information about a player in a game
type PlayerInfo struct {
	ID             uuid.UUID
	Username       string
	IsBot          bool
	Connected      bool
	DisconnectedAt *time.Time
}

// Move represents a single move in the game
type Move struct {
	Player    Cell      `json:"player"`
	Column    int       `json:"column"`
	Row       int       `json:"row"`
	MoveNum   int       `json:"move_num"`
	Timestamp time.Time `json:"timestamp"`
}

// Game represents an active game session
type Game struct {
	ID           uuid.UUID
	Player1      *PlayerInfo
	Player2      *PlayerInfo
	Board        *Board
	CurrentTurn  Cell
	Moves        []Move
	Status       GameStatus
	Result       GameResult
	Winner       Cell
	WinningCells [][2]int // Coordinates of winning cells
	StartedAt    time.Time
	EndedAt      *time.Time

	mu sync.RWMutex
}

// NewGame creates a new game session
func NewGame(player1, player2 *PlayerInfo) *Game {
	return &Game{
		ID:          uuid.New(),
		Player1:     player1,
		Player2:     player2,
		Board:       NewBoard(),
		CurrentTurn: Player1, // Player 1 always goes first
		Moves:       make([]Move, 0),
		Status:      GameStatusInProgress,
		StartedAt:   time.Now(),
	}
}

// MakeMove attempts to make a move in the specified column
// Returns the row where disc landed, or error message
func (g *Game) MakeMove(player Cell, col int) (int, string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Validate game status
	if g.Status != GameStatusInProgress {
		return -1, "game is not in progress"
	}

	// Validate turn
	if g.CurrentTurn != player {
		return -1, "not your turn"
	}

	// Validate column
	if col < 0 || col >= Columns {
		return -1, "invalid column"
	}

	// Drop disc
	row := g.Board.DropDisc(col, player)
	if row == -1 {
		return -1, "column is full"
	}

	// Record move
	move := Move{
		Player:    player,
		Column:    col,
		Row:       row,
		MoveNum:   len(g.Moves) + 1,
		Timestamp: time.Now(),
	}
	g.Moves = append(g.Moves, move)

	// Check for win
	if won, cells := g.checkWin(row, col, player); won {
		g.Status = GameStatusFinished
		g.Winner = player
		g.WinningCells = cells
		now := time.Now()
		g.EndedAt = &now
		if player == Player1 {
			g.Result = ResultPlayer1Win
		} else {
			g.Result = ResultPlayer2Win
		}
		return row, ""
	}

	// Check for draw
	if g.Board.IsBoardFull() {
		g.Status = GameStatusFinished
		g.Result = ResultDraw
		now := time.Now()
		g.EndedAt = &now
		return row, ""
	}

	// Switch turn
	if g.CurrentTurn == Player1 {
		g.CurrentTurn = Player2
	} else {
		g.CurrentTurn = Player1
	}

	return row, ""
}

// checkWin checks if the last move at (row, col) creates a win
// Returns true and the winning cells if a win is found
func (g *Game) checkWin(row, col int, player Cell) (bool, [][2]int) {
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal (down-right)
		{1, -1}, // Diagonal (down-left)
	}

	for _, dir := range directions {
		cells := g.countLine(row, col, dir[0], dir[1], player)
		if len(cells) >= 4 {
			return true, cells
		}
	}

	return false, nil
}

// countLine counts connected cells in both directions along a line
func (g *Game) countLine(row, col, dRow, dCol int, player Cell) [][2]int {
	cells := [][2]int{{row, col}}

	// Count in positive direction
	for i := 1; i < 4; i++ {
		r, c := row+dRow*i, col+dCol*i
		if r < 0 || r >= Rows || c < 0 || c >= Columns {
			break
		}
		if g.Board.GetCell(r, c) != player {
			break
		}
		cells = append(cells, [2]int{r, c})
	}

	// Count in negative direction
	for i := 1; i < 4; i++ {
		r, c := row-dRow*i, col-dCol*i
		if r < 0 || r >= Rows || c < 0 || c >= Columns {
			break
		}
		if g.Board.GetCell(r, c) != player {
			break
		}
		cells = append(cells, [2]int{r, c})
	}

	return cells
}

// IsGameOver returns true if the game has ended
func (g *Game) IsGameOver() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.Status == GameStatusFinished
}

// GetCurrentPlayer returns the player whose turn it is
func (g *Game) GetCurrentPlayer() Cell {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.CurrentTurn
}

// GetPlayerInfo returns player info based on their cell type
func (g *Game) GetPlayerInfo(player Cell) *PlayerInfo {
	if player == Player1 {
		return g.Player1
	}
	return g.Player2
}

// GetOpponentInfo returns opponent info based on player's cell type
func (g *Game) GetOpponentInfo(player Cell) *PlayerInfo {
	if player == Player1 {
		return g.Player2
	}
	return g.Player1
}

// Forfeit ends the game with the specified player as loser
func (g *Game) Forfeit(loser Cell) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Status = GameStatusFinished
	g.Result = ResultForfeit
	now := time.Now()
	g.EndedAt = &now

	if loser == Player1 {
		g.Winner = Player2
	} else {
		g.Winner = Player1
	}
}

// SetDisconnected marks a player as disconnected
func (g *Game) SetDisconnected(player Cell) {
	g.mu.Lock()
	defer g.mu.Unlock()

	info := g.GetPlayerInfo(player)
	if info != nil {
		info.Connected = false
		now := time.Now()
		info.DisconnectedAt = &now
	}
	g.Status = GameStatusDisconnected
}

// SetReconnected marks a player as reconnected
func (g *Game) SetReconnected(player Cell) {
	g.mu.Lock()
	defer g.mu.Unlock()

	info := g.GetPlayerInfo(player)
	if info != nil {
		info.Connected = true
		info.DisconnectedAt = nil
	}
	g.Status = GameStatusInProgress
}

// Duration returns the game duration in seconds
func (g *Game) Duration() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	endTime := time.Now()
	if g.EndedAt != nil {
		endTime = *g.EndedAt
	}
	return int(endTime.Sub(g.StartedAt).Seconds())
}
