package main

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/real-time-dashboard/backend/pkg/config"
	"github.com/real-time-dashboard/backend/pkg/health"
	"github.com/real-time-dashboard/backend/pkg/log"
	"github.com/real-time-dashboard/backend/pkg/middleware"
	"github.com/real-time-dashboard/backend/pkg/observability"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSService struct {
	clients map[*websocket.Conn]bool
}

func NewWSService() *WSService {
	return &WSService{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (ws *WSService) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.LogError("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ws.clients[conn] = true
	log.LogInfo("Client connected. Total: %d", len(ws.clients))

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(ws.clients, conn)
			log.LogInfo("Client disconnected. Total: %d", len(ws.clients))
			break
		}
	}
}

func (ws *WSService) GetMetrics(c *gin.Context) {
	c.JSON(200, gin.H{"connections": len(ws.clients)})
}

func main() {
	cfg := config.Load()
	wsService := NewWSService()
	
	// Initialize tracing
	tp, err := observability.InitTracing("websocket-service", "http://jaeger:14268/api/traces")
	if err != nil {
		log.LogError("Failed to initialize tracing: %v", err)
	}
	defer func() {
		if tp != nil {
			tp.Shutdown(context.Background())
		}
	}()
	
	r := gin.Default()
	
	// Apply middleware
	r.Use(middleware.TracingMiddleware("websocket-service"))
	r.Use(middleware.MetricsMiddleware())
	
	r.GET("/health", gin.WrapF(health.HealthHandler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ws", wsService.HandleWebSocket)
	r.GET("/ws-metrics", wsService.GetMetrics)

	log.LogInfo("WebSocket Service starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}