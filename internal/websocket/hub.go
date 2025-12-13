package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"connect-four/internal/game"
	"connect-four/internal/models"
)

// Hub maintains the set of active clients and manages game sessions
type Hub struct {
	// Registered clients by username
	clients map[string]*Client

	// Active games by game ID
	games map[uuid.UUID]*GameSession

	// Player to game mapping
	playerGames map[string]uuid.UUID

	// Matchmaking queue
	matchQueue chan *Client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe access
	mu sync.RWMutex

	// Configuration
	matchmakingTimeout time.Duration
	reconnectTimeout   time.Duration
	botMoveDelay       time.Duration
}

// GameSession wraps a game with its connected clients
type GameSession struct {
	Game    *game.Game
	Player1 *Client
	Player2 *Client // nil if bot game
	IsBot   bool
}

// NewHub creates a new Hub instance
func NewHub(matchmakingTimeout, reconnectTimeout, botMoveDelay time.Duration) *Hub {
	return &Hub{
		clients:            make(map[string]*Client),
		games:              make(map[uuid.UUID]*GameSession),
		playerGames:        make(map[string]uuid.UUID),
		matchQueue:         make(chan *Client, 100),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		matchmakingTimeout: matchmakingTimeout,
		reconnectTimeout:   reconnectTimeout,
		botMoveDelay:       botMoveDelay,
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)
		}
	}
}

// handleRegister adds a new client to the hub
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if player is reconnecting to an active game
	if gameID, exists := h.playerGames[client.Username]; exists {
		if session, ok := h.games[gameID]; ok {
			h.handleReconnection(client, session)
			return
		}
	}

	h.clients[client.Username] = client
	log.Info().Str("username", client.Username).Msg("Client registered")
}

// handleUnregister removes a client and handles disconnection
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.Username]; ok {
		delete(h.clients, client.Username)
		close(client.send)

		// Check if player was in a game
		if gameID, exists := h.playerGames[client.Username]; exists {
			if session, ok := h.games[gameID]; ok {
				h.handleDisconnection(client, session)
			}
		}
	}

	log.Info().Str("username", client.Username).Msg("Client unregistered")
}

// handleReconnection handles a player reconnecting to an active game
func (h *Hub) handleReconnection(client *Client, session *GameSession) {
	// Update client reference
	h.clients[client.Username] = client

	// Determine player color
	playerColor := game.Player1
	if session.Game.Player2 != nil && session.Game.Player2.Username == client.Username {
		playerColor = game.Player2
		session.Player2 = client
	} else {
		session.Player1 = client
	}

	// Mark as reconnected
	session.Game.SetReconnected(playerColor)

	// Send current game state
	client.SendMessage(models.WSTypeGameState, models.GameStatePayload{
		GameID:      session.Game.ID.String(),
		Board:       session.Game.Board.ToSlice(),
		CurrentTurn: int(session.Game.CurrentTurn),
		YourColor:   int(playerColor),
		YourTurn:    session.Game.CurrentTurn == playerColor,
		Opponent:    session.Game.GetOpponentInfo(playerColor).Username,
	})

	// Notify opponent
	opponent := session.Player1
	if playerColor == game.Player1 {
		opponent = session.Player2
	}
	if opponent != nil {
		opponent.SendMessage(models.WSTypeOpponentReconnected, nil)
	}

	log.Info().Str("username", client.Username).Str("gameId", session.Game.ID.String()).Msg("Player reconnected")
}

// handleDisconnection handles a player disconnecting from a game
func (h *Hub) handleDisconnection(client *Client, session *GameSession) {
	playerColor := game.Player1
	if session.Game.Player2 != nil && session.Game.Player2.Username == client.Username {
		playerColor = game.Player2
	}

	// Mark as disconnected
	session.Game.SetDisconnected(playerColor)

	// Notify opponent
	opponent := session.Player1
	if playerColor == game.Player1 {
		opponent = session.Player2
	}
	if opponent != nil {
		opponent.SendMessage(models.WSTypeOpponentDisconnected, models.OpponentDisconnectedPayload{
			Timeout: int(h.reconnectTimeout.Seconds()),
		})
	}

	// Start reconnection timeout
	go h.startReconnectTimer(session, playerColor)

	log.Info().Str("username", client.Username).Msg("Player disconnected from game")
}

// startReconnectTimer waits for reconnection or forfeits the game
func (h *Hub) startReconnectTimer(session *GameSession, disconnectedPlayer game.Cell) {
	time.Sleep(h.reconnectTimeout)

	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if game still exists and player still disconnected
	if session.Game.Status != game.GameStatusDisconnected {
		return
	}

	// Forfeit the game
	session.Game.Forfeit(disconnectedPlayer)

	// Notify the connected opponent
	winner := session.Player1
	winnerInfo := session.Game.Player1
	if disconnectedPlayer == game.Player1 {
		winner = session.Player2
		winnerInfo = session.Game.Player2
	}

	if winner != nil && winnerInfo != nil {
		winner.SendMessage(models.WSTypeGameForfeited, models.GameForfeitedPayload{
			Winner: winnerInfo.Username,
		})
	}

	// Cleanup
	h.cleanupGame(session)

	log.Info().Str("gameId", session.Game.ID.String()).Msg("Game forfeited due to disconnect timeout")
}

// cleanupGame removes a finished game from tracking
func (h *Hub) cleanupGame(session *GameSession) {
	delete(h.games, session.Game.ID)
	if session.Game.Player1 != nil {
		delete(h.playerGames, session.Game.Player1.Username)
	}
	if session.Game.Player2 != nil {
		delete(h.playerGames, session.Game.Player2.Username)
	}
}

// GetClient returns a client by username
func (h *Hub) GetClient(username string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[username]
}

// GetGameSession returns a game session by ID
func (h *Hub) GetGameSession(gameID uuid.UUID) *GameSession {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.games[gameID]
}

// BroadcastToGame sends a message to all players in a game
func (h *Hub) BroadcastToGame(gameID uuid.UUID, msgType models.WSMessageType, payload interface{}) {
	h.mu.RLock()
	session, exists := h.games[gameID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	if session.Player1 != nil {
		session.Player1.SendMessage(msgType, payload)
	}
	if session.Player2 != nil {
		session.Player2.SendMessage(msgType, payload)
	}
}

// CreateGame creates a new game session
func (h *Hub) CreateGame(player1, player2 *Client, isBot bool) *GameSession {
	h.mu.Lock()
	defer h.mu.Unlock()

	p1Info := &game.PlayerInfo{
		ID:        uuid.New(),
		Username:  player1.Username,
		IsBot:     false,
		Connected: true,
	}

	var p2Info *game.PlayerInfo
	if player2 != nil {
		p2Info = &game.PlayerInfo{
			ID:        uuid.New(),
			Username:  player2.Username,
			IsBot:     false,
			Connected: true,
		}
	} else {
		p2Info = &game.PlayerInfo{
			ID:       uuid.New(),
			Username: "Bot",
			IsBot:    true,
		}
	}

	g := game.NewGame(p1Info, p2Info)

	session := &GameSession{
		Game:    g,
		Player1: player1,
		Player2: player2,
		IsBot:   isBot,
	}

	h.games[g.ID] = session
	h.playerGames[player1.Username] = g.ID
	if player2 != nil {
		h.playerGames[player2.Username] = g.ID
	}

	return session
}

// SendMessage is a helper to create and serialize a WebSocket message
func CreateMessage(msgType models.WSMessageType, payload interface{}) ([]byte, error) {
	msg := models.WSMessage{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	return json.Marshal(msg)
}
