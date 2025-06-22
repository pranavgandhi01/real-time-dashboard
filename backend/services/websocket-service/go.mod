module websocket-service

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/websocket v1.5.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	github.com/prometheus/client_golang v1.17.0
)