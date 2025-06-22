package ratelimit

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	
	// Test allowing connections within limit
	if !rl.Allow("192.168.1.1:8080") {
		t.Error("Expected first connection to be allowed")
	}
	
	if !rl.Allow("192.168.1.1:8081") {
		t.Error("Expected second connection to be allowed")
	}
	
	// Test rate limiting
	if rl.Allow("192.168.1.1:8082") {
		t.Error("Expected third connection to be rate limited")
	}
	
	// Test different IP is allowed
	if !rl.Allow("192.168.1.2:8080") {
		t.Error("Expected connection from different IP to be allowed")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := NewRateLimiter(1, 100*time.Millisecond)
	
	// Use up the limit
	if !rl.Allow("192.168.1.1:8080") {
		t.Error("Expected first connection to be allowed")
	}
	
	// Should be rate limited
	if rl.Allow("192.168.1.1:8081") {
		t.Error("Expected second connection to be rate limited")
	}
	
	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)
	
	// Should be allowed again
	if !rl.Allow("192.168.1.1:8082") {
		t.Error("Expected connection after window expiry to be allowed")
	}
}

func TestRateLimiter_IPParsing(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	
	// Test with port
	if !rl.Allow("192.168.1.1:8080") {
		t.Error("Expected connection with port to be allowed")
	}
	
	// Same IP different port should be rate limited
	if rl.Allow("192.168.1.1:9090") {
		t.Error("Expected same IP different port to be rate limited")
	}
	
	// Test without port
	if !rl.Allow("192.168.1.2") {
		t.Error("Expected connection without port to be allowed")
	}
}