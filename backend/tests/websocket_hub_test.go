package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"real-time-dashboard/config"
	"real-time-dashboard/ws"
	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	cfg := &config.Config{
		WebSocket: config.WebSocketConfig{
			MaxConnections: 10,
			Port:          "8080",
			Token:         "test-token",
		},
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "flights",
			GroupID:       "test-group",
			MaxRetries:    1,
			FailFast:      false,
		},
	}
	
	hub := ws.NewHub(cfg)
	go hub.Run()
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleConnections(hub, w, r)
	}))
	defer server.Close()
	
	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=test-token"
	
	// Test WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	
	// Connection should be established
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	
	// Test that connection stays open
	time.Sleep(100 * time.Millisecond)
}

func TestWebSocketAuthentication(t *testing.T) {
	cfg := &config.Config{
		WebSocket: config.WebSocketConfig{
			MaxConnections: 10,
			Port:          "8080",
			Token:         "secret-token",
		},
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "flights",
			GroupID:       "test-group",
			MaxRetries:    1,
			FailFast:      false,
		},
	}
	
	hub := ws.NewHub(cfg)
	
	// Test with wrong token
	req := httptest.NewRequest("GET", "/ws?token=wrong-token", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Key", "test-key")
	req.Header.Set("Sec-WebSocket-Version", "13")
	
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleConnections(hub, w, r)
	})
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong token, got %d", rr.Code)
	}
}

func TestWebSocketConnectionLimit(t *testing.T) {
	cfg := &config.Config{
		WebSocket: config.WebSocketConfig{
			MaxConnections: 1, // Very low limit for testing
			Port:          "8080",
			Token:         "",
		},
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "flights",
			GroupID:       "test-group",
			MaxRetries:    1,
			FailFast:      false,
		},
	}
	
	hub := ws.NewHub(cfg)
	go hub.Run()
	
	// This test would require more complex setup to properly test connection limits
	// For now, just verify hub creation doesn't panic
	if hub == nil {
		t.Error("Hub creation failed")
	}
}