// backend/log/log.go
package log // Change package to 'log'

import (
	"log"
	"os"
	"strconv"
	"strings" // Import strings for ToLower
)

// LogLevel defines the logging verbosity.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota // 0
	LogLevelInfo                  // 1
	LogLevelWarn                  // 2
	LogLevelError                 // 3
	LogLevelFatal                 // 4
)

// currentLogLevel is the global variable that controls the minimum log level to print.
var currentLogLevel LogLevel = LogLevelInfo // Default to INFO level

func init() {
	// Set log flags for date, time, and file/line number
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Read log level from environment variable LOG_LEVEL
	logLevelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(logLevelStr) {
	case "debug":
		currentLogLevel = LogLevelDebug
	case "info":
		currentLogLevel = LogLevelInfo
	case "warn":
		currentLogLevel = LogLevelWarn
	case "error":
		currentLogLevel = LogLevelError
	case "fatal":
		currentLogLevel = LogLevelFatal
	default:
		// Try parsing as integer for backward compatibility or direct level setting
		if level, err := strconv.Atoi(logLevelStr); err == nil {
			if level >= 0 && level <= int(LogLevelFatal) {
				currentLogLevel = LogLevel(level)
			} else {
				log.Printf("WARN: Invalid LOG_LEVEL environment variable '%s', defaulting to INFO.", logLevelStr)
			}
		} else if logLevelStr != "" { // If not empty and not a valid level/int
			log.Printf("WARN: Unknown LOG_LEVEL environment variable '%s', defaulting to INFO.", logLevelStr)
		}
	}
	log.Printf("INFO: Logging level set to %v (Enum Value: %d)", currentLogLevel.String(), currentLogLevel)
}

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogDebug prints a debug message if the currentLogLevel allows.
func LogDebug(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelDebug {
		log.Printf("DEBUG: "+format, v...)
	}
}

// LogInfo prints an info message if the currentLogLevel allows.
func LogInfo(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelInfo {
		log.Printf("INFO: "+format, v...)
	}
}

// LogWarn prints a warning message if the currentLogLevel allows.
func LogWarn(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelWarn {
		log.Printf("WARN: "+format, v...)
	}
}

// LogError prints an error message if the currentLogLevel allows.
func LogError(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelError {
		log.Printf("ERROR: "+format, v...)
	}
}

// LogFatal prints a fatal message and exits the program.
func LogFatal(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelFatal {
		log.Fatalf("FATAL: "+format, v...)
	}
}