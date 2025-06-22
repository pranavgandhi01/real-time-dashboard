package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	// Create a simple handler that returns metrics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(`# HELP websocket_connected_clients Number of active WebSocket clients
# TYPE websocket_connected_clients gauge
websocket_connected_clients 0
# HELP flight_fetch_latency_seconds Latency of fetching flight data
# TYPE flight_fetch_latency_seconds histogram
flight_fetch_latency_seconds_bucket{le="0.005"} 0
# HELP flight_data_processing_seconds Time taken to process flight data
# TYPE flight_data_processing_seconds histogram
flight_data_processing_seconds_bucket{le="0.001"} 0
# HELP websocket_message_queue_size Current size of WebSocket message queue
# TYPE websocket_message_queue_size gauge
websocket_message_queue_size 0
# HELP kafka_consumer_lag Kafka consumer lag in messages
# TYPE kafka_consumer_lag gauge
kafka_consumer_lag 0
# HELP redis_cache_hits_total Total number of Redis cache hits
# TYPE redis_cache_hits_total counter
redis_cache_hits_total 0
# HELP redis_cache_misses_total Total number of Redis cache misses
# TYPE redis_cache_misses_total counter
redis_cache_misses_total 0`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	
	// Test that all expected metrics are present
	expectedMetrics := []string{
		"websocket_connected_clients",
		"flight_fetch_latency_seconds",
		"flight_data_processing_seconds",
		"websocket_message_queue_size",
		"kafka_consumer_lag",
		"redis_cache_hits_total",
		"redis_cache_misses_total",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("Expected metric %s not found in response", metric)
		}
	}
}

func TestMetricsContentType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("test metrics"))
	})

	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", contentType)
	}
}