package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Player represents a user in the system (GORM model)
type Player struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Wins      int       `gorm:"default:0"`
	Losses    int       `gorm:"default:0"`
	Draws     int       `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate generates UUID if not set
func (p *Player) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// GameResultType represents the outcome of a game
type GameResultType string

const (
	GameResultPlayer1Win GameResultType = "player1"
	GameResultPlayer2Win GameResultType = "player2"
	GameResultDraw       GameResultType = "draw"
	GameResultForfeit    GameResultType = "forfeit"
)

// GameRecord represents a game in the database (GORM model)
type GameRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Player1ID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	Player1         *Player        `gorm:"foreignKey:Player1ID"`
	Player2ID       *uuid.UUID     `gorm:"type:uuid;index"`
	Player2         *Player        `gorm:"foreignKey:Player2ID"`
	IsBotGame       bool           `gorm:"default:false"`
	WinnerID        *uuid.UUID     `gorm:"type:uuid"`
	Winner          *Player        `gorm:"foreignKey:WinnerID"`
	Result          GameResultType `gorm:"size:10"`
	Moves           string         `gorm:"type:jsonb;default:'[]'"`
	DurationSeconds int            `gorm:"default:0"`
	StartedAt       time.Time
	EndedAt         *time.Time
	CreatedAt       time.Time
}

// GameEvent represents an analytics event (GORM model)
type GameEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GameID    uuid.UUID `gorm:"type:uuid;index"`
	EventType string    `gorm:"size:50;not null;index"`
	EventData string    `gorm:"type:jsonb"`
	CreatedAt time.Time
}

// LeaderboardEntry represents a player's ranking (used for API responses)
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Draws    int    `json:"draws"`
	Games    int    `json:"games"`
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Player{}, &GameRecord{}, &GameEvent{})
}
