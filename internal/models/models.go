package models

import (
	"time"

	"github.com/google/uuid"
)

// Player represents a user in the system
type Player struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Wins      int       `json:"wins"`
	Losses    int       `json:"losses"`
	Draws     int       `json:"draws"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GameResult represents the outcome of a game
type GameResult string

const (
	GameResultPlayer1 GameResult = "player1"
	GameResultPlayer2 GameResult = "player2"
	GameResultDraw    GameResult = "draw"
	GameResultForfeit GameResult = "forfeit"
)

// Move represents a single move in a game
type Move struct {
	Player   int       `json:"player"`   // 1 or 2
	Column   int       `json:"column"`   // 0-6
	Row      int       `json:"row"`      // 0-5 (calculated after drop)
	MoveNum  int       `json:"move_num"` // Move number in sequence
	PlayedAt time.Time `json:"played_at"`
}

// Game represents a game record in the database
type Game struct {
	ID              uuid.UUID  `json:"id"`
	Player1ID       uuid.UUID  `json:"player1_id"`
	Player2ID       *uuid.UUID `json:"player2_id"` // nil if bot game
	IsBotGame       bool       `json:"is_bot_game"`
	WinnerID        *uuid.UUID `json:"winner_id"`
	Result          GameResult `json:"result"`
	Moves           []Move     `json:"moves"`
	DurationSeconds int        `json:"duration_seconds"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
}

// GameEvent represents an analytics event
type GameEvent struct {
	ID        uuid.UUID              `json:"id"`
	GameID    uuid.UUID              `json:"game_id"`
	EventType string                 `json:"event_type"`
	EventData map[string]interface{} `json:"event_data"`
	CreatedAt time.Time              `json:"created_at"`
}

// LeaderboardEntry represents a player's ranking
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Draws    int    `json:"draws"`
}

// WebSocket message types
type WSMessageType string

const (
	// Client -> Server
	WSTypeJoinQueue  WSMessageType = "join_queue"
	WSTypeMakeMove   WSMessageType = "make_move"
	WSTypeReconnect  WSMessageType = "reconnect"
	WSTypeLeaveGame  WSMessageType = "leave_game"

	// Server -> Client
	WSTypeQueueJoined          WSMessageType = "queue_joined"
	WSTypeGameStarted          WSMessageType = "game_started"
	WSTypeMoveMade             WSMessageType = "move_made"
	WSTypeInvalidMove          WSMessageType = "invalid_move"
	WSTypeGameOver             WSMessageType = "game_over"
	WSTypeOpponentDisconnected WSMessageType = "opponent_disconnected"
	WSTypeOpponentReconnected  WSMessageType = "opponent_reconnected"
	WSTypeGameForfeited        WSMessageType = "game_forfeited"
	WSTypeError                WSMessageType = "error"
	WSTypeGameState            WSMessageType = "game_state"
)

// WSMessage is the envelope for WebSocket messages
type WSMessage struct {
	Type      WSMessageType `json:"type"`
	Payload   interface{}   `json:"payload"`
	Timestamp time.Time     `json:"timestamp"`
}

// Client -> Server payloads
type JoinQueuePayload struct {
	Username string `json:"username"`
}

type MakeMovePayload struct {
	Column int `json:"column"`
}

type ReconnectPayload struct {
	GameID   string `json:"gameId"`
	Username string `json:"username"`
}

// Server -> Client payloads
type QueueJoinedPayload struct {
	Position int `json:"position"`
}

type GameStartedPayload struct {
	GameID    string `json:"gameId"`
	Opponent  string `json:"opponent"`
	YourTurn  bool   `json:"yourTurn"`
	YourColor int    `json:"yourColor"` // 1 = Red, 2 = Yellow
}

type MoveMadePayload struct {
	Column int     `json:"column"`
	Row    int     `json:"row"`
	Player int     `json:"player"`
	Board  [][]int `json:"board"`
}

type InvalidMovePayload struct {
	Reason string `json:"reason"`
}

type GameOverPayload struct {
	Winner     string  `json:"winner"` // username or "draw"
	Result     string  `json:"result"` // "win", "loss", "draw"
	FinalBoard [][]int `json:"finalBoard"`
}

type OpponentDisconnectedPayload struct {
	Timeout int `json:"timeout"` // seconds remaining
}

type GameForfeitedPayload struct {
	Winner string `json:"winner"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type GameStatePayload struct {
	GameID      string  `json:"gameId"`
	Board       [][]int `json:"board"`
	CurrentTurn int     `json:"currentTurn"`
	YourColor   int     `json:"yourColor"`
	YourTurn    bool    `json:"yourTurn"`
	Opponent    string  `json:"opponent"`
}
