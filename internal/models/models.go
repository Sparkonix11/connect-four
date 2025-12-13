// =============================================================================
// WEBSOCKET MESSAGE TYPES - Connect Four
// =============================================================================
// ⚠️  SYNC REQUIRED: These types MUST stay in sync with:
//     - shared/types/types.ts (TypeScript frontend types)
//     - shared/schema.json (JSON Schema - source of truth)
//
// When updating message types:
// 1. Update shared/schema.json first (source of truth)
// 2. Update this file to match
// 3. Update shared/types/types.ts to match
// 4. Run both `staticcheck ./...` and `npm run lint && npx tsc --noEmit`
// =============================================================================

package models

import (
	"time"
)

// WSMessageType defines the type of WebSocket message
// SYNC: shared/schema.json -> definitions.WSMessageType
type WSMessageType string

const (
	// Client -> Server
	WSTypeJoinQueue      WSMessageType = "join_queue"
	WSTypeMakeMove       WSMessageType = "make_move"
	WSTypeReconnect      WSMessageType = "reconnect"
	WSTypeLeaveGame      WSMessageType = "leave_game"
	WSTypeResumeSession  WSMessageType = "resume_session"
	WSTypeAbandonSession WSMessageType = "abandon_session"

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
	WSTypeExistingSession      WSMessageType = "existing_session"
)

// WSMessage is the envelope for WebSocket messages
// SYNC: shared/schema.json -> definitions.WSMessage
type WSMessage struct {
	Type      WSMessageType `json:"type"`
	Payload   interface{}   `json:"payload"`
	Timestamp time.Time     `json:"timestamp"`
}

// =============================================================================
// Client -> Server Payloads
// =============================================================================

// JoinQueuePayload - SYNC: shared/schema.json -> definitions.JoinQueuePayload
type JoinQueuePayload struct {
	Username string `json:"username"`
}

// MakeMovePayload - SYNC: shared/schema.json -> definitions.MakeMovePayload
type MakeMovePayload struct {
	Column int `json:"column"` // 0-6
}

// ReconnectPayload - SYNC: shared/schema.json -> definitions.ReconnectPayload
type ReconnectPayload struct {
	GameID   string `json:"gameId"`
	Username string `json:"username"`
}

// =============================================================================
// Server -> Client Payloads
// =============================================================================

// QueueJoinedPayload - SYNC: shared/schema.json -> definitions.QueueJoinedPayload
type QueueJoinedPayload struct {
	Position int `json:"position"`
}

// GameStartedPayload - SYNC: shared/schema.json -> definitions.GameStartedPayload
type GameStartedPayload struct {
	GameID    string `json:"gameId"`
	Opponent  string `json:"opponent"`
	YourTurn  bool   `json:"yourTurn"`
	YourColor int    `json:"yourColor"` // 1 = Red, 2 = Yellow
}

// MoveMadePayload - SYNC: shared/schema.json -> definitions.MoveMadePayload
type MoveMadePayload struct {
	Column int     `json:"column"`
	Row    int     `json:"row"`
	Player int     `json:"player"` // 1 or 2
	Board  [][]int `json:"board"`  // 6 rows x 7 columns, 0=empty, 1=P1, 2=P2
}

// InvalidMovePayload - SYNC: shared/schema.json -> definitions.InvalidMovePayload
type InvalidMovePayload struct {
	Reason string `json:"reason"`
}

// GameOverPayload - SYNC: shared/schema.json -> definitions.GameOverPayload
type GameOverPayload struct {
	Winner     string  `json:"winner"`     // username or "draw"
	Result     string  `json:"result"`     // "win", "loss", "draw", "forfeit"
	FinalBoard [][]int `json:"finalBoard"` // Final board state
}

// OpponentDisconnectedPayload - SYNC: shared/schema.json -> definitions.OpponentDisconnectedPayload
type OpponentDisconnectedPayload struct {
	Timeout int `json:"timeout"` // seconds remaining for reconnect
}

// GameForfeitedPayload - SYNC: shared/schema.json -> definitions.GameForfeitedPayload
type GameForfeitedPayload struct {
	Winner string `json:"winner"`
}

// ErrorPayload - SYNC: shared/schema.json -> definitions.ErrorPayload
type ErrorPayload struct {
	Message string `json:"message"`
}

// GameStatePayload - SYNC: shared/schema.json -> definitions.GameStatePayload
type GameStatePayload struct {
	GameID      string  `json:"gameId"`
	Board       [][]int `json:"board"`       // 6 rows x 7 columns
	CurrentTurn int     `json:"currentTurn"` // 1 or 2
	YourColor   int     `json:"yourColor"`   // 1 or 2
	YourTurn    bool    `json:"yourTurn"`
	Opponent    string  `json:"opponent"`
}

// ExistingSessionPayload - sent when player has an active game session
type ExistingSessionPayload struct {
	GameID   string `json:"gameId"`
	Opponent string `json:"opponent"`
	IsBot    bool   `json:"isBot"`
}
