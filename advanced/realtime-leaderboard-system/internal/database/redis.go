package database

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/pkg/config"
)

var RedisClient *redis.Client

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
	RedisClient = rdb
}
