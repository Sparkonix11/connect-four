package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"connect-four/internal/api/handlers"
	"connect-four/internal/api/middleware"
	"connect-four/internal/kafka"
	"connect-four/internal/matchmaking"
	"connect-four/internal/repository"
	ws "connect-four/internal/websocket"
	"connect-four/pkg/config"
)

// Server holds all dependencies for the API server
type Server struct {
	Router         *mux.Router
	Hub            *ws.Hub
	MessageHandler *ws.MessageHandler
	MatchQueue     *matchmaking.Queue
	upgrader       websocket.Upgrader
}

// NewServer creates a new API server with all routes configured
func NewServer(db *gorm.DB, cfg *config.Config, kafkaProducer *kafka.Producer) *Server {
	router := mux.NewRouter()

	// Create repositories
	playerRepo := repository.NewPlayerRepository(db)
	gameRepo := repository.NewGameRepository(db)
	leaderboardRepo := repository.NewLeaderboardRepository(db)

	// Create handlers
	playerHandler := handlers.NewPlayerHandler(playerRepo)
	gameHandler := handlers.NewGameHandler(gameRepo, playerRepo)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardRepo)

	// Create WebSocket infrastructure
	hub := ws.NewHub(cfg.MatchmakingTimeout, cfg.ReconnectTimeout, cfg.BotMoveDelay)
	matchQueue := matchmaking.NewQueue(cfg.MatchmakingTimeout)
	messageHandler := ws.NewMessageHandler(hub, matchQueue, playerRepo, kafkaProducer)

	// Create server
	server := &Server{
		Router:         router,
		Hub:            hub,
		MessageHandler: messageHandler,
		MatchQueue:     matchQueue,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: Configure for production
			},
		},
	}

	// Apply middleware
	router.Use(middleware.Logging)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET", "HEAD")

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Player endpoints
	api.HandleFunc("/players", playerHandler.CreateOrGet).Methods("POST")
	api.HandleFunc("/players/{id}", playerHandler.GetByID).Methods("GET")
	api.HandleFunc("/players/{id}/games", gameHandler.GetPlayerGames).Methods("GET")

	// Leaderboard endpoints
	api.HandleFunc("/leaderboard", leaderboardHandler.GetTopPlayers).Methods("GET")

	// Game endpoints
	api.HandleFunc("/games/{id}", gameHandler.GetByID).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", server.handleWebSocket).Methods("GET")

	return server
}

// Start starts the WebSocket hub and matchmaking queue
func (s *Server) Start() {
	go s.Hub.Run()
	s.MatchQueue.Start()
	log.Info().Msg("WebSocket hub and matchmaking queue started")
}

// handleWebSocket upgrades HTTP to WebSocket and registers the client
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade failed")
		return
	}

	client := ws.NewClient(s.Hub, conn, username)
	s.Hub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump(s.MessageHandler.HandleMessage)
}
