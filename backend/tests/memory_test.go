package tests

import (
	"testing"
	"time"
	"real-time-dashboard/memory"
)

func TestSlidingWindow_Add(t *testing.T) {
	sw := memory.NewSlidingWindow(5, 10, 1) // 5 minutes, max 10 items, 1 minute cleanup
	
	// Add test data
	testData := []byte("test flight data")
	sw.Add(testData)
	
	recent := sw.GetRecent()
	if len(recent) != 1 {
		t.Errorf("Expected 1 recent item, got %d", len(recent))
	}
	
	if string(recent[0]) != string(testData) {
		t.Errorf("Expected %s, got %s", string(testData), string(recent[0]))
	}
}

func TestSlidingWindow_MaxSize(t *testing.T) {
	sw := memory.NewSlidingWindow(5, 3, 1) // Max 3 items
	
	// Add more than max size
	for i := 0; i < 5; i++ {
		sw.Add([]byte("data" + string(rune('0'+i))))
	}
	
	recent := sw.GetRecent()
	if len(recent) > 3 {
		t.Errorf("Expected max 3 items, got %d", len(recent))
	}
}

func TestSlidingWindow_TimeWindow(t *testing.T) {
	sw := memory.NewSlidingWindow(0, 10, 1) // 0 minute window (immediate expiry)
	
	sw.Add([]byte("test data"))
	
	// Wait a bit for expiry
	time.Sleep(10 * time.Millisecond)
	
	recent := sw.GetRecent()
	if len(recent) != 0 {
		t.Errorf("Expected 0 items after expiry, got %d", len(recent))
	}
}

func TestSlidingWindow_Concurrent(t *testing.T) {
	sw := memory.NewSlidingWindow(5, 100, 1)
	
	// Concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				sw.Add([]byte("data"))
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	recent := sw.GetRecent()
	if len(recent) == 0 {
		t.Error("Expected some data after concurrent writes")
	}
}