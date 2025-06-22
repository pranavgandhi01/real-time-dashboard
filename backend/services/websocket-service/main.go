package main

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"../../../pkg/config"
	"../../../pkg/health"
	"../../../pkg/log"
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
	r := gin.Default()
	
	r.GET("/health", gin.WrapF(health.HealthHandler))
	r.GET("/ws", wsService.HandleWebSocket)
	r.GET("/metrics", wsService.GetMetrics)

	log.LogInfo("WebSocket Service starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}