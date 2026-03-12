package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/services"
)

func GetGlobalLeaderboard(c *gin.Context) {
	entries, err := services.GetLeaderboard("leaderboard:global", 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func GetGameLeaderboard(c *gin.Context) {
	gameID := c.Param("gameId")
	key := "leaderboard:game:" + gameID

	entries, err := services.GetLeaderboard(key, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func GetUserRank(c *gin.Context) {
	userID := c.Param("userId")
	ctx := c.Request.Context()

	rank, err := services.GetUserGlobalRank(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "rank": rank})
}

func GetTopPlayers(c *gin.Context) {
	gameID := c.Param("gameId")
	period := c.Query("period") // daily|weekly|monthly

	key := services.GetTimePeriodKey(gameID, period)

	entries, err := services.GetLeaderboard(key, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
