package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"connect-four/internal/bot"
	"connect-four/internal/game"
	"connect-four/internal/kafka"
	"connect-four/internal/matchmaking"
	"connect-four/internal/models"
	"connect-four/internal/repository"
)

// MessageHandler processes incoming WebSocket messages
type MessageHandler struct {
	hub           *Hub
	matchQueue    *matchmaking.Queue
	botEngine     *bot.Bot
	playerRepo    *repository.PlayerRepository
	kafkaProducer *kafka.Producer
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(hub *Hub, matchQueue *matchmaking.Queue, playerRepo *repository.PlayerRepository, kafkaProducer *kafka.Producer) *MessageHandler {
	return &MessageHandler{
		hub:           hub,
		matchQueue:    matchQueue,
		botEngine:     bot.NewBot(),
		playerRepo:    playerRepo,
		kafkaProducer: kafkaProducer,
	}
}

// HandleMessage routes incoming messages to appropriate handlers
func (h *MessageHandler) HandleMessage(client *Client, data []byte) {
	var msg models.WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		client.SendError("Invalid message format")
		return
	}

	switch msg.Type {
	case models.WSTypeJoinQueue:
		h.handleJoinQueue(client, msg.Payload)
	case models.WSTypeMakeMove:
		h.handleMakeMove(client, msg.Payload)
	case models.WSTypeLeaveGame:
		h.handleLeaveGame(client)
	case models.WSTypeResumeSession:
		h.handleResumeSession(client)
	case models.WSTypeAbandonSession:
		h.handleAbandonSession(client)
	default:
		client.SendError("Unknown message type")
	}
}

// handleJoinQueue adds a player to the matchmaking queue
func (h *MessageHandler) handleJoinQueue(client *Client, payload interface{}) {
	// Add to matchmaking queue with callbacks
	h.matchQueue.AddPlayer(
		client.Username,
		// On match with another player
		func(opponent *matchmaking.Player, isBot bool) {
			opponentClient := h.hub.GetClient(opponent.Username)
			if opponentClient == nil {
				log.Warn().Str("username", opponent.Username).Msg("Matched opponent not found")
				return
			}
			// client (the one who was waiting in queue) is Player 1 (first turn)
			// opponentClient (the one who just joined) is Player 2
			h.startGame(client, opponentClient, false)
		},
		// On timeout - start bot game
		func() {
			h.startBotGame(client)
		},
	)

	// Send queue position
	pos := h.matchQueue.QueuePosition(client.Username)
	client.SendMessage(models.WSTypeQueueJoined, models.QueueJoinedPayload{
		Position: pos,
	})

	log.Info().Str("username", client.Username).Int("position", pos).Msg("Player joined queue")
}

// startGame initializes a new game between two players
func (h *MessageHandler) startGame(player1, player2 *Client, isBot bool) {
	session := h.hub.CreateGame(player1, player2, isBot)

	// Notify Player 1
	player1.SendMessage(models.WSTypeGameStarted, models.GameStartedPayload{
		GameID:    session.Game.ID.String(),
		Opponent:  session.Game.Player2.Username,
		YourTurn:  true, // Player 1 always goes first
		YourColor: int(game.Player1),
	})

	// Notify Player 2
	player2.SendMessage(models.WSTypeGameStarted, models.GameStartedPayload{
		GameID:    session.Game.ID.String(),
		Opponent:  session.Game.Player1.Username,
		YourTurn:  false,
		YourColor: int(game.Player2),
	})

	log.Info().
		Str("gameId", session.Game.ID.String()).
		Str("player1", player1.Username).
		Str("player2", player2.Username).
		Msg("Game started")

	// Publish game started event to Kafka
	if h.kafkaProducer != nil {
		h.kafkaProducer.PublishGameStarted(context.Background(), session.Game.ID, player1.Username, player2.Username, isBot)
	}
}

// startBotGame initializes a game against the bot
func (h *MessageHandler) startBotGame(client *Client) {
	session := h.hub.CreateGame(client, nil, true)

	// Notify player
	client.SendMessage(models.WSTypeGameStarted, models.GameStartedPayload{
		GameID:    session.Game.ID.String(),
		Opponent:  "Bot",
		YourTurn:  true, // Player always goes first against bot
		YourColor: int(game.Player1),
	})

	log.Info().
		Str("gameId", session.Game.ID.String()).
		Str("player", client.Username).
		Msg("Bot game started")
}

// handleMakeMove processes a player's move
func (h *MessageHandler) handleMakeMove(client *Client, payload interface{}) {
	// Parse payload
	payloadBytes, _ := json.Marshal(payload)
	var movePayload models.MakeMovePayload
	if err := json.Unmarshal(payloadBytes, &movePayload); err != nil {
		client.SendError("Invalid move payload")
		return
	}

	// Find the game session
	session := h.findPlayerGame(client.Username)
	if session == nil {
		client.SendError("Not in a game")
		return
	}

	// Determine player color
	playerColor := game.Player1
	if session.Game.Player2 != nil && session.Game.Player2.Username == client.Username {
		playerColor = game.Player2
	}

	// Make the move
	row, errMsg := session.Game.MakeMove(playerColor, movePayload.Column)
	if errMsg != "" {
		client.SendMessage(models.WSTypeInvalidMove, models.InvalidMovePayload{
			Reason: errMsg,
		})
		return
	}

	// Broadcast move to all players
	boardState := session.Game.Board.ToSlice()
	moveMadePayload := models.MoveMadePayload{
		Column: movePayload.Column,
		Row:    row,
		Player: int(playerColor),
		Board:  boardState,
	}

	client.SendMessage(models.WSTypeMoveMade, moveMadePayload)
	if session.Player2 != nil {
		session.Player2.SendMessage(models.WSTypeMoveMade, moveMadePayload)
	}
	if session.Player1 != nil && session.Player1.Username != client.Username {
		session.Player1.SendMessage(models.WSTypeMoveMade, moveMadePayload)
	}

	// Publish move event to Kafka
	if h.kafkaProducer != nil {
		moveNum := len(session.Game.Moves)
		h.kafkaProducer.PublishGameMove(context.Background(), session.Game.ID, client.Username, movePayload.Column, moveNum)
	}

	// Check if game is over
	if session.Game.IsGameOver() {
		h.handleGameOver(session)
		return
	}

	// If bot game and it's bot's turn, make bot move
	if session.IsBot && session.Game.CurrentTurn == game.Player2 {
		go h.makeBotMove(session)
	}
}

// makeBotMove executes the bot's move with a small delay
func (h *MessageHandler) makeBotMove(session *GameSession) {
	// Add small delay for better UX
	time.Sleep(h.hub.botMoveDelay)

	// Get bot's move
	col := h.botEngine.SelectMove(session.Game.Board, game.Player2)
	if col == -1 {
		log.Error().Str("gameId", session.Game.ID.String()).Msg("Bot couldn't select move")
		return
	}

	// Make the move
	row, errMsg := session.Game.MakeMove(game.Player2, col)
	if errMsg != "" {
		log.Error().Str("error", errMsg).Msg("Bot made invalid move")
		return
	}

	// Send move to player
	boardState := session.Game.Board.ToSlice()
	session.Player1.SendMessage(models.WSTypeMoveMade, models.MoveMadePayload{
		Column: col,
		Row:    row,
		Player: int(game.Player2),
		Board:  boardState,
	})

	// Check if game is over
	if session.Game.IsGameOver() {
		h.handleGameOver(session)
	}
}

// handleGameOver sends game over messages and cleans up
func (h *MessageHandler) handleGameOver(session *GameSession) {
	var winnerName string
	result := "draw"

	switch session.Game.Result {
	case game.ResultPlayer1Win:
		winnerName = session.Game.Player1.Username
		result = "win"
	case game.ResultPlayer2Win:
		if session.IsBot {
			winnerName = "Bot"
		} else {
			winnerName = session.Game.Player2.Username
		}
		result = "win"
	case game.ResultDraw:
		winnerName = "draw"
		result = "draw"
	case game.ResultForfeit:
		if session.Game.Winner == game.Player1 {
			winnerName = session.Game.Player1.Username
		} else {
			winnerName = session.Game.Player2.Username
		}
		result = "forfeit"
	}

	gameOverPayload := models.GameOverPayload{
		Winner:     winnerName,
		Result:     result,
		FinalBoard: session.Game.Board.ToSlice(),
	}

	// Notify players
	if session.Player1 != nil {
		p1Result := result
		if result == "win" {
			if session.Game.Winner == game.Player1 {
				p1Result = "win"
			} else {
				p1Result = "loss"
			}
		}
		session.Player1.SendMessage(models.WSTypeGameOver, models.GameOverPayload{
			Winner:     winnerName,
			Result:     p1Result,
			FinalBoard: gameOverPayload.FinalBoard,
		})
	}

	if session.Player2 != nil {
		p2Result := result
		if result == "win" {
			if session.Game.Winner == game.Player2 {
				p2Result = "win"
			} else {
				p2Result = "loss"
			}
		}
		session.Player2.SendMessage(models.WSTypeGameOver, models.GameOverPayload{
			Winner:     winnerName,
			Result:     p2Result,
			FinalBoard: gameOverPayload.FinalBoard,
		})
	}

	log.Info().
		Str("gameId", session.Game.ID.String()).
		Str("winner", winnerName).
		Str("result", result).
		Msg("Game ended")

	// Persist game results to database
	if h.playerRepo != nil {
		// Create/get player records first
		p1, err := h.playerRepo.Create(session.Game.Player1.Username)
		if err != nil {
			log.Error().Err(err).Str("username", session.Game.Player1.Username).Msg("Failed to create/get player")
		}

		var p2 *models.Player
		if session.Game.Player2 != nil && !session.IsBot {
			p2, err = h.playerRepo.Create(session.Game.Player2.Username)
			if err != nil {
				log.Error().Err(err).Str("username", session.Game.Player2.Username).Msg("Failed to create/get player")
			}
		}

		// Update stats based on result
		switch session.Game.Result {
		case game.ResultPlayer1Win, game.ResultForfeit:
			if session.Game.Winner == game.Player1 {
				if p1 != nil {
					if err := h.playerRepo.IncrementWins(p1.ID); err != nil {
						log.Error().Err(err).Msg("Failed to increment wins")
					}
				}
				if p2 != nil {
					if err := h.playerRepo.IncrementLosses(p2.ID); err != nil {
						log.Error().Err(err).Msg("Failed to increment losses")
					}
				}
			} else if session.Game.Winner == game.Player2 {
				if p2 != nil {
					if err := h.playerRepo.IncrementWins(p2.ID); err != nil {
						log.Error().Err(err).Msg("Failed to increment wins")
					}
				}
				if p1 != nil {
					if err := h.playerRepo.IncrementLosses(p1.ID); err != nil {
						log.Error().Err(err).Msg("Failed to increment losses")
					}
				}
			}
		case game.ResultPlayer2Win:
			if p2 != nil {
				if err := h.playerRepo.IncrementWins(p2.ID); err != nil {
					log.Error().Err(err).Msg("Failed to increment wins")
				}
			}
			if p1 != nil {
				if err := h.playerRepo.IncrementLosses(p1.ID); err != nil {
					log.Error().Err(err).Msg("Failed to increment losses")
				}
			}
		case game.ResultDraw:
			if p1 != nil {
				if err := h.playerRepo.IncrementDraws(p1.ID); err != nil {
					log.Error().Err(err).Msg("Failed to increment draws")
				}
			}
			if p2 != nil {
				if err := h.playerRepo.IncrementDraws(p2.ID); err != nil {
					log.Error().Err(err).Msg("Failed to increment draws")
				}
			}
		}
		log.Info().Msg("Game stats persisted to database")
	}

	// Publish game ended event to Kafka
	if h.kafkaProducer != nil {
		h.kafkaProducer.PublishGameEnded(context.Background(), session.Game.ID, winnerName, result, 0, len(session.Game.Moves))
	}

	// Cleanup will happen after some delay or on next message
}

// handleLeaveGame handles voluntary game exit (forfeit)
func (h *MessageHandler) handleLeaveGame(client *Client) {
	session := h.findPlayerGame(client.Username)
	if session == nil {
		return
	}

	// Determine player color
	playerColor := game.Player1
	if session.Game.Player2 != nil && session.Game.Player2.Username == client.Username {
		playerColor = game.Player2
	}

	// Forfeit the game
	session.Game.Forfeit(playerColor)
	h.handleGameOver(session)
}

// findPlayerGame finds the game session for a player
func (h *MessageHandler) findPlayerGame(username string) *GameSession {
	h.hub.mu.RLock()
	defer h.hub.mu.RUnlock()

	for _, session := range h.hub.games {
		if session.Game.Player1.Username == username {
			return session
		}
		if session.Game.Player2 != nil && session.Game.Player2.Username == username {
			return session
		}
	}
	return nil
}

// handleResumeSession resumes an existing game session
func (h *MessageHandler) handleResumeSession(client *Client) {
	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()

	gameID, exists := h.hub.playerGames[client.Username]
	if !exists {
		client.SendError("No active session found")
		return
	}

	session, ok := h.hub.games[gameID]
	if !ok {
		client.SendError("Session no longer exists")
		delete(h.hub.playerGames, client.Username)
		return
	}

	// Reconnect to the game
	h.hub.handleReconnection(client, session)
	log.Info().Str("username", client.Username).Str("gameId", gameID.String()).Msg("Session resumed")
}

// handleAbandonSession abandons an existing session and allows fresh matchmaking
func (h *MessageHandler) handleAbandonSession(client *Client) {
	h.hub.mu.Lock()

	gameID, exists := h.hub.playerGames[client.Username]
	if !exists {
		h.hub.mu.Unlock()
		log.Info().Str("username", client.Username).Msg("No session to abandon")
		return
	}

	session, ok := h.hub.games[gameID]
	if ok {
		// Forfeit the game
		playerColor := game.Player1
		if session.Game.Player2 != nil && session.Game.Player2.Username == client.Username {
			playerColor = game.Player2
		}
		session.Game.Forfeit(playerColor)

		// Notify opponent if present
		var opponent *Client
		if playerColor == game.Player1 {
			opponent = session.Player2
		} else {
			opponent = session.Player1
		}
		if opponent != nil && !opponent.closed {
			opponent.SendMessage(models.WSTypeGameForfeited, models.GameForfeitedPayload{
				Winner: opponent.Username,
			})
		}

		// Clean up game
		delete(h.hub.games, gameID)
		delete(h.hub.playerGames, client.Username)
		if session.Game.Player1.Username != client.Username {
			delete(h.hub.playerGames, session.Game.Player1.Username)
		}
		if session.Game.Player2 != nil && session.Game.Player2.Username != client.Username {
			delete(h.hub.playerGames, session.Game.Player2.Username)
		}
	}

	h.hub.mu.Unlock()
	log.Info().Str("username", client.Username).Msg("Session abandoned")
}
