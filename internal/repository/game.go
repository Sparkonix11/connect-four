package repository

import (
	"connect-four/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GameRepository handles game database operations
type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// Create creates a new game record
func (r *GameRepository) Create(game *models.GameRecord) error {
	return r.db.Create(game).Error
}

// GetByID retrieves a game by ID
func (r *GameRepository) GetByID(id uuid.UUID) (*models.GameRecord, error) {
	var game models.GameRecord
	err := r.db.Preload("Player1").Preload("Player2").Preload("Winner").
		First(&game, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

// GetByPlayerID retrieves games for a player
func (r *GameRepository) GetByPlayerID(playerID uuid.UUID, limit int) ([]models.GameRecord, error) {
	var games []models.GameRecord
	err := r.db.Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Preload("Player1").Preload("Player2").
		Order("ended_at DESC NULLS LAST").
		Limit(limit).
		Find(&games).Error
	if err != nil {
		return nil, err
	}
	return games, nil
}

// Update updates a game record
func (r *GameRepository) Update(game *models.GameRecord) error {
	return r.db.Save(game).Error
}
