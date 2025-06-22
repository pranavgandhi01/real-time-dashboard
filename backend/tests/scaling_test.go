package tests

import (
	"testing"
	"real-time-dashboard/scaling"
)

func TestAutoScaler_ScaleUp(t *testing.T) {
	as := scaling.NewAutoScaler(100, 0.8, 0.3, 1, 1) // 80% up, 30% down, 1min cooldown
	
	// Simulate high connection load
	as.UpdateConnections(85) // 85% load
	
	shouldScale, action := as.ShouldScale()
	if !shouldScale {
		t.Error("Expected scale-up trigger at 85% load")
	}
	
	if action != "scale-up" {
		t.Errorf("Expected scale-up action, got %s", action)
	}
}

func TestAutoScaler_ScaleDown(t *testing.T) {
	as := scaling.NewAutoScaler(100, 0.8, 0.3, 1, 1)
	
	// Simulate low connection load
	as.UpdateConnections(25) // 25% load
	
	shouldScale, action := as.ShouldScale()
	if !shouldScale {
		t.Error("Expected scale-down trigger at 25% load")
	}
	
	if action != "scale-down" {
		t.Errorf("Expected scale-down action, got %s", action)
	}
}

func TestAutoScaler_Stable(t *testing.T) {
	as := scaling.NewAutoScaler(100, 0.8, 0.3, 1, 1)
	
	// Simulate stable load
	as.UpdateConnections(50) // 50% load
	
	shouldScale, action := as.ShouldScale()
	if shouldScale {
		t.Error("Expected no scaling at 50% load")
	}
	
	if action != "stable" {
		t.Errorf("Expected stable action, got %s", action)
	}
}

func TestAutoScaler_Cooldown(t *testing.T) {
	as := scaling.NewAutoScaler(100, 0.8, 0.3, 1, 1) // 1 minute cooldown
	
	// Trigger scaling
	as.UpdateConnections(85)
	shouldScale, _ := as.ShouldScale()
	if !shouldScale {
		t.Error("Expected initial scale trigger")
	}
	
	// Record scale action
	as.RecordScaleAction()
	
	// Immediate check should be in cooldown
	shouldScale, action := as.ShouldScale()
	if shouldScale {
		t.Error("Expected cooldown to prevent immediate scaling")
	}
	
	if action != "cooldown" {
		t.Errorf("Expected cooldown action, got %s", action)
	}
}

func TestAutoScaler_Metrics(t *testing.T) {
	as := scaling.NewAutoScaler(100, 0.8, 0.3, 1, 1)
	
	as.UpdateConnections(75)
	as.UpdateQueueDepth(150)
	
	metrics := as.GetMetrics()
	
	if metrics.ConnectionRatio != 0.75 {
		t.Errorf("Expected connection ratio 0.75, got %f", metrics.ConnectionRatio)
	}
	
	if metrics.MessageQueueDepth != 150 {
		t.Errorf("Expected queue depth 150, got %d", metrics.MessageQueueDepth)
	}
	
	if metrics.LastUpdated.IsZero() {
		t.Error("Expected LastUpdated to be set")
	}
}