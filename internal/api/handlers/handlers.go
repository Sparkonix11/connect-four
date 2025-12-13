package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"connect-four/internal/repository"
)

// PlayerHandler handles player-related HTTP requests
type PlayerHandler struct {
	repo *repository.PlayerRepository
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(repo *repository.PlayerRepository) *PlayerHandler {
	return &PlayerHandler{repo: repo}
}

// CreateOrGet handles POST /api/players
// Creates a new player or returns existing one with the same username
func (h *PlayerHandler) CreateOrGet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if len(req.Username) > 50 {
		http.Error(w, "Username too long (max 50 chars)", http.StatusBadRequest)
		return
	}

	player, err := h.repo.Create(req.Username)
	if err != nil {
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// GetByID handles GET /api/players/{id}
func (h *PlayerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	player, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get player", http.StatusInternalServerError)
		return
	}

	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// LeaderboardHandler handles leaderboard-related HTTP requests
type LeaderboardHandler struct {
	repo *repository.LeaderboardRepository
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler(repo *repository.LeaderboardRepository) *LeaderboardHandler {
	return &LeaderboardHandler{repo: repo}
}

// GetTopPlayers handles GET /api/leaderboard
func (h *LeaderboardHandler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	entries, err := h.repo.GetTopPlayers(limit)
	if err != nil {
		http.Error(w, "Failed to get leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	repo       *repository.GameRepository
	playerRepo *repository.PlayerRepository
}

// NewGameHandler creates a new game handler
func NewGameHandler(repo *repository.GameRepository, playerRepo *repository.PlayerRepository) *GameHandler {
	return &GameHandler{repo: repo, playerRepo: playerRepo}
}

// GetByID handles GET /api/games/{id}
func (h *GameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get game", http.StatusInternalServerError)
		return
	}

	if game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// GetPlayerGames handles GET /api/players/{id}/games
func (h *GameHandler) GetPlayerGames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	games, err := h.repo.GetByPlayerID(id, limit)
	if err != nil {
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}
