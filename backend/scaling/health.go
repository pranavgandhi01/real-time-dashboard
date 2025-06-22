package scaling

import (
	"sync"
	"time"
	flightlog "real-time-dashboard/log"
)

type HealthMetrics struct {
	ConnectionRatio    float64
	MemoryUsageRatio   float64
	MessageQueueDepth  int
	LastUpdated        time.Time
}

type AutoScaler struct {
	maxConnections     int
	currentConnections int
	scaleUpThreshold   float64
	scaleDownThreshold float64
	mu                 sync.RWMutex
	metrics            HealthMetrics
	scalingCooldown    time.Duration
	lastScaleAction    time.Time
}

func NewAutoScaler(maxConnections int, scaleUpThreshold, scaleDownThreshold float64, cooldownMinutes, monitorIntervalSeconds int) *AutoScaler {
	as := &AutoScaler{
		maxConnections:     maxConnections,
		scaleUpThreshold:   scaleUpThreshold,
		scaleDownThreshold: scaleDownThreshold,
		scalingCooldown:    time.Duration(cooldownMinutes) * time.Minute,
	}
	
	go as.monitor(monitorIntervalSeconds)
	return as
}

func (as *AutoScaler) UpdateConnections(current int) {
	as.mu.Lock()
	defer as.mu.Unlock()
	
	as.currentConnections = current
	as.metrics.ConnectionRatio = float64(current) / float64(as.maxConnections)
	as.metrics.LastUpdated = time.Now()
}

func (as *AutoScaler) UpdateQueueDepth(depth int) {
	as.mu.Lock()
	defer as.mu.Unlock()
	
	as.metrics.MessageQueueDepth = depth
}

func (as *AutoScaler) ShouldScale() (bool, string) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	
	// Cooldown check
	if time.Since(as.lastScaleAction) < as.scalingCooldown {
		return false, "cooldown"
	}
	
	ratio := as.metrics.ConnectionRatio
	
	if ratio > as.scaleUpThreshold {
		return true, "scale-up"
	}
	
	if ratio < as.scaleDownThreshold {
		return true, "scale-down"
	}
	
	return false, "stable"
}

func (as *AutoScaler) GetMetrics() HealthMetrics {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.metrics
}

func (as *AutoScaler) RecordScaleAction() {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.lastScaleAction = time.Now()
}

func (as *AutoScaler) monitor(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		shouldScale, action := as.ShouldScale()
		if shouldScale {
			flightlog.LogInfo("Auto-scaling trigger: %s (ratio: %.2f, connections: %d/%d)", 
				action, as.metrics.ConnectionRatio, as.currentConnections, as.maxConnections)
			
			// In production, this would trigger container orchestration
			as.triggerScaling(action)
		}
	}
}

func (as *AutoScaler) triggerScaling(action string) {
	// Placeholder for actual scaling logic
	// In Kubernetes: kubectl scale deployment flight-tracker --replicas=N
	// In Docker Swarm: docker service scale flight-tracker=N
	
	flightlog.LogInfo("Scaling action triggered: %s", action)
	as.RecordScaleAction()
}