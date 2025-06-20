package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"
	flightlog "real-time-dashboard/log"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
		flightlog.LogWarn("REDIS_URL not set, using default: %s", redisURL)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		flightlog.LogWarn("Redis connection failed: %v", err)
		return err
	}

	flightlog.LogInfo("Redis connected successfully")
	return nil
}

func SetFlightData(key string, data interface{}, expiration time.Duration) error {
	if rdb == nil {
		return nil // Skip if Redis not available
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, jsonData, expiration).Err()
}

func GetFlightData(key string, dest interface{}) error {
	if rdb == nil {
		return redis.Nil // Skip if Redis not available
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}