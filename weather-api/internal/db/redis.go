package db

import (
	"github.com/go-redis/redis/v8"
)

// RedisClient is a wrapper for the redis client
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new redis client
func NewRedisClient(addr, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Ping the Redis server to check if it's alive
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		Client: rdb,
	}
}

// SetKey sets a key-value pair in Redis
func SetKey(rc *RedisClient, key string, value any) error {
	return rc.Client.Set(rc.Client.Context(), key, value, 0).Err()
}

// GetKey gets a value from Redis
func GetKey(rc *RedisClient, key string) (string, error) {
	value, err := rc.Client.Get(rc.Client.Context(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}

	return value, nil
}
