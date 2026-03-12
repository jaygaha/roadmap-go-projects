package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/database"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/models"
)

type CreateGameRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game := models.Game{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Game already exists"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

func ListGames(c *gin.Context) {
	var games []models.Game
	if err := database.DB.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games"})
		return
	}

	c.JSON(http.StatusOK, games)
}
