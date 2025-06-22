package tests

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	flightlog "real-time-dashboard/log"
)

func TestLoggingLevels(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	
	// Test different log levels
	flightlog.LogInfo("Test info message")
	flightlog.LogWarn("Test warning message")
	flightlog.LogError("Test error message")
	flightlog.LogDebug("Test debug message")
	
	output := buf.String()
	
	// Check that messages are logged
	if !strings.Contains(output, "INFO") {
		t.Error("INFO level not found in log output")
	}
	
	if !strings.Contains(output, "WARN") {
		t.Error("WARN level not found in log output")
	}
	
	if !strings.Contains(output, "ERROR") {
		t.Error("ERROR level not found in log output")
	}
}

func TestLogFormatting(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	
	// Test formatted logging
	flightlog.LogInfo("Test message with %s and %d", "string", 42)
	
	output := buf.String()
	
	if !strings.Contains(output, "Test message with string and 42") {
		t.Error("Formatted message not found in log output")
	}
}

func TestLogLevelFiltering(t *testing.T) {
	// Test that debug messages are filtered based on log level
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	
	// This would test log level filtering if implemented
	flightlog.LogDebug("Debug message that might be filtered")
	
	output := buf.String()
	
	// In production, debug messages might be filtered
	// This test documents the expected behavior
	t.Logf("Debug log output: %s", output)
}