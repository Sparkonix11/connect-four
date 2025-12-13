package models

import (
	"time"
)

// WebSocket message types
type WSMessageType string

const (
	// Client -> Server
	WSTypeJoinQueue WSMessageType = "join_queue"
	WSTypeMakeMove  WSMessageType = "make_move"
	WSTypeReconnect WSMessageType = "reconnect"
	WSTypeLeaveGame WSMessageType = "leave_game"

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
