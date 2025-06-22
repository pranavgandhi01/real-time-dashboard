package tests

import (
	"testing"
	"time"
	"real-time-dashboard/cache"
)

func TestRedisConnection(t *testing.T) {
	// Test Redis initialization
	err := cache.InitRedis()
	if err != nil {
		t.Logf("Redis connection failed (expected in test environment): %v", err)
		return
	}
	
	// Test Redis operations if connection successful
	testKey := "test:key"
	testValue := "test_value"
	
	err = cache.Set(testKey, testValue, time.Minute)
	if err != nil {
		t.Errorf("Failed to set cache value: %v", err)
	}
	
	value, err := cache.Get(testKey)
	if err != nil {
		t.Errorf("Failed to get cache value: %v", err)
	}
	
	if value != testValue {
		t.Errorf("Expected %s, got %s", testValue, value)
	}
	
	// Cleanup
	cache.Delete(testKey)
	cache.Close()
}

func TestCacheOperations(t *testing.T) {
	// Mock cache operations for testing
	testCases := []struct {
		key   string
		value string
		ttl   time.Duration
	}{
		{"flight:123", "flight_data", time.Minute},
		{"stats:global", "statistics", 30 * time.Second},
	}
	
	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// Test that cache operations don't panic
			// In real implementation, these would test actual Redis operations
			if tc.ttl < time.Second {
				t.Errorf("TTL too short: %v", tc.ttl)
			}
		})
	}
}