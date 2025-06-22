package memory

import (
	"sync"
	"time"
)

type TimestampedData struct {
	Data      []byte
	Timestamp time.Time
}

type SlidingWindow struct {
	data     []TimestampedData
	mu       sync.RWMutex
	window   time.Duration
	maxSize  int
}

func NewSlidingWindow(windowMinutes, maxSize, cleanupIntervalMinutes int) *SlidingWindow {
	window := time.Duration(windowMinutes) * time.Minute
	sw := &SlidingWindow{
		data:    make([]TimestampedData, 0, maxSize),
		window:  window,
		maxSize: maxSize,
	}
	go sw.cleanup(cleanupIntervalMinutes)
	return sw
}

func (sw *SlidingWindow) Add(data []byte) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	
	now := time.Now()
	
	// Add new data
	sw.data = append(sw.data, TimestampedData{
		Data:      data,
		Timestamp: now,
	})
	
	// Size-based eviction
	if len(sw.data) > sw.maxSize {
		sw.data = sw.data[len(sw.data)-sw.maxSize:]
	}
}

func (sw *SlidingWindow) GetRecent() [][]byte {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	
	cutoff := time.Now().Add(-sw.window)
	recent := make([][]byte, 0)
	
	for _, item := range sw.data {
		if item.Timestamp.After(cutoff) {
			recent = append(recent, item.Data)
		}
	}
	
	return recent
}

func (sw *SlidingWindow) cleanup(cleanupIntervalMinutes int) {
	ticker := time.NewTicker(time.Duration(cleanupIntervalMinutes) * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		sw.mu.Lock()
		cutoff := time.Now().Add(-sw.window)
		
		// Remove expired data
		validData := make([]TimestampedData, 0, len(sw.data))
		for _, item := range sw.data {
			if item.Timestamp.After(cutoff) {
				validData = append(validData, item)
			}
		}
		
		sw.data = validData
		sw.mu.Unlock()
	}
}