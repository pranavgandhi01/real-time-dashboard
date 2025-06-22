package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	RedisURL     string
	KafkaBroker  string
	KafkaTopic   string
	FetchInterval time.Duration
	MaxConnections int
	RateLimitPerIP int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:32092"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "flight-events"),
		FetchInterval:  getDuration("FETCH_INTERVAL", "15s"),
		MaxConnections: getInt("MAX_CONNECTIONS", 1000),
		RateLimitPerIP: getInt("RATE_LIMIT_PER_IP", 5),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	d, _ := time.ParseDuration(defaultValue)
	return d
}