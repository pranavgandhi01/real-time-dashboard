package tests

import (
	"os"
	"testing"
	"real-time-dashboard/config"
)

func TestConfigLoad(t *testing.T) {
	// Test default values
	cfg := config.Load()
	
	if cfg.WebSocket.MaxConnections != 1000 {
		t.Errorf("Expected default MaxConnections 1000, got %d", cfg.WebSocket.MaxConnections)
	}
	
	if cfg.WebSocket.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.WebSocket.Port)
	}
	
	if cfg.Kafka.BrokerAddress != "localhost:9092" {
		t.Errorf("Expected default broker localhost:9092, got %s", cfg.Kafka.BrokerAddress)
	}
	
	// Test scaling defaults
	if cfg.Scaling.ScaleUpThreshold != 0.8 {
		t.Errorf("Expected default scale up threshold 0.8, got %f", cfg.Scaling.ScaleUpThreshold)
	}
	
	if cfg.Scaling.RateLimitPerIP != 5 {
		t.Errorf("Expected default rate limit 5, got %d", cfg.Scaling.RateLimitPerIP)
	}
	
	// Test memory defaults
	if cfg.Memory.WindowMinutes != 5 {
		t.Errorf("Expected default memory window 5 minutes, got %d", cfg.Memory.WindowMinutes)
	}
	
	if cfg.Memory.MaxSize != 1000 {
		t.Errorf("Expected default memory max size 1000, got %d", cfg.Memory.MaxSize)
	}
}

func TestConfigEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("WEBSOCKET_MAX_CONNECTIONS", "500")
	os.Setenv("WEBSOCKET_PORT", "9090")
	os.Setenv("KAFKA_BROKER_ADDRESS", "kafka:9092")
	os.Setenv("SCALE_UP_THRESHOLD", "0.9")
	os.Setenv("MEMORY_WINDOW_MINUTES", "10")
	
	defer func() {
		os.Unsetenv("WEBSOCKET_MAX_CONNECTIONS")
		os.Unsetenv("WEBSOCKET_PORT")
		os.Unsetenv("KAFKA_BROKER_ADDRESS")
		os.Unsetenv("SCALE_UP_THRESHOLD")
		os.Unsetenv("MEMORY_WINDOW_MINUTES")
	}()
	
	cfg := config.Load()
	
	if cfg.WebSocket.MaxConnections != 500 {
		t.Errorf("Expected MaxConnections 500, got %d", cfg.WebSocket.MaxConnections)
	}
	
	if cfg.WebSocket.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.WebSocket.Port)
	}
	
	if cfg.Kafka.BrokerAddress != "kafka:9092" {
		t.Errorf("Expected broker kafka:9092, got %s", cfg.Kafka.BrokerAddress)
	}
	
	if cfg.Scaling.ScaleUpThreshold != 0.9 {
		t.Errorf("Expected scale up threshold 0.9, got %f", cfg.Scaling.ScaleUpThreshold)
	}
	
	if cfg.Memory.WindowMinutes != 10 {
		t.Errorf("Expected memory window 10 minutes, got %d", cfg.Memory.WindowMinutes)
	}
}

func TestConfigInvalidValues(t *testing.T) {
	// Test invalid integer values fall back to defaults
	os.Setenv("WEBSOCKET_MAX_CONNECTIONS", "invalid")
	defer os.Unsetenv("WEBSOCKET_MAX_CONNECTIONS")
	
	cfg := config.Load()
	
	if cfg.WebSocket.MaxConnections != 1000 {
		t.Errorf("Expected fallback to default 1000, got %d", cfg.WebSocket.MaxConnections)
	}
}