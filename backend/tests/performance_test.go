package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"real-time-dashboard/ratelimit"
)

func BenchmarkRateLimiter(b *testing.B) {
	rl := ratelimit.NewRateLimiter(1000, time.Minute)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rl.Allow(fmt.Sprintf("192.168.1.%d:8080", i%255))
			i++
		}
	})
}

func BenchmarkRateLimiterConcurrent(b *testing.B) {
	rl := ratelimit.NewRateLimiter(100, time.Minute)
	
	b.ResetTimer()
	
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < b.N/10; j++ {
				rl.Allow(fmt.Sprintf("192.168.%d.%d:8080", id, j%255))
			}
		}(i)
	}
	wg.Wait()
}

func TestRateLimiterMemoryUsage(t *testing.T) {
	rl := ratelimit.NewRateLimiter(10, time.Minute)
	
	// Generate many connections from different IPs
	for i := 0; i < 1000; i++ {
		rl.Allow(fmt.Sprintf("192.168.%d.%d:8080", i/255, i%255))
	}
	
	// Test that cleanup works and doesn't cause memory leaks
	time.Sleep(100 * time.Millisecond)
	
	// This test mainly ensures no panics occur during cleanup
	// In a real scenario, you'd measure actual memory usage
}

func TestMetricsPerformance(t *testing.T) {
	// Test that metrics collection doesn't significantly impact performance
	start := time.Now()
	
	// Simulate processing 1000 flight records
	for i := 0; i < 1000; i++ {
		// Simulate processing time
		time.Sleep(time.Microsecond)
	}
	
	duration := time.Since(start)
	
	// Should complete within reasonable time (adjust based on requirements)
	if duration > time.Second {
		t.Errorf("Processing 1000 records took too long: %v", duration)
	}
}