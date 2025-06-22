package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"go.opentelemetry.io/otel/trace"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	SpanID    string `json:"span_id,omitempty"`
}

var (
	currentLogLevel LogLevel = LogLevelInfo
	serviceName     string
	useJSON         bool
)

func init() {
	serviceName = os.Getenv("SERVICE_NAME")
	useJSON = os.Getenv("LOG_FORMAT") == "json"
	
	if !useJSON {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

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
		if level, err := strconv.Atoi(logLevelStr); err == nil {
			if level >= 0 && level <= int(LogLevelFatal) {
				currentLogLevel = LogLevel(level)
			}
		}
	}
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

func logWithContext(ctx context.Context, level LogLevel, format string, v ...interface{}) {
	if currentLogLevel > level {
		return
	}
	
	if useJSON {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     level.String(),
			Message:   fmt.Sprintf(format, v...),
			Service:   serviceName,
		}
		
		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			entry.TraceID = span.SpanContext().TraceID().String()
			entry.SpanID = span.SpanContext().SpanID().String()
		}
		
		if data, err := json.Marshal(entry); err == nil {
			log.Println(string(data))
		}
	} else {
		log.Printf("%s: "+format, append([]interface{}{level.String()}, v...)...)
	}
}

func LogDebug(format string, v ...interface{}) {
	logWithContext(context.Background(), LogLevelDebug, format, v...)
}

func LogInfo(format string, v ...interface{}) {
	logWithContext(context.Background(), LogLevelInfo, format, v...)
}

func LogWarn(format string, v ...interface{}) {
	logWithContext(context.Background(), LogLevelWarn, format, v...)
}

func LogError(format string, v ...interface{}) {
	logWithContext(context.Background(), LogLevelError, format, v...)
}

func LogFatal(format string, v ...interface{}) {
	logWithContext(context.Background(), LogLevelFatal, format, v...)
	if currentLogLevel <= LogLevelFatal {
		os.Exit(1)
	}
}

// Context-aware logging functions
func LogDebugCtx(ctx context.Context, format string, v ...interface{}) {
	logWithContext(ctx, LogLevelDebug, format, v...)
}

func LogInfoCtx(ctx context.Context, format string, v ...interface{}) {
	logWithContext(ctx, LogLevelInfo, format, v...)
}

func LogErrorCtx(ctx context.Context, format string, v ...interface{}) {
	logWithContext(ctx, LogLevelError, format, v...)
}