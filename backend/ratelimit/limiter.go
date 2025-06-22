package ratelimit

import (
	"net"
	"sync"
	"time"
)

type RateLimiter struct {
	connections map[string][]time.Time
	maxPerIP    int
	window      time.Duration
	mu          sync.RWMutex
}

func NewRateLimiter(maxPerIP int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		connections: make(map[string][]time.Time),
		maxPerIP:    maxPerIP,
		window:      window,
	}
	
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(remoteAddr string) bool {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ip = remoteAddr
	}
	
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	connections := rl.connections[ip]
	validConnections := make([]time.Time, 0)
	for _, t := range connections {
		if t.After(cutoff) {
			validConnections = append(validConnections, t)
		}
	}
	
	if len(validConnections) >= rl.maxPerIP {
		return false
	}
	
	validConnections = append(validConnections, now)
	rl.connections[ip] = validConnections
	
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		
		for ip, connections := range rl.connections {
			validConnections := make([]time.Time, 0)
			for _, t := range connections {
				if t.After(cutoff) {
					validConnections = append(validConnections, t)
				}
			}
			
			if len(validConnections) == 0 {
				delete(rl.connections, ip)
			} else {
				rl.connections[ip] = validConnections
			}
		}
		rl.mu.Unlock()
	}
}