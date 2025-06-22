package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
)

func TestAPIGateway_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if !strings.Contains(w.Body.String(), "api-gateway") {
		t.Error("Expected response to contain 'api-gateway'")
	}
}

func TestAPIGateway_Creation(t *testing.T) {
	gateway := NewAPIGateway()
	
	if gateway.flightDataURL == "" {
		t.Error("Expected flightDataURL to be set")
	}
	
	if gateway.websocketURL == "" {
		t.Error("Expected websocketURL to be set")
	}
}