package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	
	// Initialize tracing
	tp, err := observability.InitTracing("api-gateway", "http://jaeger:14268/api/traces")
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
	r.Use(middleware.TracingMiddleware("api-gateway"))
	r.Use(middleware.MetricsMiddleware())
	r.Use(gateway.rateLimiter.Middleware())
	
	r.GET("/health", gin.WrapF(health.HealthHandler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Route to flight data service
	r.Any("/api/flights/*path", gateway.proxyToFlightData)
	r.Any("/flights/*path", gateway.proxyToFlightData)
	r.Any("/stats", gateway.proxyToFlightData)

	// Route to websocket service
	r.Any("/ws", gateway.proxyToWebSocket)

	log.LogInfo("API Gateway starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.LogError("Server failed: %v", err)
	}
}