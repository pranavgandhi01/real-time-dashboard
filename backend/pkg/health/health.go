package health

import (
	"encoding/json"
	"net/http"
	"time"
	"../log"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

var startTime = time.Now()

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	health := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services: map[string]string{
			"service": "active",
		},
		Uptime: time.Since(startTime).String(),
	}
	
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.LogError("Failed to encode health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Simple readiness check
	response := map[string]string{
		"status": "ready",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.LogError("Failed to encode readiness response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}