package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/database"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/models"
)

type LeaderboardEntry struct {
	Rank     int64  `json:"rank"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

func UpdateGlobalLeaderboard(userID uint, score int) error {
	ctx := context.Background()
	key := "leaderboard:global"
	member := fmt.Sprintf("%d", userID)
	z := &redis.Z{Score: float64(score), Member: member}
	return database.RedisClient.ZAdd(ctx, key, z).Err()
}

func GetUserRankInGame(gameID string, userID uint) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("leaderboard:game:%s", gameID)
	member := fmt.Sprintf("%d", userID)
	rank, err := database.RedisClient.ZRevRank(ctx, key, member).Result()
	if err != nil {
		return -1, err
	}
	return rank + 1, nil
}

func AddScoreToLeaderboards(userID uint, gameID string, score int) error {
	ctx := context.Background()
	member := strconv.Itoa(int(userID))

	// Global leaderboard (keep highest score)
	database.RedisClient.ZAdd(ctx, "leaderboard:global", &redis.Z{
		Score:  float64(score),
		Member: member,
	})

	// Game-specific leaderboard
	gameKey := "leaderboard:game:" + gameID
	database.RedisClient.ZAdd(ctx, gameKey, &redis.Z{
		Score:  float64(score),
		Member: member,
	})

	// Time-based leaderboards with date-scoped keys
	now := time.Now()
	dailyKey := fmt.Sprintf("%s:daily:%s", gameKey, now.Format("2006-01-02"))
	year, week := now.ISOWeek()
	weeklyKey := fmt.Sprintf("%s:weekly:%d-W%02d", gameKey, year, week)
	monthlyKey := fmt.Sprintf("%s:monthly:%s", gameKey, now.Format("2006-01"))

	for _, key := range []string{dailyKey, weeklyKey, monthlyKey} {
		database.RedisClient.ZAdd(ctx, key, &redis.Z{
			Score:  float64(score),
			Member: member,
		})
	}

	// Set TTL on time-based keys
	database.RedisClient.Expire(ctx, dailyKey, 48*time.Hour)
	database.RedisClient.Expire(ctx, weeklyKey, 8*24*time.Hour)
	database.RedisClient.Expire(ctx, monthlyKey, 32*24*time.Hour)

	return nil
}

func GetLeaderboard(key string, count int64) ([]LeaderboardEntry, error) {
	ctx := context.Background()
	scores, err := database.RedisClient.ZRevRangeWithScores(ctx, key, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]LeaderboardEntry, 0, len(scores))
	for i, z := range scores {
		uid, _ := strconv.ParseUint(z.Member.(string), 10, 64)

		var user models.User
		username := fmt.Sprintf("user_%d", uid)
		if err := database.DB.First(&user, uid).Error; err == nil {
			username = user.Username
		}

		entries = append(entries, LeaderboardEntry{
			Rank:     int64(i + 1),
			UserID:   uint(uid),
			Username: username,
			Score:    int(z.Score),
		})
	}
	return entries, nil
}

func GetUserGlobalRank(ctx context.Context, userID string) (int64, error) {
	rank, err := database.RedisClient.ZRevRank(ctx, "leaderboard:global", userID).Result()
	if err != nil {
		return -1, err
	}
	return rank + 1, nil
}

func GetTimePeriodKey(gameID, period string) string {
	now := time.Now()
	gameKey := "leaderboard:game:" + gameID

	switch period {
	case "daily":
		return fmt.Sprintf("%s:daily:%s", gameKey, now.Format("2006-01-02"))
	case "weekly":
		year, week := now.ISOWeek()
		return fmt.Sprintf("%s:weekly:%d-W%02d", gameKey, year, week)
	case "monthly":
		return fmt.Sprintf("%s:monthly:%s", gameKey, now.Format("2006-01"))
	default:
		return gameKey
	}
}
