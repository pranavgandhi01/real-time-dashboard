package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestWSService_GetMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	ws := NewWSService()
	
	r := gin.New()
	r.GET("/metrics", ws.GetMetrics)
	
	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if !strings.Contains(w.Body.String(), "connections") {
		t.Error("Expected response to contain 'connections'")
	}
}

func TestWSService_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "websocket"})
	})
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if !strings.Contains(w.Body.String(), "healthy") {
		t.Error("Expected response to contain 'healthy'")
	}
}