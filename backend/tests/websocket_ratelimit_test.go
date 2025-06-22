package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"real-time-dashboard/ratelimit"
)

func TestWebSocketRateLimit(t *testing.T) {
	rateLimiter := ratelimit.NewRateLimiter(2, time.Minute)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Connection allowed"))
	})

	// Test first connection - should be allowed
	req1 := httptest.NewRequest("GET", "/ws", nil)
	req1.RemoteAddr = "192.168.1.1:8080"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	
	if rr1.Code != http.StatusOK {
		t.Errorf("Expected first connection to be allowed, got status %d", rr1.Code)
	}

	// Test second connection - should be allowed
	req2 := httptest.NewRequest("GET", "/ws", nil)
	req2.RemoteAddr = "192.168.1.1:8081"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	
	if rr2.Code != http.StatusOK {
		t.Errorf("Expected second connection to be allowed, got status %d", rr2.Code)
	}

	// Test third connection - should be rate limited
	req3 := httptest.NewRequest("GET", "/ws", nil)
	req3.RemoteAddr = "192.168.1.1:8082"
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	
	if rr3.Code != http.StatusTooManyRequests {
		t.Errorf("Expected third connection to be rate limited, got status %d", rr3.Code)
	}
}

func TestWebSocketRateLimitDifferentIPs(t *testing.T) {
	rateLimiter := ratelimit.NewRateLimiter(1, time.Minute)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// First IP - should be allowed
	req1 := httptest.NewRequest("GET", "/ws", nil)
	req1.RemoteAddr = "192.168.1.1:8080"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	
	if rr1.Code != http.StatusOK {
		t.Errorf("Expected connection from first IP to be allowed, got status %d", rr1.Code)
	}

	// Different IP - should be allowed
	req2 := httptest.NewRequest("GET", "/ws", nil)
	req2.RemoteAddr = "192.168.1.2:8080"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	
	if rr2.Code != http.StatusOK {
		t.Errorf("Expected connection from different IP to be allowed, got status %d", rr2.Code)
	}
}