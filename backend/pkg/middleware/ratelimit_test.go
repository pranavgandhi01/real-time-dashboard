package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	rl := NewRateLimiter(2, time.Minute)
	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	
	// First request should pass
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.Header.Set("X-Forwarded-For", "192.168.1.1")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	
	if w1.Code != 200 {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}
	
	// Second request should pass
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.1.1")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	
	if w2.Code != 200 {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	
	// Third request should be rate limited
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.Header.Set("X-Forwarded-For", "192.168.1.1")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	
	if w3.Code != 429 {
		t.Errorf("Expected status 429, got %d", w3.Code)
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	rl := NewRateLimiter(1, time.Minute)
	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	
	// Request from IP 1
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.Header.Set("X-Forwarded-For", "192.168.1.1")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	
	if w1.Code != 200 {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}
	
	// Request from IP 2 should still pass
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.1.2")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	
	if w2.Code != 200 {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
}