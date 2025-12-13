package repository

import (
	"connect-four/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayerRepository handles player database operations
type PlayerRepository struct {
	db *gorm.DB
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// Create creates a new player or returns existing one with the same username
func (r *PlayerRepository) Create(username string) (*models.Player, error) {
	player := &models.Player{Username: username}

	// Upsert: create or update on conflict
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(player).Error

	if err != nil {
		return nil, err
	}

	// Fetch the full record
	return r.GetByUsername(username)
}

// GetByID retrieves a player by ID
func (r *PlayerRepository) GetByID(id uuid.UUID) (*models.Player, error) {
	var player models.Player
	err := r.db.First(&player, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &player, nil
}

// GetByUsername retrieves a player by username
func (r *PlayerRepository) GetByUsername(username string) (*models.Player, error) {
	var player models.Player
	err := r.db.First(&player, "username = ?", username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &player, nil
}

// IncrementWins increments a player's win count
func (r *PlayerRepository) IncrementWins(id uuid.UUID) error {
	return r.db.Model(&models.Player{}).Where("id = ?", id).
		UpdateColumn("wins", gorm.Expr("wins + 1")).Error
}

// IncrementLosses increments a player's loss count
func (r *PlayerRepository) IncrementLosses(id uuid.UUID) error {
	return r.db.Model(&models.Player{}).Where("id = ?", id).
		UpdateColumn("losses", gorm.Expr("losses + 1")).Error
}

// IncrementDraws increments a player's draw count
func (r *PlayerRepository) IncrementDraws(id uuid.UUID) error {
	return r.db.Model(&models.Player{}).Where("id = ?", id).
		UpdateColumn("draws", gorm.Expr("draws + 1")).Error
}
