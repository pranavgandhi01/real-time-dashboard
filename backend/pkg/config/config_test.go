package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test default values
	cfg := Load()
	
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}
	
	if cfg.MaxConnections != 1000 {
		t.Errorf("Expected default max connections 1000, got %d", cfg.MaxConnections)
	}
	
	if cfg.FetchInterval != 15*time.Second {
		t.Errorf("Expected default fetch interval 15s, got %v", cfg.FetchInterval)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("MAX_CONNECTIONS", "2000")
	os.Setenv("FETCH_INTERVAL", "30s")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("MAX_CONNECTIONS")
		os.Unsetenv("FETCH_INTERVAL")
	}()
	
	cfg := Load()
	
	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}
	
	if cfg.MaxConnections != 2000 {
		t.Errorf("Expected max connections 2000, got %d", cfg.MaxConnections)
	}
	
	if cfg.FetchInterval != 30*time.Second {
		t.Errorf("Expected fetch interval 30s, got %v", cfg.FetchInterval)
	}
}