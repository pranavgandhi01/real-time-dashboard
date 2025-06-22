package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"../observability"
)

func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)
	
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		defer span.End()
		
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)
		
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		
		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
	}
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		
		duration := time.Since(start).Seconds()
		status := string(rune(c.Writer.Status()))
		
		observability.RequestDuration.WithLabelValues(
			c.Request.Method, c.FullPath(), status,
		).Observe(duration)
		
		observability.RequestsTotal.WithLabelValues(
			c.Request.Method, c.FullPath(), status,
		).Inc()
	}
}