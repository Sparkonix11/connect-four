package repository

import (
	"connect-four/internal/models"

	"gorm.io/gorm"
)

// LeaderboardRepository handles leaderboard queries
type LeaderboardRepository struct {
	db *gorm.DB
}

// NewLeaderboardRepository creates a new leaderboard repository
func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

// GetTopPlayers retrieves the top players by wins
func (r *LeaderboardRepository) GetTopPlayers(limit int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	var entries []models.LeaderboardEntry

	err := r.db.Model(&models.Player{}).
		Select("ROW_NUMBER() OVER (ORDER BY wins DESC, (wins - losses) DESC) as rank, username, wins, losses, draws, (wins + losses + draws) as games").
		Where("(wins + losses + draws) > 0").
		Order("wins DESC, (wins - losses) DESC").
		Limit(limit).
		Scan(&entries).Error

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// GetPlayerRank retrieves a specific player's rank
func (r *LeaderboardRepository) GetPlayerRank(username string) (*models.LeaderboardEntry, error) {
	var entry models.LeaderboardEntry

	subQuery := r.db.Model(&models.Player{}).
		Select("ROW_NUMBER() OVER (ORDER BY wins DESC, (wins - losses) DESC) as rank, username, wins, losses, draws, (wins + losses + draws) as games").
		Where("(wins + losses + draws) > 0")

	err := r.db.Table("(?) as ranked", subQuery).
		Where("username = ?", username).
		Scan(&entry).Error

	if err != nil {
		return nil, err
	}

	if entry.Username == "" {
		return nil, nil
	}

	return &entry, nil
}
