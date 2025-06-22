package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"github.com/gin-gonic/gin"
	"../../../pkg/config"
	"../../../pkg/middleware"
	"../../../pkg/health"
	"../../../pkg/log"
)

type APIGateway struct {
	flightDataURL string
	websocketURL  string
	rateLimiter   *middleware.RateLimiter
}

func NewAPIGateway(cfg *config.Config) *APIGateway {
	return &APIGateway{
		flightDataURL: "http://flight-data-service:8081",
		websocketURL:  "http://websocket-service:8082",
		rateLimiter:   middleware.NewRateLimiter(cfg.RateLimitPerIP, time.Minute),
	}
}

func (gw *APIGateway) proxyToFlightData(c *gin.Context) {
	target, _ := url.Parse(gw.flightDataURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (gw *APIGateway) proxyToWebSocket(c *gin.Context) {
	target, _ := url.Parse(gw.websocketURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	cfg := config.Load()
	gateway := NewAPIGateway(cfg)
	r := gin.Default()
	
	// Apply rate limiting
	r.Use(gateway.rateLimiter.Middleware())
	
	r.GET("/health", gin.WrapF(health.HealthHandler))

	// Route to flight data service
	r.Any("/api/flights/*path", gateway.proxyToFlightData)
	r.Any("/flights/*path", gateway.proxyToFlightData)
	r.Any("/stats", gateway.proxyToFlightData)

	// Route to websocket service
	r.Any("/ws", gateway.proxyToWebSocket)

	log.LogInfo("API Gateway starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}