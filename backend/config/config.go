package config

import (
	"os"
	"strconv"
)

type Config struct {
	WebSocket WebSocketConfig
	Kafka     KafkaConfig
	Redis     RedisConfig
	Performance PerformanceConfig
	Scaling   ScalingConfig
	Memory    MemoryConfig
}

type WebSocketConfig struct {
	MaxConnections int
	Port           string
	Token          string
}

type KafkaConfig struct {
	BrokerAddress string
	Topic         string
	GroupID       string
	BatchSize     int
	BatchTimeout  int // milliseconds
	MaxAttempts   int
	WriteTimeout  int // seconds
	ReadTimeout   int // seconds
	MaxMessageSize int // bytes
	RetryInterval int // seconds
	MaxRetries    int
	FailFast      bool
}

type RedisConfig struct {
	URL string
}

type PerformanceConfig struct {
	MaxDisplayFlights int
	ClusterDistance   float64
	MessageBatchSize  int
}

type ScalingConfig struct {
	ScaleUpThreshold   float64
	ScaleDownThreshold float64
	CooldownMinutes    int
	MonitorInterval    int
	RateLimitPerIP     int
	RateLimitWindow    int
}

type MemoryConfig struct {
	WindowMinutes int
	MaxSize       int
	CleanupInterval int
}

func Load() *Config {
	return &Config{
		WebSocket: WebSocketConfig{
			MaxConnections: getEnvInt("WEBSOCKET_MAX_CONNECTIONS", 1000),
			Port:           getEnv("WEBSOCKET_PORT", "8080"),
			Token:          getEnv("WEBSOCKET_TOKEN", ""),
		},
		Kafka: KafkaConfig{
			BrokerAddress: getEnv("KAFKA_BROKER_ADDRESS", "localhost:9092"),
			Topic:         getEnv("KAFKA_TOPIC", "flights"),
			GroupID:       getEnv("KAFKA_GROUP_ID", "flight-websocket-group"),
			BatchSize:     getEnvInt("KAFKA_BATCH_SIZE", 1),
			BatchTimeout:  getEnvInt("KAFKA_BATCH_TIMEOUT", 10),
			MaxAttempts:   getEnvInt("KAFKA_MAX_ATTEMPTS", 3),
			WriteTimeout:  getEnvInt("KAFKA_WRITE_TIMEOUT", 10),
			ReadTimeout:   getEnvInt("KAFKA_READ_TIMEOUT", 10),
			MaxMessageSize: getEnvInt("KAFKA_MAX_MESSAGE_SIZE", 10485760),
			RetryInterval:  getEnvInt("KAFKA_RETRY_INTERVAL", 5),
			MaxRetries:     getEnvInt("KAFKA_MAX_RETRIES", 3),
			FailFast:       getEnv("KAFKA_FAIL_FAST", "false") == "true",
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "localhost:6379"),
		},
		Performance: PerformanceConfig{
			MaxDisplayFlights: getEnvInt("MAX_DISPLAY_FLIGHTS", 100),
			ClusterDistance:   getEnvFloat("CLUSTER_DISTANCE", 0.1),
			MessageBatchSize:  getEnvInt("MESSAGE_BATCH_SIZE", 50),
		},
		Scaling: ScalingConfig{
			ScaleUpThreshold:   getEnvFloat("SCALE_UP_THRESHOLD", 0.8),
			ScaleDownThreshold: getEnvFloat("SCALE_DOWN_THRESHOLD", 0.3),
			CooldownMinutes:    getEnvInt("SCALING_COOLDOWN_MINUTES", 5),
			MonitorInterval:    getEnvInt("SCALING_MONITOR_INTERVAL", 30),
			RateLimitPerIP:     getEnvInt("RATE_LIMIT_PER_IP", 5),
			RateLimitWindow:    getEnvInt("RATE_LIMIT_WINDOW_MINUTES", 1),
		},
		Memory: MemoryConfig{
			WindowMinutes:   getEnvInt("MEMORY_WINDOW_MINUTES", 5),
			MaxSize:        getEnvInt("MEMORY_MAX_SIZE", 1000),
			CleanupInterval: getEnvInt("MEMORY_CLEANUP_INTERVAL", 1),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}