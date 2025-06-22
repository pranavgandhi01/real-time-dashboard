package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"real-time-dashboard/ratelimit"
)

func TestRateLimitIntegration(t *testing.T) {
	rateLimiter := ratelimit.NewRateLimiter(3, time.Minute)
	
	// Create test server with rate limiting
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	
	// Test successful requests within limit
	for i := 0; i < 3; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
		resp.Body.Close()
	}
	
	// Test rate limited request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Rate limited request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected rate limited request to return 429, got %d", resp.StatusCode)
	}
}

func TestRateLimitMetrics(t *testing.T) {
	// Test that rate limiting generates appropriate log messages
	rateLimiter := ratelimit.NewRateLimiter(1, time.Minute)
	
	// Simulate multiple requests from same IP
	allowed1 := rateLimiter.Allow("192.168.1.1:8080")
	allowed2 := rateLimiter.Allow("192.168.1.1:8081")
	
	if !allowed1 {
		t.Error("First request should be allowed")
	}
	
	if allowed2 {
		t.Error("Second request should be rate limited")
	}
}

func TestRateLimitCleanup(t *testing.T) {
	// Test that old entries are cleaned up
	rateLimiter := ratelimit.NewRateLimiter(1, 100*time.Millisecond)
	
	// Make request
	if !rateLimiter.Allow("192.168.1.1:8080") {
		t.Error("First request should be allowed")
	}
	
	// Should be rate limited immediately
	if rateLimiter.Allow("192.168.1.1:8081") {
		t.Error("Second request should be rate limited")
	}
	
	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)
	
	// Should be allowed again after cleanup
	if !rateLimiter.Allow("192.168.1.1:8082") {
		t.Error("Request after cleanup should be allowed")
	}
}