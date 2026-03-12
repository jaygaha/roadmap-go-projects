package services

import (
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/database"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/models"
)

func GetUserScoreHistory(userID uint, gameID string) ([]models.Score, error) {
	var scores []models.Score
	query := database.DB.Where("user_id = ?", userID)
	if gameID != "" {
		query = query.Where("game_id = ?", gameID)
	}
	err := query.Order("submitted_at DESC").Find(&scores).Error
	return scores, err
}
