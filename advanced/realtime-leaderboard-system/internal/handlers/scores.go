package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/database"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/models"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/services"
)

type SubmitScoreRequest struct {
	GameID string `json:"game_id" binding:"required"`
	Score  int    `json:"score" binding:"required,min=0"`
}

func SubmitScore(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify game exists
	var game models.Game
	if err := database.DB.First(&game, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	scoreEntry := models.Score{
		UserID:      userID,
		GameID:      req.GameID,
		Score:       req.Score,
		SubmittedAt: time.Now(),
	}

	if err := database.DB.Create(&scoreEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
		return
	}

	if err := services.AddScoreToLeaderboards(userID, req.GameID, req.Score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Score submitted successfully"})
}

func GetScoreHistory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var scores []models.Score
	if err := database.DB.Where("user_id = ?", userID).Order("submitted_at DESC").Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch score history"})
		return
	}

	c.JSON(http.StatusOK, scores)
}
