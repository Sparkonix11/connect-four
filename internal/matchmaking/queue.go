package matchmaking

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Player represents a player waiting in the matchmaking queue
type Player struct {
	Username  string
	JoinedAt  time.Time
	OnMatch   func(opponent *Player, isBotGame bool) // Callback when matched
	OnTimeout func()                                 // Callback when bot assigned
}

// Queue manages matchmaking for players
type Queue struct {
	players    []*Player
	mu         sync.Mutex
	timeout    time.Duration
	addChan    chan *Player
	removeChan chan string
	stopChan   chan struct{}
}

// NewQueue creates a new matchmaking queue
func NewQueue(timeout time.Duration) *Queue {
	return &Queue{
		players:    make([]*Player, 0),
		timeout:    timeout,
		addChan:    make(chan *Player, 10),
		removeChan: make(chan string, 10),
		stopChan:   make(chan struct{}),
	}
}

// Start begins the matchmaking loop
func (q *Queue) Start() {
	go q.run()
	go q.checkTimeouts()
}

// Stop halts the matchmaking queue
func (q *Queue) Stop() {
	close(q.stopChan)
}

// AddPlayer adds a player to the matchmaking queue
func (q *Queue) AddPlayer(username string, onMatch func(*Player, bool), onTimeout func()) {
	player := &Player{
		Username:  username,
		JoinedAt:  time.Now(),
		OnMatch:   onMatch,
		OnTimeout: onTimeout,
	}
	q.addChan <- player
}

// RemovePlayer removes a player from the queue
func (q *Queue) RemovePlayer(username string) {
	q.removeChan <- username
}

// QueuePosition returns the player's position in queue (1-indexed)
func (q *Queue) QueuePosition(username string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, p := range q.players {
		if p.Username == username {
			return i + 1
		}
	}
	return 0
}

// run is the main matchmaking loop
func (q *Queue) run() {
	for {
		select {
		case <-q.stopChan:
			return

		case player := <-q.addChan:
			q.handleAdd(player)

		case username := <-q.removeChan:
			q.handleRemove(username)
		}
	}
}

// handleAdd adds a player and tries to match immediately
func (q *Queue) handleAdd(player *Player) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if player already in queue
	for _, p := range q.players {
		if p.Username == player.Username {
			log.Warn().Str("username", player.Username).Msg("Player already in queue")
			return
		}
	}

	// If there's another player waiting, match them
	if len(q.players) > 0 {
		opponent := q.players[0]
		q.players = q.players[1:]

		log.Info().
			Str("player1", opponent.Username).
			Str("player2", player.Username).
			Msg("Players matched")

		// Notify both players
		go opponent.OnMatch(player, false)
		go player.OnMatch(opponent, false)
		return
	}

	// Otherwise, add to queue
	q.players = append(q.players, player)
	log.Info().Str("username", player.Username).Int("queueSize", len(q.players)).Msg("Player added to queue")
}

// handleRemove removes a player from the queue
func (q *Queue) handleRemove(username string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, p := range q.players {
		if p.Username == username {
			q.players = append(q.players[:i], q.players[i+1:]...)
			log.Info().Str("username", username).Msg("Player removed from queue")
			return
		}
	}
}

// checkTimeouts periodically checks for players who have waited too long
func (q *Queue) checkTimeouts() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-q.stopChan:
			return

		case <-ticker.C:
			q.processTimeouts()
		}
	}
}

// processTimeouts handles players who have exceeded the wait time
func (q *Queue) processTimeouts() {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	var remaining []*Player

	for _, p := range q.players {
		if now.Sub(p.JoinedAt) >= q.timeout {
			log.Info().Str("username", p.Username).Msg("Matchmaking timeout - assigning bot")
			go p.OnTimeout()
		} else {
			remaining = append(remaining, p)
		}
	}

	q.players = remaining
}

// Size returns the current queue size
func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.players)
}
